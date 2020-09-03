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


[ $# -eq 1 ] || { echo "Supply new version number" >&2; exit 1; }

sed -i '/^replace k8s.io/d' go.mod
curl -LSs \
	https://raw.githubusercontent.com/kubernetes/kubernetes/v${1}/go.mod | \
	sed -n "s/\\(k8s.io.*\\) v0.0.0$/replace \\1 => \\1 kubernetes-${1}/p" >> go.mod
go get k8s.io/kubernetes@v$1
go mod tidy
