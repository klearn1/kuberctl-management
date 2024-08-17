/*
Copyright 2024 The Kubernetes Authors.

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

package e2enode

import (
	"context"
	"fmt"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/kubernetes/test/e2e/framework"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	admissionapi "k8s.io/pod-security-admission/api"
)

var _ = SIGDescribe(framework.WithSerial(), "Pods status phase", func() {
	f := framework.NewDefaultFramework("pods-status-phase-test-serial")
	addAfterEachForCleaningUpPods(f)
	f.NamespacePodSecurityLevel = admissionapi.LevelPrivileged

	ginkgo.It("should be pending during the execution of the init container after the node reboot", func(ctx context.Context) {
		init := "init"
		restartableInit := "restartable-init"
		regular := "regular"

		podLabels := map[string]string{
			"test":      "pods-status-phase-test-serial",
			"namespace": f.Namespace.Name,
		}
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "initialized-pod",
				Labels: podLabels,
			},
			Spec: v1.PodSpec{
				RestartPolicy: v1.RestartPolicyAlways,
				InitContainers: []v1.Container{
					{
						Name:  init,
						Image: busyboxImage,
						Command: ExecCommand(init, execCommand{
							Delay:    30,
							ExitCode: 0,
						}),
					},
					{
						Name:  restartableInit,
						Image: busyboxImage,
						Command: ExecCommand(restartableInit, execCommand{
							Delay:    300,
							ExitCode: 0,
						}),
						RestartPolicy: &containerRestartPolicyAlways,
					},
				},
				Containers: []v1.Container{
					{
						Name:  regular,
						Image: busyboxImage,
						Command: ExecCommand(regular, execCommand{
							Delay:    300,
							ExitCode: 0,
						}),
					},
				},
			},
		}
		preparePod(pod)

		client := e2epod.NewPodClient(f)
		pod = client.Create(ctx, pod)

		ginkgo.By("Waiting for the regular init container to be started")
		pod, err := waitForPodInitContainerStarted(ctx, f.ClientSet, pod.Namespace, pod.Name, init, f.Timeouts.PodStart)
		framework.ExpectNoError(err)
		if pod.Status.Phase != v1.PodPending {
			framework.Failf("pod should be pending during the execution of the init container")
		}

		ginkgo.By("Getting the current pod sandbox ID")
		rs, _, err := getCRIClient()
		framework.ExpectNoError(err)

		sandboxes, err := rs.ListPodSandbox(ctx, &runtimeapi.PodSandboxFilter{
			LabelSelector: podLabels,
		})
		framework.ExpectNoError(err)
		gomega.Expect(sandboxes).To(gomega.HaveLen(1))
		podSandboxID := sandboxes[0].Id

		ginkgo.By("Stopping the kubelet")
		restartKubelet := stopKubelet()
		gomega.Eventually(ctx, func() bool {
			return kubeletHealthCheck(kubeletHealthCheckURL)
		}, f.Timeouts.PodStart, f.Timeouts.Poll).Should(gomega.BeFalse())

		ginkgo.By("Stopping the pod sandbox to simulate the node reboot")
		err = rs.StopPodSandbox(ctx, podSandboxID)
		framework.ExpectNoError(err)

		ginkgo.By("Restarting the kubelet")
		restartKubelet()
		gomega.Eventually(ctx, func() bool {
			return kubeletHealthCheck(kubeletHealthCheckURL)
		}, f.Timeouts.PodStart, f.Timeouts.Poll).Should(gomega.BeTrue())

		ginkgo.By("Waiting for the regular init container to be started after the node reboot")
		pod, err = waitForPodInitContainerStarted(ctx, f.ClientSet, pod.Namespace, pod.Name, init, f.Timeouts.PodStart)
		framework.ExpectNoError(err)
		if pod.Status.Phase != v1.PodPending {
			framework.Failf("pod should be pending during the execution of the init container after the node reboot")
		}
	})
})

// waitForPodInitContainerStarted waits for the given Pod init container to start, after a successful run of the startupProbe.
func waitForPodInitContainerStarted(ctx context.Context, c clientset.Interface, namespace, podName, initContainer string, timeout time.Duration) (*v1.Pod, error) {
	conditionDesc := fmt.Sprintf("regular init container %s started", initContainer)
	pendingPod := &v1.Pod{}
	err := e2epod.WaitForPodCondition(ctx, c, namespace, podName, conditionDesc, timeout, func(pod *v1.Pod) (bool, error) {
		for _, status := range pod.Status.InitContainerStatuses {
			if status.Name == initContainer && *status.Started {
				pendingPod = pod.DeepCopy()
				return true, nil
			}
		}
		return false, nil
	})
	return pendingPod, err
}
