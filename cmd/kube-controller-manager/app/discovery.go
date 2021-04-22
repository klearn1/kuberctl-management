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

// Package app implements a server that runs a set of active
// components.  This includes replication controllers, service endpoints and
// nodes.
//
package app

import (
	"context"

	"k8s.io/controller-manager/controller"
	endpointslicecontroller "k8s.io/kubernetes/pkg/controller/endpointslice"
	endpointslicemirroringcontroller "k8s.io/kubernetes/pkg/controller/endpointslicemirroring"
)

func startEndpointSliceController(ctx context.Context, controllerContext ControllerContext) (controller.Interface, bool, error) {
	go endpointslicecontroller.NewController(
		controllerContext.InformerFactory.Core().V1().Pods(),
		controllerContext.InformerFactory.Core().V1().Services(),
		controllerContext.InformerFactory.Core().V1().Nodes(),
		controllerContext.InformerFactory.Discovery().V1().EndpointSlices(),
		controllerContext.ComponentConfig.EndpointSliceController.MaxEndpointsPerSlice,
		controllerContext.ClientBuilder.ClientOrDie("endpointslice-controller"),
		controllerContext.ComponentConfig.EndpointSliceController.EndpointUpdatesBatchPeriod.Duration,
	).Run(int(controllerContext.ComponentConfig.EndpointSliceController.ConcurrentServiceEndpointSyncs), controllerContext.Stop)
	return nil, true, nil
}

func startEndpointSliceMirroringController(ctx context.Context, controllerContext ControllerContext) (controller.Interface, bool, error) {
	go endpointslicemirroringcontroller.NewController(
		controllerContext.InformerFactory.Core().V1().Endpoints(),
		controllerContext.InformerFactory.Discovery().V1().EndpointSlices(),
		controllerContext.InformerFactory.Core().V1().Services(),
		controllerContext.ComponentConfig.EndpointSliceMirroringController.MirroringMaxEndpointsPerSubset,
		controllerContext.ClientBuilder.ClientOrDie("endpointslicemirroring-controller"),
		controllerContext.ComponentConfig.EndpointSliceMirroringController.MirroringEndpointUpdatesBatchPeriod.Duration,
	).Run(int(controllerContext.ComponentConfig.EndpointSliceMirroringController.MirroringConcurrentServiceEndpointSyncs), controllerContext.Stop)
	return nil, true, nil
}
