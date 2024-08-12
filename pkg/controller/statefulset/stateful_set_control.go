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

package statefulset

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/controller/history"
	"k8s.io/kubernetes/pkg/features"
)

// Realistic value for maximum in-flight requests when processing in parallel mode.
const MaxBatchSize = 500

// StatefulSetControl implements the control logic for updating StatefulSets and their children Pods. It is implemented
// as an interface to allow for extensions that provide different semantics. Currently, there is only one implementation.
type StatefulSetControlInterface interface {
	// UpdateStatefulSet implements the control logic for Pod creation, update, and deletion, and
	// persistent volume creation, update, and deletion.
	// If an implementation returns a non-nil error, the invocation will be retried using a rate-limited strategy.
	// Implementors should sink any errors that they do not wish to trigger a retry, and they may feel free to
	// exit exceptionally at any point provided they wish the update to be re-run at a later point in time.
	UpdateStatefulSet(ctx context.Context, set *apps.StatefulSet, pods []*v1.Pod) (*apps.StatefulSetStatus, error)
	// ListRevisions returns a array of the ControllerRevisions that represent the revisions of set. If the returned
	// error is nil, the returns slice of ControllerRevisions is valid.
	ListRevisions(set *apps.StatefulSet) ([]*apps.ControllerRevision, error)
	// AdoptOrphanRevisions adopts any orphaned ControllerRevisions that match set's Selector. If all adoptions are
	// successful the returned error is nil.
	AdoptOrphanRevisions(set *apps.StatefulSet, revisions []*apps.ControllerRevision) error
}

// NewDefaultStatefulSetControl returns a new instance of the default implementation StatefulSetControlInterface that
// implements the documented semantics for StatefulSets. podControl is the PodControlInterface used to create, update,
// and delete Pods and to create PersistentVolumeClaims. statusUpdater is the StatefulSetStatusUpdaterInterface used
// to update the status of StatefulSets. You should use an instance returned from NewRealStatefulPodControl() for any
// scenario other than testing.
func NewDefaultStatefulSetControl(
	podControl *StatefulPodControl,
	statusUpdater StatefulSetStatusUpdaterInterface,
	controllerHistory history.Interface) StatefulSetControlInterface {
	return &defaultStatefulSetControl{podControl, statusUpdater, controllerHistory}
}

type defaultStatefulSetControl struct {
	podControl        *StatefulPodControl
	statusUpdater     StatefulSetStatusUpdaterInterface
	controllerHistory history.Interface
}

// UpdateStatefulSet executes the core logic loop for a stateful set, applying the predictable and
// consistent monotonic update strategy by default - scale up proceeds in ordinal order, no new pod
// is created while any pod is unhealthy, and pods are terminated in descending order. The burst
// strategy allows these constraints to be relaxed - pods will be created and deleted eagerly and
// in no particular order. Clients using the burst strategy should be careful to ensure they
// understand the consistency implications of having unpredictable numbers of pods available.
func (ssc *defaultStatefulSetControl) UpdateStatefulSet(ctx context.Context, set *apps.StatefulSet, pods []*v1.Pod) (*apps.StatefulSetStatus, error) {
	set = set.DeepCopy() // set is modified when a new revision is created in performUpdate. Make a copy now to avoid mutation errors.

	// list all revisions and sort them
	revisions, err := ssc.ListRevisions(set)
	if err != nil {
		return nil, err
	}
	// get the current, and update revisions
	updateCtx, err := ssc.getStatefulSetRevisions(set, revisions)
	if err != nil {
		return nil, err
	}

	status, err := ssc.performUpdate(ctx, set, pods, updateCtx)
	// maintain the set's revision history limit
	errTrunc := ssc.truncateHistory(set, pods, revisions, updateCtx.currentRevision, updateCtx.updateRevision)
	if err != nil {
		errs := []error{err}
		if agg, ok := err.(utilerrors.Aggregate); ok {
			errs = agg.Errors()
		}
		return nil, utilerrors.NewAggregate(append(errs, errTrunc))
	}

	return status, errTrunc
}

