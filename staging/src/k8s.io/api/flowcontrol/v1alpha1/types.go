/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// These are valid flow-dingtinguisher methods.
const (
	// FlowDistinguisherMethodByUserType specifies that the flow distinguisher is the username works in the request context
	// so that the requests from the same user will enqueued into the same set of queues for processing.
	FlowDistinguisherMethodByUserType FlowDistinguisherMethodType = "ByUser"

	// FlowDistinguisherMethodByNamespaceType specifies that the flow distinguisher is the object's requested namespace
	// that the request accesses.  If the object is not namespaced or the request is a non-resource request then
	// the flow distinguisher is the empty string.
	FlowDistinguisherMethodByNamespaceType FlowDistinguisherMethodType = "ByNamespace"

	GroupKind          = "Group"
	ServiceAccountKind = "ServiceAccount"
	UserKind           = "User"

	APIGroupAll    = "*"
	ResourceAll    = "*"
	VerbAll        = "*"
	NonResourceAll = "*"
	NameAll        = "*"
)

// System preset priority level names
const (
	PriorityLevelConfigurationNameExempt = "exempt"
)

// Conditions
const (
	FlowSchemaConditionDangling = "Dangling"

	PriorityLevelConfigurationConditionConcurrencyShared = "ConcurrencyShared"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FlowSchema defines the schema of a group of flows. Note that a flow is made up of a set of inbound API requests with
// similar attributes and is identified by a pair of strings: the name of the FlowSchema and a "flow distinguisher".
type FlowSchema struct {
	metav1.TypeMeta `json:",inline"`
	// `metadata` is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// `spec` is the specification of the desired behavior of a flow-schema.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Spec FlowSchemaSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// `status` is the current status of a flow-schema.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Status FlowSchemaStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FlowSchemaList is a list of FlowSchema objects.
type FlowSchemaList struct {
	metav1.TypeMeta `json:",inline"`
	// `metadata` is the standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// `items` is a list of flow-schemas.
	Items []FlowSchema `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// FlowSchemaSpec describes how the flow-schema's specification looks like.
type FlowSchemaSpec struct {
	// `priorityLevelConfiguration` should reference a PriorityLevelConfiguration in the cluster. If the reference cannot
	// be resolved, the flow-schema will be ignored and marked as invalid in its status.
	// Required.
	PriorityLevelConfiguration PriorityLevelConfigurationReference `json:"priorityLevelConfiguration,omitempty" protobuf:"bytes,1,opt,name=priorityLevelConfiguration"`
	// `matchingPrecedence` is used to choose among the FlowSchemas that match a given request.  The chosen
	// FlowSchema is among those with the numerically lowest (which we take to be logically highest)
	// MatchingPrecedence.  Each MatchingPrecedence value must be non-negative.
	MatchingPrecedence int32 `json:"matchingPrecedence" protobuf:"varint,2,opt,name=matchingPrecedence"`
	// `distinguisherMethod` defines how to compute the flow distinguisher for requests that match this schema.
	// `nil` specifies that the distinguisher is disabled and thus will always be the empty string.
	DistinguisherMethod *FlowDistinguisherMethod `json:"distinguisherMethod,omitempty" protobuf:"bytes,3,opt,name=distinguisherMethod"`
	// `rules` describes which requests will match this flow schema. The flow-schema applies if any rule matches.
	Rules []PolicyRuleWithSubjects `json:"rules,omitempty" protobuf:"bytes,4,rep,name=rules"`
}

// FlowDistinguisherMethodType is the type of flow distinguisher method
type FlowDistinguisherMethodType string

// FlowDistinguisherMethod specifies the method of a flow distinguisher.
type FlowDistinguisherMethod struct {
	// `type` is the type of flow distinguisher method
	Type FlowDistinguisherMethodType `json:"type,omitempty" protobuf:"bytes,1,opt,name=type"`
}

// PriorityLevelConfigurationReference contains information that points to the "request-priority" being used.
type PriorityLevelConfigurationReference struct {
	// `name` is the name of resource being referenced
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
}

// PolicyRuleWithSubjects prescribes a test that applies to a request to an apiserver. The test considers the subject
// making the request, the verb being requested, and the resource to be acted upon.
type PolicyRuleWithSubjects struct {
	// `subjects` is the list of normal user, serviceaccount, or group that this rule cares about.
	// +optional
	Subjects []Subject `json:"subjects,omitempty" protobuf:"bytes,1,opt,name=subjects"`
	// `rule` is the target verb, resource or the subresource the rule cares about. APIGroups, Resources, etc.
	// Required.
	Rule PolicyRule `json:"rule" protobuf:"bytes,2,opt,name=rule"`
}

// PolicyRule holds information that describes a policy rule, but does not contain information
// about who the rule applies to or which namespace the rule applies to.
type PolicyRule struct {
	// `verbs` is a list of Verbs that apply to ALL the ResourceKinds and AttributeRestrictions contained in this rule.
	// VerbAll represents all verbs.
	Verbs []string `json:"verbs" protobuf:"bytes,1,rep,name=verbs"`
	// `apiGroups` is the name of the APIGroup that contains the resources.  If multiple API groups are specified, any action requested against one of
	// the enumerated resources in any API group will be allowed. '*' represents all api groups.
	// +optional
	APIGroups []string `json:"apiGroups,omitempty" protobuf:"bytes,2,rep,name=apiGroups"`
	// `resources` is a list of resources this rule applies to.  ResourceAll represents all resources.
	// +optional
	Resources []string `json:"resources,omitempty" protobuf:"bytes,3,rep,name=resources"`
	// `nonResourceURLs` is a set of partial urls that a user should have access to.  *s are allowed, but only as the full, final step in the path
	// Since non-resource URLs are not namespaced, this field is only applicable for ClusterRoles referenced from a ClusterRoleBinding.
	// Rules can either apply to API resources (such as "pods" or "secrets") or non-resource URL paths (such as "/api"),  but not both.
	// NonResourceAll represents all non-resource urls.
	// +optional
	NonResourceURLs []string `json:"nonResourceURLs,omitempty" protobuf:"bytes,6,rep,name=nonResourceURLs"`
}

// FlowSchemaStatus represents the current state of a flow-schema.
type FlowSchemaStatus struct {
	// Current state of flow-schema.
	Conditions []FlowSchemaCondition `json:"conditions,omitempty" protobuf:"bytes,1,rep,name=conditions"`
}

// FlowSchemaCondition describes conditions for a flow-schema.
type FlowSchemaCondition struct {
	// `type` is the type of the condition.
	Type FlowSchemaConditionType `json:"type,omitempty" protobuf:"bytes,1,opt,name=type"`
	// `status` is the status of the condition.
	// Can be True, False, Unknown.
	Status corev1.ConditionStatus `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
	// `lastTransitionTime` is the last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,3,opt,name=lastTransitionTime"`
	// `reason` is a unique, one-word, CamelCase reason for the condition's last transition.
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// `message` is a human-readable message indicating details about last transition.
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// FlowSchemaConditionType is a valid value for FlowSchemaStatusCondition.Type
type FlowSchemaConditionType string

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PriorityLevelConfiguration represents the configuration of a priority level.
type PriorityLevelConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// `metadata` is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// `spec` is the specification of the desired behavior of a "request-priority".
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Spec PriorityLevelConfigurationSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// `status` is the current status of a "request-priority".
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
	// +optional
	Status PriorityLevelConfigurationStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PriorityLevelConfigurationList is a list of PriorityLevelConfiguration objects.
type PriorityLevelConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	// `metadata` is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// `items` is a list of request-priorities.
	Items []PriorityLevelConfiguration `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// PriorityLevelConfigurationSpec is specification of a priority level
type PriorityLevelConfigurationSpec struct {
	// `assuredConcurrencyShares` is a positive number for a non-exempt priority level. The concurrency limit of an apiserver
	// is divided among the non-exempt priority levels in proportion to their assured concurrency shares. Basically this produces
	// the assured concurrency value (ACV) for each priority level:
	//
	//             ACV(l) = ceil( SCL * ACS(l) / ( sum[priority levels k] ACS(k) ) )
	//
	// +optional
	AssuredConcurrencyShares int32 `json:"assuredConcurrencyShares,omitempty" protobuf:"varint,1,opt,name=assuredConcurrencyShares"`
	// `queues` is a number of queues that belong to a non-exempt PriorityLevelConfiguration object. The queues exist
	// independently at each apiserver. The value must be positive for a non-exempt priority level and setting it to 1
	// disables shufflesharding and makes the distinguisher method irrelevant.
	// TODO: sugguest a default or a way of deciding on a value.
	// +optional
	Queues int32 `json:"queues,omitempty" protobuf:"varint,2,opt,name=queues"`
	// `handSize` is a small positive number for applying shuffle sharding. When a request arrives at an apiserver the
	// request flow identifier’s string pair is hashed and the hash value is used to shuffle the queue indices and deal
	// a hand of the size specified here. If empty, the hand size will the be set to 1.
	// NOTE: To figure out a better value for your cluster, please refer to (#76846)[https://github.com/kubernetes/kubernetes/issues/76846#issuecomment-523700960]
	// +optional
	HandSize int32 `json:"handSize,omitempty" protobuf:"varint,3,opt,name=handSize"`
	// `queueLengthLimit` is a length limit applied to each queue belongs to the priority.  The value must be positive
	// for a non-exempt priority level.
	// +optional
	QueueLengthLimit int32 `json:"queueLengthLimit,omitempty" protobuf:"varint,4,opt,name=queueLengthLimit"`
	// `exempt` defines whether the priority level is exempted or not.  There should be at most one exempt priority level.
	// Being exempt means that requests of that priority are not subject to concurrency limits (and thus are never queued)
	// and do not detract from the concurrency available for non-exempt requests. The field is default to false and only those system
	// preset priority level can be exempt.
	// +optional
	Exempt bool `json:"exempt,omitempty" protobuf:"varint,5,opt,name=exempt"`
}

// PriorityLevelConfigurationConditionType is a valid value for PriorityLevelConfigurationStatusCondition.Type
type PriorityLevelConfigurationConditionType string

// PriorityLevelConfigurationStatus represents the current state of a "request-priority".
type PriorityLevelConfigurationStatus struct {
	// `conditions` is the current state of "request-priority".
	Conditions []PriorityLevelConfigurationCondition `json:"conditions,omitempty" protobuf:"bytes,1,rep,name=conditions"`
}

// PriorityLevelConfigurationCondition defines the condition of priority level.
type PriorityLevelConfigurationCondition struct {
	// `type` is the type of the condition.
	Type PriorityLevelConfigurationConditionType `json:"type,omitempty" protobuf:"bytes,1,opt,name=type"`
	// `status` is the status of the condition.
	// Can be True, False, Unknown.
	Status corev1.ConditionStatus `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
	// `lastTransitionTime` is the last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,3,opt,name=lastTransitionTime"`
	// `reason` is a unique, one-word, CamelCase reason for the condition's last transition.
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// `message` is a human-readable message indicating details about last transition.
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// Subject matches a set of users.
// Syntactically, Subject is a general API object reference.
// Authorization produces a username and a set of groups, and we imagine special kinds of non-namespaced objects,
// User and Group in API group "flowcontrol.apiserver.k8s.io", to represent such a username or group.
// The only kind of true object reference that currently will match any users is ServiceAccount.
type Subject struct {
	// `kind` of object being referenced. Values defined by this API group are "User", "Group", and "ServiceAccount".
	// If the kind value is not recognized, the flow-control layer in api-server should report an error.
	Kind string `json:"kind" protobuf:"bytes,1,opt,name=kind"`
	// `apiGroup` holds the API group of the referenced subject.
	// Defaults to "" for ServiceAccount subjects.
	// Defaults to "flowcontrol.apiserver.k8s.io" for User and Group subjects.
	// +optional
	APIGroup string `json:"apiGroup,omitempty" protobuf:"bytes,2,opt.name=apiGroup"`
	// `name` of the object being referenced.  To match regardless of name, use NameAll.
	// Required.
	Name string `json:"name" protobuf:"bytes,3,opt,name=name"`
	// `namespace` of the referenced object.  If the object kind is non-namespace, such as "User" or "Group", and this value is not empty
	// the Authorizer should report an error.
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,4,opt,name=namespace"`
}
