/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package metrics

import (
	"k8s.io/kubernetes/pkg/api"
)

func init() {
	// Register the API.
	addKnownTypes()
}

// Adds the list of known types to api.Scheme.
func addKnownTypes() {
	api.Scheme.AddKnownTypes("",
		&MetricsMeta{},
		&RawNodeMetrics{},
		&RawNodeMetricsList{},
		&RawPodMetrics{},
		&RawPodMetricsList{},
		&RawContainerMetrics{},
		&NonLocalObjectReference{},
		&Sample{},
		&AggregateSample{},
		&PodSample{},
		&ContainerSample{},
		&NetworkMetrics{},
		&CPUMetrics{},
		&MemoryMetrics{},
		&CustomMetric{},
		&CustomMetricSample{},
		&RawMetricsOptions{},
	)
}

func (*MetricsMeta) IsAnAPIObject()             {}
func (*RawNodeMetrics) IsAnAPIObject()          {}
func (*RawNodeMetricsList) IsAnAPIObject()      {}
func (*RawPodMetrics) IsAnAPIObject()           {}
func (*RawPodMetricsList) IsAnAPIObject()       {}
func (*RawContainerMetrics) IsAnAPIObject()     {}
func (*NonLocalObjectReference) IsAnAPIObject() {}
func (*Sample) IsAnAPIObject()                  {}
func (*AggregateSample) IsAnAPIObject()         {}
func (*PodSample) IsAnAPIObject()               {}
func (*ContainerSample) IsAnAPIObject()         {}
func (*NetworkMetrics) IsAnAPIObject()          {}
func (*CPUMetrics) IsAnAPIObject()              {}
func (*MemoryMetrics) IsAnAPIObject()           {}
func (*CustomMetric) IsAnAPIObject()            {}
func (*CustomMetricSample) IsAnAPIObject()      {}
func (*RawMetricsOptions) IsAnAPIObject()       {}