func (ssc *defaultStatefulSetControl) performUpdate(
	ctx context.Context, set *apps.StatefulSet, pods []*v1.Pod, updateCtx *updateContext) (*apps.StatefulSetStatus, error) {
	logger := klog.FromContext(ctx)

	// perform the main update function and get the status
	currentStatus, err := ssc.updateStatefulSet(ctx, updateCtx, pods)
	if err != nil && currentStatus == nil {
		return nil, err
	}

	// make sure to update the latest status even if there is an error with non-nil currentStatus
	statusErr := ssc.updateStatefulSetStatus(ctx, set, currentStatus)
	if statusErr == nil {
		logger.V(4).Info("Updated status", "statefulSet", klog.KObj(set),
			"replicas", currentStatus.Replicas,
			"readyReplicas", currentStatus.ReadyReplicas,
			"currentReplicas", currentStatus.CurrentReplicas,
			"updatedReplicas", currentStatus.UpdatedReplicas)
	}

	switch {
	case err != nil && statusErr != nil:
		logger.Error(statusErr, "Could not update status", "statefulSet", klog.KObj(set))
		return currentStatus, err
	case err != nil:
		return currentStatus, err
	case statusErr != nil:
		return currentStatus, statusErr
	}

	logger.V(4).Info("StatefulSet revisions", "statefulSet", klog.KObj(set),
		"currentRevision", currentStatus.CurrentRevision,
		"updateRevision", currentStatus.UpdateRevision)

	return currentStatus, nil
}

func (ssc *defaultStatefulSetControl) ListRevisions(set *apps.StatefulSet) ([]*apps.ControllerRevision, error) {
	selector, err := metav1.LabelSelectorAsSelector(set.Spec.Selector)
	if err != nil {
		return nil, err
	}
	return ssc.controllerHistory.ListControllerRevisions(set, selector)
}

func (ssc *defaultStatefulSetControl) AdoptOrphanRevisions(
	set *apps.StatefulSet,
	revisions []*apps.ControllerRevision) error {
	for i := range revisions {
		adopted, err := ssc.controllerHistory.AdoptControllerRevision(set, controllerKind, revisions[i])
		if err != nil {
			return err
		}
		revisions[i] = adopted
	}
	return nil
}

// truncateHistory truncates any non-live ControllerRevisions in revisions from set's history. The UpdateRevision and
// CurrentRevision in set's Status are considered to be live. Any revisions associated with the Pods in pods are also
// considered to be live. Non-live revisions are deleted, starting with the revision with the lowest Revision, until
// only RevisionHistoryLimit revisions remain. If the returned error is nil the operation was successful. This method
// expects that revisions is sorted when supplied.
func (ssc *defaultStatefulSetControl) truncateHistory(
	set *apps.StatefulSet,
	pods []*v1.Pod,
	revisions []*apps.ControllerRevision,
	current, update string) error {
	history := make([]*apps.ControllerRevision, 0, len(revisions))
	// mark all live revisions
	live := sets.New(current, update)
	for i := range pods {
		live.Insert(getPodRevision(pods[i]))
	}
	// collect live revisions and historic revisions
	for i := range revisions {
		if !live.Has(revisions[i].Name) {
			history = append(history, revisions[i])
		}
	}
	historyLen := len(history)
	historyLimit := int(*set.Spec.RevisionHistoryLimit)
	if historyLen <= historyLimit {
		return nil
	}
	// delete any non-live history to maintain the revision limit.
	history = history[:(historyLen - historyLimit)]
	for i := 0; i < len(history); i++ {
		if err := ssc.controllerHistory.DeleteControllerRevision(history[i]); err != nil {
			return err
		}
	}
	return nil
}

type revisionCompare struct {
	Spec struct {
		Template             runtime.RawExtension `json:"template"`
		VolumeClaimTemplates runtime.RawExtension `json:"volumeClaimTemplates"`
	} `json:"spec"`
}

func equalRevision(lhs *revisionCompare, rhs *revisionCompare) (pod, all bool) {
	pod = bytes.Equal(lhs.Spec.Template.Raw, rhs.Spec.Template.Raw)
	all = pod && bytes.Equal(lhs.Spec.VolumeClaimTemplates.Raw, rhs.Spec.VolumeClaimTemplates.Raw)
	return
}

