# Kubernetes::IoK8sApiAutoscalingV2beta1MetricStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**external** | [**IoK8sApiAutoscalingV2beta1ExternalMetricStatus**](IoK8sApiAutoscalingV2beta1ExternalMetricStatus.md) |  | [optional] 
**object** | [**IoK8sApiAutoscalingV2beta1ObjectMetricStatus**](IoK8sApiAutoscalingV2beta1ObjectMetricStatus.md) |  | [optional] 
**pods** | [**IoK8sApiAutoscalingV2beta1PodsMetricStatus**](IoK8sApiAutoscalingV2beta1PodsMetricStatus.md) |  | [optional] 
**resource** | [**IoK8sApiAutoscalingV2beta1ResourceMetricStatus**](IoK8sApiAutoscalingV2beta1ResourceMetricStatus.md) |  | [optional] 
**type** | **String** | type is the type of metric source.  It will be one of \&quot;Object\&quot;, \&quot;Pods\&quot; or \&quot;Resource\&quot;, each corresponds to a matching field in the object. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAutoscalingV2beta1MetricStatus.new(external: null,
                                 object: null,
                                 pods: null,
                                 resource: null,
                                 type: null)
```


