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
        local version=$(git describe --always release-${release} | sed 's/^v\([0-9]*\.[0-9]*\.[0-9]*\).*$/\1/')
        local name=cloud-controller-build-$(echo $version | tr . -)

cat <<-EOF
---
apiVersion: batch/v1
kind: Job
metadata:
  name: ${name}
  labels:
    build: cloud-controller-build
spec:
  ttlSecondsAfterFinished: 600
  template:
    spec:
      containers:
      - name: ${name}
        image: gcr.io/kaniko-project/executor:latest
        resources:
          requests:
            memory: 1Gi
            cpu: 500m
          limits:
            memory: 1890Mi
            cpu: 2
        args: ["--dockerfile=Dockerfile",
                "--context=git://github.com/brightbox/brightbox-cloud-controller-manager.git#refs/heads/release-${release}",
                "--destination=brightbox/brightbox-cloud-controller-manager:${version}"]
        volumeMounts:
          - name: kaniko-secret
            mountPath: /root
      restartPolicy: Never
      volumes:
        - name: kaniko-secret
          secret:
            secretName: regcred
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