// getStatefulSetRevisions returns the current and update ControllerRevisions for set. It also
// returns a collision count that records the number of name collisions set saw when creating
// new ControllerRevisions. This count is incremented on every name collision and is used in
// building the ControllerRevision names for name collision avoidance. This method may create
// a new revision, or modify the Revision of an existing revision if an update to set is detected.
func (ssc *defaultStatefulSetControl) getStatefulSetRevisions(
	set *apps.StatefulSet,
	revisions []*apps.ControllerRevision) (*updateContext, error) {
	var currentRevision, updateRevision *apps.ControllerRevision

	revisionCount := len(revisions)
	history.SortControllerRevisions(revisions)

	// Use a local copy of set.Status.CollisionCount to avoid modifying set.Status directly.
	// This copy is returned so the value gets carried over to set.Status in updateStatefulSet.
	var collisionCount int32
	if set.Status.CollisionCount != nil {
		collisionCount = *set.Status.CollisionCount
	}

	// create a new revision from the current set
	nv := nextRevision(revisions)
	updateRevision, err := newRevision(set, nv, &collisionCount)
	if err != nil {
		return nil, err
	}
	var updateCmp revisionCompare
	utilruntime.Must(json.Unmarshal(updateRevision.Data.Raw, &updateCmp))

	podUpToDateRevisions := sets.New[string]()
	var retrievedRevision *apps.ControllerRevision
	for _, re := range revisions {
		var reCmp revisionCompare
		err := json.Unmarshal(re.Data.Raw, &reCmp)
		if err != nil {
			return nil, fmt.Errorf("unable to decode ControllerRevision %v: %v", klog.KObj(re), err)
		}
		pod, all := equalRevision(&reCmp, &updateCmp)
		if all {
			retrievedRevision = re
		}
		if pod {
			podUpToDateRevisions.Insert(re.Name)
		}
		if re.Name == set.Status.CurrentRevision {
			currentRevision = re
		}
	}

	if retrievedRevision == nil {
		//if there is no equivalent revision we create a new one
		updateRevision, err = ssc.controllerHistory.CreateControllerRevision(set, updateRevision, &collisionCount)
		if err != nil {
			return nil, err
		}
	} else if retrievedRevision != revisions[revisionCount-1] {
		// if the equivalent revision is not immediately prior we will roll back by incrementing the
		// Revision of the equivalent revision
		updateRevision, err = ssc.controllerHistory.UpdateControllerRevision(retrievedRevision, nv)
		if err != nil {
			return nil, err
		}
	} else {
		updateRevision = retrievedRevision
	}

	// if the current revision is nil we initialize the history by setting it to the update revision
	if currentRevision == nil {
		currentRevision = updateRevision
	}

	return newUpdateContext(set, currentRevision, updateRevision, podUpToDateRevisions, collisionCount)
}

func slowStartBatch(initialBatchSize int, remaining int, fn func(int) (bool, error)) (int, error) {
	successes := 0
	j := 0
	for batchSize := min(remaining, initialBatchSize); batchSize > 0; batchSize = min(min(2*batchSize, remaining), MaxBatchSize) {
		errCh := make(chan error, batchSize)
		var wg sync.WaitGroup
		wg.Add(batchSize)
		for i := 0; i < batchSize; i++ {
			go func(k int) {
				defer wg.Done()
				// Ignore the first parameter - relevant for monotonic only.
				if _, err := fn(k); err != nil {
					errCh <- err
				}
			}(j)
			j++
		}
		wg.Wait()
		successes += batchSize - len(errCh)
		close(errCh)
		if len(errCh) > 0 {
			errs := make([]error, 0)
			for err := range errCh {
				errs = append(errs, err)
			}
			return successes, utilerrors.NewAggregate(errs)
		}
		remaining -= batchSize
	}
	return successes, nil
}

func updateStatus(status *apps.StatefulSetStatus, minReadySeconds int32, updateCtx *updateContext, podLists ...[]*v1.Pod) {
	status.Replicas = 0
	status.ReadyReplicas = 0
	status.AvailableReplicas = 0
	status.CurrentReplicas = 0
	status.UpdatedReplicas = 0
	for _, list := range podLists {
		for _, pod := range list {
			if pod == nil {
				continue
			}
			status.Replicas++

			// count the number of running and ready replicas
			if isRunningAndReady(pod) {
				status.ReadyReplicas++
				// count the number of running and available replicas
				if isRunningAndAvailable(pod, minReadySeconds) {
					status.AvailableReplicas++
				}

			}

			// count the number of current and update replicas
			if !isTerminating(pod) {
				revision := getPodRevision(pod)
				if revision == updateCtx.currentRevision {
					status.CurrentReplicas++
				}
				if revision == updateCtx.updateRevision {
					status.UpdatedReplicas++
				}
			}
		}
	}
}

