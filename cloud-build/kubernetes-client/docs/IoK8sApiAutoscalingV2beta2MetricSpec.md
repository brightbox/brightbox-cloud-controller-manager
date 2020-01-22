# Kubernetes::IoK8sApiAutoscalingV2beta2MetricSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**external** | [**IoK8sApiAutoscalingV2beta2ExternalMetricSource**](IoK8sApiAutoscalingV2beta2ExternalMetricSource.md) |  | [optional] 
**object** | [**IoK8sApiAutoscalingV2beta2ObjectMetricSource**](IoK8sApiAutoscalingV2beta2ObjectMetricSource.md) |  | [optional] 
**pods** | [**IoK8sApiAutoscalingV2beta2PodsMetricSource**](IoK8sApiAutoscalingV2beta2PodsMetricSource.md) |  | [optional] 
**resource** | [**IoK8sApiAutoscalingV2beta2ResourceMetricSource**](IoK8sApiAutoscalingV2beta2ResourceMetricSource.md) |  | [optional] 
**type** | **String** | type is the type of metric source.  It should be one of \&quot;Object\&quot;, \&quot;Pods\&quot; or \&quot;Resource\&quot;, each mapping to a matching field in the object. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAutoscalingV2beta2MetricSpec.new(external: null,
                                 object: null,
                                 pods: null,
                                 resource: null,
                                 type: null)
```


