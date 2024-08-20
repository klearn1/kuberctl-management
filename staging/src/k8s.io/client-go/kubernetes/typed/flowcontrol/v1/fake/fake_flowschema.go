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
	v1 "k8s.io/api/flowcontrol/v1"
	flowcontrolv1 "k8s.io/client-go/applyconfigurations/flowcontrol/v1"
	gentype "k8s.io/client-go/gentype"
	clientset "k8s.io/client-go/kubernetes/typed/flowcontrol/v1"
)

// fakeFlowSchemas implements FlowSchemaInterface
type fakeFlowSchemas struct {
	*gentype.FakeClientWithListAndApply[*v1.FlowSchema, *v1.FlowSchemaList, *flowcontrolv1.FlowSchemaApplyConfiguration]
	Fake *FakeFlowcontrolV1
}

func newFakeFlowSchemas(fake *FakeFlowcontrolV1) clientset.FlowSchemaInterface {
	return &fakeFlowSchemas{
		gentype.NewFakeClientWithListAndApply[*v1.FlowSchema, *v1.FlowSchemaList, *flowcontrolv1.FlowSchemaApplyConfiguration](
			fake.Fake,
			"",
			v1.SchemeGroupVersion.WithResource("flowschemas"),
			v1.SchemeGroupVersion.WithKind("FlowSchema"),
			func() *v1.FlowSchema { return &v1.FlowSchema{} },
			func() *v1.FlowSchemaList { return &v1.FlowSchemaList{} },
			func(dst, src *v1.FlowSchemaList) { dst.ListMeta = src.ListMeta },
			func(list *v1.FlowSchemaList) []*v1.FlowSchema { return gentype.ToPointerSlice(list.Items) },
			func(list *v1.FlowSchemaList, items []*v1.FlowSchema) { list.Items = gentype.FromPointerSlice(items) },
		),
		fake,
	}
}