func (ssc *defaultStatefulSetControl) processReplica(
	ctx context.Context,
	updateCtx *updateContext,
	replicas []*v1.Pod,
	i int) (bool, error) {
	logger := klog.FromContext(ctx)
	set, revision := chooseRevision(updateCtx, i)

	// Note that pods with phase Succeeded will also trigger this event. This is
	// because final pod phase of evicted or otherwise forcibly stopped pods
	// (e.g. terminated on node reboot) is determined by the exit code of the
	// container, not by the reason for pod termination. We should restart the pod
	// regardless of the exit code.
	if replicas[i] != nil && (isFailed(replicas[i]) || isSucceeded(replicas[i])) {
		if replicas[i].DeletionTimestamp == nil {
			if err := ssc.podControl.DeleteStatefulPod(set, replicas[i]); err != nil {
				return true, err
			}
		}
		// New pod should be generated on the next sync after the current pod is removed from etcd.
		return true, nil
	}
	// If we find a Pod that has not been created we create the Pod
	if replicas[i] == nil {
		newReplica := newStatefulSetPod(set, getStartOrdinal(set)+i)
		setPodRevision(newReplica, revision)
		if utilfeature.DefaultFeatureGate.Enabled(features.StatefulSetAutoDeletePVC) {
			if isStale, err := ssc.podControl.PodClaimIsStale(set, newReplica); err != nil {
				return true, err
			} else if isStale {
				// If a pod has a stale PVC, no more work can be done this round.
				return true, err
			}
		}
		if utilfeature.DefaultFeatureGate.Enabled(features.UpdateVolumeClaimTemplate) &&
			set.Spec.VolumeClaimUpdatePolicy == apps.InPlaceStatefulSetVolumeClaimUpdatePolicy {
			// update the PVCs before creating the pod to maintain the invariant.
			if err := ssc.podControl.applyPersistentVolumeClaims(ctx, set, newReplica, false); err != nil {
				return true, err
			}
		}
		if err := ssc.podControl.CreateStatefulPod(ctx, set, newReplica); err != nil {
			return true, err
		}
		replicas[i] = newReplica
		if updateCtx.monotonic {
			// if the set does not allow bursting, return immediately
			return true, nil
		}
	}

	// If the Pod is in pending state then trigger PVC creation to create missing PVCs
	if isPending(replicas[i]) {
		claimSet := updateCtx.updateSet
		if utilfeature.DefaultFeatureGate.Enabled(features.UpdateVolumeClaimTemplate) {
			claimSet = set // create the PVCs using the same revision as Pod.
		}
		logger.V(4).Info(
			"StatefulSet is triggering PVC creation for pending Pod",
			"statefulSet", klog.KObj(claimSet), "pod", klog.KObj(replicas[i]))
		if err := ssc.podControl.createMissingPersistentVolumeClaims(ctx, claimSet, replicas[i]); err != nil {
			return true, err
		}
	}

	if updateCtx.monotonic {
		// If we find a Pod that is currently terminating, we must wait until graceful deletion
		// completes before we continue to make progress.
		if isTerminating(replicas[i]) {
			logger.V(4).Info("StatefulSet is waiting for Pod to Terminate",
				"statefulSet", klog.KObj(set), "pod", klog.KObj(replicas[i]))
			return true, nil
		}

		// If we have a Pod that has been created but is not running and ready we can not make progress.
		// We must ensure that all for each Pod, when we create it, all of its predecessors, with respect to its
		// ordinal, are Running and Ready.
		if !isRunningAndReady(replicas[i]) {
			logger.V(4).Info("StatefulSet is waiting for Pod to be Running and Ready",
				"statefulSet", klog.KObj(set), "pod", klog.KObj(replicas[i]))
			return true, nil
		}

		// If we have a Pod that has been created but is not available we can not make progress.
		// We must ensure that all for each Pod, when we create it, all of its predecessors, with respect to its
		// ordinal, are Available.
		if !isRunningAndAvailable(replicas[i], set.Spec.MinReadySeconds) {
			logger.V(4).Info("StatefulSet is waiting for Pod to be Available",
				"statefulSet", klog.KObj(set), "pod", klog.KObj(replicas[i]))
			return true, nil
		}
	}

	// Enforce the StatefulSet invariants
	retentionMatch := true
	if utilfeature.DefaultFeatureGate.Enabled(features.StatefulSetAutoDeletePVC) {
		var err error
		retentionMatch, err = ssc.podControl.ClaimsMatchRetentionPolicy(ctx, set, replicas[i])
		// An error is expected if the pod is not yet fully updated, and so return is treated as matching.
		if err != nil {
			retentionMatch = true
		}
	}

	if identityMatches(set, replicas[i]) && storageMatches(set, replicas[i]) && retentionMatch {
		return false, nil
	}

	// Make a deep copy so we don't mutate the shared cache
	replica := replicas[i].DeepCopy()
	if err := ssc.podControl.UpdateStatefulPod(ctx, set, replica); err != nil {
		return true, err
	}

	return false, nil
}

