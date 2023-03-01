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

FROM golang:1.19 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN GOCACHE=/tmp go mod download

COPY . .

RUN GOCACHE=/tmp make brightbox-cloud-controller-manager

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/brightbox-cloud-controller-manager /bin/

ENTRYPOINT ["/bin/brightbox-cloud-controller-manager"]
