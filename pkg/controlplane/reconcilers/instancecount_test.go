/*
Copyright 2014 The Kubernetes Authors.

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

package reconcilers

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	netutils "k8s.io/utils/net"
)

func TestMasterCountEndpointReconciler(t *testing.T) {
	reconcileTests := []struct {
		testName          string
		serviceName       string
		ip                string
		endpointPorts     []corev1.EndpointPort
		additionalMasters int
		initialState      []runtime.Object
		expectUpdate      []runtime.Object
		expectCreate      []runtime.Object
	}{
		{
			testName:      "no existing endpoints",
			serviceName:   "foo",
			ip:            "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			initialState:  nil,
			expectCreate:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:      "existing endpoints satisfy",
			serviceName:   "foo",
			ip:            "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			initialState:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:      "existing endpoints satisfy but too many",
			serviceName:   "foo",
			ip:            "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			initialState:  makeEndpointsArray("foo", []string{"1.2.3.4", "4.3.2.1"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectUpdate:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:          "existing endpoints satisfy but too many + extra masters",
			serviceName:       "foo",
			ip:                "1.2.3.4",
			endpointPorts:     []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			additionalMasters: 3,
			initialState:      makeEndpointsArray("foo", []string{"1.2.3.4", "4.3.2.1", "4.3.2.2", "4.3.2.3", "4.3.2.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectUpdate:      makeEndpointsArray("foo", []string{"1.2.3.4", "4.3.2.2", "4.3.2.3", "4.3.2.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:          "existing endpoints satisfy but too many + extra masters + delete first",
			serviceName:       "foo",
			ip:                "4.3.2.4",
			endpointPorts:     []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			additionalMasters: 3,
			initialState:      makeEndpointsArray("foo", []string{"1.2.3.4", "4.3.2.1", "4.3.2.2", "4.3.2.3", "4.3.2.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectUpdate:      makeEndpointsArray("foo", []string{"4.3.2.1", "4.3.2.2", "4.3.2.3", "4.3.2.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:          "existing endpoints satisfy and endpoint addresses length less than master count",
			serviceName:       "foo",
			ip:                "4.3.2.2",
			endpointPorts:     []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			additionalMasters: 3,
			initialState:      makeEndpointsArray("foo", []string{"4.3.2.1", "4.3.2.2"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectUpdate:      nil,
		},
		{
			testName:          "existing endpoints current IP missing and address length less than master count",
			serviceName:       "foo",
			ip:                "4.3.2.2",
			endpointPorts:     []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			additionalMasters: 3,
			initialState:      makeEndpointsArray("foo", []string{"4.3.2.1"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectUpdate:      makeEndpointsArray("foo", []string{"4.3.2.1", "4.3.2.2"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:      "existing endpoints wrong name",
			serviceName:   "foo",
			ip:            "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			initialState:  makeEndpointsArray("bar", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectCreate:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:      "existing endpoints wrong IP",
			serviceName:   "foo",
			ip:            "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			initialState:  makeEndpointsArray("foo", []string{"4.3.2.1"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectUpdate:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:      "existing endpoints wrong port",
			serviceName:   "foo",
			ip:            "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			initialState:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 9090, Protocol: "TCP"}}),
			expectUpdate:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:      "existing endpoints wrong protocol",
			serviceName:   "foo",
			ip:            "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			initialState:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "UDP"}}),
			expectUpdate:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:      "existing endpoints wrong port name",
			serviceName:   "foo",
			ip:            "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{{Name: "baz", Port: 8080, Protocol: "TCP"}},
			initialState:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectUpdate:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "baz", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:    "existing endpoints extra service ports satisfy",
			serviceName: "foo",
			ip:          "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{
				{Name: "foo", Port: 8080, Protocol: "TCP"},
				{Name: "bar", Port: 1000, Protocol: "TCP"},
				{Name: "baz", Port: 1010, Protocol: "TCP"},
			},
			initialState: makeEndpointsArray("foo", []string{"1.2.3.4"},
				[]corev1.EndpointPort{
					{Name: "foo", Port: 8080, Protocol: "TCP"},
					{Name: "bar", Port: 1000, Protocol: "TCP"},
					{Name: "baz", Port: 1010, Protocol: "TCP"},
				},
			),
		},
		{
			testName:    "existing endpoints extra service ports missing port",
			serviceName: "foo",
			ip:          "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{
				{Name: "foo", Port: 8080, Protocol: "TCP"},
				{Name: "bar", Port: 1000, Protocol: "TCP"},
			},
			initialState: makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectUpdate: makeEndpointsArray("foo", []string{"1.2.3.4"},
				[]corev1.EndpointPort{
					{Name: "foo", Port: 8080, Protocol: "TCP"},
					{Name: "bar", Port: 1000, Protocol: "TCP"},
				},
			),
		},
		{
			testName:      "no existing sctp endpoints",
			serviceName:   "boo",
			ip:            "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{{Name: "boo", Port: 7777, Protocol: "SCTP"}},
			initialState:  nil,
			expectCreate:  makeEndpointsArray("boo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "boo", Port: 7777, Protocol: "SCTP"}}),
		},
	}
	for _, test := range reconcileTests {
		t.Run(test.testName, func(t *testing.T) {
			fakeClient := fake.NewSimpleClientset(test.initialState...)
			epAdapter := NewEndpointsAdapter(fakeClient.CoreV1(), nil)
			reconciler := NewMasterCountEndpointReconciler(test.additionalMasters+1, epAdapter)
			err := reconciler.ReconcileEndpoints(test.serviceName, netutils.ParseIPSloppy(test.ip), test.endpointPorts, true)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			updates := []core.UpdateAction{}
			for _, action := range fakeClient.Actions() {
				if action.GetVerb() != "update" {
					continue
				}
				updates = append(updates, action.(core.UpdateAction))
			}
			if test.expectUpdate != nil {
				if len(updates) != 1 {
					t.Errorf("unexpected updates: %v", updates)
				} else if e, a := test.expectUpdate[0], updates[0].GetObject(); !reflect.DeepEqual(e, a) {
					t.Errorf("expected update:\n%#v\ngot:\n%#v\n", e, a)
				}
			}
			if test.expectUpdate == nil && len(updates) > 0 {
				t.Errorf("no update expected, yet saw: %v", updates)
			}

			creates := []core.CreateAction{}
			for _, action := range fakeClient.Actions() {
				if action.GetVerb() != "create" {
					continue
				}
				creates = append(creates, action.(core.CreateAction))
			}
			if test.expectCreate != nil {
				if len(creates) != 1 {
					t.Errorf("unexpected creates: %v", creates)
				} else if e, a := test.expectCreate[0], creates[0].GetObject(); !reflect.DeepEqual(e, a) {
					t.Errorf("expected create:\n%#v\ngot:\n%#v\n", e, a)
				}
			}
			if test.expectCreate == nil && len(creates) > 0 {
				t.Errorf("no create expected, yet saw: %v", creates)
			}
		})
	}

	nonReconcileTests := []struct {
		testName          string
		serviceName       string
		ip                string
		endpointPorts     []corev1.EndpointPort
		additionalMasters int
		initialState      []runtime.Object
		expectUpdate      []runtime.Object
		expectCreate      []runtime.Object
	}{
		{
			testName:    "existing endpoints extra service ports missing port no update",
			serviceName: "foo",
			ip:          "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{
				{Name: "foo", Port: 8080, Protocol: "TCP"},
				{Name: "bar", Port: 1000, Protocol: "TCP"},
			},
			initialState: makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectUpdate: nil,
		},
		{
			testName:    "existing endpoints extra service ports, wrong ports, wrong IP",
			serviceName: "foo",
			ip:          "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{
				{Name: "foo", Port: 8080, Protocol: "TCP"},
				{Name: "bar", Port: 1000, Protocol: "TCP"},
			},
			initialState: makeEndpointsArray("foo", []string{"4.3.2.1"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
			expectUpdate: makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
		{
			testName:      "no existing endpoints",
			serviceName:   "foo",
			ip:            "1.2.3.4",
			endpointPorts: []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}},
			initialState:  nil,
			expectCreate:  makeEndpointsArray("foo", []string{"1.2.3.4"}, []corev1.EndpointPort{{Name: "foo", Port: 8080, Protocol: "TCP"}}),
		},
	}
	for _, test := range nonReconcileTests {
		t.Run(test.testName, func(t *testing.T) {
			fakeClient := fake.NewSimpleClientset(test.initialState...)
			epAdapter := NewEndpointsAdapter(fakeClient.CoreV1(), nil)
			reconciler := NewMasterCountEndpointReconciler(test.additionalMasters+1, epAdapter)
			err := reconciler.ReconcileEndpoints(test.serviceName, netutils.ParseIPSloppy(test.ip), test.endpointPorts, false)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			updates := []core.UpdateAction{}
			for _, action := range fakeClient.Actions() {
				if action.GetVerb() != "update" {
					continue
				}
				updates = append(updates, action.(core.UpdateAction))
			}
			if test.expectUpdate != nil {
				if len(updates) != 1 {
					t.Errorf("unexpected updates: %v", updates)
				} else if e, a := test.expectUpdate[0], updates[0].GetObject(); !reflect.DeepEqual(e, a) {
					t.Errorf("expected update:\n%#v\ngot:\n%#v\n", e, a)
				}
			}
			if test.expectUpdate == nil && len(updates) > 0 {
				t.Errorf("no update expected, yet saw: %v", updates)
			}

			creates := []core.CreateAction{}
			for _, action := range fakeClient.Actions() {
				if action.GetVerb() != "create" {
					continue
				}
				creates = append(creates, action.(core.CreateAction))
			}
			if test.expectCreate != nil {
				if len(creates) != 1 {
					t.Errorf("unexpected creates: %v", creates)
				} else if e, a := test.expectCreate[0], creates[0].GetObject(); !reflect.DeepEqual(e, a) {
					t.Errorf("expected create:\n%#v\ngot:\n%#v\n", e, a)
				}
			}
			if test.expectCreate == nil && len(creates) > 0 {
				t.Errorf("no create expected, yet saw: %v", creates)
			}
		})
	}

}

func TestEmptySubsets(t *testing.T) {
	endpoints := makeEndpoints("foo", nil, nil)
	fakeClient := fake.NewSimpleClientset(endpoints)
	epAdapter := NewEndpointsAdapter(fakeClient.CoreV1(), nil)
	reconciler := NewMasterCountEndpointReconciler(1, epAdapter)
	endpointPorts := []corev1.EndpointPort{
		{Name: "foo", Port: 8080, Protocol: "TCP"},
	}
	err := reconciler.RemoveEndpoints("foo", netutils.ParseIPSloppy("1.2.3.4"), endpointPorts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
