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

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

kube::log::status "Restoring kubernetes godeps"

kube::util::godep_restored_prepare

if kube::util::godep_restored 2>&1; then
    kube::log::status "Dependencies appear to be current - skipping download"
    exit 0
fi

kube::util::ensure_godep_version

kube::log::status "Downloading dependencies - this might take a while"
GOPATH="${GOPATH}:${KUBE_ROOT}/staging" ${KUBE_GODEP:?} restore "$@"
kube::log::status "Done"
