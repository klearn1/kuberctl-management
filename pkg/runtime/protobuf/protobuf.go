/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package protobuf

import (
	"encoding/hex"
	"fmt"
	"io"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"

	"k8s.io/kubernetes/pkg/runtime"
)

// NewCodec
func NewCodec(version string, creater runtime.ObjectCreater, typer runtime.ObjectTyper, convertor runtime.ObjectConvertor) runtime.Codec {
	return &codec{
		version:   version,
		creater:   creater,
		typer:     typer,
		convertor: convertor,
	}
}

// codec decodes protobuf objects
type codec struct {
	version   string
	creater   runtime.ObjectCreater
	typer     runtime.ObjectTyper
	convertor runtime.ObjectConvertor
}

var _ runtime.Codec = codec{}

func (c codec) Decode(data []byte) (runtime.Object, error) {
	unknown := &runtime.Unknown{}
	if err := proto.Unmarshal(data, unknown); err != nil {
		return nil, err
	}
	obj, err := c.creater.New(unknown.APIVersion, unknown.Kind)
	if err != nil {
		return nil, err
	}
	pobj, ok := obj.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("runtime object is not a proto.Message: %v", reflect.TypeOf(obj))
	}
	glog.Infof("hex for %s %s: %s", unknown.APIVersion, unknown.Kind, hex.Dump(unknown.RawJSON))
	if err := proto.Unmarshal(unknown.RawJSON, pobj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (c codec) DecodeToVersion(data []byte, version string) (runtime.Object, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (c codec) DecodeInto(data []byte, obj runtime.Object) error {
	pobj, ok := obj.(proto.Message)
	if !ok {
		return fmt.Errorf("runtime object is not a proto.Message: %v", reflect.TypeOf(obj))
	}
	return proto.Unmarshal(data, pobj)
}

func (c codec) DecodeIntoWithSpecifiedVersionKind(data []byte, obj runtime.Object, kind, version string) error {
	return fmt.Errorf("unimplemented")
}

func (c codec) Encode(obj runtime.Object) (data []byte, err error) {
	version, kind, err := c.typer.ObjectVersionAndKind(obj)
	if err != nil {
		return nil, err
	}
	if len(version) == 0 {
		converted, err := c.convertor.ConvertToVersion(obj, c.version)
		if err != nil {
			return nil, err
		}
		obj = converted
	}
	m, ok := obj.(proto.Marshaler)
	if !ok {
		return nil, fmt.Errorf("object %v (kind: %s in version: %s) does not implement ProtoBuf marshalling", reflect.TypeOf(obj), kind, c.version)
	}
	b, err := m.Marshal()
	if err != nil {
		return nil, err
	}
	glog.Infof("marshaled %#v\ninto\n%s", obj, hex.Dump(b))
	return (&runtime.Unknown{
		TypeMeta: runtime.TypeMeta{
			Kind:       kind,
			APIVersion: version,
		},
		RawJSON: b,
	}).Marshal()
}
func (c codec) EncodeToStream(obj runtime.Object, stream io.Writer) error {
	return fmt.Errorf("unimplemented")
}
