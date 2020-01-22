# Kubernetes::IoK8sApiAutoscalingV2beta2ObjectMetricSource

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**described_object** | [**IoK8sApiAutoscalingV2beta2CrossVersionObjectReference**](IoK8sApiAutoscalingV2beta2CrossVersionObjectReference.md) |  | 
**metric** | [**IoK8sApiAutoscalingV2beta2MetricIdentifier**](IoK8sApiAutoscalingV2beta2MetricIdentifier.md) |  | 
**target** | [**IoK8sApiAutoscalingV2beta2MetricTarget**](IoK8sApiAutoscalingV2beta2MetricTarget.md) |  | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAutoscalingV2beta2ObjectMetricSource.new(described_object: null,
                                 metric: null,
                                 target: null)
```


