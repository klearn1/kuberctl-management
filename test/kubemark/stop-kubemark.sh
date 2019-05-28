#!/usr/bin/env bash

# Copyright 2015 The Kubernetes Authors.
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

# Script that destroys Kubemark cluster and deletes all master resources.
set -x 
KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
# shellcheck source=./bootstrap-kubemark.sh
source "${KUBE_ROOT}/test/kubemark/bootstrap-kubemark.sh"

if [[ -f "${KUBE_ROOT}/test/kubemark/${CLOUD_PROVIDER}/shutdown.sh" ]] ; then
# shellcheck disable=SC1090
  source "${KUBE_ROOT}/test/kubemark/${CLOUD_PROVIDER}/shutdown.sh" 
fi
# shellcheck source=../../cluster/kubemark/util.sh
source "${KUBE_ROOT}/cluster/kubemark/util.sh"

KUBECTL="${KUBE_ROOT}/cluster/kubectl.sh"
KUBEMARK_DIRECTORY="${KUBE_ROOT}/test/kubemark"
RESOURCE_DIRECTORY="${KUBEMARK_DIRECTORY}/resources"

detect-project &> /dev/null

"${KUBECTL}" delete -f "${RESOURCE_DIRECTORY}/addons" &> /dev/null || true
"${KUBECTL}" delete -f "${RESOURCE_DIRECTORY}/hollow-node.yaml" &> /dev/null || true
"${KUBECTL}" delete -f "${RESOURCE_DIRECTORY}/kubemark-ns.json" &> /dev/null || true

rm -rf "${RESOURCE_DIRECTORY}/addons" \
	"${RESOURCE_DIRECTORY}/kubeconfig.kubemark" \
	"${RESOURCE_DIRECTORY}/hollow-node.yaml" \
	"${RESOURCE_DIRECTORY}/kubemark-master-env.sh"  &> /dev/null || true

delete-master-instance-and-resources
set +x
