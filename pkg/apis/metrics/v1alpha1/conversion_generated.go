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

// DO NOT EDIT. THIS FILE IS AUTO-GENERATED BY $KUBEROOT/hack/update-generated-conversions.sh

package v1alpha1

import (
	reflect "reflect"

	api "k8s.io/kubernetes/pkg/api"
	metrics "k8s.io/kubernetes/pkg/apis/metrics"
	conversion "k8s.io/kubernetes/pkg/conversion"
)

func autoconvert_metrics_AggregateSample_To_v1alpha1_AggregateSample(in *metrics.AggregateSample, out *AggregateSample, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.AggregateSample))(in)
	}
	if err := convert_metrics_Sample_To_v1alpha1_Sample(&in.Sample, &out.Sample, s); err != nil {
		return err
	}
	if in.CPU != nil {
		out.CPU = new(CPUMetrics)
		if err := convert_metrics_CPUMetrics_To_v1alpha1_CPUMetrics(in.CPU, out.CPU, s); err != nil {
			return err
		}
	} else {
		out.CPU = nil
	}
	if in.Memory != nil {
		out.Memory = new(MemoryMetrics)
		if err := convert_metrics_MemoryMetrics_To_v1alpha1_MemoryMetrics(in.Memory, out.Memory, s); err != nil {
			return err
		}
	} else {
		out.Memory = nil
	}
	if in.Network != nil {
		out.Network = new(NetworkMetrics)
		if err := convert_metrics_NetworkMetrics_To_v1alpha1_NetworkMetrics(in.Network, out.Network, s); err != nil {
			return err
		}
	} else {
		out.Network = nil
	}
	if in.Filesystem != nil {
		out.Filesystem = make([]FilesystemMetrics, len(in.Filesystem))
		for i := range in.Filesystem {
			if err := convert_metrics_FilesystemMetrics_To_v1alpha1_FilesystemMetrics(&in.Filesystem[i], &out.Filesystem[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Filesystem = nil
	}
	return nil
}

func convert_metrics_AggregateSample_To_v1alpha1_AggregateSample(in *metrics.AggregateSample, out *AggregateSample, s conversion.Scope) error {
	return autoconvert_metrics_AggregateSample_To_v1alpha1_AggregateSample(in, out, s)
}

func autoconvert_metrics_CPUMetrics_To_v1alpha1_CPUMetrics(in *metrics.CPUMetrics, out *CPUMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.CPUMetrics))(in)
	}
	if in.TotalCores != nil {
		if err := s.Convert(&in.TotalCores, &out.TotalCores, 0); err != nil {
			return err
		}
	} else {
		out.TotalCores = nil
	}
	return nil
}

func convert_metrics_CPUMetrics_To_v1alpha1_CPUMetrics(in *metrics.CPUMetrics, out *CPUMetrics, s conversion.Scope) error {
	return autoconvert_metrics_CPUMetrics_To_v1alpha1_CPUMetrics(in, out, s)
}

func autoconvert_metrics_ContainerSample_To_v1alpha1_ContainerSample(in *metrics.ContainerSample, out *ContainerSample, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.ContainerSample))(in)
	}
	if err := convert_metrics_Sample_To_v1alpha1_Sample(&in.Sample, &out.Sample, s); err != nil {
		return err
	}
	if in.CPU != nil {
		out.CPU = new(CPUMetrics)
		if err := convert_metrics_CPUMetrics_To_v1alpha1_CPUMetrics(in.CPU, out.CPU, s); err != nil {
			return err
		}
	} else {
		out.CPU = nil
	}
	if in.Memory != nil {
		out.Memory = new(MemoryMetrics)
		if err := convert_metrics_MemoryMetrics_To_v1alpha1_MemoryMetrics(in.Memory, out.Memory, s); err != nil {
			return err
		}
	} else {
		out.Memory = nil
	}
	if in.Filesystem != nil {
		out.Filesystem = make([]FilesystemMetrics, len(in.Filesystem))
		for i := range in.Filesystem {
			if err := convert_metrics_FilesystemMetrics_To_v1alpha1_FilesystemMetrics(&in.Filesystem[i], &out.Filesystem[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Filesystem = nil
	}
	return nil
}

func convert_metrics_ContainerSample_To_v1alpha1_ContainerSample(in *metrics.ContainerSample, out *ContainerSample, s conversion.Scope) error {
	return autoconvert_metrics_ContainerSample_To_v1alpha1_ContainerSample(in, out, s)
}

func autoconvert_metrics_FilesystemMetrics_To_v1alpha1_FilesystemMetrics(in *metrics.FilesystemMetrics, out *FilesystemMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.FilesystemMetrics))(in)
	}
	out.Device = in.Device
	if in.UsageBytes != nil {
		if err := s.Convert(&in.UsageBytes, &out.UsageBytes, 0); err != nil {
			return err
		}
	} else {
		out.UsageBytes = nil
	}
	if in.LimitBytes != nil {
		if err := s.Convert(&in.LimitBytes, &out.LimitBytes, 0); err != nil {
			return err
		}
	} else {
		out.LimitBytes = nil
	}
	return nil
}

