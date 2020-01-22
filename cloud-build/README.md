# Ruby client

Generating a Kuberenetes OpenAPI ruby client

Get the latest api spec from your cluster

```
$ kubectl get --raw /openapi/v2 > kubeapi.json
```

Run the OpenAPI generator to recreate the client

```
$ docker run --rm -v ${PWD}:/local -v /tmp:/tmp -u $(id -u ${USER}):$(id -g ${USER}) openapitools/openapi-generator-cli generate -i /local/kubeapi.json  -c /local/config.json -g ruby -o /local/kubernetes-client
```
