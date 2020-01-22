# Kubernetes::IoK8sApiDiscoveryV1beta1EndpointConditions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ready** | **Boolean** | ready indicates that this endpoint is prepared to receive traffic, according to whatever system is managing the endpoint. A nil value indicates an unknown state. In most cases consumers should interpret this unknown state as ready. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiDiscoveryV1beta1EndpointConditions.new(ready: null)
```


