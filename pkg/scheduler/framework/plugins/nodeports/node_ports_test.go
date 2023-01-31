/*
Copyright 2019 The Kubernetes Authors.

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

package nodeports

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	st "k8s.io/kubernetes/pkg/scheduler/testing"
)

func newPod(host string, hostPortInfos ...string) *v1.Pod {
	networkPorts := []v1.ContainerPort{}
	for _, portInfo := range hostPortInfos {
		splited := strings.Split(portInfo, "/")
		hostPort, _ := strconv.Atoi(splited[2])

		networkPorts = append(networkPorts, v1.ContainerPort{
			HostIP:   splited[1],
			HostPort: int32(hostPort),
			Protocol: v1.Protocol(splited[0]),
		})
	}
	return st.MakePod().Node(host).ContainerPort(networkPorts).Obj()
}

func TestNodePorts(t *testing.T) {
	tests := []struct {
		pod                 *v1.Pod
		nodeInfo            *framework.NodeInfo
		name                string
		wantPreFilterStatus *framework.Status
		wantFilterStatus    *framework.Status
	}{
		{
			pod:                 &v1.Pod{},
			nodeInfo:            framework.NewNodeInfo(),
			name:                "skip scheduling",
			wantPreFilterStatus: framework.NewStatus(framework.Skip),
		},
		{
			pod: newPod("m1", "UDP/127.0.0.1/8080"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "UDP/127.0.0.1/9090")),
			name: "other port",
		},
		{
			pod: newPod("m1", "UDP/127.0.0.1/8080"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "UDP/127.0.0.1/8080")),
			name:             "same udp port",
			wantFilterStatus: framework.NewStatus(framework.Unschedulable, ErrReason),
		},
		{
			pod: newPod("m1", "TCP/127.0.0.1/8080"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "TCP/127.0.0.1/8080")),
			name:             "same tcp port",
			wantFilterStatus: framework.NewStatus(framework.Unschedulable, ErrReason),
		},
		{
			pod: newPod("m1", "TCP/127.0.0.1/8080"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "TCP/127.0.0.2/8080")),
			name: "different host ip",
		},
		{
			pod: newPod("m1", "UDP/127.0.0.1/8080"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "TCP/127.0.0.1/8080")),
			name: "different protocol",
		},
		{
			pod: newPod("m1", "UDP/127.0.0.1/8000", "UDP/127.0.0.1/8080"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "UDP/127.0.0.1/8080")),
			name:             "second udp port conflict",
			wantFilterStatus: framework.NewStatus(framework.Unschedulable, ErrReason),
		},
		{
			pod: newPod("m1", "TCP/127.0.0.1/8001", "UDP/127.0.0.1/8080"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "TCP/127.0.0.1/8001", "UDP/127.0.0.1/8081")),
			name:             "first tcp port conflict",
			wantFilterStatus: framework.NewStatus(framework.Unschedulable, ErrReason),
		},
		{
			pod: newPod("m1", "TCP/0.0.0.0/8001"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "TCP/127.0.0.1/8001")),
			name:             "first tcp port conflict due to 0.0.0.0 hostIP",
			wantFilterStatus: framework.NewStatus(framework.Unschedulable, ErrReason),
		},
		{
			pod: newPod("m1", "TCP/10.0.10.10/8001", "TCP/0.0.0.0/8001"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "TCP/127.0.0.1/8001")),
			name:             "TCP hostPort conflict due to 0.0.0.0 hostIP",
			wantFilterStatus: framework.NewStatus(framework.Unschedulable, ErrReason),
		},
		{
			pod: newPod("m1", "TCP/127.0.0.1/8001"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "TCP/0.0.0.0/8001")),
			name:             "second tcp port conflict to 0.0.0.0 hostIP",
			wantFilterStatus: framework.NewStatus(framework.Unschedulable, ErrReason),
		},
		{
			pod: newPod("m1", "UDP/127.0.0.1/8001"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "TCP/0.0.0.0/8001")),
			name: "second different protocol",
		},
		{
			pod: newPod("m1", "UDP/127.0.0.1/8001"),
			nodeInfo: framework.NewNodeInfo(
				newPod("m1", "TCP/0.0.0.0/8001", "UDP/0.0.0.0/8001")),
			name:             "UDP hostPort conflict due to 0.0.0.0 hostIP",
			wantFilterStatus: framework.NewStatus(framework.Unschedulable, ErrReason),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, _ := New(nil, nil)
			cycleState := framework.NewCycleState()
			_, preFilterStatus := p.(framework.PreFilterPlugin).PreFilter(context.Background(), cycleState, test.pod)
			if test.wantPreFilterStatus != nil {
				if diff := cmp.Diff(test.wantPreFilterStatus, preFilterStatus); diff != "" {
					t.Errorf("preFilter: status does not match (-want,+got):\n%s", diff)
				}
				return
			}
			if !preFilterStatus.IsSuccess() {
				t.Errorf("prefilter failed with status: %v", preFilterStatus)
				return
			}
			filterStatus := p.(framework.FilterPlugin).Filter(context.Background(), cycleState, test.pod, test.nodeInfo)
			if diff := cmp.Diff(test.wantFilterStatus, filterStatus); diff != "" {
				t.Errorf("filter: status does not match (-want,+got):\n%s", diff)
				return
			}
		})
	}
}

func TestPreFilterDisabled(t *testing.T) {
	pod := &v1.Pod{}
	nodeInfo := framework.NewNodeInfo()
	node := v1.Node{}
	nodeInfo.SetNode(&node)
	p, _ := New(nil, nil)
	cycleState := framework.NewCycleState()
	gotStatus := p.(framework.FilterPlugin).Filter(context.Background(), cycleState, pod, nodeInfo)
	wantStatus := framework.AsStatus(fmt.Errorf(`reading "PreFilterNodePorts" from cycleState: %w`, framework.ErrNotFound))
	if !reflect.DeepEqual(*wantStatus, *gotStatus) {
		t.Errorf("status does not match: %v, want: %v", gotStatus, wantStatus)
		return
	}
}

func TestGetContainerPorts(t *testing.T) {
	tests := []struct {
		pod1     *v1.Pod
		pod2     *v1.Pod
		expected []*v1.ContainerPort
	}{
		{
			pod1: st.MakePod().ContainerPort([]v1.ContainerPort{
				{
					ContainerPort: 8001,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8002,
					Protocol:      v1.ProtocolTCP,
				},
			}).ContainerPort([]v1.ContainerPort{
				{
					ContainerPort: 8003,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8004,
					Protocol:      v1.ProtocolTCP,
				},
			}).Obj(),
			pod2: st.MakePod().ContainerPort([]v1.ContainerPort{
				{
					ContainerPort: 8011,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8012,
					Protocol:      v1.ProtocolTCP,
				},
			}).ContainerPort([]v1.ContainerPort{
				{
					ContainerPort: 8013,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8014,
					Protocol:      v1.ProtocolTCP,
				},
			}).Obj(),
			expected: []*v1.ContainerPort{
				{
					ContainerPort: 8001,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8002,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8003,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8004,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8011,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8012,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8013,
					Protocol:      v1.ProtocolTCP,
				},
				{
					ContainerPort: 8014,
					Protocol:      v1.ProtocolTCP,
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			result := getContainerPorts(test.pod1, test.pod2)
			if diff := cmp.Diff(test.expected, result); diff != "" {
				t.Errorf("container ports: container ports does not match (-want,+got):\n%s", diff)
				return
			}
		})
	}
}
