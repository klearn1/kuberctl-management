/*
Copyright 2023 The Kubernetes Authors.

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

package client

// Args are the arguments for the client generator.
type Args struct {
	InputDir string `doc:"The root package under which to search for Go files which request clients to be generated. This must be a local path, not a Go package."`

	// TODO: mirror other arguments from kube::codegen::gen_client function
}