func (ssc *defaultStatefulSetControl) processCondemned(ctx context.Context, set *apps.StatefulSet, firstUnhealthyOrdinal int, monotonic bool, condemned []*v1.Pod, i int) (bool, error) {
	logger := klog.FromContext(ctx)
	if isTerminating(condemned[i]) {
		// if we are in monotonic mode, block and wait for terminating pods to expire
		if monotonic {
			logger.V(4).Info("StatefulSet is waiting for Pod to Terminate prior to scale down",
				"statefulSet", klog.KObj(set), "pod", klog.KObj(condemned[i]))
			return true, nil
		}
		return false, nil
	}
	// if we are in monotonic mode and the condemned target is not the first unhealthy Pod, block
	if monotonic && getOrdinal(condemned[i]) != firstUnhealthyOrdinal {
		if !isRunningAndReady(condemned[i]) {
			logger.V(4).Info("StatefulSet is waiting for Pod to be Running and Ready prior to scale down",
				"statefulSet", klog.KObj(set), "ordinal", firstUnhealthyOrdinal)
			return true, nil
		}
		if !isRunningAndAvailable(condemned[i], set.Spec.MinReadySeconds) {
			logger.V(4).Info("StatefulSet is waiting for Pod to be Available prior to scale down",
				"statefulSet", klog.KObj(set), "ordinal", firstUnhealthyOrdinal)
			return true, nil
		}
	}

	logger.V(2).Info("Pod of StatefulSet is terminating for scale down",
		"statefulSet", klog.KObj(set), "pod", klog.KObj(condemned[i]))
	return true, ssc.podControl.DeleteStatefulPod(set, condemned[i])
}

type updateContext struct {
	currentSet *apps.StatefulSet
	updateSet  *apps.StatefulSet

	currentRevision string
	updateRevision  string

	monotonic      bool
	collisionCount int32

	podUpToDateRevisions sets.Set[string]
}

func newUpdateContext(set *apps.StatefulSet, currentRevision, updateRevision *apps.ControllerRevision,
	podUpToDateRevisions sets.Set[string], collisionCount int32) (*updateContext, error) {

	c := &updateContext{
		currentRevision: currentRevision.Name,
		updateRevision:  updateRevision.Name,
		monotonic:       !allowsBurst(set),
		collisionCount:  collisionCount,

		podUpToDateRevisions: podUpToDateRevisions,
	}
	var err error
	// get the current and update revisions of the set.
	c.currentSet, err = ApplyRevision(set, currentRevision)
	if err != nil {
		return nil, err
	}
	c.updateSet, err = ApplyRevision(set, updateRevision)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *updateContext) runForAll(ctx context.Context, pods []*v1.Pod, fn func(ctx context.Context, updateCtx *updateContext, pods []*v1.Pod, i int) (bool, error)) (bool, error) {
	if c.monotonic {
		for i := range pods {
			if shouldExit, err := fn(ctx, c, pods, i); shouldExit || err != nil {
				return true, err
			}
		}
	} else {
		if _, err := slowStartBatch(1, len(pods), func(i int) (bool, error) {
			return fn(ctx, c, pods, i)
		}); err != nil {
			return true, err
		}
	}
	return false, nil
}

