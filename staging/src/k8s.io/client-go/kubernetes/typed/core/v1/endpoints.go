/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	gentype "k8s.io/client-go/gentype"
	scheme "k8s.io/client-go/kubernetes/scheme"
)

// EndpointsGetter has a method to return a EndpointsInterface.
// A group's client should implement this interface.
type EndpointsGetter interface {
	Endpoints(namespace string) EndpointsInterface
}

// EndpointsInterface has methods to work with Endpoints resources.
type EndpointsInterface interface {
	Create(ctx context.Context, endpoints *v1.Endpoints, opts metav1.CreateOptions) (*v1.Endpoints, error)
	Update(ctx context.Context, endpoints *v1.Endpoints, opts metav1.UpdateOptions) (*v1.Endpoints, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Endpoints, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.EndpointsList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Endpoints, err error)
	Apply(ctx context.Context, endpoints *corev1.EndpointsApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Endpoints, err error)
	EndpointsExpansion
}

// endpoints implements EndpointsInterface
type endpoints struct {
	*gentype.ClientWithListAndApply[*v1.Endpoints, *v1.EndpointsList, *corev1.EndpointsApplyConfiguration]
}

// newEndpoints returns a Endpoints
func newEndpoints(c *CoreV1Client, namespace string) *endpoints {
	return &endpoints{
		gentype.NewClientWithListAndApply[*v1.Endpoints, *v1.EndpointsList, *corev1.EndpointsApplyConfiguration](
			"endpoints",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *v1.Endpoints { return &v1.Endpoints{} },
			func() *v1.EndpointsList { return &v1.EndpointsList{} },
			true,
		),
	}
}
