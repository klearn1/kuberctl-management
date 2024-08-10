/*
Copyright 2015 The Kubernetes Authors.

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

package pod

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	podtest "k8s.io/kubernetes/pkg/kubelet/pod/testing"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

// Stub out mirror client for testing purpose.
func newTestManager() (*basicManager, *podtest.FakeMirrorClient) {
	fakeMirrorClient := podtest.NewFakeMirrorClient()
	manager := NewBasicPodManager().(*basicManager)
	return manager, fakeMirrorClient
}

// Tests that pods/maps are properly set after the pod update, and the basic
// methods work correctly.
func TestGetSetPods(t *testing.T) {
	var (
		mirrorPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("mirror-pod-uid"),
				Name:      "mirror-static-pod-name",
				Namespace: metav1.NamespaceDefault,
				Annotations: map[string]string{
					kubetypes.ConfigSourceAnnotationKey: "api",
					kubetypes.ConfigMirrorAnnotationKey: "mirror",
				},
			},
		}
		staticPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:         types.UID("static-pod-uid"),
				Name:        "mirror-static-pod-name",
				Namespace:   metav1.NamespaceDefault,
				Annotations: map[string]string{kubetypes.ConfigSourceAnnotationKey: "file"},
			},
		}
		normalPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("normal-pod-uid"),
				Name:      "normal-pod-name",
				Namespace: metav1.NamespaceDefault,
			},
		}
	)

	testCase := []struct {
		name         string
		podList      []*v1.Pod
		wantPod      *v1.Pod
		expectPods   []*v1.Pod
		expectGetPod *v1.Pod
		expectUID    types.UID
	}{
		{
			name:         "Get normal pod",
			podList:      []*v1.Pod{mirrorPod, staticPod, normalPod},
			wantPod:      normalPod,
			expectPods:   []*v1.Pod{staticPod, normalPod},
			expectGetPod: normalPod,
			expectUID:    types.UID("static-pod-uid"),
		},
		{
			name:         "Get static pod",
			podList:      []*v1.Pod{mirrorPod, staticPod, normalPod},
			wantPod:      staticPod,
			expectPods:   []*v1.Pod{staticPod, normalPod},
			expectGetPod: staticPod,
			expectUID:    types.UID("static-pod-uid"),
		},
	}
	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			podManager, _ := newTestManager()
			podManager.SetPods(test.podList)
			actualPods := podManager.GetPods()

			// Tests that all regular pods, except for mirror pods, are recorded correctly.

			sortOpt := cmpopts.SortSlices(func(a, b *v1.Pod) bool { return a.Name < b.Name })
			if diff := cmp.Diff(actualPods, test.expectPods, sortOpt); diff != "" {
				t.Errorf("actualPods and expectPods are not equal: %s", diff)
			}

			// Tests UID translation works as expected. Convert static pod UID for comparison only.
			if uid := podManager.TranslatePodUID(mirrorPod.UID); uid != kubetypes.ResolvedPodUID(test.expectUID) {
				t.Errorf("unable to translate UID %q to the static POD's UID %q; %#v",
					mirrorPod.UID, staticPod.UID, podManager.mirrorPodByUID)
			}
			fullName := fmt.Sprintf("%s_%s", test.wantPod.Name, test.wantPod.Namespace)

			// Test the basic Get methods.
			actualPod, ok := podManager.GetPodByFullName(fullName)
			if !ok {
				if diff := cmp.Diff(actualPod, test.expectGetPod); diff != "" {
					t.Errorf("unexpected to get pod by full name: %s", diff)
				}
			}

			actualPod, ok = podManager.GetPodByName(test.wantPod.Namespace, test.wantPod.Name)
			if !ok {
				if diff := cmp.Diff(actualPod, test.expectGetPod); diff != "" {
					t.Errorf("unexpected to get pod by name: %s", diff)
				}
			}
		})
	}
}

func TestRemovePods(t *testing.T) {
	var (
		mirrorPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("mirror-pod-uid"),
				Name:      "mirror-static-pod-name",
				Namespace: metav1.NamespaceDefault,
				Annotations: map[string]string{
					kubetypes.ConfigSourceAnnotationKey: "api",
					kubetypes.ConfigMirrorAnnotationKey: "mirror",
				},
			},
		}
		staticPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:         types.UID("static-pod-uid"),
				Name:        "mirror-static-pod-name",
				Namespace:   metav1.NamespaceDefault,
				Annotations: map[string]string{kubetypes.ConfigSourceAnnotationKey: "file"},
			},
		}
		normalPod = &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				UID:       types.UID("normal-pod-uid"),
				Name:      "normal-pod-name",
				Namespace: metav1.NamespaceDefault,
			},
		}
	)

	testCase := []struct {
		name             string
		podList          []*v1.Pod
		needToRemovePod  *v1.Pod
		expectPods       []*v1.Pod
		expectMirrorPods []*v1.Pod
	}{
		{
			name:             "Remove mirror pod",
			podList:          []*v1.Pod{mirrorPod, staticPod, normalPod},
			needToRemovePod:  mirrorPod,
			expectPods:       []*v1.Pod{normalPod, staticPod},
			expectMirrorPods: []*v1.Pod{},
		},
		{
			name:             "Remove static pod",
			podList:          []*v1.Pod{mirrorPod, staticPod, normalPod},
			needToRemovePod:  staticPod,
			expectPods:       []*v1.Pod{normalPod},
			expectMirrorPods: []*v1.Pod{mirrorPod},
		},
		{
			name:             "Remove normal pod",
			podList:          []*v1.Pod{mirrorPod, staticPod, normalPod},
			needToRemovePod:  normalPod,
			expectPods:       []*v1.Pod{staticPod},
			expectMirrorPods: []*v1.Pod{mirrorPod},
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			podManager, _ := newTestManager()
			podManager.SetPods(test.podList)
			podManager.RemovePod(test.needToRemovePod)
			actualPods1 := podManager.GetPods()
			actualPods2, actualMirrorPods, _ := podManager.GetPodsAndMirrorPods()

			sortOpt := cmpopts.SortSlices(func(a, b *v1.Pod) bool { return a.Name < b.Name })

			// Check if the actual pods and mirror pods match the expected ones.
			if diff := cmp.Diff(actualPods1, actualPods2, sortOpt); diff != "" {
				t.Errorf("actualPods1 and actualPods2 are not equal: %s", diff)
			}
			if diff := cmp.Diff(actualPods1, test.expectPods, sortOpt); diff != "" {
				t.Errorf("actualPods1 and expectPods are not equal: %s", diff)

			}
			if diff := cmp.Diff(actualMirrorPods, test.expectMirrorPods, sortOpt); diff != "" {
				t.Errorf("actualMirrorPods and expectMirrorPods are not equal: %s", diff)
			}
		})
	}
}
