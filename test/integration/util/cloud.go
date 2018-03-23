/*
Copyright 2018 The Kubernetes Authors.

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

package util

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/golang/glog"

	"golang.org/x/oauth2"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/gce"
)

const (
	// TestProjectID is the project id used for creating NewMockGCECloud
	TestProjectID = "test-project"
	// TestNetworkProjectID is the network project id for creating NewMockGCECloud
	TestNetworkProjectID = "net-test-project"
	// TestRegion is the region for creating NewMockGCECloud
	TestRegion = "test-region"
	// TestZone is the zone for creating NewMockGCECloud
	TestZone = "test-zone"
	// TestNetworkName is the network name for creating NewMockGCECloud
	TestNetworkName = "test-network"
	// TestSubnetworkName is the sub network name for creating NewMockGCECloud
	TestSubnetworkName = "test-sub-network"
	// TestSecondaryRangeName is the secondary range name for creating NewMockGCECloud
	TestSecondaryRangeName = "test-secondary-range"
)

type mockTokenSource struct{}

func (*mockTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  "access",
		TokenType:    "Bearer",
		RefreshToken: "refresh",
		Expiry:       time.Now().Add(1 * time.Hour),
	}, nil
}

// NewMockGCECloud returns a handle to a GCECloud instance that is
// served by a mock http server
func NewMockGCECloud(handler http.Handler) (ShutdownFunc, *gce.GCECloud, error) {
	h := func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("[%s] %s", r.Method, r.RequestURI)
		handler.ServeHTTP(w, r)
	}
	s := httptest.NewServer(http.HandlerFunc(h))
	baseURL := s.URL + "/compute/v1/"
	glog.Infof("Mock server running at: %s", s.URL)

	config := &gce.CloudConfig{
		ApiEndpoint:        baseURL,
		ProjectID:          TestProjectID,
		NetworkProjectID:   TestNetworkProjectID,
		Region:             TestRegion,
		Zone:               TestZone,
		ManagedZones:       []string{TestZone},
		NetworkName:        TestNetworkName,
		SubnetworkName:     TestSubnetworkName,
		SecondaryRangeName: TestSecondaryRangeName,
		NodeTags:           []string{},
		UseMetadataServer:  false,
		TokenSource:        &mockTokenSource{},
	}
	cloud, err := gce.CreateGCECloud(config)

	shutdownFunc := func() {
		glog.Infof("Shutting down mock cloud server")
		s.Close()
	}
	return shutdownFunc, cloud, err
}
