/*
Copyright 2017 The Kubernetes Authors.

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

package hugepagehandler

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"

	units "github.com/docker/go-units"
	libcontainercgroups "github.com/opencontainers/runc/libcontainer/cgroups"
	cgroupfs "github.com/opencontainers/runc/libcontainer/cgroups/fs"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

type handler struct {
	// containerRuntime is the container runtime service interface needed
	// to make UpdateContainerResources() calls against the containers.
	containerRuntime runtimeService
}

var _ Handler = &handler{}

// NewHandler creates new handler
func NewHandler() (Handler, error) {
	handler := &handler{}
	return handler, nil
}

func (m *handler) Start(containerRuntime runtimeService) {
	m.containerRuntime = containerRuntime
}

func (m *handler) AddContainer(p *v1.Pod, c *v1.Container, containerID string) error {
	hugepageLimits := []*runtimeapi.HugepageLimit{}
	pageSizes := sets.NewString()

	for resourceObj, amountObj := range c.Resources.Limits {
		// Check if resourceObj is "hugepages-<pageSize>". api.ResourceHugePagesPrefix = "hugepages-"
		if v1helper.IsHugePageResourceName(resourceObj) {
			pageSize, err := v1helper.HugePageSizeFromResourceName(resourceObj)
			if err != nil {
				klog.Infof("[hugepagehandler] Fail to parse hugepage size")
				continue
			}
			sizeString := units.CustomSize("%g%s", float64(pageSize.Value()), 1024.0, libcontainercgroups.HugePageSizeUnitList)
			hugepageLimits = append(hugepageLimits, &runtimeapi.HugepageLimit{
				PageSize: sizeString,
				Limit:    uint64(amountObj.Value()),
			})
			pageSizes.Insert(sizeString)
		}
	}

	// For each page size omitted, limit to 0
	for _, pageSize := range cgroupfs.HugePageSizes {
		if pageSizes.Has(pageSize) {
			continue
		}
		hugepageLimits = append(hugepageLimits, &runtimeapi.HugepageLimit{
			PageSize: pageSize,
			Limit:    uint64(0),
		})
	}

	if err := m.updateContainerHugepageLimit(containerID, hugepageLimits); err != nil {
		return fmt.Errorf("[hugepagehandler] AddContainer error: %v", err)
	}

	return nil
}

func (m *handler) updateContainerHugepageLimit(containerID string, hugepageLimit []*runtimeapi.HugepageLimit) error {
	return m.containerRuntime.UpdateContainerResources(
		containerID,
		&runtimeapi.LinuxContainerResources{
			HugepageLimits: hugepageLimit,
		})
}
