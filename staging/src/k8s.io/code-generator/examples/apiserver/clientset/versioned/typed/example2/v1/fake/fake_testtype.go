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
	gentype "k8s.io/client-go/gentype"
	v1 "k8s.io/code-generator/examples/apiserver/apis/example2/v1"
	clientset "k8s.io/code-generator/examples/apiserver/clientset/versioned/typed/example2/v1"
)

// fakeTestTypes implements TestTypeInterface
type fakeTestTypes struct {
	*gentype.FakeClientWithList[*v1.TestType, *v1.TestTypeList]
	Fake *FakeSecondExampleV1
}

func newFakeTestTypes(fake *FakeSecondExampleV1, namespace string) clientset.TestTypeInterface {
	return &fakeTestTypes{
		gentype.NewFakeClientWithList[*v1.TestType, *v1.TestTypeList](
			fake.Fake,
			namespace,
			v1.SchemeGroupVersion.WithResource("testtypes"),
			v1.SchemeGroupVersion.WithKind("TestType"),
			func() *v1.TestType { return &v1.TestType{} },
			func() *v1.TestTypeList { return &v1.TestTypeList{} },
			func(dst, src *v1.TestTypeList) { dst.ListMeta = src.ListMeta },
			func(list *v1.TestTypeList) []*v1.TestType { return gentype.ToPointerSlice(list.Items) },
			func(list *v1.TestTypeList, items []*v1.TestType) { list.Items = gentype.FromPointerSlice(items) },
		),
		fake,
	}
}
