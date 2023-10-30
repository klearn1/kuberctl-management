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

// ReplicationControllerSpecApplyConfiguration represents an declarative configuration of the ReplicationControllerSpec type for use
// with apply.
type ReplicationControllerSpecApplyConfiguration struct {
	Replicas        *int32                             `json:"replicas,omitempty"`
	MinReadySeconds *int32                             `json:"minReadySeconds,omitempty"`
	Selector        map[string]string                  `json:"selector,omitempty"`
	Template        *PodTemplateSpecApplyConfiguration `json:"template,omitempty"`
}

// ReplicationControllerSpecApplyConfiguration constructs an declarative configuration of the ReplicationControllerSpec type for use with
// apply.
func ReplicationControllerSpec() *ReplicationControllerSpecApplyConfiguration {
	return &ReplicationControllerSpecApplyConfiguration{}
}

// WithReplicas sets the Replicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Replicas field is set to the value of the last call.
func (b *ReplicationControllerSpecApplyConfiguration) WithReplicas(value int32) *ReplicationControllerSpecApplyConfiguration {
	b.Replicas = &value
	return b
}

// WithMinReadySeconds sets the MinReadySeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MinReadySeconds field is set to the value of the last call.
func (b *ReplicationControllerSpecApplyConfiguration) WithMinReadySeconds(value int32) *ReplicationControllerSpecApplyConfiguration {
	b.MinReadySeconds = &value
	return b
}

// WithSelector puts the entries into the Selector field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Selector field,
// overwriting an existing map entries in Selector field with the same key.
// Deprecated: WithSelector does not replace existing map for atomic map type. Use WithNewSelector instead.
func (b *ReplicationControllerSpecApplyConfiguration) WithSelector(entries map[string]string) *ReplicationControllerSpecApplyConfiguration {
	if b.Selector == nil && len(entries) > 0 {
		b.Selector = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Selector[k] = v
	}
	return b
}

// WithNewSelector replaces the Selector field in the declarative configuration with given entries
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the Selector field is set to the entries of the last call.
func (b *ReplicationControllerSpecApplyConfiguration) WithNewSelector(entries map[string]string) *ReplicationControllerSpecApplyConfiguration {
	b.Selector = make(map[string]string, len(entries))
	for k, v := range entries {
		b.Selector[k] = v
	}
	return b
}

// WithTemplate sets the Template field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Template field is set to the value of the last call.
func (b *ReplicationControllerSpecApplyConfiguration) WithTemplate(value *PodTemplateSpecApplyConfiguration) *ReplicationControllerSpecApplyConfiguration {
	b.Template = value
	return b
}
