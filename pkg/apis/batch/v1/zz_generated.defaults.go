//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by defaulter-gen. DO NOT EDIT.

package v1

import (
	"encoding/json"

	v1 "k8s.io/api/batch/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	corev1 "k8s.io/kubernetes/pkg/apis/core/v1"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&v1.CronJob{}, func(obj interface{}) { SetObjectDefaults_CronJob(obj.(*v1.CronJob)) })
	scheme.AddTypeDefaultingFunc(&v1.CronJobList{}, func(obj interface{}) { SetObjectDefaults_CronJobList(obj.(*v1.CronJobList)) })
	scheme.AddTypeDefaultingFunc(&v1.Job{}, func(obj interface{}) { SetObjectDefaults_Job(obj.(*v1.Job)) })
	scheme.AddTypeDefaultingFunc(&v1.JobList{}, func(obj interface{}) { SetObjectDefaults_JobList(obj.(*v1.JobList)) })
	return nil
}

func SetObjectDefaults_CronJob(in *v1.CronJob) {
	SetDefaults_CronJob(in)
	corev1.SetDefaults_PodSpec(&in.Spec.JobTemplate.Spec.Template.Spec)
	for i := range in.Spec.JobTemplate.Spec.Template.Spec.Volumes {
		a := &in.Spec.JobTemplate.Spec.Template.Spec.Volumes[i]
		corev1.SetDefaults_Volume(a)
		if a.VolumeSource.HostPath != nil {
			corev1.SetDefaults_HostPathVolumeSource(a.VolumeSource.HostPath)
		}
		if a.VolumeSource.Secret != nil {
			if a.VolumeSource.Secret.DefaultMode == nil {
				var ptrVar1 int32 = SecretVolumeSourceDefaultMode
				a.VolumeSource.Secret.DefaultMode = &ptrVar1
			}
		}
		if a.VolumeSource.ISCSI != nil {
			corev1.SetDefaults_ISCSIVolumeSource(a.VolumeSource.ISCSI)
		}
		if a.VolumeSource.RBD != nil {
			corev1.SetDefaults_RBDVolumeSource(a.VolumeSource.RBD)
		}
		if a.VolumeSource.DownwardAPI != nil {
			for j := range a.VolumeSource.DownwardAPI.Items {
				b := &a.VolumeSource.DownwardAPI.Items[j]
				if b.FieldRef != nil {
					corev1.SetDefaults_ObjectFieldSelector(b.FieldRef)
				}
			}
			if a.VolumeSource.DownwardAPI.DefaultMode == nil {
				var ptrVar1 int32 = DownwardAPIVolumeSourceDefaultMode
				a.VolumeSource.DownwardAPI.DefaultMode = &ptrVar1
			}
		}
		if a.VolumeSource.ConfigMap != nil {
			if a.VolumeSource.ConfigMap.DefaultMode == nil {
				var ptrVar1 int32 = ConfigMapVolumeSourceDefaultMode
				a.VolumeSource.ConfigMap.DefaultMode = &ptrVar1
			}
		}
		if a.VolumeSource.AzureDisk != nil {
			corev1.SetDefaults_AzureDiskVolumeSource(a.VolumeSource.AzureDisk)
		}
		if a.VolumeSource.Projected != nil {
			for j := range a.VolumeSource.Projected.Sources {
				b := &a.VolumeSource.Projected.Sources[j]
				if b.DownwardAPI != nil {
					for k := range b.DownwardAPI.Items {
						c := &b.DownwardAPI.Items[k]
						if c.FieldRef != nil {
							corev1.SetDefaults_ObjectFieldSelector(c.FieldRef)
						}
					}
				}
				if b.ServiceAccountToken != nil {
					if b.ServiceAccountToken.ExpirationSeconds == nil {
						var ptrVar1 int64 = 3600
						b.ServiceAccountToken.ExpirationSeconds = &ptrVar1
					}
				}
			}
			if a.VolumeSource.Projected.DefaultMode == nil {
				var ptrVar1 int32 = ProjectedVolumeSourceDefaultMode
				a.VolumeSource.Projected.DefaultMode = &ptrVar1
			}
		}
		if a.VolumeSource.ScaleIO != nil {
			corev1.SetDefaults_ScaleIOVolumeSource(a.VolumeSource.ScaleIO)
		}
		if a.VolumeSource.Ephemeral != nil {
			if a.VolumeSource.Ephemeral.VolumeClaimTemplate != nil {
				corev1.SetDefaults_PersistentVolumeClaimSpec(&a.VolumeSource.Ephemeral.VolumeClaimTemplate.Spec)
				corev1.SetDefaults_ResourceList(&a.VolumeSource.Ephemeral.VolumeClaimTemplate.Spec.Resources.Limits)
				corev1.SetDefaults_ResourceList(&a.VolumeSource.Ephemeral.VolumeClaimTemplate.Spec.Resources.Requests)
			}
		}
	}
	for i := range in.Spec.JobTemplate.Spec.Template.Spec.InitContainers {
		a := &in.Spec.JobTemplate.Spec.Template.Spec.InitContainers[i]
		corev1.SetDefaults_Container(a)
		for j := range a.Ports {
			b := &a.Ports[j]
			if b.Protocol == "" {
				b.Protocol = "TCP"
			}
		}
		for j := range a.Env {
			b := &a.Env[j]
			if b.ValueFrom != nil {
				if b.ValueFrom.FieldRef != nil {
					corev1.SetDefaults_ObjectFieldSelector(b.ValueFrom.FieldRef)
				}
			}
		}
		corev1.SetDefaults_ResourceList(&a.Resources.Limits)
		corev1.SetDefaults_ResourceList(&a.Resources.Requests)
		if a.LivenessProbe != nil {
			if a.LivenessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.LivenessProbe.ProbeHandler.HTTPGet)
			}
			if a.LivenessProbe.ProbeHandler.GRPC != nil {
				if a.LivenessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.LivenessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.LivenessProbe.TimeoutSeconds == 0 {
				a.LivenessProbe.TimeoutSeconds = 0
			}
			if a.LivenessProbe.PeriodSeconds == 0 {
				a.LivenessProbe.PeriodSeconds = 10
			}
			if a.LivenessProbe.SuccessThreshold == 0 {
				a.LivenessProbe.SuccessThreshold = 1
			}
			if a.LivenessProbe.FailureThreshold == 0 {
				a.LivenessProbe.FailureThreshold = 3
			}
		}
		if a.ReadinessProbe != nil {
			if a.ReadinessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.ReadinessProbe.ProbeHandler.HTTPGet)
			}
			if a.ReadinessProbe.ProbeHandler.GRPC != nil {
				if a.ReadinessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.ReadinessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.ReadinessProbe.TimeoutSeconds == 0 {
				a.ReadinessProbe.TimeoutSeconds = 0
			}
			if a.ReadinessProbe.PeriodSeconds == 0 {
				a.ReadinessProbe.PeriodSeconds = 10
			}
			if a.ReadinessProbe.SuccessThreshold == 0 {
				a.ReadinessProbe.SuccessThreshold = 1
			}
			if a.ReadinessProbe.FailureThreshold == 0 {
				a.ReadinessProbe.FailureThreshold = 3
			}
		}
		if a.StartupProbe != nil {
			if a.StartupProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.StartupProbe.ProbeHandler.HTTPGet)
			}
			if a.StartupProbe.ProbeHandler.GRPC != nil {
				if a.StartupProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.StartupProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.StartupProbe.TimeoutSeconds == 0 {
				a.StartupProbe.TimeoutSeconds = 0
			}
			if a.StartupProbe.PeriodSeconds == 0 {
				a.StartupProbe.PeriodSeconds = 10
			}
			if a.StartupProbe.SuccessThreshold == 0 {
				a.StartupProbe.SuccessThreshold = 1
			}
			if a.StartupProbe.FailureThreshold == 0 {
				a.StartupProbe.FailureThreshold = 3
			}
		}
		if a.Lifecycle != nil {
			if a.Lifecycle.PostStart != nil {
				if a.Lifecycle.PostStart.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.Lifecycle.PostStart.HTTPGet)
				}
			}
			if a.Lifecycle.PreStop != nil {
				if a.Lifecycle.PreStop.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.Lifecycle.PreStop.HTTPGet)
				}
			}
		}
		if a.TerminationMessagePath == "" {
			a.TerminationMessagePath = TerminationMessagePathDefault
		}
		if a.TerminationMessagePolicy == "" {
			a.TerminationMessagePolicy = TerminationMessageReadFile
		}
	}
	for i := range in.Spec.JobTemplate.Spec.Template.Spec.Containers {
		a := &in.Spec.JobTemplate.Spec.Template.Spec.Containers[i]
		corev1.SetDefaults_Container(a)
		for j := range a.Ports {
			b := &a.Ports[j]
			if b.Protocol == "" {
				b.Protocol = "TCP"
			}
		}
		for j := range a.Env {
			b := &a.Env[j]
			if b.ValueFrom != nil {
				if b.ValueFrom.FieldRef != nil {
					corev1.SetDefaults_ObjectFieldSelector(b.ValueFrom.FieldRef)
				}
			}
		}
		corev1.SetDefaults_ResourceList(&a.Resources.Limits)
		corev1.SetDefaults_ResourceList(&a.Resources.Requests)
		if a.LivenessProbe != nil {
			if a.LivenessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.LivenessProbe.ProbeHandler.HTTPGet)
			}
			if a.LivenessProbe.ProbeHandler.GRPC != nil {
				if a.LivenessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.LivenessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.LivenessProbe.TimeoutSeconds == 0 {
				a.LivenessProbe.TimeoutSeconds = 0
			}
			if a.LivenessProbe.PeriodSeconds == 0 {
				a.LivenessProbe.PeriodSeconds = 10
			}
			if a.LivenessProbe.SuccessThreshold == 0 {
				a.LivenessProbe.SuccessThreshold = 1
			}
			if a.LivenessProbe.FailureThreshold == 0 {
				a.LivenessProbe.FailureThreshold = 3
			}
		}
		if a.ReadinessProbe != nil {
			if a.ReadinessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.ReadinessProbe.ProbeHandler.HTTPGet)
			}
			if a.ReadinessProbe.ProbeHandler.GRPC != nil {
				if a.ReadinessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.ReadinessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.ReadinessProbe.TimeoutSeconds == 0 {
				a.ReadinessProbe.TimeoutSeconds = 0
			}
			if a.ReadinessProbe.PeriodSeconds == 0 {
				a.ReadinessProbe.PeriodSeconds = 10
			}
			if a.ReadinessProbe.SuccessThreshold == 0 {
				a.ReadinessProbe.SuccessThreshold = 1
			}
			if a.ReadinessProbe.FailureThreshold == 0 {
				a.ReadinessProbe.FailureThreshold = 3
			}
		}
		if a.StartupProbe != nil {
			if a.StartupProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.StartupProbe.ProbeHandler.HTTPGet)
			}
			if a.StartupProbe.ProbeHandler.GRPC != nil {
				if a.StartupProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.StartupProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.StartupProbe.TimeoutSeconds == 0 {
				a.StartupProbe.TimeoutSeconds = 0
			}
			if a.StartupProbe.PeriodSeconds == 0 {
				a.StartupProbe.PeriodSeconds = 10
			}
			if a.StartupProbe.SuccessThreshold == 0 {
				a.StartupProbe.SuccessThreshold = 1
			}
			if a.StartupProbe.FailureThreshold == 0 {
				a.StartupProbe.FailureThreshold = 3
			}
		}
		if a.Lifecycle != nil {
			if a.Lifecycle.PostStart != nil {
				if a.Lifecycle.PostStart.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.Lifecycle.PostStart.HTTPGet)
				}
			}
			if a.Lifecycle.PreStop != nil {
				if a.Lifecycle.PreStop.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.Lifecycle.PreStop.HTTPGet)
				}
			}
		}
		if a.TerminationMessagePath == "" {
			a.TerminationMessagePath = TerminationMessagePathDefault
		}
		if a.TerminationMessagePolicy == "" {
			a.TerminationMessagePolicy = TerminationMessageReadFile
		}
	}
	for i := range in.Spec.JobTemplate.Spec.Template.Spec.EphemeralContainers {
		a := &in.Spec.JobTemplate.Spec.Template.Spec.EphemeralContainers[i]
		corev1.SetDefaults_EphemeralContainer(a)
		for j := range a.EphemeralContainerCommon.Ports {
			b := &a.EphemeralContainerCommon.Ports[j]
			if b.Protocol == "" {
				b.Protocol = "TCP"
			}
		}
		for j := range a.EphemeralContainerCommon.Env {
			b := &a.EphemeralContainerCommon.Env[j]
			if b.ValueFrom != nil {
				if b.ValueFrom.FieldRef != nil {
					corev1.SetDefaults_ObjectFieldSelector(b.ValueFrom.FieldRef)
				}
			}
		}
		corev1.SetDefaults_ResourceList(&a.EphemeralContainerCommon.Resources.Limits)
		corev1.SetDefaults_ResourceList(&a.EphemeralContainerCommon.Resources.Requests)
		if a.EphemeralContainerCommon.LivenessProbe != nil {
			if a.EphemeralContainerCommon.LivenessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.EphemeralContainerCommon.LivenessProbe.ProbeHandler.HTTPGet)
			}
			if a.EphemeralContainerCommon.LivenessProbe.ProbeHandler.GRPC != nil {
				if a.EphemeralContainerCommon.LivenessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.EphemeralContainerCommon.LivenessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.EphemeralContainerCommon.LivenessProbe.TimeoutSeconds == 0 {
				a.EphemeralContainerCommon.LivenessProbe.TimeoutSeconds = 0
			}
			if a.EphemeralContainerCommon.LivenessProbe.PeriodSeconds == 0 {
				a.EphemeralContainerCommon.LivenessProbe.PeriodSeconds = 10
			}
			if a.EphemeralContainerCommon.LivenessProbe.SuccessThreshold == 0 {
				a.EphemeralContainerCommon.LivenessProbe.SuccessThreshold = 1
			}
			if a.EphemeralContainerCommon.LivenessProbe.FailureThreshold == 0 {
				a.EphemeralContainerCommon.LivenessProbe.FailureThreshold = 3
			}
		}
		if a.EphemeralContainerCommon.ReadinessProbe != nil {
			if a.EphemeralContainerCommon.ReadinessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.EphemeralContainerCommon.ReadinessProbe.ProbeHandler.HTTPGet)
			}
			if a.EphemeralContainerCommon.ReadinessProbe.ProbeHandler.GRPC != nil {
				if a.EphemeralContainerCommon.ReadinessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.EphemeralContainerCommon.ReadinessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.EphemeralContainerCommon.ReadinessProbe.TimeoutSeconds == 0 {
				a.EphemeralContainerCommon.ReadinessProbe.TimeoutSeconds = 0
			}
			if a.EphemeralContainerCommon.ReadinessProbe.PeriodSeconds == 0 {
				a.EphemeralContainerCommon.ReadinessProbe.PeriodSeconds = 10
			}
			if a.EphemeralContainerCommon.ReadinessProbe.SuccessThreshold == 0 {
				a.EphemeralContainerCommon.ReadinessProbe.SuccessThreshold = 1
			}
			if a.EphemeralContainerCommon.ReadinessProbe.FailureThreshold == 0 {
				a.EphemeralContainerCommon.ReadinessProbe.FailureThreshold = 3
			}
		}
		if a.EphemeralContainerCommon.StartupProbe != nil {
			if a.EphemeralContainerCommon.StartupProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.EphemeralContainerCommon.StartupProbe.ProbeHandler.HTTPGet)
			}
			if a.EphemeralContainerCommon.StartupProbe.ProbeHandler.GRPC != nil {
				if a.EphemeralContainerCommon.StartupProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.EphemeralContainerCommon.StartupProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.EphemeralContainerCommon.StartupProbe.TimeoutSeconds == 0 {
				a.EphemeralContainerCommon.StartupProbe.TimeoutSeconds = 0
			}
			if a.EphemeralContainerCommon.StartupProbe.PeriodSeconds == 0 {
				a.EphemeralContainerCommon.StartupProbe.PeriodSeconds = 10
			}
			if a.EphemeralContainerCommon.StartupProbe.SuccessThreshold == 0 {
				a.EphemeralContainerCommon.StartupProbe.SuccessThreshold = 1
			}
			if a.EphemeralContainerCommon.StartupProbe.FailureThreshold == 0 {
				a.EphemeralContainerCommon.StartupProbe.FailureThreshold = 3
			}
		}
		if a.EphemeralContainerCommon.Lifecycle != nil {
			if a.EphemeralContainerCommon.Lifecycle.PostStart != nil {
				if a.EphemeralContainerCommon.Lifecycle.PostStart.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.EphemeralContainerCommon.Lifecycle.PostStart.HTTPGet)
				}
			}
			if a.EphemeralContainerCommon.Lifecycle.PreStop != nil {
				if a.EphemeralContainerCommon.Lifecycle.PreStop.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.EphemeralContainerCommon.Lifecycle.PreStop.HTTPGet)
				}
			}
		}
	}
	if in.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy == "" {
		in.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy = RestartPolicyAlways
	}
	if in.Spec.JobTemplate.Spec.Template.Spec.TerminationGracePeriodSeconds == nil {
		var ptrVar1 int64 = DefaultTerminationGracePeriodSeconds
		in.Spec.JobTemplate.Spec.Template.Spec.TerminationGracePeriodSeconds = &ptrVar1
	}
	if in.Spec.JobTemplate.Spec.Template.Spec.DNSPolicy == "" {
		in.Spec.JobTemplate.Spec.Template.Spec.DNSPolicy = DNSClusterFirst
	}
	if in.Spec.JobTemplate.Spec.Template.Spec.SecurityContext == nil {
		if err := json.Unmarshal([]byte(`{}`), &in.Spec.JobTemplate.Spec.Template.Spec.SecurityContext); err != nil {
			panic(err)
		}
	}
	if in.Spec.JobTemplate.Spec.Template.Spec.SchedulerName == "" {
		in.Spec.JobTemplate.Spec.Template.Spec.SchedulerName = DefaultSchedulerName
	}
	if in.Spec.JobTemplate.Spec.Template.Spec.EnableServiceLinks == nil {
		var ptrVar1 bool = DefaultEnableServiceLinks
		in.Spec.JobTemplate.Spec.Template.Spec.EnableServiceLinks = &ptrVar1
	}
	corev1.SetDefaults_ResourceList(&in.Spec.JobTemplate.Spec.Template.Spec.Overhead)
}

