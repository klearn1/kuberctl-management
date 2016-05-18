/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package util

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/runtime"
)

// DynamicVersionIntefaces provides an object converter and metadata
// accessor appropriate for use with unstructured objects.
func DynamicVersionInterfaces(unversioned.GroupVersion) (*meta.VersionInterfaces, error) {
	return &meta.VersionInterfaces{
		ObjectConvertor:  &runtime.UnstructuredObjectConverter{},
		MetadataAccessor: meta.NewAccessor(),
	}, nil
}

// NewDiscoveryRESTMapper returns a RESTMapper based on discovery information.
func NewDiscoveryRESTMapper(resources []*unversioned.APIResourceList, versionFunc meta.VersionInterfacesFunc) (*meta.DefaultRESTMapper, error) {
	rm := meta.NewDefaultRESTMapper(nil, versionFunc)
	for _, resourceList := range resources {
		gv, err := unversioned.ParseGroupVersion(resourceList.GroupVersion)
		if err != nil {
			return nil, err
		}

		for _, resource := range resourceList.APIResources {
			gvk := gv.WithKind(resource.Kind)
			scope := meta.RESTScopeRoot
			if resource.Namespaced {
				scope = meta.RESTScopeNamespace
			}
			rm.Add(gvk, scope)
		}
	}
	return rm, nil
}

// DynamicObjectTyper provides an ObjectTyper implmentation for
// runtime.Unstructured object based on discovery information.
type DynamicObjectTyper struct {
	registered map[unversioned.GroupVersionKind]bool
}

func NewDynamicObjectTyper(resources []*unversioned.APIResourceList) (runtime.ObjectTyper, error) {
	dot := &DynamicObjectTyper{registered: make(map[unversioned.GroupVersionKind]bool)}
	for _, resourceList := range resources {
		gv, err := unversioned.ParseGroupVersion(resourceList.GroupVersion)
		if err != nil {
			return nil, err
		}

		for _, resource := range resourceList.APIResources {
			dot.registered[gv.WithKind(resource.Kind)] = true
		}
	}
	return dot, nil
}

// ObjectKind returns the group,version,kind of the provided object,
// or an error if the object in not *runtime.Unstructured or has no
// group,version,kind information.
func (dot *DynamicObjectTyper) ObjectKind(obj runtime.Object) (unversioned.GroupVersionKind, error) {
	if _, ok := obj.(*runtime.Unstructured); !ok {
		return unversioned.GroupVersionKind{}, fmt.Errorf("type %T is invalid for dynamic object typer", obj)
	}

	return obj.GetObjectKind().GroupVersionKind(), nil
}

// ObjectKinds returns a slice of one element with the
// group,version,kind of the provided object, or an error if the
// object is not *runtime.Unstructured or has no group,version,kind
// information.
func (dot *DynamicObjectTyper) ObjectKinds(obj runtime.Object) ([]unversioned.GroupVersionKind, error) {
	gvk, err := dot.ObjectKind(obj)
	if err != nil {
		return nil, err
	}

	return []unversioned.GroupVersionKind{gvk}, nil
}

// Recognizes returns true if the provided group,version,kind was in
// the discovery information.
func (dot *DynamicObjectTyper) Recognizes(gvk unversioned.GroupVersionKind) bool {
	return dot.registered[gvk]
}

// IsUnversioned returns false always because *runtime.Unstructured
// objects should always have group,version,kind information set. ok
// will be true if the object's group,version,kind is registered.
func (dot *DynamicObjectTyper) IsUnversioned(obj runtime.Object) (unversioned bool, ok bool) {
	gvk, err := dot.ObjectKind(obj)
	if err != nil {
		return false, false
	}

	return false, dot.registered[gvk]
}
