# Kubernetes::IoK8sApiAutoscalingV2beta2ObjectMetricStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**current** | [**IoK8sApiAutoscalingV2beta2MetricValueStatus**](IoK8sApiAutoscalingV2beta2MetricValueStatus.md) |  | 
**described_object** | [**IoK8sApiAutoscalingV2beta2CrossVersionObjectReference**](IoK8sApiAutoscalingV2beta2CrossVersionObjectReference.md) |  | 
**metric** | [**IoK8sApiAutoscalingV2beta2MetricIdentifier**](IoK8sApiAutoscalingV2beta2MetricIdentifier.md) |  | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAutoscalingV2beta2ObjectMetricStatus.new(current: null,
                                 described_object: null,
                                 metric: null)
```


