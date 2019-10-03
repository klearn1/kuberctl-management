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
	v1 "k8s.io/api/storage/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	corev1 "k8s.io/internal-api/apis/core/v1"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&v1.StorageClass{}, func(obj interface{}) { SetObjectDefaults_StorageClass(obj.(*v1.StorageClass)) })
	scheme.AddTypeDefaultingFunc(&v1.StorageClassList{}, func(obj interface{}) { SetObjectDefaults_StorageClassList(obj.(*v1.StorageClassList)) })
	scheme.AddTypeDefaultingFunc(&v1.VolumeAttachment{}, func(obj interface{}) { SetObjectDefaults_VolumeAttachment(obj.(*v1.VolumeAttachment)) })
	scheme.AddTypeDefaultingFunc(&v1.VolumeAttachmentList{}, func(obj interface{}) { SetObjectDefaults_VolumeAttachmentList(obj.(*v1.VolumeAttachmentList)) })
	return nil
}

func SetObjectDefaults_StorageClass(in *v1.StorageClass) {
	SetDefaults_StorageClass(in)
}

func SetObjectDefaults_StorageClassList(in *v1.StorageClassList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_StorageClass(a)
	}
}

func SetObjectDefaults_VolumeAttachment(in *v1.VolumeAttachment) {
	if in.Spec.Source.InlineVolumeSpec != nil {
		corev1.SetDefaults_ResourceList(&in.Spec.Source.InlineVolumeSpec.Capacity)
		if in.Spec.Source.InlineVolumeSpec.PersistentVolumeSource.HostPath != nil {
			corev1.SetDefaults_HostPathVolumeSource(in.Spec.Source.InlineVolumeSpec.PersistentVolumeSource.HostPath)
		}
		if in.Spec.Source.InlineVolumeSpec.PersistentVolumeSource.RBD != nil {
			corev1.SetDefaults_RBDPersistentVolumeSource(in.Spec.Source.InlineVolumeSpec.PersistentVolumeSource.RBD)
		}
		if in.Spec.Source.InlineVolumeSpec.PersistentVolumeSource.ISCSI != nil {
			corev1.SetDefaults_ISCSIPersistentVolumeSource(in.Spec.Source.InlineVolumeSpec.PersistentVolumeSource.ISCSI)
		}
		if in.Spec.Source.InlineVolumeSpec.PersistentVolumeSource.AzureDisk != nil {
			corev1.SetDefaults_AzureDiskVolumeSource(in.Spec.Source.InlineVolumeSpec.PersistentVolumeSource.AzureDisk)
		}
		if in.Spec.Source.InlineVolumeSpec.PersistentVolumeSource.ScaleIO != nil {
			corev1.SetDefaults_ScaleIOPersistentVolumeSource(in.Spec.Source.InlineVolumeSpec.PersistentVolumeSource.ScaleIO)
		}
	}
}

func SetObjectDefaults_VolumeAttachmentList(in *v1.VolumeAttachmentList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetObjectDefaults_VolumeAttachment(a)
	}
}
