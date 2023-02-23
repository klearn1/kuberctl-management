/*
Copyright 2023 The Kubernetes Authors.

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

package apply

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const defaultApplySetParentResource = "secrets"
const defaultApplySetParentVersion = "v1"

// ApplySet tracks the information about an applyset apply/prune
type ApplySet struct {
	// ParentRef is the reference to the parent object that is used to track the applyset.
	ParentRef *ApplySetParentRef

	// resources is the set of all the resources that (might) be part of this applyset.
	resources map[schema.GroupVersionResource]struct{}

	// namespaces is the set of all namespaces that (might) contain objects that are part of this applyset.
	namespaces map[string]struct{}
}

// ApplySetParentRef stores object and type meta for the parent object that is used to track the applyset.
type ApplySetParentRef struct {
	Name        string
	Namespace   string
	RESTMapping *meta.RESTMapping
}

const invalidParentRefFmt = "invalid parent reference %q: %w"

// NewApplySet creates a new ApplySet object from a parent reference in the format [RESOURCE][.GROUP]/NAME
func NewApplySet(parentRefStr string, namespace string, mapper meta.RESTMapper) (*ApplySet, error) {
	var parent *ApplySetParentRef
	var err error

	if parent, err = parentRefFromStr(parentRefStr, mapper); err != nil {
		return nil, fmt.Errorf(invalidParentRefFmt, parentRefStr, err)
	}
	parent.Namespace = namespace
	if err := parent.Validate(); err != nil {
		return nil, fmt.Errorf(invalidParentRefFmt, parentRefStr, err)
	}

	return &ApplySet{
		resources:  make(map[schema.GroupVersionResource]struct{}),
		namespaces: make(map[string]struct{}),
		ParentRef:  parent,
	}, nil
}

// ID is the label value that we are using to identify this applyset.
func (a ApplySet) ID() string {
	// TODO: base64(sha256(gknn))
	return "placeholder-todo"
}

func (a *ApplySet) LabelsForMember() map[string]string {
	return map[string]string{
		"applyset.k8s.io/part-of": a.ID(),
	}
}

func (p *ApplySetParentRef) IsNamespaced() bool {
	return p.RESTMapping.Scope.Name() == meta.RESTScopeNameNamespace
}

func (p *ApplySetParentRef) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name cannot be blank")
	}
	if p.IsNamespaced() && p.Namespace == "" {
		// TODO: Can this actually happen? We seem to get the value 'default' here when the namespace flag unspecified or even blank. May need KEP update.
		return fmt.Errorf("namespace is required to use namespace-scoped ApplySet")
	}
	return nil
}

// parentRefFromStr creates a new ApplySetParentRef from a parent reference in the format [RESOURCE][.GROUP]/NAME
func parentRefFromStr(parentRefStr string, mapper meta.RESTMapper) (*ApplySetParentRef, error) {
	var gvr schema.GroupVersionResource
	var name string

	if groupRes, nameSuffix, hasTypeInfo := strings.Cut(parentRefStr, "/"); hasTypeInfo {
		name = nameSuffix
		gvr = schema.ParseGroupResource(groupRes).WithVersion("")
	} else {
		name = parentRefStr
		gvr = schema.GroupVersionResource{Version: defaultApplySetParentVersion, Resource: defaultApplySetParentResource}
	}

	gvk, err := mapper.KindFor(gvr)
	if err != nil {
		return nil, err
	}
	mapping, err := mapper.RESTMapping(gvk.GroupKind())
	if err != nil {
		return nil, err
	}
	return &ApplySetParentRef{Name: name, RESTMapping: mapping}, nil
}