func convert_metrics_FilesystemMetrics_To_v1alpha1_FilesystemMetrics(in *metrics.FilesystemMetrics, out *FilesystemMetrics, s conversion.Scope) error {
	return autoconvert_metrics_FilesystemMetrics_To_v1alpha1_FilesystemMetrics(in, out, s)
}

func autoconvert_metrics_MemoryMetrics_To_v1alpha1_MemoryMetrics(in *metrics.MemoryMetrics, out *MemoryMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.MemoryMetrics))(in)
	}
	if in.TotalBytes != nil {
		if err := s.Convert(&in.TotalBytes, &out.TotalBytes, 0); err != nil {
			return err
		}
	} else {
		out.TotalBytes = nil
	}
	if in.UsageBytes != nil {
		if err := s.Convert(&in.UsageBytes, &out.UsageBytes, 0); err != nil {
			return err
		}
	} else {
		out.UsageBytes = nil
	}
	if in.PageFaults != nil {
		out.PageFaults = new(int64)
		*out.PageFaults = *in.PageFaults
	} else {
		out.PageFaults = nil
	}
	if in.MajorPageFaults != nil {
		out.MajorPageFaults = new(int64)
		*out.MajorPageFaults = *in.MajorPageFaults
	} else {
		out.MajorPageFaults = nil
	}
	return nil
}

func convert_metrics_MemoryMetrics_To_v1alpha1_MemoryMetrics(in *metrics.MemoryMetrics, out *MemoryMetrics, s conversion.Scope) error {
	return autoconvert_metrics_MemoryMetrics_To_v1alpha1_MemoryMetrics(in, out, s)
}

func autoconvert_metrics_MetricsMeta_To_v1alpha1_MetricsMeta(in *metrics.MetricsMeta, out *MetricsMeta, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.MetricsMeta))(in)
	}
	out.SelfLink = in.SelfLink
	return nil
}

func convert_metrics_MetricsMeta_To_v1alpha1_MetricsMeta(in *metrics.MetricsMeta, out *MetricsMeta, s conversion.Scope) error {
	return autoconvert_metrics_MetricsMeta_To_v1alpha1_MetricsMeta(in, out, s)
}

func autoconvert_metrics_NetworkMetrics_To_v1alpha1_NetworkMetrics(in *metrics.NetworkMetrics, out *NetworkMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.NetworkMetrics))(in)
	}
	if in.RxBytes != nil {
		if err := s.Convert(&in.RxBytes, &out.RxBytes, 0); err != nil {
			return err
		}
	} else {
		out.RxBytes = nil
	}
	if in.RxErrors != nil {
		out.RxErrors = new(int64)
		*out.RxErrors = *in.RxErrors
	} else {
		out.RxErrors = nil
	}
	if in.TxBytes != nil {
		if err := s.Convert(&in.TxBytes, &out.TxBytes, 0); err != nil {
			return err
		}
	} else {
		out.TxBytes = nil
	}
	if in.TxErrors != nil {
		out.TxErrors = new(int64)
		*out.TxErrors = *in.TxErrors
	} else {
		out.TxErrors = nil
	}
	return nil
}

func convert_metrics_NetworkMetrics_To_v1alpha1_NetworkMetrics(in *metrics.NetworkMetrics, out *NetworkMetrics, s conversion.Scope) error {
	return autoconvert_metrics_NetworkMetrics_To_v1alpha1_NetworkMetrics(in, out, s)
}

func autoconvert_metrics_NonLocalObjectReference_To_v1alpha1_NonLocalObjectReference(in *metrics.NonLocalObjectReference, out *NonLocalObjectReference, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.NonLocalObjectReference))(in)
	}
	out.Name = in.Name
	out.Namespace = in.Namespace
	out.UID = in.UID
	return nil
}

func convert_metrics_NonLocalObjectReference_To_v1alpha1_NonLocalObjectReference(in *metrics.NonLocalObjectReference, out *NonLocalObjectReference, s conversion.Scope) error {
	return autoconvert_metrics_NonLocalObjectReference_To_v1alpha1_NonLocalObjectReference(in, out, s)
}

