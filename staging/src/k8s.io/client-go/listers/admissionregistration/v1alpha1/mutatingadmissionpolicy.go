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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// MutatingAdmissionPolicyLister helps list MutatingAdmissionPolicies.
// All objects returned here must be treated as read-only.
type MutatingAdmissionPolicyLister interface {
	// List lists all MutatingAdmissionPolicies in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.MutatingAdmissionPolicy, err error)
	// Get retrieves the MutatingAdmissionPolicy from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.MutatingAdmissionPolicy, error)
	MutatingAdmissionPolicyListerExpansion
}

// mutatingAdmissionPolicyLister implements the MutatingAdmissionPolicyLister interface.
type mutatingAdmissionPolicyLister struct {
	listers.ResourceIndexer[*v1alpha1.MutatingAdmissionPolicy]
}

// NewMutatingAdmissionPolicyLister returns a new MutatingAdmissionPolicyLister.
func NewMutatingAdmissionPolicyLister(indexer cache.Indexer) MutatingAdmissionPolicyLister {
	return &mutatingAdmissionPolicyLister{listers.New[*v1alpha1.MutatingAdmissionPolicy](indexer, v1alpha1.Resource("mutatingadmissionpolicy"))}
}
