/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubernetesservice

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	"k8s.io/kubernetes/pkg/controlplane/reconcilers"
)

const (
	kubernetesServiceName = "kubernetes"
)

// Controller is the controller manager for the core bootstrap Kubernetes
// controller loops, which manage creating the "kubernetes" service and
// provide the IP repair check on service IPs
type Controller struct {
	Config

	client kubernetes.Interface
}

type Config struct {
	PublicIP net.IP

	EndpointReconciler reconcilers.EndpointReconciler
	EndpointInterval   time.Duration

	// ServiceIP indicates where the kubernetes service will live.  It may not be nil.
	ServiceIP                 net.IP
	ServicePort               int
	PublicServicePort         int
	KubernetesServiceNodePort int
}

// New returns a controller for watching the kubernetes service endpoints.
func New(config Config, client kubernetes.Interface) *Controller {
	return &Controller{
		Config: config,
		client: client,
	}
}

// Start begins the core controller loops that must exist for bootstrapping
// a cluster.
func (c *Controller) Start(stopCh <-chan struct{}) {
	// Reconcile during first run removing itself until server is ready.
	endpointPorts := createEndpointPortSpec(c.PublicServicePort, "https")
	if err := c.EndpointReconciler.RemoveEndpoints(kubernetesServiceName, c.PublicIP, endpointPorts); err == nil {
		klog.Error("Found stale data, removed previous endpoints on kubernetes service, apiserver didn't exit successfully previously")
	} else if !storage.IsNotFound(err) {
		klog.Errorf("Error removing old endpoints from kubernetes service: %v", err)
	}

	go c.Run(stopCh)
}

// Stop cleans up this API Servers endpoint reconciliation leases so another master can take over more quickly.
func (c *Controller) Stop() {
	endpointPorts := createEndpointPortSpec(c.PublicServicePort, "https")
	finishedReconciling := make(chan struct{})
	go func() {
		defer close(finishedReconciling)
		klog.Infof("Shutting down kubernetes service endpoint reconciler")
		c.EndpointReconciler.StopReconciling()
		if err := c.EndpointReconciler.RemoveEndpoints(kubernetesServiceName, c.PublicIP, endpointPorts); err != nil {
			klog.Errorf("Unable to remove endpoints from kubernetes service: %v", err)
		}
		c.EndpointReconciler.Destroy()
	}()

	select {
	case <-finishedReconciling:
		// done
	case <-time.After(2 * c.EndpointInterval):
		// don't block server shutdown forever if we can't reach etcd to remove ourselves
		klog.Warning("RemoveEndpoints() timed out")
	}
}

// Run periodically updates the kubernetes service
func (c *Controller) Run(ch <-chan struct{}) {
	// wait until process is ready
	wait.PollImmediateUntil(100*time.Millisecond, func() (bool, error) {
		var code int
		c.client.CoreV1().RESTClient().Get().AbsPath("/readyz").Do(context.TODO()).StatusCode(&code)
		return code == http.StatusOK, nil
	}, ch)

	wait.NonSlidingUntil(func() {
		// Service definition is not reconciled after first
		// run, ports and type will be corrected only during
		// start.
		if err := c.UpdateKubernetesService(false); err != nil {
			runtime.HandleError(fmt.Errorf("unable to sync kubernetes service: %v", err))
		}
	}, c.EndpointInterval, ch)
}

// UpdateKubernetesService attempts to update the default Kube service.
func (c *Controller) UpdateKubernetesService(reconcile bool) error {
	// Update service & endpoint records.
	// TODO: when it becomes possible to change this stuff,
	// stop polling and start watching.
	// TODO: add endpoints of all replicas, not just the elected master.
	if _, err := c.client.CoreV1().Namespaces().Get(context.TODO(), metav1.NamespaceDefault, metav1.GetOptions{}); err != nil {
		if _, err := c.client.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      metav1.NamespaceDefault,
				Namespace: "",
			},
		}, metav1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}

	servicePorts, serviceType := createPortAndServiceSpec(c.ServicePort, c.PublicServicePort, c.KubernetesServiceNodePort, "https")
	if err := c.CreateOrUpdateMasterServiceIfNeeded(kubernetesServiceName, c.ServiceIP, servicePorts, serviceType, reconcile); err != nil {
		return err
	}
	endpointPorts := createEndpointPortSpec(c.PublicServicePort, "https")
	if err := c.EndpointReconciler.ReconcileEndpoints(kubernetesServiceName, c.PublicIP, endpointPorts, reconcile); err != nil {
		return err
	}
	return nil
}

