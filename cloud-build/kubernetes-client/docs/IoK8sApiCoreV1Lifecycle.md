# Kubernetes::IoK8sApiCoreV1Lifecycle

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**post_start** | [**IoK8sApiCoreV1Handler**](IoK8sApiCoreV1Handler.md) |  | [optional] 
**pre_stop** | [**IoK8sApiCoreV1Handler**](IoK8sApiCoreV1Handler.md) |  | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1Lifecycle.new(post_start: null,
                                 pre_stop: null)
```


