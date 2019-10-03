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

package discovery

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "k8s.io/internal-api/apis/core"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EndpointSlice represents a subset of the endpoints that implement a service.
// For a given service there may be multiple EndpointSlice objects, selected by
// labels, which must be joined to produce the full set of endpoints.
type EndpointSlice struct {
	metav1.TypeMeta
	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta
	// addressType specifies the type of address carried by this EndpointSlice.
	// All addresses in this slice must be the same type.
	// +optional
	AddressType *AddressType
	// endpoints is a list of unique endpoints in this slice. Each slice may
	// include a maximum of 1000 endpoints.
	// +listType=atomic
	Endpoints []Endpoint
	// ports specifies the list of network ports exposed by each endpoint in
	// this slice. Each port must have a unique name. When ports is empty, it
	// indicates that there are no defined ports. When a port is defined with a
	// nil port value, it indicates "all ports". Each slice may include a
	// maximum of 100 ports.
	// +optional
	// +listType=atomic
	Ports []EndpointPort
}

// AddressType represents the type of address referred to by an endpoint.
type AddressType string

const (
	// AddressTypeIP represents an IP Address.
	AddressTypeIP = AddressType("IP")
)

// Endpoint represents a single logical "backend" implementing a service.
type Endpoint struct {
	// addresses of this endpoint. The contents of this field are interpreted
	// according to the corresponding EndpointSlice addressType field. This
	// allows for cases like dual-stack (IPv4 and IPv6) networking. Consumers
	// (e.g. kube-proxy) must handle different types of addresses in the context
	// of their own capabilities. This must contain at least one address but no
	// more than 100.
	// +listType=set
	Addresses []string
	// conditions contains information about the current status of the endpoint.
	Conditions EndpointConditions
	// hostname of this endpoint. This field may be used by consumers of
	// endpoints to distinguish endpoints from each other (e.g. in DNS names).
	// Multiple endpoints which use the same hostname should be considered
	// fungible (e.g. multiple A values in DNS). Must pass DNS Label (RFC 1123)
	// validation.
	// +optional
	Hostname *string
	// targetRef is a reference to a Kubernetes object that represents this
	// endpoint.
	// +optional
	TargetRef *api.ObjectReference
	// topology contains arbitrary topology information associated with the
	// endpoint. These key/value pairs must conform with the label format.
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels
	// Topology may include a maximum of 16 key/value pairs. This includes, but
	// is not limited to the following well known keys:
	// * kubernetes.io/hostname: the value indicates the hostname of the node
	//   where the endpoint is located. This should match the corresponding
	//   node label.
	// * topology.kubernetes.io/zone: the value indicates the zone where the
	//   endpoint is located. This should match the corresponding node label.
	// * topology.kubernetes.io/region: the value indicates the region where the
	//   endpoint is located. This should match the corresponding node label.
	// +optional
	Topology map[string]string
}

// EndpointConditions represents the current condition of an endpoint.
type EndpointConditions struct {
	// ready indicates that this endpoint is prepared to receive traffic,
	// according to whatever system is managing the endpoint. A nil value
	// indicates an unknown state. In most cases consumers should interpret this
	// unknown state as ready.
	Ready *bool
}

// EndpointPort represents a Port used by an EndpointSlice.
type EndpointPort struct {
	// The name of this port. All ports in an EndpointSlice must have a unique
	// name. If the EndpointSlice is dervied from a Kubernetes service, this
	// corresponds to the Service.ports[].name.
	// Name must either be an empty string or pass IANA_SVC_NAME validation:
	// * must be no more than 15 characters long
	// * may contain only [-a-z0-9]
	// * must contain at least one letter [a-z]
	// * it must not start or end with a hyphen, nor contain adjacent hyphens
	Name *string
	// The IP protocol for this port.
	// Must be UDP, TCP, or SCTP.
	Protocol *api.Protocol
	// The port number of the endpoint.
	// If this is not specified, ports are not restricted and must be
	// interpreted in the context of the specific consumer.
	Port *int32
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EndpointSliceList represents a list of endpoint slices.
type EndpointSliceList struct {
	metav1.TypeMeta
	// Standard list metadata.
	// +optional
	metav1.ListMeta
	// List of endpoint slices
	// +listType=set
	Items []EndpointSlice
}
