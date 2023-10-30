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

// PolicyRuleApplyConfiguration represents an declarative configuration of the PolicyRule type for use
// with apply.
type PolicyRuleApplyConfiguration struct {
	Verbs           []string `json:"verbs,omitempty"`
	APIGroups       []string `json:"apiGroups,omitempty"`
	Resources       []string `json:"resources,omitempty"`
	ResourceNames   []string `json:"resourceNames,omitempty"`
	NonResourceURLs []string `json:"nonResourceURLs,omitempty"`
}

// PolicyRuleApplyConfiguration constructs an declarative configuration of the PolicyRule type for use with
// apply.
func PolicyRule() *PolicyRuleApplyConfiguration {
	return &PolicyRuleApplyConfiguration{}
}

// WithVerbs adds the given value to the Verbs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Verbs field.
// Deprecated: WithVerbs does not replace existing list for atomic list type. Use WithNewVerbs instead.
func (b *PolicyRuleApplyConfiguration) WithVerbs(values ...string) *PolicyRuleApplyConfiguration {
	for i := range values {
		b.Verbs = append(b.Verbs, values[i])
	}
	return b
}

// WithNewVerbs replaces the Verbs field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the Verbs field is set to the values of the last call.
func (b *PolicyRuleApplyConfiguration) WithNewVerbs(values ...string) *PolicyRuleApplyConfiguration {
	b.Verbs = make([]string, len(values))
	copy(b.Verbs, values)
	return b
}

// WithAPIGroups adds the given value to the APIGroups field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the APIGroups field.
// Deprecated: WithAPIGroups does not replace existing list for atomic list type. Use WithNewAPIGroups instead.
func (b *PolicyRuleApplyConfiguration) WithAPIGroups(values ...string) *PolicyRuleApplyConfiguration {
	for i := range values {
		b.APIGroups = append(b.APIGroups, values[i])
	}
	return b
}

// WithNewAPIGroups replaces the APIGroups field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the APIGroups field is set to the values of the last call.
func (b *PolicyRuleApplyConfiguration) WithNewAPIGroups(values ...string) *PolicyRuleApplyConfiguration {
	b.APIGroups = make([]string, len(values))
	copy(b.APIGroups, values)
	return b
}

// WithResources adds the given value to the Resources field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Resources field.
// Deprecated: WithResources does not replace existing list for atomic list type. Use WithNewResources instead.
func (b *PolicyRuleApplyConfiguration) WithResources(values ...string) *PolicyRuleApplyConfiguration {
	for i := range values {
		b.Resources = append(b.Resources, values[i])
	}
	return b
}

// WithNewResources replaces the Resources field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the Resources field is set to the values of the last call.
func (b *PolicyRuleApplyConfiguration) WithNewResources(values ...string) *PolicyRuleApplyConfiguration {
	b.Resources = make([]string, len(values))
	copy(b.Resources, values)
	return b
}

// WithResourceNames adds the given value to the ResourceNames field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the ResourceNames field.
// Deprecated: WithResourceNames does not replace existing list for atomic list type. Use WithNewResourceNames instead.
func (b *PolicyRuleApplyConfiguration) WithResourceNames(values ...string) *PolicyRuleApplyConfiguration {
	for i := range values {
		b.ResourceNames = append(b.ResourceNames, values[i])
	}
	return b
}

// WithNewResourceNames replaces the ResourceNames field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the ResourceNames field is set to the values of the last call.
func (b *PolicyRuleApplyConfiguration) WithNewResourceNames(values ...string) *PolicyRuleApplyConfiguration {
	b.ResourceNames = make([]string, len(values))
	copy(b.ResourceNames, values)
	return b
}

// WithNonResourceURLs adds the given value to the NonResourceURLs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the NonResourceURLs field.
// Deprecated: WithNonResourceURLs does not replace existing list for atomic list type. Use WithNewNonResourceURLs instead.
func (b *PolicyRuleApplyConfiguration) WithNonResourceURLs(values ...string) *PolicyRuleApplyConfiguration {
	for i := range values {
		b.NonResourceURLs = append(b.NonResourceURLs, values[i])
	}
	return b
}

// WithNewNonResourceURLs replaces the NonResourceURLs field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the NonResourceURLs field is set to the values of the last call.
func (b *PolicyRuleApplyConfiguration) WithNewNonResourceURLs(values ...string) *PolicyRuleApplyConfiguration {
	b.NonResourceURLs = make([]string, len(values))
	copy(b.NonResourceURLs, values)
	return b
}
