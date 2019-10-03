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

package v1

import (
	v1 "k8s.io/api/authentication/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	authentication "k8s.io/internal-api/apis/authentication"
)

func addConversionFuncs(scheme *runtime.Scheme) error {
	// Add non-generated conversion functions
	return scheme.AddConversionFuncs()
}

// Convert_v1_UserInfo_To_authentication_UserInfo is an autogenerated conversion function.
func Convert_v1_UserInfo_To_authentication_UserInfo(in *v1.UserInfo, out *authentication.UserInfo, s conversion.Scope) error {
	return autoConvert_v1_UserInfo_To_authentication_UserInfo(in, out, s)
}

// Convert_authentication_UserInfo_To_v1_UserInfo is an autogenerated conversion function.
func Convert_authentication_UserInfo_To_v1_UserInfo(in *authentication.UserInfo, out *v1.UserInfo, s conversion.Scope) error {
	return autoConvert_authentication_UserInfo_To_v1_UserInfo(in, out, s)
}
