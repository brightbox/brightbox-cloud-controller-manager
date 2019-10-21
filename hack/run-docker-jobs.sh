#!/bin/sh

set -e
first_release=16
last_release=16

launch_job() {
	local release=${1}
	local version=$(git describe --always release-${release} | egrep -o '^v[0-9]+\.[0-9]+\.[0-9]+' | sed 's/^v//')
	local name=cloud-controller-build-$(echo $version | tr . -) 

kubectl apply -f - <<-EOF
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
for word in $(seq ${first_release} ${last_release})
do
	launch_job 1.${word}
done

