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

// IPBlockApplyConfiguration represents an declarative configuration of the IPBlock type for use
// with apply.
type IPBlockApplyConfiguration struct {
	CIDR   *string  `json:"cidr,omitempty"`
	Except []string `json:"except,omitempty"`
}

// IPBlockApplyConfiguration constructs an declarative configuration of the IPBlock type for use with
// apply.
func IPBlock() *IPBlockApplyConfiguration {
	return &IPBlockApplyConfiguration{}
}

// WithCIDR sets the CIDR field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CIDR field is set to the value of the last call.
func (b *IPBlockApplyConfiguration) WithCIDR(value string) *IPBlockApplyConfiguration {
	b.CIDR = &value
	return b
}

// WithExcept adds the given value to the Except field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Except field.
// Deprecated: WithExcept does not replace existing list for atomic list type. Use WithNewExcept instead.
func (b *IPBlockApplyConfiguration) WithExcept(values ...string) *IPBlockApplyConfiguration {
	for i := range values {
		b.Except = append(b.Except, values[i])
	}
	return b
}

// WithNewExcept replaces the Except field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the Except field is set to the values of the last call.
func (b *IPBlockApplyConfiguration) WithNewExcept(values ...string) *IPBlockApplyConfiguration {
	b.Except = make([]string, len(values))
	copy(b.Except, values)
	return b
}