func autoconvert_metrics_PodSample_To_v1alpha1_PodSample(in *metrics.PodSample, out *PodSample, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.PodSample))(in)
	}
	if err := convert_metrics_Sample_To_v1alpha1_Sample(&in.Sample, &out.Sample, s); err != nil {
		return err
	}
	if in.Network != nil {
		out.Network = new(NetworkMetrics)
		if err := convert_metrics_NetworkMetrics_To_v1alpha1_NetworkMetrics(in.Network, out.Network, s); err != nil {
			return err
		}
	} else {
		out.Network = nil
	}
	return nil
}

func convert_metrics_PodSample_To_v1alpha1_PodSample(in *metrics.PodSample, out *PodSample, s conversion.Scope) error {
	return autoconvert_metrics_PodSample_To_v1alpha1_PodSample(in, out, s)
}

func autoconvert_metrics_RawContainerMetrics_To_v1alpha1_RawContainerMetrics(in *metrics.RawContainerMetrics, out *RawContainerMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.RawContainerMetrics))(in)
	}
	out.Name = in.Name
	if in.Labels != nil {
		out.Labels = make(map[string]string)
		for key, val := range in.Labels {
			out.Labels[key] = val
		}
	} else {
		out.Labels = nil
	}
	if in.Samples != nil {
		out.Samples = make([]ContainerSample, len(in.Samples))
		for i := range in.Samples {
			if err := convert_metrics_ContainerSample_To_v1alpha1_ContainerSample(&in.Samples[i], &out.Samples[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Samples = nil
	}
	return nil
}

func convert_metrics_RawContainerMetrics_To_v1alpha1_RawContainerMetrics(in *metrics.RawContainerMetrics, out *RawContainerMetrics, s conversion.Scope) error {
	return autoconvert_metrics_RawContainerMetrics_To_v1alpha1_RawContainerMetrics(in, out, s)
}

func autoconvert_metrics_RawMetricsOptions_To_v1alpha1_RawMetricsOptions(in *metrics.RawMetricsOptions, out *RawMetricsOptions, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.RawMetricsOptions))(in)
	}
	out.MaxSamples = in.MaxSamples
	return nil
}

func convert_metrics_RawMetricsOptions_To_v1alpha1_RawMetricsOptions(in *metrics.RawMetricsOptions, out *RawMetricsOptions, s conversion.Scope) error {
	return autoconvert_metrics_RawMetricsOptions_To_v1alpha1_RawMetricsOptions(in, out, s)
}

