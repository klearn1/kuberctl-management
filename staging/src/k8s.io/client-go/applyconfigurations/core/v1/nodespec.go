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

// NodeSpecApplyConfiguration represents an declarative configuration of the NodeSpec type for use
// with apply.
type NodeSpecApplyConfiguration struct {
	PodCIDR            *string                             `json:"podCIDR,omitempty"`
	PodCIDRs           []string                            `json:"podCIDRs,omitempty"`
	ProviderID         *string                             `json:"providerID,omitempty"`
	Unschedulable      *bool                               `json:"unschedulable,omitempty"`
	Taints             []TaintApplyConfiguration           `json:"taints,omitempty"`
	ConfigSource       *NodeConfigSourceApplyConfiguration `json:"configSource,omitempty"`
	DoNotUseExternalID *string                             `json:"externalID,omitempty"`
}

// NodeSpecApplyConfiguration constructs an declarative configuration of the NodeSpec type for use with
// apply.
func NodeSpec() *NodeSpecApplyConfiguration {
	return &NodeSpecApplyConfiguration{}
}

// WithPodCIDR sets the PodCIDR field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PodCIDR field is set to the value of the last call.
func (b *NodeSpecApplyConfiguration) WithPodCIDR(value string) *NodeSpecApplyConfiguration {
	b.PodCIDR = &value
	return b
}

// WithPodCIDRs adds the given value to the PodCIDRs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the PodCIDRs field.
// Deprecated: WithPodCIDRs does not replace existing list for atomic list type. Use WithNewPodCIDRs instead.
func (b *NodeSpecApplyConfiguration) WithPodCIDRs(values ...string) *NodeSpecApplyConfiguration {
	for i := range values {
		b.PodCIDRs = append(b.PodCIDRs, values[i])
	}
	return b
}

// WithNewPodCIDRs replaces the PodCIDRs field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the PodCIDRs field is set to the values of the last call.
func (b *NodeSpecApplyConfiguration) WithNewPodCIDRs(values ...string) *NodeSpecApplyConfiguration {
	b.PodCIDRs = make([]string, len(values))
	copy(b.PodCIDRs, values)
	return b
}

// WithProviderID sets the ProviderID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ProviderID field is set to the value of the last call.
func (b *NodeSpecApplyConfiguration) WithProviderID(value string) *NodeSpecApplyConfiguration {
	b.ProviderID = &value
	return b
}

// WithUnschedulable sets the Unschedulable field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Unschedulable field is set to the value of the last call.
func (b *NodeSpecApplyConfiguration) WithUnschedulable(value bool) *NodeSpecApplyConfiguration {
	b.Unschedulable = &value
	return b
}

// WithTaints adds the given value to the Taints field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Taints field.
// Deprecated: WithTaints does not replace existing list for atomic list type. Use WithNewTaints instead.
func (b *NodeSpecApplyConfiguration) WithTaints(values ...*TaintApplyConfiguration) *NodeSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithTaints")
		}
		b.Taints = append(b.Taints, *values[i])
	}
	return b
}

// WithNewTaints replaces the Taints field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the Taints field is set to the values of the last call.
func (b *NodeSpecApplyConfiguration) WithNewTaints(values ...*TaintApplyConfiguration) *NodeSpecApplyConfiguration {
	b.Taints = nil
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNewTaints")
		}
		b.Taints = append(b.Taints, *values[i])
	}
	return b
}

// WithConfigSource sets the ConfigSource field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ConfigSource field is set to the value of the last call.
func (b *NodeSpecApplyConfiguration) WithConfigSource(value *NodeConfigSourceApplyConfiguration) *NodeSpecApplyConfiguration {
	b.ConfigSource = value
	return b
}

// WithDoNotUseExternalID sets the DoNotUseExternalID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DoNotUseExternalID field is set to the value of the last call.
func (b *NodeSpecApplyConfiguration) WithDoNotUseExternalID(value string) *NodeSpecApplyConfiguration {
	b.DoNotUseExternalID = &value
	return b
}
