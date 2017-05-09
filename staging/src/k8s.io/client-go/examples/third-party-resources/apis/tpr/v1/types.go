/*
Copyright 2017 The Kubernetes Authors.

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

package v1

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const ExampleResourcePlural = "examples"

type Example struct {
	metav1.TypeMeta `json:",inline"`
	// WARNING: Don't call the field below ObjectMeta to avoid the issue with ugorji
	// https://github.com/kubernetes/client-go/issues/8#issuecomment-285333502
	Metadata metav1.ObjectMeta `json:"metadata"`

	Spec ExampleSpec `json:"spec"`

	Status ExampleStatus `json:"status,omitempty"`
}

type ExampleSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
}

type ExampleStatus struct {
	State   ExampleState `json:"state,omitempty"`
	Message string       `json:"message,omitempty"`
}

type ExampleState string

const (
	ExampleStateCreated   ExampleState = "Created"
	ExampleStateProcessed ExampleState = "Processed"
)

type ExampleList struct {
	metav1.TypeMeta `json:",inline"`
	// WARNING: Don't call the field below ListMeta to avoid the issue with ugorji
	// https://github.com/kubernetes/client-go/issues/8#issuecomment-285333502
	Metadata metav1.ListMeta `json:"metadata"`

	Items []Example `json:"items"`
}

// The code below is used only to work around a known problem with third-party
// resources and ugorji. If/when these issues are resolved, the code below
// should no longer be required.

// Required to satisfy Object interface
func (e *Example) GetObjectKind() schema.ObjectKind {
	return &e.TypeMeta
}

// Required to satisfy ObjectMetaAccessor interface
func (e *Example) GetObjectMeta() metav1.Object {
	return &e.Metadata
}

// Required to satisfy Object interface
func (el *ExampleList) GetObjectKind() schema.ObjectKind {
	return &el.TypeMeta
}

// Required to satisfy ListMetaAccessor interface
func (el *ExampleList) GetListMeta() metav1.List {
	return &el.Metadata
}

type ExampleListCopy ExampleList
type ExampleCopy Example

func (e *Example) UnmarshalJSON(data []byte) error {
	tmp := ExampleCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := Example(tmp)
	*e = tmp2
	return nil
}

func (el *ExampleList) UnmarshalJSON(data []byte) error {
	tmp := ExampleListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := ExampleList(tmp)
	*el = tmp2
	return nil
}