func autoconvert_metrics_RawNodeMetrics_To_v1alpha1_RawNodeMetrics(in *metrics.RawNodeMetrics, out *RawNodeMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.RawNodeMetrics))(in)
	}
	if err := s.Convert(&in.TypeMeta, &out.TypeMeta, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.ListMeta, &out.ListMeta, 0); err != nil {
		return err
	}
	out.NodeName = in.NodeName
	if in.Total != nil {
		out.Total = make([]AggregateSample, len(in.Total))
		for i := range in.Total {
			if err := convert_metrics_AggregateSample_To_v1alpha1_AggregateSample(&in.Total[i], &out.Total[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Total = nil
	}
	if in.SystemContainers != nil {
		out.SystemContainers = make([]RawContainerMetrics, len(in.SystemContainers))
		for i := range in.SystemContainers {
			if err := convert_metrics_RawContainerMetrics_To_v1alpha1_RawContainerMetrics(&in.SystemContainers[i], &out.SystemContainers[i], s); err != nil {
				return err
			}
		}
	} else {
		out.SystemContainers = nil
	}
	return nil
}

func convert_metrics_RawNodeMetrics_To_v1alpha1_RawNodeMetrics(in *metrics.RawNodeMetrics, out *RawNodeMetrics, s conversion.Scope) error {
	return autoconvert_metrics_RawNodeMetrics_To_v1alpha1_RawNodeMetrics(in, out, s)
}

func autoconvert_metrics_RawNodeMetricsList_To_v1alpha1_RawNodeMetricsList(in *metrics.RawNodeMetricsList, out *RawNodeMetricsList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.RawNodeMetricsList))(in)
	}
	if err := s.Convert(&in.TypeMeta, &out.TypeMeta, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.ListMeta, &out.ListMeta, 0); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]RawNodeMetrics, len(in.Items))
		for i := range in.Items {
			if err := convert_metrics_RawNodeMetrics_To_v1alpha1_RawNodeMetrics(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_metrics_RawNodeMetricsList_To_v1alpha1_RawNodeMetricsList(in *metrics.RawNodeMetricsList, out *RawNodeMetricsList, s conversion.Scope) error {
	return autoconvert_metrics_RawNodeMetricsList_To_v1alpha1_RawNodeMetricsList(in, out, s)
}

func autoconvert_metrics_RawPodMetrics_To_v1alpha1_RawPodMetrics(in *metrics.RawPodMetrics, out *RawPodMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.RawPodMetrics))(in)
	}
	if err := s.Convert(&in.TypeMeta, &out.TypeMeta, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.ListMeta, &out.ListMeta, 0); err != nil {
		return err
	}
	if err := convert_metrics_NonLocalObjectReference_To_v1alpha1_NonLocalObjectReference(&in.PodRef, &out.PodRef, s); err != nil {
		return err
	}
	if in.Containers != nil {
		out.Containers = make([]RawContainerMetrics, len(in.Containers))
		for i := range in.Containers {
			if err := convert_metrics_RawContainerMetrics_To_v1alpha1_RawContainerMetrics(&in.Containers[i], &out.Containers[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Containers = nil
	}
	if in.Samples != nil {
		out.Samples = make([]PodSample, len(in.Samples))
		for i := range in.Samples {
			if err := convert_metrics_PodSample_To_v1alpha1_PodSample(&in.Samples[i], &out.Samples[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Samples = nil
	}
	return nil
}

func convert_metrics_RawPodMetrics_To_v1alpha1_RawPodMetrics(in *metrics.RawPodMetrics, out *RawPodMetrics, s conversion.Scope) error {
	return autoconvert_metrics_RawPodMetrics_To_v1alpha1_RawPodMetrics(in, out, s)
}

func autoconvert_metrics_RawPodMetricsList_To_v1alpha1_RawPodMetricsList(in *metrics.RawPodMetricsList, out *RawPodMetricsList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.RawPodMetricsList))(in)
	}
	if err := s.Convert(&in.TypeMeta, &out.TypeMeta, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.ListMeta, &out.ListMeta, 0); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]RawPodMetrics, len(in.Items))
		for i := range in.Items {
			if err := convert_metrics_RawPodMetrics_To_v1alpha1_RawPodMetrics(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_metrics_RawPodMetricsList_To_v1alpha1_RawPodMetricsList(in *metrics.RawPodMetricsList, out *RawPodMetricsList, s conversion.Scope) error {
	return autoconvert_metrics_RawPodMetricsList_To_v1alpha1_RawPodMetricsList(in, out, s)
}

func autoconvert_metrics_Sample_To_v1alpha1_Sample(in *metrics.Sample, out *Sample, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*metrics.Sample))(in)
	}
	if err := s.Convert(&in.SampleTime, &out.SampleTime, 0); err != nil {
		return err
	}
	return nil
}

func convert_metrics_Sample_To_v1alpha1_Sample(in *metrics.Sample, out *Sample, s conversion.Scope) error {
	return autoconvert_metrics_Sample_To_v1alpha1_Sample(in, out, s)
}

func autoconvert_v1alpha1_AggregateSample_To_metrics_AggregateSample(in *AggregateSample, out *metrics.AggregateSample, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*AggregateSample))(in)
	}
	if err := convert_v1alpha1_Sample_To_metrics_Sample(&in.Sample, &out.Sample, s); err != nil {
		return err
	}
	if in.CPU != nil {
		out.CPU = new(metrics.CPUMetrics)
		if err := convert_v1alpha1_CPUMetrics_To_metrics_CPUMetrics(in.CPU, out.CPU, s); err != nil {
			return err
		}
	} else {
		out.CPU = nil
	}
	if in.Memory != nil {
		out.Memory = new(metrics.MemoryMetrics)
		if err := convert_v1alpha1_MemoryMetrics_To_metrics_MemoryMetrics(in.Memory, out.Memory, s); err != nil {
			return err
		}
	} else {
		out.Memory = nil
	}
	if in.Network != nil {
		out.Network = new(metrics.NetworkMetrics)
		if err := convert_v1alpha1_NetworkMetrics_To_metrics_NetworkMetrics(in.Network, out.Network, s); err != nil {
			return err
		}
	} else {
		out.Network = nil
	}
	if in.Filesystem != nil {
		out.Filesystem = make([]metrics.FilesystemMetrics, len(in.Filesystem))
		for i := range in.Filesystem {
			if err := convert_v1alpha1_FilesystemMetrics_To_metrics_FilesystemMetrics(&in.Filesystem[i], &out.Filesystem[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Filesystem = nil
	}
	return nil
}

func convert_v1alpha1_AggregateSample_To_metrics_AggregateSample(in *AggregateSample, out *metrics.AggregateSample, s conversion.Scope) error {
	return autoconvert_v1alpha1_AggregateSample_To_metrics_AggregateSample(in, out, s)
}

func autoconvert_v1alpha1_CPUMetrics_To_metrics_CPUMetrics(in *CPUMetrics, out *metrics.CPUMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*CPUMetrics))(in)
	}
	if in.TotalCores != nil {
		if err := s.Convert(&in.TotalCores, &out.TotalCores, 0); err != nil {
			return err
		}
	} else {
		out.TotalCores = nil
	}
	return nil
}

