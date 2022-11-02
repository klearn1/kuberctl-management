#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script checks whether updating of device plugin API is needed or not. We
# should run `hack/update-generated-device-plugin.sh` if device plugin API is
# out of date.
# Usage: `hack/verify-generated-device-plugin.sh`.

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
ERROR="Device plugin api is out of date. Please run hack/update-generated-device-plugin.sh"
DEVICE_PLUGIN_ALPHA="${KUBE_ROOT}/staging/src/k8s.io/kubelet/pkg/apis/deviceplugin/v1alpha/"
DEVICE_PLUGIN_V1BETA1="${KUBE_ROOT}/staging/src/k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1/"
DEVICE_PLUGIN_V1="${KUBE_ROOT}/staging/src/k8s.io/kubelet/pkg/apis/deviceplugin/v1/"

source "${KUBE_ROOT}/hack/lib/protoc.sh"
kube::golang::setup_env

function cleanup {
	rm -rf "${DEVICE_PLUGIN_ALPHA}/_tmp/"
	rm -rf "${DEVICE_PLUGIN_V1BETA1}/_tmp/"
	rm -rf "${DEVICE_PLUGIN_V1}/_tmp/"
}

trap cleanup EXIT

mkdir -p "${DEVICE_PLUGIN_ALPHA}/_tmp"
cp "${DEVICE_PLUGIN_ALPHA}/api.pb.go" "${DEVICE_PLUGIN_ALPHA}/_tmp/"
mkdir -p "${DEVICE_PLUGIN_V1BETA1}/_tmp"
cp "${DEVICE_PLUGIN_V1BETA1}/api.pb.go" "${DEVICE_PLUGIN_V1BETA1}/_tmp/"
mkdir -p "${DEVICE_PLUGIN_V1}/_tmp"
cp "${DEVICE_PLUGIN_V1}/api.pb.go" "${DEVICE_PLUGIN_V1}/_tmp/"

KUBE_VERBOSE=3 "${KUBE_ROOT}/hack/update-generated-device-plugin.sh"
kube::protoc::diff "${DEVICE_PLUGIN_ALPHA}/api.pb.go" "${DEVICE_PLUGIN_ALPHA}/_tmp/api.pb.go" "${ERROR}"
echo "Generated device plugin alpha api is up to date."
kube::protoc::diff "${DEVICE_PLUGIN_V1BETA1}/api.pb.go" "${DEVICE_PLUGIN_V1BETA1}/_tmp/api.pb.go" "${ERROR}"
echo "Generated device plugin beta api is up to date."
kube::protoc::diff "${DEVICE_PLUGIN_V1}/api.pb.go" "${DEVICE_PLUGIN_V1}/_tmp/api.pb.go" "${ERROR}"
echo "Generated device plugin beta api is up to date."
