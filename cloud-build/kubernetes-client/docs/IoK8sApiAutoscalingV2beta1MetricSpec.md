# Kubernetes::IoK8sApiAutoscalingV2beta1MetricSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**external** | [**IoK8sApiAutoscalingV2beta1ExternalMetricSource**](IoK8sApiAutoscalingV2beta1ExternalMetricSource.md) |  | [optional] 
**object** | [**IoK8sApiAutoscalingV2beta1ObjectMetricSource**](IoK8sApiAutoscalingV2beta1ObjectMetricSource.md) |  | [optional] 
**pods** | [**IoK8sApiAutoscalingV2beta1PodsMetricSource**](IoK8sApiAutoscalingV2beta1PodsMetricSource.md) |  | [optional] 
**resource** | [**IoK8sApiAutoscalingV2beta1ResourceMetricSource**](IoK8sApiAutoscalingV2beta1ResourceMetricSource.md) |  | [optional] 
**type** | **String** | type is the type of metric source.  It should be one of \&quot;Object\&quot;, \&quot;Pods\&quot; or \&quot;Resource\&quot;, each mapping to a matching field in the object. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAutoscalingV2beta1MetricSpec.new(external: null,
                                 object: null,
                                 pods: null,
                                 resource: null,
                                 type: null)
```


