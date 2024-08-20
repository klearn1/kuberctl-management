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

package v1alpha1

import (
	admissionregistrationv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
)

// MutatingAdmissionPolicySpecApplyConfiguration represents a declarative configuration of the MutatingAdmissionPolicySpec type for use
// with apply.
type MutatingAdmissionPolicySpecApplyConfiguration struct {
	ParamKind        *ParamKindApplyConfiguration                     `json:"paramKind,omitempty"`
	MatchConstraints *MatchResourcesApplyConfiguration                `json:"matchConstraints,omitempty"`
	Mutations        []MutationApplyConfiguration                     `json:"mutations,omitempty"`
	FailurePolicy    *admissionregistrationv1alpha1.FailurePolicyType `json:"failurePolicy,omitempty"`
	MatchConditions  []MatchConditionApplyConfiguration               `json:"matchConditions,omitempty"`
}

// MutatingAdmissionPolicySpecApplyConfiguration constructs a declarative configuration of the MutatingAdmissionPolicySpec type for use with
// apply.
func MutatingAdmissionPolicySpec() *MutatingAdmissionPolicySpecApplyConfiguration {
	return &MutatingAdmissionPolicySpecApplyConfiguration{}
}

// WithParamKind sets the ParamKind field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ParamKind field is set to the value of the last call.
func (b *MutatingAdmissionPolicySpecApplyConfiguration) WithParamKind(value *ParamKindApplyConfiguration) *MutatingAdmissionPolicySpecApplyConfiguration {
	b.ParamKind = value
	return b
}

// WithMatchConstraints sets the MatchConstraints field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MatchConstraints field is set to the value of the last call.
func (b *MutatingAdmissionPolicySpecApplyConfiguration) WithMatchConstraints(value *MatchResourcesApplyConfiguration) *MutatingAdmissionPolicySpecApplyConfiguration {
	b.MatchConstraints = value
	return b
}

// WithMutations adds the given value to the Mutations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Mutations field.
func (b *MutatingAdmissionPolicySpecApplyConfiguration) WithMutations(values ...*MutationApplyConfiguration) *MutatingAdmissionPolicySpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithMutations")
		}
		b.Mutations = append(b.Mutations, *values[i])
	}
	return b
}

// WithFailurePolicy sets the FailurePolicy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FailurePolicy field is set to the value of the last call.
func (b *MutatingAdmissionPolicySpecApplyConfiguration) WithFailurePolicy(value admissionregistrationv1alpha1.FailurePolicyType) *MutatingAdmissionPolicySpecApplyConfiguration {
	b.FailurePolicy = &value
	return b
}

// WithMatchConditions adds the given value to the MatchConditions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the MatchConditions field.
func (b *MutatingAdmissionPolicySpecApplyConfiguration) WithMatchConditions(values ...*MatchConditionApplyConfiguration) *MutatingAdmissionPolicySpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithMatchConditions")
		}
		b.MatchConditions = append(b.MatchConditions, *values[i])
	}
	return b
}