func convert_v1alpha1_CPUMetrics_To_metrics_CPUMetrics(in *CPUMetrics, out *metrics.CPUMetrics, s conversion.Scope) error {
	return autoconvert_v1alpha1_CPUMetrics_To_metrics_CPUMetrics(in, out, s)
}

func autoconvert_v1alpha1_ContainerSample_To_metrics_ContainerSample(in *ContainerSample, out *metrics.ContainerSample, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*ContainerSample))(in)
	}
	if err := convert_v1alpha1_Sample_To_metrics_Sample(&in.Sample, &out.Sample, s); err != nil {
		return err
	}
	if in.CPU != nil {
		out.CPU = new(metrics.CPUMetrics)
		if err := convert_v1alpha1_CPUMetrics_To_metrics_CPUMetrics(in.CPU, out.CPU, s); err != nil {
			return err
		}
	} else {
		out.CPU = nil
	}
	if in.Memory != nil {
		out.Memory = new(metrics.MemoryMetrics)
		if err := convert_v1alpha1_MemoryMetrics_To_metrics_MemoryMetrics(in.Memory, out.Memory, s); err != nil {
			return err
		}
	} else {
		out.Memory = nil
	}
	if in.Filesystem != nil {
		out.Filesystem = make([]metrics.FilesystemMetrics, len(in.Filesystem))
		for i := range in.Filesystem {
			if err := convert_v1alpha1_FilesystemMetrics_To_metrics_FilesystemMetrics(&in.Filesystem[i], &out.Filesystem[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Filesystem = nil
	}
	return nil
}

func convert_v1alpha1_ContainerSample_To_metrics_ContainerSample(in *ContainerSample, out *metrics.ContainerSample, s conversion.Scope) error {
	return autoconvert_v1alpha1_ContainerSample_To_metrics_ContainerSample(in, out, s)
}

func autoconvert_v1alpha1_FilesystemMetrics_To_metrics_FilesystemMetrics(in *FilesystemMetrics, out *metrics.FilesystemMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*FilesystemMetrics))(in)
	}
	out.Device = in.Device
	if in.UsageBytes != nil {
		if err := s.Convert(&in.UsageBytes, &out.UsageBytes, 0); err != nil {
			return err
		}
	} else {
		out.UsageBytes = nil
	}
	if in.LimitBytes != nil {
		if err := s.Convert(&in.LimitBytes, &out.LimitBytes, 0); err != nil {
			return err
		}
	} else {
		out.LimitBytes = nil
	}
	return nil
}

func convert_v1alpha1_FilesystemMetrics_To_metrics_FilesystemMetrics(in *FilesystemMetrics, out *metrics.FilesystemMetrics, s conversion.Scope) error {
	return autoconvert_v1alpha1_FilesystemMetrics_To_metrics_FilesystemMetrics(in, out, s)
}

func autoconvert_v1alpha1_MemoryMetrics_To_metrics_MemoryMetrics(in *MemoryMetrics, out *metrics.MemoryMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*MemoryMetrics))(in)
	}
	if in.TotalBytes != nil {
		if err := s.Convert(&in.TotalBytes, &out.TotalBytes, 0); err != nil {
			return err
		}
	} else {
		out.TotalBytes = nil
	}
	if in.UsageBytes != nil {
		if err := s.Convert(&in.UsageBytes, &out.UsageBytes, 0); err != nil {
			return err
		}
	} else {
		out.UsageBytes = nil
	}
	if in.PageFaults != nil {
		out.PageFaults = new(int64)
		*out.PageFaults = *in.PageFaults
	} else {
		out.PageFaults = nil
	}
	if in.MajorPageFaults != nil {
		out.MajorPageFaults = new(int64)
		*out.MajorPageFaults = *in.MajorPageFaults
	} else {
		out.MajorPageFaults = nil
	}
	return nil
}

