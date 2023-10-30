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

package v1beta1

// FlowSchemaSpecApplyConfiguration represents an declarative configuration of the FlowSchemaSpec type for use
// with apply.
type FlowSchemaSpecApplyConfiguration struct {
	PriorityLevelConfiguration *PriorityLevelConfigurationReferenceApplyConfiguration `json:"priorityLevelConfiguration,omitempty"`
	MatchingPrecedence         *int32                                                 `json:"matchingPrecedence,omitempty"`
	DistinguisherMethod        *FlowDistinguisherMethodApplyConfiguration             `json:"distinguisherMethod,omitempty"`
	Rules                      []PolicyRulesWithSubjectsApplyConfiguration            `json:"rules,omitempty"`
}

// FlowSchemaSpecApplyConfiguration constructs an declarative configuration of the FlowSchemaSpec type for use with
// apply.
func FlowSchemaSpec() *FlowSchemaSpecApplyConfiguration {
	return &FlowSchemaSpecApplyConfiguration{}
}

// WithPriorityLevelConfiguration sets the PriorityLevelConfiguration field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PriorityLevelConfiguration field is set to the value of the last call.
func (b *FlowSchemaSpecApplyConfiguration) WithPriorityLevelConfiguration(value *PriorityLevelConfigurationReferenceApplyConfiguration) *FlowSchemaSpecApplyConfiguration {
	b.PriorityLevelConfiguration = value
	return b
}

// WithMatchingPrecedence sets the MatchingPrecedence field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MatchingPrecedence field is set to the value of the last call.
func (b *FlowSchemaSpecApplyConfiguration) WithMatchingPrecedence(value int32) *FlowSchemaSpecApplyConfiguration {
	b.MatchingPrecedence = &value
	return b
}

// WithDistinguisherMethod sets the DistinguisherMethod field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DistinguisherMethod field is set to the value of the last call.
func (b *FlowSchemaSpecApplyConfiguration) WithDistinguisherMethod(value *FlowDistinguisherMethodApplyConfiguration) *FlowSchemaSpecApplyConfiguration {
	b.DistinguisherMethod = value
	return b
}

// WithRules adds the given value to the Rules field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Rules field.
// Deprecated: WithRules does not replace existing list for atomic list type. Use WithNewRules instead.
func (b *FlowSchemaSpecApplyConfiguration) WithRules(values ...*PolicyRulesWithSubjectsApplyConfiguration) *FlowSchemaSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithRules")
		}
		b.Rules = append(b.Rules, *values[i])
	}
	return b
}

// WithNewRules replaces the Rules field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the Rules field is set to the values of the last call.
func (b *FlowSchemaSpecApplyConfiguration) WithNewRules(values ...*PolicyRulesWithSubjectsApplyConfiguration) *FlowSchemaSpecApplyConfiguration {
	b.Rules = nil
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNewRules")
		}
		b.Rules = append(b.Rules, *values[i])
	}
	return b
}
