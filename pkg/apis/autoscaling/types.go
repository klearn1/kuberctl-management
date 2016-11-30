/*
Copyright 2016 The Kubernetes Authors.

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

package autoscaling

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
)

// Scale represents a scaling request for a resource.
type Scale struct {
	metav1.TypeMeta
	// Standard object metadata; More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta

	// defines the behavior of the scale. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status.
	// +optional
	Spec ScaleSpec

	// current status of the scale. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status. Read-only.
	// +optional
	Status ScaleStatus
}

// ScaleSpec describes the attributes of a scale subresource.
type ScaleSpec struct {
	// desired number of instances for the scaled object.
	// +optional
	Replicas int32
}

// ScaleStatus represents the current status of a scale subresource.
type ScaleStatus struct {
	// actual number of observed instances of the scaled object.
	Replicas int32

	// label query over pods that should match the replicas count. This is same
	// as the label selector but in the string format to avoid introspection
	// by clients. The string will be in the same format as the query-param syntax.
	// More info: http://kubernetes.io/docs/user-guide/labels#label-selectors
	// +optional
	Selector string
}

// CrossVersionObjectReference contains enough information to let you identify the referred resource.
type CrossVersionObjectReference struct {
	// Kind of the referent; More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#types-kinds"
	Kind string
	// Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string
	// API version of the referent
	// +optional
	APIVersion string
}

// specification of a horizontal pod autoscaler.
type HorizontalPodAutoscalerSpec struct {
	// reference to scaled resource; horizontal pod autoscaler will learn the current resource consumption
	// and will set the desired number of pods by using its Scale subresource.
	ScaleTargetRef CrossVersionObjectReference
	// lower limit for the number of pods that can be set by the autoscaler, default 1.
	// +optional
	MinReplicas *int32
	// upper limit for the number of pods that can be set by the autoscaler. It cannot be smaller than MinReplicas.
	MaxReplicas int32
	// the metrics to use to calculate the desired replica count (the
	// maximum replica count across all metrics will be used).	The
	// desired replica count is calculated multiplying the ratio between
	// the target value and the current value by the current number of
	// pods.  Ergo, metrics used must decrease as the pod count is
	// increased, and vice-versa.  See the individual metric source
	// types for more information about how each type of metric
	// must respond.
	// +optional
	Metrics []MetricSpec
}

// a type of metric source
type MetricSourceType string

var (
	// a metric describing a kubernetes object (for example, hits-per-second on an Ingress object)
	ObjectSourceType MetricSourceType = "Object"
	// a metric describing each pod in the current scale target (for example, transactions-processed-per-second).
	// The values will be averaged together before being compared to the target value
	PodsSourceType MetricSourceType = "Pods"
	// a resource metric known to Kubernetes, as specified in requests and limits, describing each pod
	// in the current scale target (e.g. CPU or memory).  Such metrics are built in to Kubernetes,
	// and have special scaling options on top of those available to normal per-pod metrics (the "pods" source)
	ResourceSourceType MetricSourceType = "Resource"
)

// a specification for how to scale based on a single metric
// (only `type` and one other matching field should be set at once)
type MetricSpec struct {
	// the type of metric source (should match one of the fields below)
	Type MetricSourceType

	// a metric describing a single kubernetes object (for example, hits-per-second on an Ingress object)
	// +optional
	Object *ObjectMetricSource
	// a metric describing each pod in the current scale target (for example, transactions-processed-per-second).
	// The values will be averaged together before being compared to the target value
	// +optional
	Pods *PodsMetricSource
	// a resource metric (such as those specified in requests and limits) known to Kubernetes
	// describing each pod in the current scale target (e.g. CPU or memory). Such metrics are
	// built in to Kubernetes, and have special scaling options on top of those available to
	// normal per-pod metrics using the "pods" source.
	// +optional
	Resource *ResourceMetricSource
}

// a metric describing a single kubernetes object (for example, hits-per-second on an Ingress object)
type ObjectMetricSource struct {
	// the described Kubernetes object
	Target CrossVersionObjectReference

	// the name of the metric in question
	MetricName string
	// the target value of the metric (as a quantity)
	TargetValue resource.Quantity
}

// a metric describing each pod in the current scale target (for example, transactions-processed-per-second).
// The values will be averaged together before being compared to the target value
type PodsMetricSource struct {
	// the name of the metric in question
	MetricName string
	// the target value of the metric (as a quantity)
	TargetAverageValue resource.Quantity
}

// a resource metric known to Kubernetes, as specified in requests and limits, describing each pod
// in the current scale target (e.g. CPU or memory).  The values will be averaged together before
// being compared to the target.  Such metrics are built in to Kubernetes, and have special
// scaling options on top of those available to normal per-pod metrics using the "pods" source.
// Only one "target" type should be set.
type ResourceMetricSource struct {
	// the name of the resource in question
	Name api.ResourceName
	// the target value of the resource metric, represented as
	// a percentage of the requested value of the resource on the pods.
	// +optional
	TargetAverageUtilization *int32
	// the target value of the resource metric as a raw value, similarly
	// to the "pods" metric source type.
	// +optional
	TargetAverageValue *resource.Quantity
}

// current status of a horizontal pod autoscaler
type HorizontalPodAutoscalerStatus struct {
	// most recent generation observed by this autoscaler.
	// +optional
	ObservedGeneration *int64

	// last time the HorizontalPodAutoscaler scaled the number of pods;
	// used by the autoscaler to control how often the number of pods is changed.
	// +optional
	LastScaleTime *metav1.Time

	// current number of replicas of pods managed by this autoscaler.
	CurrentReplicas int32

	// desired number of replicas of pods managed by this autoscaler.
	DesiredReplicas int32

	// the last read state of the metrics used by this autoscaler
	CurrentMetrics []MetricStatus
}

// the status of a single metric
type MetricStatus struct {
	// the type of metric source
	Type MetricSourceType

	// a metric describing a single kubernetes object (for example, hits-per-second on an Ingress object)
	// +optional
	Object *ObjectMetricStatus
	// a metric describing each pod in the current scale target (for example, transactions-processed-per-second).
	// The values will be averaged together before being compared to the target value
	// +optional
	Pods *PodsMetricStatus
	// a resource metric known to Kubernetes, as specified in requests and limits, describing each pod
	// in the current scale target (e.g. CPU or memory).  Such metrics are built in to Kubernetes,
	// and have special scaling options on top of those available to normal per-pod metrics using the "pods" source.
	// +optional
	Resource *ResourceMetricStatus
}

// a metric describing a single kubernetes object (for example, hits-per-second on an Ingress object)
type ObjectMetricStatus struct {
	// the described Kubernetes object
	Target CrossVersionObjectReference

	// the name of the metric in question
	MetricName string
	// the current value of the metric (as a quantity)
	CurrentValue resource.Quantity
}

// a metric describing each pod in the current scale target (for example, transactions-processed-per-second).
// The values will be averaged together before being compared to the target value
type PodsMetricStatus struct {
	// the name of the metric in question
	MetricName string
	// the current value of the metric (as a quantity)
	CurrentAverageValue resource.Quantity
}

// a resource metric known to Kubernetes, as specified in requests and limits, describing each pod
// in the current scale target (e.g. CPU or memory).  The values will be averaged together before
// being compared to the target.  Such metrics are built in to Kubernetes, and have special
// scaling options on top of those available to normal per-pod metrics using the "pods" source.
// Only one "target" type should be set.  Note that the current raw value is always displayed
// (even when the current values as request utilization is also displayed).
type ResourceMetricStatus struct {
	// the name of the resource in question
	Name api.ResourceName
	// the target value of the resource metric, represented as
	// a percentage of the requested value of the resource on the pods
	// (only populated if the corresponding request target was set)
	// +optional
	CurrentAverageUtilization *int32
	// the current value of the resource metric as a raw value
	CurrentAverageValue resource.Quantity
}

// +genclient=true

// configuration of a horizontal pod autoscaler.
type HorizontalPodAutoscaler struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// behaviour of autoscaler. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status.
	// +optional
	Spec HorizontalPodAutoscalerSpec

	// current information about the autoscaler.
	// +optional
	Status HorizontalPodAutoscalerStatus
}

// list of horizontal pod autoscaler objects.
type HorizontalPodAutoscalerList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// list of horizontal pod autoscaler objects.
	Items []HorizontalPodAutoscaler
}
