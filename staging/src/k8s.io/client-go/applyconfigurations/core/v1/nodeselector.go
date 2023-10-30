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

// NodeSelectorApplyConfiguration represents an declarative configuration of the NodeSelector type for use
// with apply.
type NodeSelectorApplyConfiguration struct {
	NodeSelectorTerms []NodeSelectorTermApplyConfiguration `json:"nodeSelectorTerms,omitempty"`
}

// NodeSelectorApplyConfiguration constructs an declarative configuration of the NodeSelector type for use with
// apply.
func NodeSelector() *NodeSelectorApplyConfiguration {
	return &NodeSelectorApplyConfiguration{}
}

// WithNodeSelectorTerms adds the given value to the NodeSelectorTerms field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the NodeSelectorTerms field.
// Deprecated: WithNodeSelectorTerms does not replace existing list for atomic list type. Use WithNewNodeSelectorTerms instead.
func (b *NodeSelectorApplyConfiguration) WithNodeSelectorTerms(values ...*NodeSelectorTermApplyConfiguration) *NodeSelectorApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNodeSelectorTerms")
		}
		b.NodeSelectorTerms = append(b.NodeSelectorTerms, *values[i])
	}
	return b
}

// WithNewNodeSelectorTerms replaces the NodeSelectorTerms field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the NodeSelectorTerms field is set to the values of the last call.
func (b *NodeSelectorApplyConfiguration) WithNewNodeSelectorTerms(values ...*NodeSelectorTermApplyConfiguration) *NodeSelectorApplyConfiguration {
	b.NodeSelectorTerms = nil
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNewNodeSelectorTerms")
		}
		b.NodeSelectorTerms = append(b.NodeSelectorTerms, *values[i])
	}
	return b
}
