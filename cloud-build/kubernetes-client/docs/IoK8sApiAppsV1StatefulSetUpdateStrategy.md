# Kubernetes::IoK8sApiAppsV1StatefulSetUpdateStrategy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**rolling_update** | [**IoK8sApiAppsV1RollingUpdateStatefulSetStrategy**](IoK8sApiAppsV1RollingUpdateStatefulSetStrategy.md) |  | [optional] 
**type** | **String** | Type indicates the type of the StatefulSetUpdateStrategy. Default is RollingUpdate. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAppsV1StatefulSetUpdateStrategy.new(rolling_update: null,
                                 type: null)
```


