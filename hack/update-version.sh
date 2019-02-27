#!/bin/bash
# Copyright 2018 Brightbox Systems Ltd
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


[ $# -eq 1 ] || { echo "Supply new version number" >&2; exit 1; }

go get k8s.io/kubernetes@v$1 \
	k8s.io/cloud-provider@kubernetes-$1\
	k8s.io/api@kubernetes-$1\
	k8s.io/apimachinery@kubernetes-$1\
	k8s.io/apiserver@kubernetes-$1\
	k8s.io/apiextensions-apiserver@kubernetes-$1\
	k8s.io/cloud-provider@kubernetes-$1\
	k8s.io/csi-api@kubernetes-$1\
	k8s.io/kube-controller-manager@kubernetes-$1 \
	k8s.io/client-go@kubernetes-$1