func convert_v1alpha1_MemoryMetrics_To_metrics_MemoryMetrics(in *MemoryMetrics, out *metrics.MemoryMetrics, s conversion.Scope) error {
	return autoconvert_v1alpha1_MemoryMetrics_To_metrics_MemoryMetrics(in, out, s)
}

func autoconvert_v1alpha1_MetricsMeta_To_metrics_MetricsMeta(in *MetricsMeta, out *metrics.MetricsMeta, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*MetricsMeta))(in)
	}
	out.SelfLink = in.SelfLink
	return nil
}

func convert_v1alpha1_MetricsMeta_To_metrics_MetricsMeta(in *MetricsMeta, out *metrics.MetricsMeta, s conversion.Scope) error {
	return autoconvert_v1alpha1_MetricsMeta_To_metrics_MetricsMeta(in, out, s)
}

func autoconvert_v1alpha1_NetworkMetrics_To_metrics_NetworkMetrics(in *NetworkMetrics, out *metrics.NetworkMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*NetworkMetrics))(in)
	}
	if in.RxBytes != nil {
		if err := s.Convert(&in.RxBytes, &out.RxBytes, 0); err != nil {
			return err
		}
	} else {
		out.RxBytes = nil
	}
	if in.RxErrors != nil {
		out.RxErrors = new(int64)
		*out.RxErrors = *in.RxErrors
	} else {
		out.RxErrors = nil
	}
	if in.TxBytes != nil {
		if err := s.Convert(&in.TxBytes, &out.TxBytes, 0); err != nil {
			return err
		}
	} else {
		out.TxBytes = nil
	}
	if in.TxErrors != nil {
		out.TxErrors = new(int64)
		*out.TxErrors = *in.TxErrors
	} else {
		out.TxErrors = nil
	}
	return nil
}

func convert_v1alpha1_NetworkMetrics_To_metrics_NetworkMetrics(in *NetworkMetrics, out *metrics.NetworkMetrics, s conversion.Scope) error {
	return autoconvert_v1alpha1_NetworkMetrics_To_metrics_NetworkMetrics(in, out, s)
}

func autoconvert_v1alpha1_NonLocalObjectReference_To_metrics_NonLocalObjectReference(in *NonLocalObjectReference, out *metrics.NonLocalObjectReference, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*NonLocalObjectReference))(in)
	}
	out.Name = in.Name
	out.Namespace = in.Namespace
	out.UID = in.UID
	return nil
}

func convert_v1alpha1_NonLocalObjectReference_To_metrics_NonLocalObjectReference(in *NonLocalObjectReference, out *metrics.NonLocalObjectReference, s conversion.Scope) error {
	return autoconvert_v1alpha1_NonLocalObjectReference_To_metrics_NonLocalObjectReference(in, out, s)
}

func autoconvert_v1alpha1_PodSample_To_metrics_PodSample(in *PodSample, out *metrics.PodSample, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*PodSample))(in)
	}
	if err := convert_v1alpha1_Sample_To_metrics_Sample(&in.Sample, &out.Sample, s); err != nil {
		return err
	}
	if in.Network != nil {
		out.Network = new(metrics.NetworkMetrics)
		if err := convert_v1alpha1_NetworkMetrics_To_metrics_NetworkMetrics(in.Network, out.Network, s); err != nil {
			return err
		}
	} else {
		out.Network = nil
	}
	return nil
}

func convert_v1alpha1_PodSample_To_metrics_PodSample(in *PodSample, out *metrics.PodSample, s conversion.Scope) error {
	return autoconvert_v1alpha1_PodSample_To_metrics_PodSample(in, out, s)
}

func autoconvert_v1alpha1_RawContainerMetrics_To_metrics_RawContainerMetrics(in *RawContainerMetrics, out *metrics.RawContainerMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*RawContainerMetrics))(in)
	}
	out.Name = in.Name
	if in.Labels != nil {
		out.Labels = make(map[string]string)
		for key, val := range in.Labels {
			out.Labels[key] = val
		}
	} else {
		out.Labels = nil
	}
	if in.Samples != nil {
		out.Samples = make([]metrics.ContainerSample, len(in.Samples))
		for i := range in.Samples {
			if err := convert_v1alpha1_ContainerSample_To_metrics_ContainerSample(&in.Samples[i], &out.Samples[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Samples = nil
	}
	return nil
}

func convert_v1alpha1_RawContainerMetrics_To_metrics_RawContainerMetrics(in *RawContainerMetrics, out *metrics.RawContainerMetrics, s conversion.Scope) error {
	return autoconvert_v1alpha1_RawContainerMetrics_To_metrics_RawContainerMetrics(in, out, s)
}

func autoconvert_v1alpha1_RawMetricsOptions_To_metrics_RawMetricsOptions(in *RawMetricsOptions, out *metrics.RawMetricsOptions, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*RawMetricsOptions))(in)
	}
	out.MaxSamples = in.MaxSamples
	return nil
}

