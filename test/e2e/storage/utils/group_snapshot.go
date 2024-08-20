/*
Copyright 2024 The Kubernetes Authors.

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

package utils

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/storage/names"
)

const (
	// VolumeGroupSnapshot is the group snapshot api
	VolumeGroupSnapshotAPIGroup = "groupsnapshot.storage.k8s.io"
	// VolumeGroupSnapshotAPIVersion is the group snapshot api version
	VolumeGroupSnapshotAPIVersion = "groupsnapshot.storage.k8s.io/v1alpha1"
)

var (

	// VolumeGroupSnapshotGVR is GroupVersionResource for groupvolumesnapshots
	VolumeGroupSnapshotGVR = schema.GroupVersionResource{Group: VolumeGroupSnapshotAPIGroup, Version: "v1alpha1", Resource: "volumegroupsnapshots"}
	// VolumeGroupSnapshotClassGVR is GroupVersionResource for groupvolumesnapshotclasses
	VolumeGroupSnapshotClassGVR = schema.GroupVersionResource{Group: VolumeGroupSnapshotAPIGroup, Version: "v1alpha1", Resource: "volumegroupsnapshotclasses"}
)

func GenerateVolumeGroupSnapshotClassSpec(
	snapshotter string,
	parameters map[string]string,
	ns string,
) *unstructured.Unstructured {
	volumeGroupSnapshotClass := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "VolumeGroupSnapshotClass",
			"apiVersion": VolumeGroupSnapshotAPIVersion,
			"metadata": map[string]interface{}{
				// Name must be unique, so let's base it on namespace name and use GenerateName
				"name": names.SimpleNameGenerator.GenerateName(ns),
			},
			"driver":         snapshotter,
			"parameters":     parameters,
			"deletionPolicy": "Delete",
		},
	}

	return volumeGroupSnapshotClass
}
