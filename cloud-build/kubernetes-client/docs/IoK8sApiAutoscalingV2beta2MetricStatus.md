# Kubernetes::IoK8sApiAutoscalingV2beta2MetricStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**external** | [**IoK8sApiAutoscalingV2beta2ExternalMetricStatus**](IoK8sApiAutoscalingV2beta2ExternalMetricStatus.md) |  | [optional] 
**object** | [**IoK8sApiAutoscalingV2beta2ObjectMetricStatus**](IoK8sApiAutoscalingV2beta2ObjectMetricStatus.md) |  | [optional] 
**pods** | [**IoK8sApiAutoscalingV2beta2PodsMetricStatus**](IoK8sApiAutoscalingV2beta2PodsMetricStatus.md) |  | [optional] 
**resource** | [**IoK8sApiAutoscalingV2beta2ResourceMetricStatus**](IoK8sApiAutoscalingV2beta2ResourceMetricStatus.md) |  | [optional] 
**type** | **String** | type is the type of metric source.  It will be one of \&quot;Object\&quot;, \&quot;Pods\&quot; or \&quot;Resource\&quot;, each corresponds to a matching field in the object. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAutoscalingV2beta2MetricStatus.new(external: null,
                                 object: null,
                                 pods: null,
                                 resource: null,
                                 type: null)
```


