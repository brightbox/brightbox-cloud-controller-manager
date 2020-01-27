vendor: kubernetes-client
	bundle install --clean --without development test --deployment

kubernetes-client: kubeapi.json
	docker run --rm -v $(CURDIR):/local \
		-u $(shell id -u $(USER)):$(shell id -g $(USER)) \
		openapitools/openapi-generator-cli generate \
		-c /local/config.json \
		-i /local/$< \
		-g ruby \
		-o /local/$@

kubeapi.json:
	kubectl get --raw /openapi/v2 > $@

.PHONY: clean
clean:
	rm -rf kubernetes-client kubeapi.json vendor