func SetObjectDefaults_CronJobList(in *v1.CronJobList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_CronJob(a)
	}
}

func SetObjectDefaults_Job(in *v1.Job) {
	SetDefaults_Job(in)
	corev1.SetDefaults_PodSpec(&in.Spec.Template.Spec)
	for i := range in.Spec.Template.Spec.Volumes {
		a := &in.Spec.Template.Spec.Volumes[i]
		corev1.SetDefaults_Volume(a)
		if a.VolumeSource.HostPath != nil {
			corev1.SetDefaults_HostPathVolumeSource(a.VolumeSource.HostPath)
		}
		if a.VolumeSource.Secret != nil {
			if a.VolumeSource.Secret.DefaultMode == nil {
				var ptrVar1 int32 = SecretVolumeSourceDefaultMode
				a.VolumeSource.Secret.DefaultMode = &ptrVar1
			}
		}
		if a.VolumeSource.ISCSI != nil {
			corev1.SetDefaults_ISCSIVolumeSource(a.VolumeSource.ISCSI)
		}
		if a.VolumeSource.RBD != nil {
			corev1.SetDefaults_RBDVolumeSource(a.VolumeSource.RBD)
		}
		if a.VolumeSource.DownwardAPI != nil {
			for j := range a.VolumeSource.DownwardAPI.Items {
				b := &a.VolumeSource.DownwardAPI.Items[j]
				if b.FieldRef != nil {
					corev1.SetDefaults_ObjectFieldSelector(b.FieldRef)
				}
			}
			if a.VolumeSource.DownwardAPI.DefaultMode == nil {
				var ptrVar1 int32 = DownwardAPIVolumeSourceDefaultMode
				a.VolumeSource.DownwardAPI.DefaultMode = &ptrVar1
			}
		}
		if a.VolumeSource.ConfigMap != nil {
			if a.VolumeSource.ConfigMap.DefaultMode == nil {
				var ptrVar1 int32 = ConfigMapVolumeSourceDefaultMode
				a.VolumeSource.ConfigMap.DefaultMode = &ptrVar1
			}
		}
		if a.VolumeSource.AzureDisk != nil {
			corev1.SetDefaults_AzureDiskVolumeSource(a.VolumeSource.AzureDisk)
		}
		if a.VolumeSource.Projected != nil {
			for j := range a.VolumeSource.Projected.Sources {
				b := &a.VolumeSource.Projected.Sources[j]
				if b.DownwardAPI != nil {
					for k := range b.DownwardAPI.Items {
						c := &b.DownwardAPI.Items[k]
						if c.FieldRef != nil {
							corev1.SetDefaults_ObjectFieldSelector(c.FieldRef)
						}
					}
				}
				if b.ServiceAccountToken != nil {
					if b.ServiceAccountToken.ExpirationSeconds == nil {
						var ptrVar1 int64 = 3600
						b.ServiceAccountToken.ExpirationSeconds = &ptrVar1
					}
				}
			}
			if a.VolumeSource.Projected.DefaultMode == nil {
				var ptrVar1 int32 = ProjectedVolumeSourceDefaultMode
				a.VolumeSource.Projected.DefaultMode = &ptrVar1
			}
		}
		if a.VolumeSource.ScaleIO != nil {
			corev1.SetDefaults_ScaleIOVolumeSource(a.VolumeSource.ScaleIO)
		}
		if a.VolumeSource.Ephemeral != nil {
			if a.VolumeSource.Ephemeral.VolumeClaimTemplate != nil {
				corev1.SetDefaults_PersistentVolumeClaimSpec(&a.VolumeSource.Ephemeral.VolumeClaimTemplate.Spec)
				corev1.SetDefaults_ResourceList(&a.VolumeSource.Ephemeral.VolumeClaimTemplate.Spec.Resources.Limits)
				corev1.SetDefaults_ResourceList(&a.VolumeSource.Ephemeral.VolumeClaimTemplate.Spec.Resources.Requests)
			}
		}
	}
	for i := range in.Spec.Template.Spec.InitContainers {
		a := &in.Spec.Template.Spec.InitContainers[i]
		corev1.SetDefaults_Container(a)
		for j := range a.Ports {
			b := &a.Ports[j]
			if b.Protocol == "" {
				b.Protocol = "TCP"
			}
		}
		for j := range a.Env {
			b := &a.Env[j]
			if b.ValueFrom != nil {
				if b.ValueFrom.FieldRef != nil {
					corev1.SetDefaults_ObjectFieldSelector(b.ValueFrom.FieldRef)
				}
			}
		}
		corev1.SetDefaults_ResourceList(&a.Resources.Limits)
		corev1.SetDefaults_ResourceList(&a.Resources.Requests)
		if a.LivenessProbe != nil {
			if a.LivenessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.LivenessProbe.ProbeHandler.HTTPGet)
			}
			if a.LivenessProbe.ProbeHandler.GRPC != nil {
				if a.LivenessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.LivenessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.LivenessProbe.TimeoutSeconds == 0 {
				a.LivenessProbe.TimeoutSeconds = 0
			}
			if a.LivenessProbe.PeriodSeconds == 0 {
				a.LivenessProbe.PeriodSeconds = 10
			}
			if a.LivenessProbe.SuccessThreshold == 0 {
				a.LivenessProbe.SuccessThreshold = 1
			}
			if a.LivenessProbe.FailureThreshold == 0 {
				a.LivenessProbe.FailureThreshold = 3
			}
		}
		if a.ReadinessProbe != nil {
			if a.ReadinessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.ReadinessProbe.ProbeHandler.HTTPGet)
			}
			if a.ReadinessProbe.ProbeHandler.GRPC != nil {
				if a.ReadinessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.ReadinessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.ReadinessProbe.TimeoutSeconds == 0 {
				a.ReadinessProbe.TimeoutSeconds = 0
			}
			if a.ReadinessProbe.PeriodSeconds == 0 {
				a.ReadinessProbe.PeriodSeconds = 10
			}
			if a.ReadinessProbe.SuccessThreshold == 0 {
				a.ReadinessProbe.SuccessThreshold = 1
			}
			if a.ReadinessProbe.FailureThreshold == 0 {
				a.ReadinessProbe.FailureThreshold = 3
			}
		}
		if a.StartupProbe != nil {
			if a.StartupProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.StartupProbe.ProbeHandler.HTTPGet)
			}
			if a.StartupProbe.ProbeHandler.GRPC != nil {
				if a.StartupProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.StartupProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.StartupProbe.TimeoutSeconds == 0 {
				a.StartupProbe.TimeoutSeconds = 0
			}
			if a.StartupProbe.PeriodSeconds == 0 {
				a.StartupProbe.PeriodSeconds = 10
			}
			if a.StartupProbe.SuccessThreshold == 0 {
				a.StartupProbe.SuccessThreshold = 1
			}
			if a.StartupProbe.FailureThreshold == 0 {
				a.StartupProbe.FailureThreshold = 3
			}
		}
		if a.Lifecycle != nil {
			if a.Lifecycle.PostStart != nil {
				if a.Lifecycle.PostStart.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.Lifecycle.PostStart.HTTPGet)
				}
			}
			if a.Lifecycle.PreStop != nil {
				if a.Lifecycle.PreStop.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.Lifecycle.PreStop.HTTPGet)
				}
			}
		}
		if a.TerminationMessagePath == "" {
			a.TerminationMessagePath = TerminationMessagePathDefault
		}
		if a.TerminationMessagePolicy == "" {
			a.TerminationMessagePolicy = TerminationMessageReadFile
		}
	}
	for i := range in.Spec.Template.Spec.Containers {
		a := &in.Spec.Template.Spec.Containers[i]
		corev1.SetDefaults_Container(a)
		for j := range a.Ports {
			b := &a.Ports[j]
			if b.Protocol == "" {
				b.Protocol = "TCP"
			}
		}
		for j := range a.Env {
			b := &a.Env[j]
			if b.ValueFrom != nil {
				if b.ValueFrom.FieldRef != nil {
					corev1.SetDefaults_ObjectFieldSelector(b.ValueFrom.FieldRef)
				}
			}
		}
		corev1.SetDefaults_ResourceList(&a.Resources.Limits)
		corev1.SetDefaults_ResourceList(&a.Resources.Requests)
		if a.LivenessProbe != nil {
			if a.LivenessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.LivenessProbe.ProbeHandler.HTTPGet)
			}
			if a.LivenessProbe.ProbeHandler.GRPC != nil {
				if a.LivenessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.LivenessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.LivenessProbe.TimeoutSeconds == 0 {
				a.LivenessProbe.TimeoutSeconds = 0
			}
			if a.LivenessProbe.PeriodSeconds == 0 {
				a.LivenessProbe.PeriodSeconds = 10
			}
			if a.LivenessProbe.SuccessThreshold == 0 {
				a.LivenessProbe.SuccessThreshold = 1
			}
			if a.LivenessProbe.FailureThreshold == 0 {
				a.LivenessProbe.FailureThreshold = 3
			}
		}
		if a.ReadinessProbe != nil {
			if a.ReadinessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.ReadinessProbe.ProbeHandler.HTTPGet)
			}
			if a.ReadinessProbe.ProbeHandler.GRPC != nil {
				if a.ReadinessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.ReadinessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.ReadinessProbe.TimeoutSeconds == 0 {
				a.ReadinessProbe.TimeoutSeconds = 0
			}
			if a.ReadinessProbe.PeriodSeconds == 0 {
				a.ReadinessProbe.PeriodSeconds = 10
			}
			if a.ReadinessProbe.SuccessThreshold == 0 {
				a.ReadinessProbe.SuccessThreshold = 1
			}
			if a.ReadinessProbe.FailureThreshold == 0 {
				a.ReadinessProbe.FailureThreshold = 3
			}
		}
		if a.StartupProbe != nil {
			if a.StartupProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.StartupProbe.ProbeHandler.HTTPGet)
			}
			if a.StartupProbe.ProbeHandler.GRPC != nil {
				if a.StartupProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.StartupProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.StartupProbe.TimeoutSeconds == 0 {
				a.StartupProbe.TimeoutSeconds = 0
			}
			if a.StartupProbe.PeriodSeconds == 0 {
				a.StartupProbe.PeriodSeconds = 10
			}
			if a.StartupProbe.SuccessThreshold == 0 {
				a.StartupProbe.SuccessThreshold = 1
			}
			if a.StartupProbe.FailureThreshold == 0 {
				a.StartupProbe.FailureThreshold = 3
			}
		}
		if a.Lifecycle != nil {
			if a.Lifecycle.PostStart != nil {
				if a.Lifecycle.PostStart.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.Lifecycle.PostStart.HTTPGet)
				}
			}
			if a.Lifecycle.PreStop != nil {
				if a.Lifecycle.PreStop.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.Lifecycle.PreStop.HTTPGet)
				}
			}
		}
		if a.TerminationMessagePath == "" {
			a.TerminationMessagePath = TerminationMessagePathDefault
		}
		if a.TerminationMessagePolicy == "" {
			a.TerminationMessagePolicy = TerminationMessageReadFile
		}
	}
	for i := range in.Spec.Template.Spec.EphemeralContainers {
		a := &in.Spec.Template.Spec.EphemeralContainers[i]
		corev1.SetDefaults_EphemeralContainer(a)
		for j := range a.EphemeralContainerCommon.Ports {
			b := &a.EphemeralContainerCommon.Ports[j]
			if b.Protocol == "" {
				b.Protocol = "TCP"
			}
		}
		for j := range a.EphemeralContainerCommon.Env {
			b := &a.EphemeralContainerCommon.Env[j]
			if b.ValueFrom != nil {
				if b.ValueFrom.FieldRef != nil {
					corev1.SetDefaults_ObjectFieldSelector(b.ValueFrom.FieldRef)
				}
			}
		}
		corev1.SetDefaults_ResourceList(&a.EphemeralContainerCommon.Resources.Limits)
		corev1.SetDefaults_ResourceList(&a.EphemeralContainerCommon.Resources.Requests)
		if a.EphemeralContainerCommon.LivenessProbe != nil {
			if a.EphemeralContainerCommon.LivenessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.EphemeralContainerCommon.LivenessProbe.ProbeHandler.HTTPGet)
			}
			if a.EphemeralContainerCommon.LivenessProbe.ProbeHandler.GRPC != nil {
				if a.EphemeralContainerCommon.LivenessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.EphemeralContainerCommon.LivenessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.EphemeralContainerCommon.LivenessProbe.TimeoutSeconds == 0 {
				a.EphemeralContainerCommon.LivenessProbe.TimeoutSeconds = 0
			}
			if a.EphemeralContainerCommon.LivenessProbe.PeriodSeconds == 0 {
				a.EphemeralContainerCommon.LivenessProbe.PeriodSeconds = 10
			}
			if a.EphemeralContainerCommon.LivenessProbe.SuccessThreshold == 0 {
				a.EphemeralContainerCommon.LivenessProbe.SuccessThreshold = 1
			}
			if a.EphemeralContainerCommon.LivenessProbe.FailureThreshold == 0 {
				a.EphemeralContainerCommon.LivenessProbe.FailureThreshold = 3
			}
		}
		if a.EphemeralContainerCommon.ReadinessProbe != nil {
			if a.EphemeralContainerCommon.ReadinessProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.EphemeralContainerCommon.ReadinessProbe.ProbeHandler.HTTPGet)
			}
			if a.EphemeralContainerCommon.ReadinessProbe.ProbeHandler.GRPC != nil {
				if a.EphemeralContainerCommon.ReadinessProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.EphemeralContainerCommon.ReadinessProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.EphemeralContainerCommon.ReadinessProbe.TimeoutSeconds == 0 {
				a.EphemeralContainerCommon.ReadinessProbe.TimeoutSeconds = 0
			}
			if a.EphemeralContainerCommon.ReadinessProbe.PeriodSeconds == 0 {
				a.EphemeralContainerCommon.ReadinessProbe.PeriodSeconds = 10
			}
			if a.EphemeralContainerCommon.ReadinessProbe.SuccessThreshold == 0 {
				a.EphemeralContainerCommon.ReadinessProbe.SuccessThreshold = 1
			}
			if a.EphemeralContainerCommon.ReadinessProbe.FailureThreshold == 0 {
				a.EphemeralContainerCommon.ReadinessProbe.FailureThreshold = 3
			}
		}
		if a.EphemeralContainerCommon.StartupProbe != nil {
			if a.EphemeralContainerCommon.StartupProbe.ProbeHandler.HTTPGet != nil {
				corev1.SetDefaults_HTTPGetAction(a.EphemeralContainerCommon.StartupProbe.ProbeHandler.HTTPGet)
			}
			if a.EphemeralContainerCommon.StartupProbe.ProbeHandler.GRPC != nil {
				if a.EphemeralContainerCommon.StartupProbe.ProbeHandler.GRPC.Service == nil {
					var ptrVar1 string = ""
					a.EphemeralContainerCommon.StartupProbe.ProbeHandler.GRPC.Service = &ptrVar1
				}
			}
			if a.EphemeralContainerCommon.StartupProbe.TimeoutSeconds == 0 {
				a.EphemeralContainerCommon.StartupProbe.TimeoutSeconds = 0
			}
			if a.EphemeralContainerCommon.StartupProbe.PeriodSeconds == 0 {
				a.EphemeralContainerCommon.StartupProbe.PeriodSeconds = 10
			}
			if a.EphemeralContainerCommon.StartupProbe.SuccessThreshold == 0 {
				a.EphemeralContainerCommon.StartupProbe.SuccessThreshold = 1
			}
			if a.EphemeralContainerCommon.StartupProbe.FailureThreshold == 0 {
				a.EphemeralContainerCommon.StartupProbe.FailureThreshold = 3
			}
		}
		if a.EphemeralContainerCommon.Lifecycle != nil {
			if a.EphemeralContainerCommon.Lifecycle.PostStart != nil {
				if a.EphemeralContainerCommon.Lifecycle.PostStart.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.EphemeralContainerCommon.Lifecycle.PostStart.HTTPGet)
				}
			}
			if a.EphemeralContainerCommon.Lifecycle.PreStop != nil {
				if a.EphemeralContainerCommon.Lifecycle.PreStop.HTTPGet != nil {
					corev1.SetDefaults_HTTPGetAction(a.EphemeralContainerCommon.Lifecycle.PreStop.HTTPGet)
				}
			}
		}
	}
	if in.Spec.Template.Spec.RestartPolicy == "" {
		in.Spec.Template.Spec.RestartPolicy = RestartPolicyAlways
	}
	if in.Spec.Template.Spec.TerminationGracePeriodSeconds == nil {
		var ptrVar1 int64 = DefaultTerminationGracePeriodSeconds
		in.Spec.Template.Spec.TerminationGracePeriodSeconds = &ptrVar1
	}
	if in.Spec.Template.Spec.DNSPolicy == "" {
		in.Spec.Template.Spec.DNSPolicy = DNSClusterFirst
	}
	if in.Spec.Template.Spec.SecurityContext == nil {
		if err := json.Unmarshal([]byte(`{}`), &in.Spec.Template.Spec.SecurityContext); err != nil {
			panic(err)
		}
	}
	if in.Spec.Template.Spec.SchedulerName == "" {
		in.Spec.Template.Spec.SchedulerName = DefaultSchedulerName
	}
	if in.Spec.Template.Spec.EnableServiceLinks == nil {
		var ptrVar1 bool = DefaultEnableServiceLinks
		in.Spec.Template.Spec.EnableServiceLinks = &ptrVar1
	}
	corev1.SetDefaults_ResourceList(&in.Spec.Template.Spec.Overhead)
}

func SetObjectDefaults_JobList(in *v1.JobList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_Job(a)
	}
}
