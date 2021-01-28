// +build linux

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

package cm

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	featuregatetesting "k8s.io/component-base/featuregate/testing"
	pkgfeatures "k8s.io/kubernetes/pkg/features"
)

// getResourceList returns a ResourceList with the
// specified cpu and memory resource values
func getResourceList(cpu, memory string) v1.ResourceList {
	res := v1.ResourceList{}
	if cpu != "" {
		res[v1.ResourceCPU] = resource.MustParse(cpu)
	}
	if memory != "" {
		res[v1.ResourceMemory] = resource.MustParse(memory)
	}
	return res
}

// getResourceRequirements returns a ResourceRequirements object
func getResourceRequirements(requests, limits v1.ResourceList) v1.ResourceRequirements {
	res := v1.ResourceRequirements{}
	res.Requests = requests
	res.Limits = limits
	return res
}

func TestResourceConfigForPod(t *testing.T) {
	defaultQuotaPeriod := uint64(100 * time.Millisecond / time.Microsecond)
	tunedQuotaPeriod := uint64(5 * time.Millisecond / time.Microsecond)

	minShares := uint64(MinShares)
	burstableShares := MilliCPUToShares(100)
	memoryQuantity := resource.MustParse("200Mi")
	burstableMemory := memoryQuantity.Value()
	burstablePartialShares := MilliCPUToShares(200)
	burstableQuota := MilliCPUToQuota(200, int64(defaultQuotaPeriod))
	guaranteedShares := MilliCPUToShares(100)
	guaranteedQuota := MilliCPUToQuota(100, int64(defaultQuotaPeriod))
	guaranteedTunedQuota := MilliCPUToQuota(100, int64(tunedQuotaPeriod))
	memoryQuantity = resource.MustParse("100Mi")
	cpuNoLimit := int64(-1)
	guaranteedMemory := memoryQuantity.Value()
	testCases := map[string]struct {
		pod              *v1.Pod
		expected         *ResourceConfig
		enforceCPULimits bool
		quotaPeriod      uint64
		cpuManagerPolicy string
	}{
		"besteffort": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("", ""), getResourceList("", "")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &minShares},
			cpuManagerPolicy: "none",
		},
		"burstable-no-limits": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("", "")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstableShares},
			cpuManagerPolicy: "none",
		},
		"burstable-with-limits": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstableShares, CpuQuota: &burstableQuota, CpuPeriod: &defaultQuotaPeriod, Memory: &burstableMemory},
			cpuManagerPolicy: "none",
		},
		"burstable-with-limits-no-cpu-enforcement": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
					},
				},
			},
			enforceCPULimits: false,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstableShares, CpuQuota: &cpuNoLimit, CpuPeriod: &defaultQuotaPeriod, Memory: &burstableMemory},
			cpuManagerPolicy: "none",
		},
		"burstable-partial-limits": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("", "")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstablePartialShares},
			cpuManagerPolicy: "none",
		},
		"burstable-with-limits-with-tuned-quota": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstableShares, CpuQuota: &burstableQuota, CpuPeriod: &tunedQuotaPeriod, Memory: &burstableMemory},
			cpuManagerPolicy: "none",
		},
		"burstable-with-limits-no-cpu-enforcement-with-tuned-quota": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
					},
				},
			},
			enforceCPULimits: false,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstableShares, CpuQuota: &cpuNoLimit, CpuPeriod: &tunedQuotaPeriod, Memory: &burstableMemory},
			cpuManagerPolicy: "none",
		},
		"burstable-partial-limits-with-tuned-quota": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("", "")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstablePartialShares},
			cpuManagerPolicy: "none",
		},
		"guaranteed": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &guaranteedQuota, CpuPeriod: &defaultQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "none",
		},
		"guaranteed-with-cpu-manager-policy-static": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &cpuNoLimit, CpuPeriod: &defaultQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "static",
		},
		"guaranteed-no-cpu-enforcement": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: false,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &cpuNoLimit, CpuPeriod: &defaultQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "none",
		},
		"guaranteed-with-tuned-quota-and-cpu-manager-policy-static": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &cpuNoLimit, CpuPeriod: &tunedQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "static",
		},
		"guaranteed-with-tuned-quota": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &guaranteedTunedQuota, CpuPeriod: &tunedQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "none",
		},
		"guaranteed-no-cpu-enforcement-with-tuned-quota": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: false,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &cpuNoLimit, CpuPeriod: &tunedQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "none",
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			if testCase.cpuManagerPolicy == "static" {
				defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, pkgfeatures.CPUManager, true)()
			}
			actual := ResourceConfigForPod(testCase.pod, testCase.enforceCPULimits, testCase.quotaPeriod, testCase.cpuManagerPolicy)

			if !reflect.DeepEqual(actual.CpuPeriod, testCase.expected.CpuPeriod) {
				t.Errorf("unexpected result, test: %v, cpu period not as expected", testName)
			}
			if !reflect.DeepEqual(actual.CpuQuota, testCase.expected.CpuQuota) {
				t.Errorf("unexpected result, test: %v, cpu quota not as expected", testName)
			}
			if !reflect.DeepEqual(actual.CpuShares, testCase.expected.CpuShares) {
				t.Errorf("unexpected result, test: %v, cpu shares not as expected", testName)
			}
			if !reflect.DeepEqual(actual.Memory, testCase.expected.Memory) {
				t.Errorf("unexpected result, test: %v, memory not as expected", testName)
			}
		})
	}
}

