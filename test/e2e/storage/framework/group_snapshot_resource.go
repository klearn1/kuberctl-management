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

package framework

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo/v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/storage/utils"
)

func getGroupSnapshot(groupName string, ns, snapshotClassName string) *unstructured.Unstructured {
	snapshot := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "VolumeGroupSnapshot",
			"apiVersion": utils.VolumeGroupSnapshotAPIVersion,
			"metadata": map[string]interface{}{
				"generateName": "group-snapshot-",
				"namespace":    ns,
			},
			"spec": map[string]interface{}{
				"volumeGroupSnapshotClassName": snapshotClassName,
				"source": map[string]interface{}{
					"selector": map[string]interface{}{
						"matchLabels": map[string]interface{}{
							"group": groupName,
						},
					},
				},
			},
		},
	}

	return snapshot
}

// SnapshotResource represents a snapshot class, a snapshot and its bound snapshot contents for a specific test case
type GroupSnapshotResource struct {
	Config  *PerTestConfig
	Pattern TestPattern

	Vgs        *unstructured.Unstructured
	Vgscontent *unstructured.Unstructured
	Vgsclass   *unstructured.Unstructured
}

// CreateSnapshot creates a VolumeSnapshotClass with given SnapshotDeletionPolicy and a VolumeSnapshot
// from the VolumeSnapshotClass using a dynamic client.
// Returns the unstructured VolumeSnapshotClass and VolumeSnapshot objects.
func CreateGroupSnapshot(ctx context.Context, sDriver VoulmeGroupSnapshottableTestDriver, config *PerTestConfig, pattern TestPattern, groupName string, pvcNamespace string, timeouts *framework.TimeoutContext, parameters map[string]string) (*unstructured.Unstructured, *unstructured.Unstructured) {
	defer ginkgo.GinkgoRecover()
	var err error
	if pattern.SnapshotType != GroupVolumeSnapshot {
		err = fmt.Errorf("SnapshotType must be set to GroupVolumeSnapshot")
		framework.ExpectNoError(err)
	}
	dc := config.Framework.DynamicClient

	ginkgo.By("creating a GroupSnapshotClass")
	gsclass := sDriver.GetVolumeGroupSnapshotClass(ctx, config, parameters)
	if gsclass == nil {
		framework.Failf("Failed to get group snapshot class based on test config")
	}
	gsclass.Object["deletionPolicy"] = pattern.SnapshotDeletionPolicy.String()

	gsclass, err = dc.Resource(utils.VolumeGroupSnapshotClassGVR).Create(ctx, gsclass, metav1.CreateOptions{})
	framework.ExpectNoError(err)
	gsclass, err = dc.Resource(utils.VolumeGroupSnapshotClassGVR).Get(ctx, gsclass.GetName(), metav1.GetOptions{})
	framework.ExpectNoError(err)

	ginkgo.By("creating a dynamic GroupVolumeSnapshot")
	// Prepare a dynamically provisioned group volume snapshot with certain data
	groupSnapshot := getGroupSnapshot(groupName, pvcNamespace, gsclass.GetName())

	groupSnapshot, err = dc.Resource(utils.VolumeGroupSnapshotGVR).Namespace(groupSnapshot.GetNamespace()).Create(ctx, groupSnapshot, metav1.CreateOptions{})
	framework.ExpectNoError(err)
	ginkgo.By("Waiting for snapshot to be ready")
	err = utils.WaitForGroupSnapshotReady(ctx, dc, groupSnapshot.GetNamespace(), groupSnapshot.GetName(), framework.Poll, timeouts.SnapshotCreate)
	framework.ExpectNoError(err)
	ginkgo.By("Getting group snapshot and content")
	groupSnapshot, err = dc.Resource(utils.VolumeGroupSnapshotGVR).Namespace(groupSnapshot.GetNamespace()).Get(ctx, groupSnapshot.GetName(), metav1.GetOptions{})
	framework.ExpectNoError(err)

	return gsclass, groupSnapshot
}

func (r *GroupSnapshotResource) CleanupResource(ctx context.Context, timeouts *framework.TimeoutContext) error {
	defer ginkgo.GinkgoRecover()
	dc := r.Config.Framework.DynamicClient
	err := dc.Resource(utils.VolumeGroupSnapshotClassGVR).Delete(ctx, r.Vgsclass.GetName(), metav1.DeleteOptions{})
	framework.ExpectNoError(err)
	return nil
}

func CreateGroupSnapshotResource(ctx context.Context, sDriver VoulmeGroupSnapshottableTestDriver, config *PerTestConfig, pattern TestPattern, pvcName string, pvcNamespace string, timeouts *framework.TimeoutContext, parameters map[string]string) *GroupSnapshotResource {
	vgsclass, snapshot := CreateGroupSnapshot(ctx, sDriver, config, pattern, pvcName, pvcNamespace, timeouts, parameters)
	vgs := &GroupSnapshotResource{
		Config:     config,
		Pattern:    pattern,
		Vgs:        snapshot,
		Vgsclass:   vgsclass,
		Vgscontent: nil,
	}
	return vgs
}
