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

package garbagecollector

import (
	"fmt"
)

const tprMessage = `If %s is a thirdparty resource (tpr), please note that the garbage collector doesn't support tpr yet. Once tpr is supported, object with ownerReferences referring non-existing tpr objects will be deleted by the garbage collector.`

type restMappingError struct {
	kind    string
	version string
}

func (r *restMappingError) Error() string {
	versionKind := fmt.Sprintf("%s/%s", r.version, r.kind)
	return fmt.Sprintf("unable to get REST mapping for %s.", versionKind)
}

// Message prints more details
func (r *restMappingError) Message() string {
	versionKind := fmt.Sprintf("%s/%s", r.version, r.kind)
	errMsg := fmt.Sprintf("unable to get REST mapping for %s. ", versionKind)
	errMsg += fmt.Sprintf(tprMessage, versionKind)
	errMsg += fmt.Sprintf(" If %s is not a tpr, then you should remove ownerReferences that refer %s objects manually.", versionKind, versionKind)
	return errMsg
}

func newRESTMappingError(kind, version string) *restMappingError {
	return &restMappingError{kind: kind, version: version}
}
