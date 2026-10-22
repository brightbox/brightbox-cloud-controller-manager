#!/bin/bash
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

[ $# -eq 1 ] || { echo "Supply new version number" >&2; exit 1; }

m=v0.${1#1.}
go get k8s.io/kubernetes@v$1 k8s.io/api@$m k8s.io/apimachinery@$m k8s.io/cloud-provider@$m \
  k8s.io/component-base@$m k8s.io/controller-manager@$m
go mod tidy
