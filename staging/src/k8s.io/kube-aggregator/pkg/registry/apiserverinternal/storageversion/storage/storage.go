/*
Copyright 2018 The Kubernetes Authors.

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

package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/apis/apiserverinternal"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	strategy "k8s.io/kube-aggregator/pkg/registry/apiserverinternal/storageversion"
)

// REST implements a RESTStorage for storage version against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against storageVersions
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:     func() runtime.Object { return &apiserverinternal.StorageVersion{} },
		NewListFunc: func() runtime.Object { return &apiserverinternal.StorageVersionList{} },
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*apiserverinternal.StorageVersion).Name, nil
		},
		DefaultQualifiedResource: apiserverinternal.Resource("storageversions"),

		CreateStrategy: strategy.Strategy,
		UpdateStrategy: strategy.Strategy,
		DeleteStrategy: strategy.Strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}
