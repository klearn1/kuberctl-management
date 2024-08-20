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

package v1beta2

import (
	"context"

	v1beta2 "k8s.io/api/apps/v1beta2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	appsv1beta2 "k8s.io/client-go/applyconfigurations/apps/v1beta2"
	gentype "k8s.io/client-go/gentype"
	scheme "k8s.io/client-go/kubernetes/scheme"
)

// ControllerRevisionsGetter has a method to return a ControllerRevisionInterface.
// A group's client should implement this interface.
type ControllerRevisionsGetter interface {
	ControllerRevisions(namespace string) ControllerRevisionInterface
}

// ControllerRevisionInterface has methods to work with ControllerRevision resources.
type ControllerRevisionInterface interface {
	Create(ctx context.Context, controllerRevision *v1beta2.ControllerRevision, opts v1.CreateOptions) (*v1beta2.ControllerRevision, error)
	Update(ctx context.Context, controllerRevision *v1beta2.ControllerRevision, opts v1.UpdateOptions) (*v1beta2.ControllerRevision, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta2.ControllerRevision, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta2.ControllerRevisionList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta2.ControllerRevision, err error)
	Apply(ctx context.Context, controllerRevision *appsv1beta2.ControllerRevisionApplyConfiguration, opts v1.ApplyOptions) (result *v1beta2.ControllerRevision, err error)
	ControllerRevisionExpansion
}

// controllerRevisions implements ControllerRevisionInterface
type controllerRevisions struct {
	*gentype.ClientWithListAndApply[*v1beta2.ControllerRevision, *v1beta2.ControllerRevisionList, *appsv1beta2.ControllerRevisionApplyConfiguration]
}

// newControllerRevisions returns a ControllerRevisions
func newControllerRevisions(c *AppsV1beta2Client, namespace string) *controllerRevisions {
	return &controllerRevisions{
		gentype.NewClientWithListAndApply[*v1beta2.ControllerRevision, *v1beta2.ControllerRevisionList, *appsv1beta2.ControllerRevisionApplyConfiguration](
			"controllerrevisions",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *v1beta2.ControllerRevision { return &v1beta2.ControllerRevision{} },
			func() *v1beta2.ControllerRevisionList { return &v1beta2.ControllerRevisionList{} },
			true,
		),
	}
}
