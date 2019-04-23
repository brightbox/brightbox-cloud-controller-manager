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

REGISTRY ?= brightbox
BUILD ?= $(shell git describe --always --dirty)
VERSION ?= $(shell git describe --always --dirty | egrep -o '^v[0-9]+\.[0-9]+\.[0-9]+')
GOOS ?= linux
ARCH ?= amd64
SRC := $(git ls-files "*.go" | grep -v vendor)
BIN := brightbox-cloud-controller-manager
PKG := github.com/brightbox/${BIN}
LDFLAGS := $(shell KUBE_ROOT="." KUBE_GO_PACKAGE=${PKG} hack/version.sh)

.PHONY: all
all: clean check build

.PHONY: clean
clean:
	GOOS=${GOOS} GOARCH=${ARCH} go clean -i -x ./...

.PHONY: compile
compile: check-headers gofmt ${BIN}
${BIN}:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${ARCH} go build \
	    -ldflags "-s -w ${LDFLAGS}" \
	    -o ${BIN}

.PHONY: version
version:
	@echo ${VERSION}:${BUILD}

.PHONY: gofmt
gofmt:
	@./hack/gofmt.sh ${SRC}

.PHONY: golint
golint:
	golint ${SRC}

.PHONY: govet
govet:
	go vet ${SRC}

.PHONY: check-headers
check-headers: 
	@./hack/verify-boilerplate.sh

.PHONY: check
check: check-headers gofmt govet golint

.PHONY: container
container: test compile
	docker build -t ${REGISTRY}/${BIN}:${VERSION} .

.PHONY: push
push: container
	docker push ${REGISTRY}/${BIN}:${VERSION}

.PHONY: test
test: check-headers gofmt
	go test -v ./...

main.go:
	@./hack/create-main.sh