func TestResourceConfigForPodWithCustomCPUCFSQuotaPeriod(t *testing.T) {
	defaultQuotaPeriod := uint64(100 * time.Millisecond / time.Microsecond)
	tunedQuotaPeriod := uint64(5 * time.Millisecond / time.Microsecond)
	tunedQuota := int64(1 * time.Millisecond / time.Microsecond)

	defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, pkgfeatures.CPUCFSQuotaPeriod, true)()

	minShares := uint64(MinShares)
	burstableShares := MilliCPUToShares(100)
	memoryQuantity := resource.MustParse("200Mi")
	burstableMemory := memoryQuantity.Value()
	burstablePartialShares := MilliCPUToShares(200)
	burstableQuota := MilliCPUToQuota(200, int64(defaultQuotaPeriod))
	guaranteedShares := MilliCPUToShares(100)
	guaranteedQuota := MilliCPUToQuota(100, int64(defaultQuotaPeriod))
	guaranteedTunedQuota := MilliCPUToQuota(100, int64(tunedQuotaPeriod))
	memoryQuantity = resource.MustParse("100Mi")
	cpuNoLimit := int64(-1)
	guaranteedMemory := memoryQuantity.Value()
	testCases := map[string]struct {
		pod              *v1.Pod
		expected         *ResourceConfig
		enforceCPULimits bool
		quotaPeriod      uint64
		cpuManagerPolicy string
	}{
		"besteffort": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("", ""), getResourceList("", "")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &minShares},
			cpuManagerPolicy: "none",
		},
		"burstable-no-limits": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("", "")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstableShares},
			cpuManagerPolicy: "none",
		},
		"burstable-with-limits": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstableShares, CpuQuota: &burstableQuota, CpuPeriod: &defaultQuotaPeriod, Memory: &burstableMemory},
			cpuManagerPolicy: "none",
		},
		"burstable-with-limits-no-cpu-enforcement": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
					},
				},
			},
			enforceCPULimits: false,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstableShares, CpuQuota: &cpuNoLimit, CpuPeriod: &defaultQuotaPeriod, Memory: &burstableMemory},
			cpuManagerPolicy: "none",
		},
		"burstable-partial-limits": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("", "")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstablePartialShares},
			cpuManagerPolicy: "none",
		},
		"burstable-with-limits-with-tuned-quota": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstableShares, CpuQuota: &tunedQuota, CpuPeriod: &tunedQuotaPeriod, Memory: &burstableMemory},
			cpuManagerPolicy: "none",
		},
		"burstable-with-limits-no-cpu-enforcement-with-tuned-quota": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
					},
				},
			},
			enforceCPULimits: false,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstableShares, CpuQuota: &cpuNoLimit, CpuPeriod: &tunedQuotaPeriod, Memory: &burstableMemory},
			cpuManagerPolicy: "none",
		},
		"burstable-partial-limits-with-tuned-quota": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("200m", "200Mi")),
						},
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("", "")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &burstablePartialShares},
			cpuManagerPolicy: "none",
		},
		"guaranteed": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &guaranteedQuota, CpuPeriod: &defaultQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "none",
		},
		"guaranteed-with-cpu-manager-policy-static": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &cpuNoLimit, CpuPeriod: &defaultQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "static",
		},
		"guaranteed-no-cpu-enforcement": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: false,
			quotaPeriod:      defaultQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &cpuNoLimit, CpuPeriod: &defaultQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "none",
		},
		"guaranteed-with-tuned-quota": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &guaranteedTunedQuota, CpuPeriod: &tunedQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "none",
		},
		"guaranteed-with-tuned-quota-and-cpu-manager-policy-static": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: true,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &cpuNoLimit, CpuPeriod: &tunedQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "static",
		},
		"guaranteed-no-cpu-enforcement-with-tuned-quota": {
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: getResourceRequirements(getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
						},
					},
				},
			},
			enforceCPULimits: false,
			quotaPeriod:      tunedQuotaPeriod,
			expected:         &ResourceConfig{CpuShares: &guaranteedShares, CpuQuota: &cpuNoLimit, CpuPeriod: &tunedQuotaPeriod, Memory: &guaranteedMemory},
			cpuManagerPolicy: "none",
		},
	}

	for testName, testCase := range testCases {

		t.Run(testName, func(t *testing.T) {
			if testCase.cpuManagerPolicy == "static" {
				defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, pkgfeatures.CPUManager, true)()
			}
			actual := ResourceConfigForPod(testCase.pod, testCase.enforceCPULimits, testCase.quotaPeriod, testCase.cpuManagerPolicy)

			if !reflect.DeepEqual(actual.CpuPeriod, testCase.expected.CpuPeriod) {
				t.Errorf("unexpected result, test: %v, cpu period not as expected", testName)
			}
			if !reflect.DeepEqual(actual.CpuQuota, testCase.expected.CpuQuota) {
				t.Errorf("unexpected result, test: %v, cpu quota not as expected", testName)
			}
			if !reflect.DeepEqual(actual.CpuShares, testCase.expected.CpuShares) {
				t.Errorf("unexpected result, test: %v, cpu shares not as expected", testName)
			}
			if !reflect.DeepEqual(actual.Memory, testCase.expected.Memory) {
				t.Errorf("unexpected result, test: %v, memory not as expected", testName)
			}
		})
	}
}

