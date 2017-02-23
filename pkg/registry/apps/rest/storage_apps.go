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

package rest

import (
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/apps"
	appsapiv1alpha1 "k8s.io/kubernetes/pkg/apis/apps/v1alpha1"
	appsapiv1beta1 "k8s.io/kubernetes/pkg/apis/apps/v1beta1"
	statefulsetstore "k8s.io/kubernetes/pkg/registry/apps/petset/storage"
	pipstore "k8s.io/kubernetes/pkg/registry/apps/podinjectionpolicy/storage"
)

type RESTStorageProvider struct{}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(apps.GroupName, api.Registry, api.Scheme, api.ParameterCodec, api.Codecs)

	if apiResourceConfigSource.AnyResourcesForVersionEnabled(appsapiv1beta1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[appsapiv1beta1.SchemeGroupVersion.Version] = p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter)
		apiGroupInfo.GroupMeta.GroupVersion = appsapiv1beta1.SchemeGroupVersion
	}

	if apiResourceConfigSource.AnyResourcesForVersionEnabled(appsapiv1alpha1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[appsapiv1alpha1.SchemeGroupVersion.Version] = p.v1alpha1Storage(apiResourceConfigSource, restOptionsGetter)
		apiGroupInfo.GroupMeta.GroupVersion = appsapiv1alpha1.SchemeGroupVersion
	}

	return apiGroupInfo, true
}

func (p RESTStorageProvider) v1beta1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) map[string]rest.Storage {
	version := appsapiv1beta1.SchemeGroupVersion

	storage := map[string]rest.Storage{}
	if apiResourceConfigSource.ResourceEnabled(version.WithResource("statefulsets")) {
		statefulsetStorage, statefulsetStatusStorage := statefulsetstore.NewREST(restOptionsGetter)
		storage["statefulsets"] = statefulsetStorage
		storage["statefulsets/status"] = statefulsetStatusStorage
	}
	return storage
}

func (p RESTStorageProvider) v1alpha1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) map[string]rest.Storage {
	version := appsapiv1beta1.SchemeGroupVersion

	storage := map[string]rest.Storage{}
	if apiResourceConfigSource.ResourceEnabled(version.WithResource("podinjectionpolicies")) {
		podInjectionPolicyStorage := pipstore.NewREST(restOptionsGetter)
		storage["podinjectionpolicies"] = podInjectionPolicyStorage
	}
	return storage
}

func (p RESTStorageProvider) GroupName() string {
	return apps.GroupName
}
