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
	"sync"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/api/v1"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/secret"
)

// Manager stores and manages access to pods, maintaining the mappings
// between static pods and mirror pods.
//
// The kubelet discovers pod updates from 3 sources: file, http, and
// apiserver. Pods from non-apiserver sources are called static pods, and API
// server is not aware of the existence of static pods. In order to monitor
// the status of such pods, the kubelet creates a mirror pod for each static
// pod via the API server.
//
// A mirror pod has the same pod full name (name and namespace) as its static
// counterpart (albeit different metadata such as UID, etc). By leveraging the
// fact that the kubelet reports the pod status using the pod full name, the
// status of the mirror pod always reflects the actual status of the static
// pod. When a static pod gets deleted, the associated orphaned mirror pod
// will also be removed.
type Manager interface {
	// GetPods returns the regular pods bound to the kubelet and their spec.
	GetPods() []*Pod
	// GetPodByFullName returns the (non-mirror) pod that matches full name, as well as
	// whether the pod was found.
	GetPodByFullName(podFullName string) (*Pod, bool)
	// GetPodByName provides the (non-mirror) pod that matches namespace and
	// name, as well as whether the pod was found.
	GetPodByName(namespace, name string) (*Pod, bool)
	// GetPodByUID provides the (non-mirror) pod that matches pod UID, as well as
	// whether the pod is found.
	GetPodByUID(types.UID) (*Pod, bool)
	// GetPodByMirrorPod returns the static pod for the given mirror pod and
	// whether it was known to the pod manger.
	GetPodByMirrorPod(*Pod) (*Pod, bool)
	// GetMirrorPodByPod returns the mirror pod for the given static pod and
	// whether it was known to the pod manager.
	GetMirrorPodByPod(*Pod) (*Pod, bool)
	// GetPodsAndMirrorPods returns the both regular and mirror pods.
	GetPodsAndMirrorPods() ([]*Pod, []*Pod)
	// AddPod adds the given pod to the manager.
	AddPod(pod *v1.Pod) *Pod
	// UpdatePod updates the given pod in the manager.
	UpdatePod(pod *Pod)
	// DeletePod deletes the given pod from the manager.  For mirror pods,
	// this means deleting the mappings related to mirror pods.  For non-
	// mirror pods, this means deleting from indexes for all non-mirror pods.
	DeletePod(pod *Pod)
	// DeleteOrphanedMirrorPods deletes all mirror pods which do not have
	// associated static pods. This method sends deletion requests to the API
	// server, but does NOT modify the internal pod storage in BasicManager.
	DeleteOrphanedMirrorPods()
	// TranslatePodUID returns the actual UID of a pod. If the UID belongs to
	// a mirror pod, returns the UID of its static pod. Otherwise, returns the
	// original UID.
	//
	// All public-facing functions should perform this translation for UIDs
	// because user may provide a mirror pod UID, which is not recognized by
	// internal Kubelet functions.
	TranslatePodUID(uid types.UID) types.UID
	// GetUIDTranslations returns the mappings of static pod UIDs to mirror pod
	// UIDs and mirror pod UIDs to static pod UIDs.
	GetUIDTranslations() (podToMirror, mirrorToPod map[types.UID]types.UID)

	MirrorClient
}

// BasicManager is a functional Manger.
//
// All fields in BasicManager are read-only and are updated calling SetPods,
// AddPod, UpdatePod, or DeletePod.
type BasicManager struct {
	// Protects all internal maps.
	lock sync.RWMutex

	// Regular pods indexed by UID.
	podByUID map[types.UID]*Pod
	// Mirror pods indexed by UID.
	mirrorPodByUID map[types.UID]*Pod

	// Pods indexed by full name for easy access.
	podByFullName       map[string]*Pod
	mirrorPodByFullName map[string]*Pod

	// Mirror pod UID to pod UID map.
	translationByUID map[types.UID]types.UID

	// BasicManager is keeping secretManager up-to-date.
	secretManager secret.Manager

	// A mirror pod client to create/delete mirror pods.
	MirrorClient
}

// NewBasicPodManager returns a functional Manager.
func NewBasicPodManager(client MirrorClient, secretManager secret.Manager) *BasicManager {
	pm := &BasicManager{}
	pm.secretManager = secretManager
	pm.MirrorClient = client
	pm.SetPods(nil)
	return pm
}

// Set the internal pods based on the new pods.
func (pm *BasicManager) SetPods(newPods []*v1.Pod) []*Pod {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	pm.podByUID = make(map[types.UID]*Pod)
	pm.podByFullName = make(map[string]*Pod)
	pm.mirrorPodByUID = make(map[types.UID]*Pod)
	pm.mirrorPodByFullName = make(map[string]*Pod)
	pm.translationByUID = make(map[types.UID]types.UID)

	kubepods := FromAPIPods(newPods)
	pm.updatePodsInternal(kubepods...)
	return kubepods
}

func (pm *BasicManager) AddPod(pod *v1.Pod) *Pod {
	p := NewPod(pod)
	pm.UpdatePod(p)
	return p
}

