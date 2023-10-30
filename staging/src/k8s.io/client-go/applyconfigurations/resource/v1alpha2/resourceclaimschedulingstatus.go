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

package v1alpha2

// ResourceClaimSchedulingStatusApplyConfiguration represents an declarative configuration of the ResourceClaimSchedulingStatus type for use
// with apply.
type ResourceClaimSchedulingStatusApplyConfiguration struct {
	Name            *string  `json:"name,omitempty"`
	UnsuitableNodes []string `json:"unsuitableNodes,omitempty"`
}

// ResourceClaimSchedulingStatusApplyConfiguration constructs an declarative configuration of the ResourceClaimSchedulingStatus type for use with
// apply.
func ResourceClaimSchedulingStatus() *ResourceClaimSchedulingStatusApplyConfiguration {
	return &ResourceClaimSchedulingStatusApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *ResourceClaimSchedulingStatusApplyConfiguration) WithName(value string) *ResourceClaimSchedulingStatusApplyConfiguration {
	b.Name = &value
	return b
}

// WithUnsuitableNodes adds the given value to the UnsuitableNodes field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the UnsuitableNodes field.
// Deprecated: WithUnsuitableNodes does not replace existing list for atomic list type. Use WithNewUnsuitableNodes instead.
func (b *ResourceClaimSchedulingStatusApplyConfiguration) WithUnsuitableNodes(values ...string) *ResourceClaimSchedulingStatusApplyConfiguration {
	for i := range values {
		b.UnsuitableNodes = append(b.UnsuitableNodes, values[i])
	}
	return b
}

// WithNewUnsuitableNodes replaces the UnsuitableNodes field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the UnsuitableNodes field is set to the values of the last call.
func (b *ResourceClaimSchedulingStatusApplyConfiguration) WithNewUnsuitableNodes(values ...string) *ResourceClaimSchedulingStatusApplyConfiguration {
	b.UnsuitableNodes = make([]string, len(values))
	copy(b.UnsuitableNodes, values)
	return b
}
