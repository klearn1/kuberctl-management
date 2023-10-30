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

import (
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	v1 "k8s.io/client-go/applyconfigurations/admissionregistration/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

// MutatingWebhookApplyConfiguration represents an declarative configuration of the MutatingWebhook type for use
// with apply.
type MutatingWebhookApplyConfiguration struct {
	Name                    *string                                              `json:"name,omitempty"`
	ClientConfig            *WebhookClientConfigApplyConfiguration               `json:"clientConfig,omitempty"`
	Rules                   []v1.RuleWithOperationsApplyConfiguration            `json:"rules,omitempty"`
	FailurePolicy           *admissionregistrationv1beta1.FailurePolicyType      `json:"failurePolicy,omitempty"`
	MatchPolicy             *admissionregistrationv1beta1.MatchPolicyType        `json:"matchPolicy,omitempty"`
	NamespaceSelector       *metav1.LabelSelectorApplyConfiguration              `json:"namespaceSelector,omitempty"`
	ObjectSelector          *metav1.LabelSelectorApplyConfiguration              `json:"objectSelector,omitempty"`
	SideEffects             *admissionregistrationv1beta1.SideEffectClass        `json:"sideEffects,omitempty"`
	TimeoutSeconds          *int32                                               `json:"timeoutSeconds,omitempty"`
	AdmissionReviewVersions []string                                             `json:"admissionReviewVersions,omitempty"`
	ReinvocationPolicy      *admissionregistrationv1beta1.ReinvocationPolicyType `json:"reinvocationPolicy,omitempty"`
	MatchConditions         []MatchConditionApplyConfiguration                   `json:"matchConditions,omitempty"`
}

// MutatingWebhookApplyConfiguration constructs an declarative configuration of the MutatingWebhook type for use with
// apply.
func MutatingWebhook() *MutatingWebhookApplyConfiguration {
	return &MutatingWebhookApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *MutatingWebhookApplyConfiguration) WithName(value string) *MutatingWebhookApplyConfiguration {
	b.Name = &value
	return b
}

// WithClientConfig sets the ClientConfig field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ClientConfig field is set to the value of the last call.
func (b *MutatingWebhookApplyConfiguration) WithClientConfig(value *WebhookClientConfigApplyConfiguration) *MutatingWebhookApplyConfiguration {
	b.ClientConfig = value
	return b
}

// WithRules adds the given value to the Rules field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Rules field.
// Deprecated: WithRules does not replace existing list for atomic list type. Use WithNewRules instead.
func (b *MutatingWebhookApplyConfiguration) WithRules(values ...*v1.RuleWithOperationsApplyConfiguration) *MutatingWebhookApplyConfiguration {
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
func (b *MutatingWebhookApplyConfiguration) WithNewRules(values ...*v1.RuleWithOperationsApplyConfiguration) *MutatingWebhookApplyConfiguration {
	b.Rules = nil
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNewRules")
		}
		b.Rules = append(b.Rules, *values[i])
	}
	return b
}

// WithFailurePolicy sets the FailurePolicy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FailurePolicy field is set to the value of the last call.
func (b *MutatingWebhookApplyConfiguration) WithFailurePolicy(value admissionregistrationv1beta1.FailurePolicyType) *MutatingWebhookApplyConfiguration {
	b.FailurePolicy = &value
	return b
}

// WithMatchPolicy sets the MatchPolicy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MatchPolicy field is set to the value of the last call.
func (b *MutatingWebhookApplyConfiguration) WithMatchPolicy(value admissionregistrationv1beta1.MatchPolicyType) *MutatingWebhookApplyConfiguration {
	b.MatchPolicy = &value
	return b
}

// WithNamespaceSelector sets the NamespaceSelector field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NamespaceSelector field is set to the value of the last call.
func (b *MutatingWebhookApplyConfiguration) WithNamespaceSelector(value *metav1.LabelSelectorApplyConfiguration) *MutatingWebhookApplyConfiguration {
	b.NamespaceSelector = value
	return b
}

// WithObjectSelector sets the ObjectSelector field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ObjectSelector field is set to the value of the last call.
func (b *MutatingWebhookApplyConfiguration) WithObjectSelector(value *metav1.LabelSelectorApplyConfiguration) *MutatingWebhookApplyConfiguration {
	b.ObjectSelector = value
	return b
}

// WithSideEffects sets the SideEffects field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SideEffects field is set to the value of the last call.
func (b *MutatingWebhookApplyConfiguration) WithSideEffects(value admissionregistrationv1beta1.SideEffectClass) *MutatingWebhookApplyConfiguration {
	b.SideEffects = &value
	return b
}

// WithTimeoutSeconds sets the TimeoutSeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TimeoutSeconds field is set to the value of the last call.
func (b *MutatingWebhookApplyConfiguration) WithTimeoutSeconds(value int32) *MutatingWebhookApplyConfiguration {
	b.TimeoutSeconds = &value
	return b
}

// WithAdmissionReviewVersions adds the given value to the AdmissionReviewVersions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the AdmissionReviewVersions field.
// Deprecated: WithAdmissionReviewVersions does not replace existing list for atomic list type. Use WithNewAdmissionReviewVersions instead.
func (b *MutatingWebhookApplyConfiguration) WithAdmissionReviewVersions(values ...string) *MutatingWebhookApplyConfiguration {
	for i := range values {
		b.AdmissionReviewVersions = append(b.AdmissionReviewVersions, values[i])
	}
	return b
}

// WithNewAdmissionReviewVersions replaces the AdmissionReviewVersions field in the declarative configuration with given values
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the AdmissionReviewVersions field is set to the values of the last call.
func (b *MutatingWebhookApplyConfiguration) WithNewAdmissionReviewVersions(values ...string) *MutatingWebhookApplyConfiguration {
	b.AdmissionReviewVersions = make([]string, len(values))
	copy(b.AdmissionReviewVersions, values)
	return b
}

// WithReinvocationPolicy sets the ReinvocationPolicy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ReinvocationPolicy field is set to the value of the last call.
func (b *MutatingWebhookApplyConfiguration) WithReinvocationPolicy(value admissionregistrationv1beta1.ReinvocationPolicyType) *MutatingWebhookApplyConfiguration {
	b.ReinvocationPolicy = &value
	return b
}

// WithMatchConditions adds the given value to the MatchConditions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the MatchConditions field.
func (b *MutatingWebhookApplyConfiguration) WithMatchConditions(values ...*MatchConditionApplyConfiguration) *MutatingWebhookApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithMatchConditions")
		}
		b.MatchConditions = append(b.MatchConditions, *values[i])
	}
	return b
}