func convert_v1alpha1_RawMetricsOptions_To_metrics_RawMetricsOptions(in *RawMetricsOptions, out *metrics.RawMetricsOptions, s conversion.Scope) error {
	return autoconvert_v1alpha1_RawMetricsOptions_To_metrics_RawMetricsOptions(in, out, s)
}

func autoconvert_v1alpha1_RawNodeMetrics_To_metrics_RawNodeMetrics(in *RawNodeMetrics, out *metrics.RawNodeMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*RawNodeMetrics))(in)
	}
	if err := s.Convert(&in.TypeMeta, &out.TypeMeta, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.ListMeta, &out.ListMeta, 0); err != nil {
		return err
	}
	out.NodeName = in.NodeName
	if in.Total != nil {
		out.Total = make([]metrics.AggregateSample, len(in.Total))
		for i := range in.Total {
			if err := convert_v1alpha1_AggregateSample_To_metrics_AggregateSample(&in.Total[i], &out.Total[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Total = nil
	}
	if in.SystemContainers != nil {
		out.SystemContainers = make([]metrics.RawContainerMetrics, len(in.SystemContainers))
		for i := range in.SystemContainers {
			if err := convert_v1alpha1_RawContainerMetrics_To_metrics_RawContainerMetrics(&in.SystemContainers[i], &out.SystemContainers[i], s); err != nil {
				return err
			}
		}
	} else {
		out.SystemContainers = nil
	}
	return nil
}

func convert_v1alpha1_RawNodeMetrics_To_metrics_RawNodeMetrics(in *RawNodeMetrics, out *metrics.RawNodeMetrics, s conversion.Scope) error {
	return autoconvert_v1alpha1_RawNodeMetrics_To_metrics_RawNodeMetrics(in, out, s)
}

func autoconvert_v1alpha1_RawNodeMetricsList_To_metrics_RawNodeMetricsList(in *RawNodeMetricsList, out *metrics.RawNodeMetricsList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*RawNodeMetricsList))(in)
	}
	if err := s.Convert(&in.TypeMeta, &out.TypeMeta, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.ListMeta, &out.ListMeta, 0); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]metrics.RawNodeMetrics, len(in.Items))
		for i := range in.Items {
			if err := convert_v1alpha1_RawNodeMetrics_To_metrics_RawNodeMetrics(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1alpha1_RawNodeMetricsList_To_metrics_RawNodeMetricsList(in *RawNodeMetricsList, out *metrics.RawNodeMetricsList, s conversion.Scope) error {
	return autoconvert_v1alpha1_RawNodeMetricsList_To_metrics_RawNodeMetricsList(in, out, s)
}

func autoconvert_v1alpha1_RawPodMetrics_To_metrics_RawPodMetrics(in *RawPodMetrics, out *metrics.RawPodMetrics, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*RawPodMetrics))(in)
	}
	if err := s.Convert(&in.TypeMeta, &out.TypeMeta, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.ListMeta, &out.ListMeta, 0); err != nil {
		return err
	}
	if err := convert_v1alpha1_NonLocalObjectReference_To_metrics_NonLocalObjectReference(&in.PodRef, &out.PodRef, s); err != nil {
		return err
	}
	if in.Containers != nil {
		out.Containers = make([]metrics.RawContainerMetrics, len(in.Containers))
		for i := range in.Containers {
			if err := convert_v1alpha1_RawContainerMetrics_To_metrics_RawContainerMetrics(&in.Containers[i], &out.Containers[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Containers = nil
	}
	if in.Samples != nil {
		out.Samples = make([]metrics.PodSample, len(in.Samples))
		for i := range in.Samples {
			if err := convert_v1alpha1_PodSample_To_metrics_PodSample(&in.Samples[i], &out.Samples[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Samples = nil
	}
	return nil
}

func convert_v1alpha1_RawPodMetrics_To_metrics_RawPodMetrics(in *RawPodMetrics, out *metrics.RawPodMetrics, s conversion.Scope) error {
	return autoconvert_v1alpha1_RawPodMetrics_To_metrics_RawPodMetrics(in, out, s)
}

func autoconvert_v1alpha1_RawPodMetricsList_To_metrics_RawPodMetricsList(in *RawPodMetricsList, out *metrics.RawPodMetricsList, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*RawPodMetricsList))(in)
	}
	if err := s.Convert(&in.TypeMeta, &out.TypeMeta, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.ListMeta, &out.ListMeta, 0); err != nil {
		return err
	}
	if in.Items != nil {
		out.Items = make([]metrics.RawPodMetrics, len(in.Items))
		for i := range in.Items {
			if err := convert_v1alpha1_RawPodMetrics_To_metrics_RawPodMetrics(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func convert_v1alpha1_RawPodMetricsList_To_metrics_RawPodMetricsList(in *RawPodMetricsList, out *metrics.RawPodMetricsList, s conversion.Scope) error {
	return autoconvert_v1alpha1_RawPodMetricsList_To_metrics_RawPodMetricsList(in, out, s)
}

func autoconvert_v1alpha1_Sample_To_metrics_Sample(in *Sample, out *metrics.Sample, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*Sample))(in)
	}
	if err := s.Convert(&in.SampleTime, &out.SampleTime, 0); err != nil {
		return err
	}
	return nil
}

func convert_v1alpha1_Sample_To_metrics_Sample(in *Sample, out *metrics.Sample, s conversion.Scope) error {
	return autoconvert_v1alpha1_Sample_To_metrics_Sample(in, out, s)
}

func init() {
	err := api.Scheme.AddGeneratedConversionFuncs(
		autoconvert_metrics_AggregateSample_To_v1alpha1_AggregateSample,
		autoconvert_metrics_CPUMetrics_To_v1alpha1_CPUMetrics,
		autoconvert_metrics_ContainerSample_To_v1alpha1_ContainerSample,
		autoconvert_metrics_FilesystemMetrics_To_v1alpha1_FilesystemMetrics,
		autoconvert_metrics_MemoryMetrics_To_v1alpha1_MemoryMetrics,
		autoconvert_metrics_MetricsMeta_To_v1alpha1_MetricsMeta,
		autoconvert_metrics_NetworkMetrics_To_v1alpha1_NetworkMetrics,
		autoconvert_metrics_NonLocalObjectReference_To_v1alpha1_NonLocalObjectReference,
		autoconvert_metrics_PodSample_To_v1alpha1_PodSample,
		autoconvert_metrics_RawContainerMetrics_To_v1alpha1_RawContainerMetrics,
		autoconvert_metrics_RawMetricsOptions_To_v1alpha1_RawMetricsOptions,
		autoconvert_metrics_RawNodeMetricsList_To_v1alpha1_RawNodeMetricsList,
		autoconvert_metrics_RawNodeMetrics_To_v1alpha1_RawNodeMetrics,
		autoconvert_metrics_RawPodMetricsList_To_v1alpha1_RawPodMetricsList,
		autoconvert_metrics_RawPodMetrics_To_v1alpha1_RawPodMetrics,
		autoconvert_metrics_Sample_To_v1alpha1_Sample,
		autoconvert_v1alpha1_AggregateSample_To_metrics_AggregateSample,
		autoconvert_v1alpha1_CPUMetrics_To_metrics_CPUMetrics,
		autoconvert_v1alpha1_ContainerSample_To_metrics_ContainerSample,
		autoconvert_v1alpha1_FilesystemMetrics_To_metrics_FilesystemMetrics,
		autoconvert_v1alpha1_MemoryMetrics_To_metrics_MemoryMetrics,
		autoconvert_v1alpha1_MetricsMeta_To_metrics_MetricsMeta,
		autoconvert_v1alpha1_NetworkMetrics_To_metrics_NetworkMetrics,
		autoconvert_v1alpha1_NonLocalObjectReference_To_metrics_NonLocalObjectReference,
		autoconvert_v1alpha1_PodSample_To_metrics_PodSample,
		autoconvert_v1alpha1_RawContainerMetrics_To_metrics_RawContainerMetrics,
		autoconvert_v1alpha1_RawMetricsOptions_To_metrics_RawMetricsOptions,
		autoconvert_v1alpha1_RawNodeMetricsList_To_metrics_RawNodeMetricsList,
		autoconvert_v1alpha1_RawNodeMetrics_To_metrics_RawNodeMetrics,
		autoconvert_v1alpha1_RawPodMetricsList_To_metrics_RawPodMetricsList,
		autoconvert_v1alpha1_RawPodMetrics_To_metrics_RawPodMetrics,
		autoconvert_v1alpha1_Sample_To_metrics_Sample,
	)
	if err != nil {
		// If one of the conversion functions is malformed, detect it immediately.
		panic(err)
	}
}