func TestHugePageLimits(t *testing.T) {
	Mi := int64(1024 * 1024)
	type inputStruct struct {
		key   string
		input string
	}

	testCases := []struct {
		name     string
		inputs   []inputStruct
		expected map[int64]int64
	}{
		{
			name: "no valid hugepages",
			inputs: []inputStruct{
				{
					key:   "2Mi",
					input: "128",
				},
			},
			expected: map[int64]int64{},
		},
		{
			name: "2Mi only",
			inputs: []inputStruct{
				{
					key:   v1.ResourceHugePagesPrefix + "2Mi",
					input: "128",
				},
			},
			expected: map[int64]int64{2 * Mi: 128},
		},
		{
			name: "2Mi and 4Mi",
			inputs: []inputStruct{
				{
					key:   v1.ResourceHugePagesPrefix + "2Mi",
					input: "128",
				},
				{
					key:   v1.ResourceHugePagesPrefix + strconv.FormatInt(2*Mi, 10),
					input: "256",
				},
				{
					key:   v1.ResourceHugePagesPrefix + "4Mi",
					input: "512",
				},
				{
					key:   "4Mi",
					input: "1024",
				},
			},
			expected: map[int64]int64{2 * Mi: 384, 4 * Mi: 512},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			resourceList := v1.ResourceList{}

			for _, input := range testcase.inputs {
				value, err := resource.ParseQuantity(input.input)
				if err != nil {
					t.Fatalf("error in parsing hugepages, value: %s", input.input)
				} else {
					resourceList[v1.ResourceName(input.key)] = value
				}
			}

			resultValue := HugePageLimits(resourceList)

			if !reflect.DeepEqual(testcase.expected, resultValue) {
				t.Errorf("unexpected result, expected: %v, actual: %v", testcase.expected, resultValue)
			}
		})

	}
}