func (c *updateContext) isPodUpToDate(pod *v1.Pod) bool {
	return c.podUpToDateRevisions.Has(getPodRevision(pod))
}

// updateStatefulSet performs the update function for a StatefulSet. This method creates, updates, and deletes Pods in
// the set in order to conform the system to the target state for the set. The target state always contains
// set.Spec.Replicas Pods with a Ready Condition. If the UpdateStrategy.Type for the set is
// RollingUpdateStatefulSetStrategyType then all Pods in the set must be at set.Status.CurrentRevision.
// If the UpdateStrategy.Type for the set is OnDeleteStatefulSetStrategyType, the target state implies nothing about
// the revisions of Pods in the set. If the UpdateStrategy.Type for the set is PartitionStatefulSetStrategyType, then
// all Pods with ordinal less than UpdateStrategy.Partition.Ordinal must be at Status.CurrentRevision and all other
// Pods must be at Status.UpdateRevision. If the returned error is nil, the returned StatefulSetStatus is valid and the
// update must be recorded. If the error is not nil, the method should be retried until successful.
func (ssc *defaultStatefulSetControl) updateStatefulSet(
	ctx context.Context,
	updateCtx *updateContext,
	pods []*v1.Pod) (*apps.StatefulSetStatus, error) {
	logger := klog.FromContext(ctx)
	set := updateCtx.updateSet

	// set the generation, and revisions in the returned status
	status := apps.StatefulSetStatus{}
	status.ObservedGeneration = set.Generation
	status.CurrentRevision = updateCtx.currentRevision
	status.UpdateRevision = updateCtx.updateRevision
	status.CollisionCount = new(int32)
	*status.CollisionCount = updateCtx.collisionCount

	updateStatus(&status, set.Spec.MinReadySeconds, updateCtx, pods)

	replicaCount := int(*set.Spec.Replicas)
	// slice that will contain all Pods such that getStartOrdinal(set) <= getOrdinal(pod) <= getEndOrdinal(set)
	replicas := make([]*v1.Pod, replicaCount)
	// slice that will contain all Pods such that getOrdinal(pod) < getStartOrdinal(set) OR getOrdinal(pod) > getEndOrdinal(set)
	condemned := make([]*v1.Pod, 0, len(pods))
	unhealthy := 0
	firstUnhealthyOrdinal := -1

	// First we partition pods into two lists valid replicas and condemned Pods
	for _, pod := range pods {
		if podInOrdinalRange(pod, set) {
			// if the ordinal of the pod is within the range of the current number of replicas,
			// insert it at the indirection of its ordinal
			replicas[getOrdinal(pod)-getStartOrdinal(set)] = pod
		} else if getOrdinal(pod) >= 0 {
			// if the ordinal is valid, but not within the range add it to the condemned list
			condemned = append(condemned, pod)
		}
		// If the ordinal could not be parsed (ord < 0), ignore the Pod.
	}

	// sort the condemned Pods by their ordinals
	sort.Sort(descendingOrdinal(condemned))

	// find the first unhealthy Pod
	for i, pod := range replicas {
		if pod == nil || !isHealthy(pod) {
			unhealthy++
			if firstUnhealthyOrdinal < 0 {
				firstUnhealthyOrdinal = getStartOrdinal(set) + i
			}
		}
	}

	// or the first unhealthy condemned Pod (condemned are sorted in descending order for ease of use)
	for i := len(condemned) - 1; i >= 0; i-- {
		if !isHealthy(condemned[i]) {
			unhealthy++
			if firstUnhealthyOrdinal < 0 {
				firstUnhealthyOrdinal = getOrdinal(condemned[i])
			}
		}
	}

	if unhealthy > 0 {
		logger.V(4).Info("StatefulSet has unhealthy Pods", "statefulSet", klog.KObj(set), "unhealthyReplicas", unhealthy, "ordinal", firstUnhealthyOrdinal)
	}

	// If the StatefulSet is being deleted, don't do anything other than updating
	// status.
	if set.DeletionTimestamp != nil {
		return &status, nil
	}

	// First, process each living replica. Exit if we run into an error or something blocking in monotonic mode.
	if shouldExit, err := updateCtx.runForAll(ctx, replicas, ssc.processReplica); shouldExit || err != nil {
		updateStatus(&status, set.Spec.MinReadySeconds, updateCtx, replicas, condemned)
		return &status, err
	}

	// Fix pod claims for condemned pods, if necessary.
	if utilfeature.DefaultFeatureGate.Enabled(features.StatefulSetAutoDeletePVC) {
		if shouldExit, err := updateCtx.runForAll(ctx, condemned, ssc.fixPodClaim); shouldExit || err != nil {
			updateStatus(&status, set.Spec.MinReadySeconds, updateCtx, replicas, condemned)
			return &status, err
		}
	}

	// At this point, in monotonic mode all of the current Replicas are Running, Ready and Available,
	// and we can consider termination.
	// We will wait for all predecessors to be Running and Ready prior to attempting a deletion.
	// We will terminate Pods in a monotonically decreasing order.
	// Note that we do not resurrect Pods in this interval. Also note that scaling will take precedence over
	// updates.
	processCondemnedFn := func(ctx context.Context, updateCtx *updateContext, condemned []*v1.Pod, i int) (bool, error) {
		return ssc.processCondemned(ctx, updateCtx.updateSet, firstUnhealthyOrdinal, updateCtx.monotonic, condemned, i)
	}
	if shouldExit, err := updateCtx.runForAll(ctx, condemned, processCondemnedFn); shouldExit || err != nil {
		updateStatus(&status, set.Spec.MinReadySeconds, updateCtx, replicas, condemned)
		return &status, err
	}

	updateStatus(&status, set.Spec.MinReadySeconds, updateCtx, replicas, condemned)

	// for the OnDelete strategy we short circuit. Pods will be updated when they are manually deleted.
	if set.Spec.UpdateStrategy.Type == apps.OnDeleteStatefulSetStrategyType {
		return &status, nil
	}

	if utilfeature.DefaultFeatureGate.Enabled(features.MaxUnavailableStatefulSet) {
		return updateStatefulSetAfterInvariantEstablished(ctx,
			ssc,
			set,
			replicas,
			updateCtx.updateRevision,
			status,
		)
	}

	// we compute the minimum ordinal of the target sequence for a destructive update based on the strategy.
	updateMin := 0
	if set.Spec.UpdateStrategy.RollingUpdate != nil {
		updateMin = int(*set.Spec.UpdateStrategy.RollingUpdate.Partition)
	}
	// we update the Pod with the largest ordinal that does not match the update revision.
	for target := len(replicas) - 1; target >= updateMin; target-- {

		if getPodRevision(replicas[target]) != updateCtx.updateRevision && !isTerminating(replicas[target]) {
			if updateCtx.isPodUpToDate(replicas[target]) {
				// Pod template unchanged, update PVCs then Pod revision label
				logger.V(2).Info("Pod of StatefulSet is unchanged",
					"statefulSet", klog.KObj(set), "pod", klog.KObj(replicas[target]))

				if utilfeature.DefaultFeatureGate.Enabled(features.UpdateVolumeClaimTemplate) &&
					set.Spec.VolumeClaimUpdatePolicy == apps.InPlaceStatefulSetVolumeClaimUpdatePolicy {

					err := ssc.podControl.applyPersistentVolumeClaims(ctx, set, replicas[target], false)
					if err != nil {
						return &status, err
					}
				}
				setPodRevision(replicas[target], updateCtx.updateRevision)
				if err := ssc.podControl.objectMgr.UpdatePod(replicas[target]); err != nil {
					if !errors.IsNotFound(err) {
						return &status, err
					}
				}
			} else {
				// delete the Pod if it is not already terminating and does not match the update revision.
				// PVCs will be updated before creating new Pod.
				logger.V(2).Info("Pod of StatefulSet is terminating for update",
					"statefulSet", klog.KObj(set), "pod", klog.KObj(replicas[target]))
				if err := ssc.podControl.DeleteStatefulPod(set, replicas[target]); err != nil {
					if !errors.IsNotFound(err) {
						return &status, err
					}
				}
			}
			status.CurrentReplicas--
			return &status, nil
		}

		if ready, err := ssc.podControl.readyForUpdate(ctx, set, replicas[target]); err != nil || !ready {
			return &status, err
		}

	}
	return &status, nil
}

