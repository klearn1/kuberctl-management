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

package v1beta2

import (
	v1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	v1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

// StatefulSetSpecApplyConfiguration represents a declarative configuration of the StatefulSetSpec type for use
// with apply.
type StatefulSetSpecApplyConfiguration struct {
	Replicas                             *int32                                                             `json:"replicas,omitempty"`
	Selector                             *v1.LabelSelectorApplyConfiguration                                `json:"selector,omitempty"`
	Template                             *corev1.PodTemplateSpecApplyConfiguration                          `json:"template,omitempty"`
	VolumeClaimTemplates                 []corev1.PersistentVolumeClaimApplyConfiguration                   `json:"volumeClaimTemplates,omitempty"`
	ServiceName                          *string                                                            `json:"serviceName,omitempty"`
	PodManagementPolicy                  *v1beta2.PodManagementPolicyType                                   `json:"podManagementPolicy,omitempty"`
	UpdateStrategy                       *StatefulSetUpdateStrategyApplyConfiguration                       `json:"updateStrategy,omitempty"`
	RevisionHistoryLimit                 *int32                                                             `json:"revisionHistoryLimit,omitempty"`
	MinReadySeconds                      *int32                                                             `json:"minReadySeconds,omitempty"`
	PersistentVolumeClaimRetentionPolicy *StatefulSetPersistentVolumeClaimRetentionPolicyApplyConfiguration `json:"persistentVolumeClaimRetentionPolicy,omitempty"`
	Ordinals                             *StatefulSetOrdinalsApplyConfiguration                             `json:"ordinals,omitempty"`
	VolumeClaimUpdatePolicy              *v1beta2.StatefulSetVolumeClaimUpdatePolicyType                    `json:"volumeClaimUpdatePolicy,omitempty"`
}

// StatefulSetSpecApplyConfiguration constructs a declarative configuration of the StatefulSetSpec type for use with
// apply.
func StatefulSetSpec() *StatefulSetSpecApplyConfiguration {
	return &StatefulSetSpecApplyConfiguration{}
}

// WithReplicas sets the Replicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Replicas field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithReplicas(value int32) *StatefulSetSpecApplyConfiguration {
	b.Replicas = &value
	return b
}

// WithSelector sets the Selector field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Selector field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithSelector(value *v1.LabelSelectorApplyConfiguration) *StatefulSetSpecApplyConfiguration {
	b.Selector = value
	return b
}

// WithTemplate sets the Template field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Template field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithTemplate(value *corev1.PodTemplateSpecApplyConfiguration) *StatefulSetSpecApplyConfiguration {
	b.Template = value
	return b
}

// WithVolumeClaimTemplates adds the given value to the VolumeClaimTemplates field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the VolumeClaimTemplates field.
func (b *StatefulSetSpecApplyConfiguration) WithVolumeClaimTemplates(values ...*corev1.PersistentVolumeClaimApplyConfiguration) *StatefulSetSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithVolumeClaimTemplates")
		}
		b.VolumeClaimTemplates = append(b.VolumeClaimTemplates, *values[i])
	}
	return b
}

// WithServiceName sets the ServiceName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ServiceName field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithServiceName(value string) *StatefulSetSpecApplyConfiguration {
	b.ServiceName = &value
	return b
}

// WithPodManagementPolicy sets the PodManagementPolicy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PodManagementPolicy field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithPodManagementPolicy(value v1beta2.PodManagementPolicyType) *StatefulSetSpecApplyConfiguration {
	b.PodManagementPolicy = &value
	return b
}

// WithUpdateStrategy sets the UpdateStrategy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UpdateStrategy field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithUpdateStrategy(value *StatefulSetUpdateStrategyApplyConfiguration) *StatefulSetSpecApplyConfiguration {
	b.UpdateStrategy = value
	return b
}

// WithRevisionHistoryLimit sets the RevisionHistoryLimit field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RevisionHistoryLimit field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithRevisionHistoryLimit(value int32) *StatefulSetSpecApplyConfiguration {
	b.RevisionHistoryLimit = &value
	return b
}

// WithMinReadySeconds sets the MinReadySeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MinReadySeconds field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithMinReadySeconds(value int32) *StatefulSetSpecApplyConfiguration {
	b.MinReadySeconds = &value
	return b
}

// WithPersistentVolumeClaimRetentionPolicy sets the PersistentVolumeClaimRetentionPolicy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PersistentVolumeClaimRetentionPolicy field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithPersistentVolumeClaimRetentionPolicy(value *StatefulSetPersistentVolumeClaimRetentionPolicyApplyConfiguration) *StatefulSetSpecApplyConfiguration {
	b.PersistentVolumeClaimRetentionPolicy = value
	return b
}

// WithOrdinals sets the Ordinals field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Ordinals field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithOrdinals(value *StatefulSetOrdinalsApplyConfiguration) *StatefulSetSpecApplyConfiguration {
	b.Ordinals = value
	return b
}

// WithVolumeClaimUpdatePolicy sets the VolumeClaimUpdatePolicy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VolumeClaimUpdatePolicy field is set to the value of the last call.
func (b *StatefulSetSpecApplyConfiguration) WithVolumeClaimUpdatePolicy(value v1beta2.StatefulSetVolumeClaimUpdatePolicyType) *StatefulSetSpecApplyConfiguration {
	b.VolumeClaimUpdatePolicy = &value
	return b
}
