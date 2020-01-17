#!/bin/sh
# Copyright 2019 Brightbox Systems Ltd
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


set -e
first_release=14
last_release=17

create_job_manifest() {
        local release=${1}
	local prefix="cloud-controller-build"
        local version=$(git describe --always release-${release} | sed 's/^v\([0-9]*\.[0-9]*\.[0-9]*\).*$/\1/')
        local name=${prefix}-$(echo $version | tr . -)
	local image="gcr.io/kaniko-project/executor:latest"
        local requestsCpu="1700m"
        local requestsMemory="1Gi"
        local limitsCpu="2"
        local limitsMemory="1890Mi"
        local gitRepo="git://github.com/brightbox/brightbox-cloud-controller-manager.git"
        local dockerTarget="brightbox/brightbox-cloud-controller-manager"
        local secretName="regcred"
	local holdTime=600

cat <<-EOF
---
apiVersion: batch/v1
kind: Job
metadata:
  name: ${name}
  labels:
    build: ${prefix}
spec:
  ttlSecondsAfterFinished: ${holdTime}
  template:
    spec:
      containers:
      - name: ${name}
        image: ${image}
        resources:
          requests:
            memory: ${requestsMemory}
            cpu: ${requestsCpu}
          limits:
            memory: ${limitsMemory}
            cpu: ${limitsCpu}
        args: ["--dockerfile=Dockerfile",
                "--context=${gitRepo}#refs/heads/release-${release}",
                "--destination=${dockerTarget}:${version}"]
        volumeMounts:
          - name: ${secretName}
            mountPath: /root
      restartPolicy: Never
      volumes:
        - name: ${secretName}
          secret:
            secretName: ${secretName}
            items:
              - key: .dockerconfigjson
                path: .docker/config.json
EOF
}

echo "Running jobs on k8s"
for word in $(seq ${last_release} -1 ${first_release})
do
        create_job_manifest 1.${word}
done | kubectl apply -f -