func (ssc *defaultStatefulSetControl) fixPodClaim(ctx context.Context, updateCtx *updateContext, condemned []*v1.Pod, i int) (bool, error) {
	if matchPolicy, err := ssc.podControl.ClaimsMatchRetentionPolicy(ctx, updateCtx.updateSet, condemned[i]); err != nil {
		return true, err
	} else if !matchPolicy {
		if err := ssc.podControl.UpdatePodClaimForRetentionPolicy(ctx, updateCtx.updateSet, condemned[i]); err != nil {
			return true, err
		}
	}
	return false, nil
}

func updateStatefulSetAfterInvariantEstablished(
	ctx context.Context,
	ssc *defaultStatefulSetControl,
	set *apps.StatefulSet,
	replicas []*v1.Pod,
	updateRevision string,
	status apps.StatefulSetStatus,
) (*apps.StatefulSetStatus, error) {

	logger := klog.FromContext(ctx)
	replicaCount := int(*set.Spec.Replicas)

	// we compute the minimum ordinal of the target sequence for a destructive update based on the strategy.
	updateMin := 0
	maxUnavailable := 1
	if set.Spec.UpdateStrategy.RollingUpdate != nil {
		updateMin = int(*set.Spec.UpdateStrategy.RollingUpdate.Partition)

		// if the feature was enabled and then later disabled, MaxUnavailable may have a value
		// more than 1. Ignore the passed in value and Use maxUnavailable as 1 to enforce
		// expected behavior when feature gate is not enabled.
		var err error
		maxUnavailable, err = getStatefulSetMaxUnavailable(set.Spec.UpdateStrategy.RollingUpdate.MaxUnavailable, replicaCount)
		if err != nil {
			return &status, err
		}
	}

	// Collect all targets in the range between getStartOrdinal(set) and getEndOrdinal(set). Count any targets in that range
	// that are unhealthy i.e. terminated or not running and ready as unavailable). Select the
	// (MaxUnavailable - Unavailable) Pods, in order with respect to their ordinal for termination. Delete
	// those pods and count the successful deletions. Update the status with the correct number of deletions.
	unavailablePods := 0
	for target := len(replicas) - 1; target >= 0; target-- {
		if !isHealthy(replicas[target]) {
			unavailablePods++
		}
	}

	if unavailablePods >= maxUnavailable {
		logger.V(2).Info("StatefulSet found unavailablePods, more than or equal to allowed maxUnavailable",
			"statefulSet", klog.KObj(set),
			"unavailablePods", unavailablePods,
			"maxUnavailable", maxUnavailable)
		return &status, nil
	}

	// Now we need to delete MaxUnavailable- unavailablePods
	// start deleting one by one starting from the highest ordinal first
	podsToDelete := maxUnavailable - unavailablePods

	deletedPods := 0
	for target := len(replicas) - 1; target >= updateMin && deletedPods < podsToDelete; target-- {

		// delete the Pod if it is healthy and the revision doesnt match the target
		if getPodRevision(replicas[target]) != updateRevision && !isTerminating(replicas[target]) {
			// delete the Pod if it is healthy and the revision doesnt match the target
			logger.V(2).Info("StatefulSet terminating Pod for update",
				"statefulSet", klog.KObj(set),
				"pod", klog.KObj(replicas[target]))
			if err := ssc.podControl.DeleteStatefulPod(set, replicas[target]); err != nil {
				if !errors.IsNotFound(err) {
					return &status, err
				}
			}
			deletedPods++
			status.CurrentReplicas--
		}
	}
	return &status, nil
}

// updateStatefulSetStatus updates set's Status to be equal to status. If status indicates a complete update, it is
// mutated to indicate completion. If status is semantically equivalent to set's Status no update is performed. If the
// returned error is nil, the update is successful.
func (ssc *defaultStatefulSetControl) updateStatefulSetStatus(
	ctx context.Context,
	set *apps.StatefulSet,
	status *apps.StatefulSetStatus) error {
	// complete any in progress rolling update if necessary
	completeRollingUpdate(set, status)

	// if the status is not inconsistent do not perform an update
	if !inconsistentStatus(set, status) {
		return nil
	}

	// copy set and update its status
	set = set.DeepCopy()
	if err := ssc.statusUpdater.UpdateStatefulSetStatus(ctx, set, status); err != nil {
		return err
	}

	return nil
}

var _ StatefulSetControlInterface = &defaultStatefulSetControl{}
