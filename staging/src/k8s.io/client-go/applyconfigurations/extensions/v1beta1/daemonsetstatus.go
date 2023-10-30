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

// DaemonSetStatusApplyConfiguration represents an declarative configuration of the DaemonSetStatus type for use
// with apply.
type DaemonSetStatusApplyConfiguration struct {
	CurrentNumberScheduled *int32                                 `json:"currentNumberScheduled,omitempty"`
	NumberMisscheduled     *int32                                 `json:"numberMisscheduled,omitempty"`
	DesiredNumberScheduled *int32                                 `json:"desiredNumberScheduled,omitempty"`
	NumberReady            *int32                                 `json:"numberReady,omitempty"`
	ObservedGeneration     *int64                                 `json:"observedGeneration,omitempty"`
	UpdatedNumberScheduled *int32                                 `json:"updatedNumberScheduled,omitempty"`
	NumberAvailable        *int32                                 `json:"numberAvailable,omitempty"`
	NumberUnavailable      *int32                                 `json:"numberUnavailable,omitempty"`
	CollisionCount         *int32                                 `json:"collisionCount,omitempty"`
	Conditions             []DaemonSetConditionApplyConfiguration `json:"conditions,omitempty"`
}

// DaemonSetStatusApplyConfiguration constructs an declarative configuration of the DaemonSetStatus type for use with
// apply.
func DaemonSetStatus() *DaemonSetStatusApplyConfiguration {
	return &DaemonSetStatusApplyConfiguration{}
}

// WithCurrentNumberScheduled sets the CurrentNumberScheduled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CurrentNumberScheduled field is set to the value of the last call.
func (b *DaemonSetStatusApplyConfiguration) WithCurrentNumberScheduled(value int32) *DaemonSetStatusApplyConfiguration {
	b.CurrentNumberScheduled = &value
	return b
}

// WithNumberMisscheduled sets the NumberMisscheduled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NumberMisscheduled field is set to the value of the last call.
func (b *DaemonSetStatusApplyConfiguration) WithNumberMisscheduled(value int32) *DaemonSetStatusApplyConfiguration {
	b.NumberMisscheduled = &value
	return b
}

// WithDesiredNumberScheduled sets the DesiredNumberScheduled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DesiredNumberScheduled field is set to the value of the last call.
func (b *DaemonSetStatusApplyConfiguration) WithDesiredNumberScheduled(value int32) *DaemonSetStatusApplyConfiguration {
	b.DesiredNumberScheduled = &value
	return b
}

// WithNumberReady sets the NumberReady field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NumberReady field is set to the value of the last call.
func (b *DaemonSetStatusApplyConfiguration) WithNumberReady(value int32) *DaemonSetStatusApplyConfiguration {
	b.NumberReady = &value
	return b
}

// WithObservedGeneration sets the ObservedGeneration field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ObservedGeneration field is set to the value of the last call.
func (b *DaemonSetStatusApplyConfiguration) WithObservedGeneration(value int64) *DaemonSetStatusApplyConfiguration {
	b.ObservedGeneration = &value
	return b
}

// WithUpdatedNumberScheduled sets the UpdatedNumberScheduled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UpdatedNumberScheduled field is set to the value of the last call.
func (b *DaemonSetStatusApplyConfiguration) WithUpdatedNumberScheduled(value int32) *DaemonSetStatusApplyConfiguration {
	b.UpdatedNumberScheduled = &value
	return b
}

// WithNumberAvailable sets the NumberAvailable field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NumberAvailable field is set to the value of the last call.
func (b *DaemonSetStatusApplyConfiguration) WithNumberAvailable(value int32) *DaemonSetStatusApplyConfiguration {
	b.NumberAvailable = &value
	return b
}

// WithNumberUnavailable sets the NumberUnavailable field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NumberUnavailable field is set to the value of the last call.
func (b *DaemonSetStatusApplyConfiguration) WithNumberUnavailable(value int32) *DaemonSetStatusApplyConfiguration {
	b.NumberUnavailable = &value
	return b
}

// WithCollisionCount sets the CollisionCount field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CollisionCount field is set to the value of the last call.
func (b *DaemonSetStatusApplyConfiguration) WithCollisionCount(value int32) *DaemonSetStatusApplyConfiguration {
	b.CollisionCount = &value
	return b
}

// WithConditions adds the given value to the Conditions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Conditions field.
// Deprecated: WithConditions does not replace existing list for atomic list type. Use WithNewConditions instead.
func (b *DaemonSetStatusApplyConfiguration) WithConditions(values ...*DaemonSetConditionApplyConfiguration) *DaemonSetStatusApplyConfiguration {
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
func (b *DaemonSetStatusApplyConfiguration) WithNewConditions(values ...*DaemonSetConditionApplyConfiguration) *DaemonSetStatusApplyConfiguration {
	b.Conditions = nil
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNewConditions")
		}
		b.Conditions = append(b.Conditions, *values[i])
	}
	return b
}