// createPortAndServiceSpec creates an array of service ports.
// If the NodePort value is 0, just the servicePort is used, otherwise, a node port is exposed.
func createPortAndServiceSpec(servicePort int, targetServicePort int, nodePort int, servicePortName string) ([]corev1.ServicePort, corev1.ServiceType) {
	// Use the Cluster IP type for the service port if NodePort isn't provided.
	// Otherwise, we will be binding the master service to a NodePort.
	servicePorts := []corev1.ServicePort{{
		Protocol:   corev1.ProtocolTCP,
		Port:       int32(servicePort),
		Name:       servicePortName,
		TargetPort: intstr.FromInt(targetServicePort),
	}}
	serviceType := corev1.ServiceTypeClusterIP
	if nodePort > 0 {
		servicePorts[0].NodePort = int32(nodePort)
		serviceType = corev1.ServiceTypeNodePort
	}
	return servicePorts, serviceType
}

// createEndpointPortSpec creates the endpoint ports
func createEndpointPortSpec(endpointPort int, endpointPortName string) []corev1.EndpointPort {
	return []corev1.EndpointPort{{
		Protocol: corev1.ProtocolTCP,
		Port:     int32(endpointPort),
		Name:     endpointPortName,
	}}
}

// CreateOrUpdateMasterServiceIfNeeded will create the specified service if it
// doesn't already exist.
func (c *Controller) CreateOrUpdateMasterServiceIfNeeded(serviceName string, serviceIP net.IP, servicePorts []corev1.ServicePort, serviceType corev1.ServiceType, reconcile bool) error {
	if s, err := c.client.CoreV1().Services(metav1.NamespaceDefault).Get(context.TODO(), serviceName, metav1.GetOptions{}); err == nil {
		// The service already exists.
		if reconcile {
			if svc, updated := getMasterServiceUpdateIfNeeded(s, servicePorts, serviceType); updated {
				klog.Warningf("Resetting master service %q to %#v", serviceName, svc)
				_, err := c.client.CoreV1().Services(metav1.NamespaceDefault).Update(context.TODO(), svc, metav1.UpdateOptions{})
				return err
			}
		}
		return nil
	}
	singleStack := corev1.IPFamilyPolicySingleStack
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: metav1.NamespaceDefault,
			Labels:    map[string]string{"provider": "kubernetes", "component": "apiserver"},
		},
		Spec: corev1.ServiceSpec{
			Ports: servicePorts,
			// maintained by this code, not by the pod selector
			Selector:        nil,
			ClusterIP:       serviceIP.String(),
			IPFamilyPolicy:  &singleStack,
			SessionAffinity: corev1.ServiceAffinityNone,
			Type:            serviceType,
		},
	}

	_, err := c.client.CoreV1().Services(metav1.NamespaceDefault).Create(context.TODO(), svc, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		return c.CreateOrUpdateMasterServiceIfNeeded(serviceName, serviceIP, servicePorts, serviceType, reconcile)
	}
	return err
}

// getMasterServiceUpdateIfNeeded sets service attributes for the given apiserver service.
func getMasterServiceUpdateIfNeeded(svc *corev1.Service, servicePorts []corev1.ServicePort, serviceType corev1.ServiceType) (s *corev1.Service, updated bool) {
	// Determine if the service is in the format we expect
	// (servicePorts are present and service type matches)
	formatCorrect := checkServiceFormat(svc, servicePorts, serviceType)
	if formatCorrect {
		return svc, false
	}
	svc.Spec.Ports = servicePorts
	svc.Spec.Type = serviceType
	return svc, true
}

// Determine if the service is in the correct format
// getMasterServiceUpdateIfNeeded expects (servicePorts are correct
// and service type matches).
func checkServiceFormat(s *corev1.Service, ports []corev1.ServicePort, serviceType corev1.ServiceType) (formatCorrect bool) {
	if s.Spec.Type != serviceType {
		return false
	}
	if len(ports) != len(s.Spec.Ports) {
		return false
	}
	for i, port := range ports {
		if port != s.Spec.Ports[i] {
			return false
		}
	}
	return true
}
