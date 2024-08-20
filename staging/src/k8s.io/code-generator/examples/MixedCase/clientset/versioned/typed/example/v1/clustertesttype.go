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

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
	v1 "k8s.io/code-generator/examples/MixedCase/apis/example/v1"
	examplev1 "k8s.io/code-generator/examples/MixedCase/applyconfiguration/example/v1"
	scheme "k8s.io/code-generator/examples/MixedCase/clientset/versioned/scheme"
)

// ClusterTestTypesGetter has a method to return a ClusterTestTypeInterface.
// A group's client should implement this interface.
type ClusterTestTypesGetter interface {
	ClusterTestTypes() ClusterTestTypeInterface
}

// ClusterTestTypeInterface has methods to work with ClusterTestType resources.
type ClusterTestTypeInterface interface {
	Create(ctx context.Context, clusterTestType *v1.ClusterTestType, opts metav1.CreateOptions) (*v1.ClusterTestType, error)
	Update(ctx context.Context, clusterTestType *v1.ClusterTestType, opts metav1.UpdateOptions) (*v1.ClusterTestType, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, clusterTestType *v1.ClusterTestType, opts metav1.UpdateOptions) (*v1.ClusterTestType, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.ClusterTestType, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ClusterTestTypeList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.ClusterTestType, err error)
	Apply(ctx context.Context, clusterTestType *examplev1.ClusterTestTypeApplyConfiguration, opts metav1.ApplyOptions) (result *v1.ClusterTestType, err error)
	// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
	ApplyStatus(ctx context.Context, clusterTestType *examplev1.ClusterTestTypeApplyConfiguration, opts metav1.ApplyOptions) (result *v1.ClusterTestType, err error)
	GetScale(ctx context.Context, clusterTestTypeName string, options metav1.GetOptions) (*autoscalingv1.Scale, error)
	UpdateScale(ctx context.Context, clusterTestTypeName string, scale *autoscalingv1.Scale, opts metav1.UpdateOptions) (*autoscalingv1.Scale, error)
	CreateScale(ctx context.Context, clusterTestTypeName string, scale *autoscalingv1.Scale, opts metav1.CreateOptions) (*autoscalingv1.Scale, error)

	ClusterTestTypeExpansion
}

// clusterTestTypes implements ClusterTestTypeInterface
type clusterTestTypes struct {
	*gentype.ClientWithListAndApply[*v1.ClusterTestType, *v1.ClusterTestTypeList, *examplev1.ClusterTestTypeApplyConfiguration]
}

// newClusterTestTypes returns a ClusterTestTypes
func newClusterTestTypes(c *ExampleV1Client) *clusterTestTypes {
	return &clusterTestTypes{
		gentype.NewClientWithListAndApply[*v1.ClusterTestType, *v1.ClusterTestTypeList, *examplev1.ClusterTestTypeApplyConfiguration](
			"clustertesttypes",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *v1.ClusterTestType { return &v1.ClusterTestType{} },
			func() *v1.ClusterTestTypeList { return &v1.ClusterTestTypeList{} },
			false,
		),
	}
}

// GetScale takes name of the clusterTestType, and returns the corresponding autoscalingv1.Scale object, and an error if there is any.
func (c *clusterTestTypes) GetScale(ctx context.Context, clusterTestTypeName string, options metav1.GetOptions) (result *autoscalingv1.Scale, err error) {
	result = &autoscalingv1.Scale{}
	err = c.GetClient().Get().
		Resource("clustertesttypes").
		Name(clusterTestTypeName).
		SubResource("scale").
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// UpdateScale takes the top resource name and the representation of a scale and updates it. Returns the server's representation of the scale, and an error, if there is any.
func (c *clusterTestTypes) UpdateScale(ctx context.Context, clusterTestTypeName string, scale *autoscalingv1.Scale, opts metav1.UpdateOptions) (result *autoscalingv1.Scale, err error) {
	result = &autoscalingv1.Scale{}
	err = c.GetClient().Put().
		Resource("clustertesttypes").
		Name(clusterTestTypeName).
		SubResource("scale").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(scale).
		Do(ctx).
		Into(result)
	return
}

// CreateScale takes the representation of a scale and creates it.  Returns the server's representation of the scale, and an error, if there is any.
func (c *clusterTestTypes) CreateScale(ctx context.Context, clusterTestTypeName string, scale *autoscalingv1.Scale, opts metav1.CreateOptions) (result *autoscalingv1.Scale, err error) {
	result = &autoscalingv1.Scale{}
	err = c.GetClient().Post().
		Resource("clustertesttypes").
		Name(clusterTestTypeName).
		SubResource("scale").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(scale).
		Do(ctx).
		Into(result)
	return
}
