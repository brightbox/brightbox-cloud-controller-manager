# Copyright 2020 Brightbox Systems Ltd
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

vendor: kubernetes-client
	bundle config set --local clean 'true'
	bundle config set --local without 'development test'
	bundle config set --local deployment 'true'
	bundle config set --local path 'vendor/bundle'
	bundle install

kubernetes-client: kubeapi.json
	docker run --rm -v $(CURDIR):/local \
		-u $(shell id -u $(USER)):$(shell id -g $(USER)) \
		openapitools/openapi-generator-cli generate \
		-g ruby \
		-p "moduleName=Kubernetes" \
		-i /local/$< \
		-o /local/$@

kubeapi.json:
	kubectl get --raw /openapi/v2 > $@

.PHONY: clean
clean:
	rm -rf kubernetes-client kubeapi.json vendor
