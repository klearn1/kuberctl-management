#!/bin/bash

# Copyright 2015 The Kubernetes Authors All rights reserved.
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

set -o pipefail

# TODO: Switch this example to get the IP address from DNS instead of using the
# kubernetes API.

IP=""
if [[ -n "${KUBERNETES_SERVICE_HOST}" ]]; then
  : ${NAMESPACE:=rethinkdb} # TODO: get namespace from downward API
  # try to pick up first different ip from endpoints
  MYHOST=$(ip addr | grep 'state UP' -A2 | tail -n1 | awk '{print $2}' | cut -f1  -d'/')
  URL="https://${KUBERNETES_SERVICE_HOST}:${KUBERNETES_SERVICE_PORT}/api/v1/namespaces/${NAMESPACE}/endpoints/rethinkdb-driver"
  token=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
  IP=$(curl -s ${URL} --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt --header "Authorization: Bearer ${token}" \
    | jq -s -r --arg h "${MYHOST}" '.[0].subsets | .[].addresses | [ .[].IP ] | map(select(. != $h)) | .[0]') || exit 1
  [[ "${IP}" == null ]] && IP=""
fi

if [[ -n "${IP}" ]]; then
  ENDPOINT="${IP}:29015"
  echo "Join to ${ENDPOINT}"
  exec rethinkdb --bind all  --join ${ENDPOINT}
else
  echo "Start single instance"
  exec rethinkdb --bind all
fi
