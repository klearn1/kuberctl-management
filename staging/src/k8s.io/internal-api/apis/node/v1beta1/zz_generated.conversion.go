// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by conversion-gen. DO NOT EDIT.

package v1beta1

import (
	unsafe "unsafe"

	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/node/v1beta1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	core "k8s.io/internal-api/apis/core"
	node "k8s.io/internal-api/apis/node"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1beta1.Overhead)(nil), (*node.Overhead)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Overhead_To_node_Overhead(a.(*v1beta1.Overhead), b.(*node.Overhead), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*node.Overhead)(nil), (*v1beta1.Overhead)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_node_Overhead_To_v1beta1_Overhead(a.(*node.Overhead), b.(*v1beta1.Overhead), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RuntimeClass)(nil), (*node.RuntimeClass)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RuntimeClass_To_node_RuntimeClass(a.(*v1beta1.RuntimeClass), b.(*node.RuntimeClass), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*node.RuntimeClass)(nil), (*v1beta1.RuntimeClass)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_node_RuntimeClass_To_v1beta1_RuntimeClass(a.(*node.RuntimeClass), b.(*v1beta1.RuntimeClass), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RuntimeClassList)(nil), (*node.RuntimeClassList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RuntimeClassList_To_node_RuntimeClassList(a.(*v1beta1.RuntimeClassList), b.(*node.RuntimeClassList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*node.RuntimeClassList)(nil), (*v1beta1.RuntimeClassList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_node_RuntimeClassList_To_v1beta1_RuntimeClassList(a.(*node.RuntimeClassList), b.(*v1beta1.RuntimeClassList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.Scheduling)(nil), (*node.Scheduling)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Scheduling_To_node_Scheduling(a.(*v1beta1.Scheduling), b.(*node.Scheduling), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*node.Scheduling)(nil), (*v1beta1.Scheduling)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_node_Scheduling_To_v1beta1_Scheduling(a.(*node.Scheduling), b.(*v1beta1.Scheduling), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1beta1_Overhead_To_node_Overhead(in *v1beta1.Overhead, out *node.Overhead, s conversion.Scope) error {
	out.PodFixed = *(*core.ResourceList)(unsafe.Pointer(&in.PodFixed))
	return nil
}

// Convert_v1beta1_Overhead_To_node_Overhead is an autogenerated conversion function.
func Convert_v1beta1_Overhead_To_node_Overhead(in *v1beta1.Overhead, out *node.Overhead, s conversion.Scope) error {
	return autoConvert_v1beta1_Overhead_To_node_Overhead(in, out, s)
}

func autoConvert_node_Overhead_To_v1beta1_Overhead(in *node.Overhead, out *v1beta1.Overhead, s conversion.Scope) error {
	out.PodFixed = *(*v1.ResourceList)(unsafe.Pointer(&in.PodFixed))
	return nil
}

// Convert_node_Overhead_To_v1beta1_Overhead is an autogenerated conversion function.
func Convert_node_Overhead_To_v1beta1_Overhead(in *node.Overhead, out *v1beta1.Overhead, s conversion.Scope) error {
	return autoConvert_node_Overhead_To_v1beta1_Overhead(in, out, s)
}

func autoConvert_v1beta1_RuntimeClass_To_node_RuntimeClass(in *v1beta1.RuntimeClass, out *node.RuntimeClass, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.Handler = in.Handler
	out.Overhead = (*node.Overhead)(unsafe.Pointer(in.Overhead))
	out.Scheduling = (*node.Scheduling)(unsafe.Pointer(in.Scheduling))
	return nil
}

// Convert_v1beta1_RuntimeClass_To_node_RuntimeClass is an autogenerated conversion function.
func Convert_v1beta1_RuntimeClass_To_node_RuntimeClass(in *v1beta1.RuntimeClass, out *node.RuntimeClass, s conversion.Scope) error {
	return autoConvert_v1beta1_RuntimeClass_To_node_RuntimeClass(in, out, s)
}

func autoConvert_node_RuntimeClass_To_v1beta1_RuntimeClass(in *node.RuntimeClass, out *v1beta1.RuntimeClass, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.Handler = in.Handler
	out.Overhead = (*v1beta1.Overhead)(unsafe.Pointer(in.Overhead))
	out.Scheduling = (*v1beta1.Scheduling)(unsafe.Pointer(in.Scheduling))
	return nil
}

// Convert_node_RuntimeClass_To_v1beta1_RuntimeClass is an autogenerated conversion function.
func Convert_node_RuntimeClass_To_v1beta1_RuntimeClass(in *node.RuntimeClass, out *v1beta1.RuntimeClass, s conversion.Scope) error {
	return autoConvert_node_RuntimeClass_To_v1beta1_RuntimeClass(in, out, s)
}

func autoConvert_v1beta1_RuntimeClassList_To_node_RuntimeClassList(in *v1beta1.RuntimeClassList, out *node.RuntimeClassList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]node.RuntimeClass)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1beta1_RuntimeClassList_To_node_RuntimeClassList is an autogenerated conversion function.
func Convert_v1beta1_RuntimeClassList_To_node_RuntimeClassList(in *v1beta1.RuntimeClassList, out *node.RuntimeClassList, s conversion.Scope) error {
	return autoConvert_v1beta1_RuntimeClassList_To_node_RuntimeClassList(in, out, s)
}

func autoConvert_node_RuntimeClassList_To_v1beta1_RuntimeClassList(in *node.RuntimeClassList, out *v1beta1.RuntimeClassList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]v1beta1.RuntimeClass)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_node_RuntimeClassList_To_v1beta1_RuntimeClassList is an autogenerated conversion function.
func Convert_node_RuntimeClassList_To_v1beta1_RuntimeClassList(in *node.RuntimeClassList, out *v1beta1.RuntimeClassList, s conversion.Scope) error {
	return autoConvert_node_RuntimeClassList_To_v1beta1_RuntimeClassList(in, out, s)
}

func autoConvert_v1beta1_Scheduling_To_node_Scheduling(in *v1beta1.Scheduling, out *node.Scheduling, s conversion.Scope) error {
	out.NodeSelector = *(*map[string]string)(unsafe.Pointer(&in.NodeSelector))
	out.Tolerations = *(*[]core.Toleration)(unsafe.Pointer(&in.Tolerations))
	return nil
}

// Convert_v1beta1_Scheduling_To_node_Scheduling is an autogenerated conversion function.
func Convert_v1beta1_Scheduling_To_node_Scheduling(in *v1beta1.Scheduling, out *node.Scheduling, s conversion.Scope) error {
	return autoConvert_v1beta1_Scheduling_To_node_Scheduling(in, out, s)
}

func autoConvert_node_Scheduling_To_v1beta1_Scheduling(in *node.Scheduling, out *v1beta1.Scheduling, s conversion.Scope) error {
	out.NodeSelector = *(*map[string]string)(unsafe.Pointer(&in.NodeSelector))
	out.Tolerations = *(*[]v1.Toleration)(unsafe.Pointer(&in.Tolerations))
	return nil
}

// Convert_node_Scheduling_To_v1beta1_Scheduling is an autogenerated conversion function.
func Convert_node_Scheduling_To_v1beta1_Scheduling(in *node.Scheduling, out *v1beta1.Scheduling, s conversion.Scope) error {
	return autoConvert_node_Scheduling_To_v1beta1_Scheduling(in, out, s)
}
