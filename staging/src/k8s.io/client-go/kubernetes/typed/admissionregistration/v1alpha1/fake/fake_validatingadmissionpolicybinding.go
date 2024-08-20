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

package fake

import (
	v1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	admissionregistrationv1alpha1 "k8s.io/client-go/applyconfigurations/admissionregistration/v1alpha1"
	gentype "k8s.io/client-go/gentype"
	clientset "k8s.io/client-go/kubernetes/typed/admissionregistration/v1alpha1"
)

// fakeValidatingAdmissionPolicyBindings implements ValidatingAdmissionPolicyBindingInterface
type fakeValidatingAdmissionPolicyBindings struct {
	*gentype.FakeClientWithListAndApply[*v1alpha1.ValidatingAdmissionPolicyBinding, *v1alpha1.ValidatingAdmissionPolicyBindingList, *admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingApplyConfiguration]
	Fake *FakeAdmissionregistrationV1alpha1
}

func newFakeValidatingAdmissionPolicyBindings(fake *FakeAdmissionregistrationV1alpha1) clientset.ValidatingAdmissionPolicyBindingInterface {
	return &fakeValidatingAdmissionPolicyBindings{
		gentype.NewFakeClientWithListAndApply[*v1alpha1.ValidatingAdmissionPolicyBinding, *v1alpha1.ValidatingAdmissionPolicyBindingList, *admissionregistrationv1alpha1.ValidatingAdmissionPolicyBindingApplyConfiguration](
			fake.Fake,
			"",
			v1alpha1.SchemeGroupVersion.WithResource("validatingadmissionpolicybindings"),
			v1alpha1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicyBinding"),
			func() *v1alpha1.ValidatingAdmissionPolicyBinding { return &v1alpha1.ValidatingAdmissionPolicyBinding{} },
			func() *v1alpha1.ValidatingAdmissionPolicyBindingList {
				return &v1alpha1.ValidatingAdmissionPolicyBindingList{}
			},
			func(dst, src *v1alpha1.ValidatingAdmissionPolicyBindingList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.ValidatingAdmissionPolicyBindingList) []*v1alpha1.ValidatingAdmissionPolicyBinding {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.ValidatingAdmissionPolicyBindingList, items []*v1alpha1.ValidatingAdmissionPolicyBinding) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