func (pm *BasicManager) UpdatePod(pod *Pod) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.updatePodsInternal(pod)
}

// updatePodsInternal replaces the given pods in the current state of the
// manager, updating the various indices.  The caller is assumed to hold the
// lock.
func (pm *BasicManager) updatePodsInternal(pods ...*Pod) {
	for _, pod := range pods {
		if pm.secretManager != nil {
			// TODO: Consider detecting only status update and in such case do
			// not register pod, as it doesn't really matter.
			pm.secretManager.RegisterPod(pod.GetAPIPod())
		}
		podFullName := pod.GetFullName()
		if pod.IsMirror() {
			pm.mirrorPodByUID[pod.UID()] = pod
			pm.mirrorPodByFullName[podFullName] = pod
			if p, ok := pm.podByFullName[podFullName]; ok {
				pm.translationByUID[pod.UID()] = p.UID()
			}
		} else {
			pm.podByUID[pod.UID()] = pod
			pm.podByFullName[podFullName] = pod
			if mirror, ok := pm.mirrorPodByFullName[podFullName]; ok {
				pm.translationByUID[mirror.UID()] = pod.UID()
			}
		}
	}
}

func (pm *BasicManager) DeletePod(pod *Pod) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	if pm.secretManager != nil {
		pm.secretManager.UnregisterPod(pod.GetAPIPod())
	}
	podFullName := pod.GetFullName()
	if pod.IsMirror() {
		delete(pm.mirrorPodByUID, pod.UID())
		delete(pm.mirrorPodByFullName, podFullName)
		delete(pm.translationByUID, pod.UID())
	} else {
		delete(pm.podByUID, pod.UID())
		delete(pm.podByFullName, podFullName)
	}
}

func (pm *BasicManager) GetPods() []*Pod {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	return podsMapToPods(pm.podByUID)
}

func (pm *BasicManager) GetPodsAndMirrorPods() ([]*Pod, []*Pod) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pods := podsMapToPods(pm.podByUID)
	mirrorPods := podsMapToPods(pm.mirrorPodByUID)
	return pods, mirrorPods
}

func (pm *BasicManager) GetPodByUID(uid types.UID) (*Pod, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pod, ok := pm.podByUID[uid]
	return pod, ok
}

func (pm *BasicManager) GetPodByName(namespace, name string) (*Pod, bool) {
	podFullName := kubecontainer.BuildPodFullName(name, namespace)
	return pm.GetPodByFullName(podFullName)
}

func (pm *BasicManager) GetPodByFullName(podFullName string) (*Pod, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pod, ok := pm.podByFullName[podFullName]
	return pod, ok
}

func (pm *BasicManager) TranslatePodUID(uid types.UID) types.UID {
	if uid == "" {
		return uid
	}

	pm.lock.RLock()
	defer pm.lock.RUnlock()
	if translated, ok := pm.translationByUID[uid]; ok {
		return translated
	}
	return uid
}

func (pm *BasicManager) GetUIDTranslations() (podToMirror, mirrorToPod map[types.UID]types.UID) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	podToMirror = make(map[types.UID]types.UID, len(pm.translationByUID))
	mirrorToPod = make(map[types.UID]types.UID, len(pm.translationByUID))
	// Insert empty translation mapping for all static pods.
	for uid, pod := range pm.podByUID {
		if !pod.IsStatic() {
			continue
		}
		podToMirror[uid] = ""
	}
	// Fill in translations. Notice that if there is no mirror pod for a
	// static pod, its uid will be translated into empty string "". This
	// is WAI, from the caller side we can know that the static pod doesn't
	// have a corresponding mirror pod instead of using static pod uid directly.
	for k, v := range pm.translationByUID {
		mirrorToPod[k] = v
		podToMirror[v] = k
	}
	return podToMirror, mirrorToPod
}

func (pm *BasicManager) getOrphanedMirrorPodNames() []string {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	var podFullNames []string
	for podFullName := range pm.mirrorPodByFullName {
		if _, ok := pm.podByFullName[podFullName]; !ok {
			podFullNames = append(podFullNames, podFullName)
		}
	}
	return podFullNames
}

func (pm *BasicManager) DeleteOrphanedMirrorPods() {
	podFullNames := pm.getOrphanedMirrorPodNames()
	for _, podFullName := range podFullNames {
		pm.MirrorClient.DeleteMirrorPod(podFullName)
	}
}

func podsMapToPods(UIDMap map[types.UID]*Pod) []*Pod {
	pods := make([]*Pod, 0, len(UIDMap))
	for _, pod := range UIDMap {
		pods = append(pods, pod)
	}
	return pods
}

func (pm *BasicManager) GetMirrorPodByPod(pod *Pod) (*Pod, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	mirrorPod, ok := pm.mirrorPodByFullName[pod.GetFullName()]
	return mirrorPod, ok
}

func (pm *BasicManager) GetPodByMirrorPod(mirrorPod *Pod) (*Pod, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pod, ok := pm.podByFullName[mirrorPod.GetFullName()]
	return pod, ok
}
