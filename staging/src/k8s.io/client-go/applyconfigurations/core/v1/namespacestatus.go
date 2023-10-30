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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	v1 "k8s.io/api/core/v1"
)

// NamespaceStatusApplyConfiguration represents an declarative configuration of the NamespaceStatus type for use
// with apply.
type NamespaceStatusApplyConfiguration struct {
	Phase      *v1.NamespacePhase                     `json:"phase,omitempty"`
	Conditions []NamespaceConditionApplyConfiguration `json:"conditions,omitempty"`
}

// NamespaceStatusApplyConfiguration constructs an declarative configuration of the NamespaceStatus type for use with
// apply.
func NamespaceStatus() *NamespaceStatusApplyConfiguration {
	return &NamespaceStatusApplyConfiguration{}
}

// WithPhase sets the Phase field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Phase field is set to the value of the last call.
func (b *NamespaceStatusApplyConfiguration) WithPhase(value v1.NamespacePhase) *NamespaceStatusApplyConfiguration {
	b.Phase = &value
	return b
}

// WithConditions adds the given value to the Conditions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Conditions field.
// Deprecated: WithConditions does not replace existing list for atomic list type. Use WithNewConditions instead.
func (b *NamespaceStatusApplyConfiguration) WithConditions(values ...*NamespaceConditionApplyConfiguration) *NamespaceStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithConditions")
		}
		b.Conditions = append(b.Conditions, *values[i])
	}
	return b
}

// WithNewConditions replaces the Conditions field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the Conditions field is set to the values of the last call.
func (b *NamespaceStatusApplyConfiguration) WithNewConditions(values ...*NamespaceConditionApplyConfiguration) *NamespaceStatusApplyConfiguration {
	b.Conditions = nil
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNewConditions")
		}
		b.Conditions = append(b.Conditions, *values[i])
	}
	return b
}
