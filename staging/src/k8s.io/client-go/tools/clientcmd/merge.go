/*
Copyright 2024 The Kubernetes Authors.

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

package clientcmd

import (
	"fmt"
	"reflect"
	"strings"
)

// recursively merges src into dst:
// - non-pointer struct fields are recursively merged
// - maps are shallow merged with src keys taking priority over dst
// - non-zero src fields encountered during recursion that are not maps or structs overwrite and recursion stops
func merge[T any](dst, src *T) error {
	if dst == nil {
		return fmt.Errorf("cannot merge into nil pointer")
	}
	if src == nil {
		return nil
	}
	return mergeValues(nil, reflect.ValueOf(dst).Elem(), reflect.ValueOf(src).Elem())
}

func mergeValues(fieldNames []string, dst, src reflect.Value) error {
	dstType := dst.Type()
	// sanity check types match
	if srcType := src.Type(); dstType != srcType {
		return fmt.Errorf("cannot merge mismatched types (%s, %s) at %s", dstType, srcType, strings.Join(fieldNames, "."))
	}
	// sanity check dst can be set
	if !dst.CanSet() {
		return fmt.Errorf("unsettable value at %s", strings.Join(fieldNames, "."))
	}
	// if src is zero, nothing to do
	if src.IsZero() {
		return nil
	}
	// if dst is zero, just set, don't bother merging
	if dst.IsZero() {
		dst.Set(src)
		return nil
	}

	switch dstType.Kind() {
	case reflect.Struct:
		// recursively merge exported struct fields with src overwriting
		for i := 0; i < dstType.NumField(); i++ {
			if fieldInfo := dstType.Field(i); fieldInfo.IsExported() {
				if err := mergeValues(append(fieldNames, fieldInfo.Name), dst.Field(i), src.Field(i)); err != nil {
					return err
				}
			}
		}
	case reflect.Map:
		// shallow-merge maps with src overwriting
		for _, mapKey := range src.MapKeys() {
			dst.SetMapIndex(mapKey, src.MapIndex(mapKey))
		}
	default:
		// overwrite dst for other types
		dst.Set(src)
	}

	return nil
}
