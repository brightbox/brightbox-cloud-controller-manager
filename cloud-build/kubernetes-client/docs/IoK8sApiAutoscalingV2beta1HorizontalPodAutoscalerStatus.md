# Kubernetes::IoK8sApiAutoscalingV2beta1HorizontalPodAutoscalerStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**conditions** | [**Array&lt;IoK8sApiAutoscalingV2beta1HorizontalPodAutoscalerCondition&gt;**](IoK8sApiAutoscalingV2beta1HorizontalPodAutoscalerCondition.md) | conditions is the set of conditions required for this autoscaler to scale its target, and indicates whether or not those conditions are met. | 
**current_metrics** | [**Array&lt;IoK8sApiAutoscalingV2beta1MetricStatus&gt;**](IoK8sApiAutoscalingV2beta1MetricStatus.md) | currentMetrics is the last read state of the metrics used by this autoscaler. | [optional] 
**current_replicas** | **Integer** | currentReplicas is current number of replicas of pods managed by this autoscaler, as last seen by the autoscaler. | 
**desired_replicas** | **Integer** | desiredReplicas is the desired number of replicas of pods managed by this autoscaler, as last calculated by the autoscaler. | 
**last_scale_time** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**observed_generation** | **Integer** | observedGeneration is the most recent generation observed by this autoscaler. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAutoscalingV2beta1HorizontalPodAutoscalerStatus.new(conditions: null,
                                 current_metrics: null,
                                 current_replicas: null,
                                 desired_replicas: null,
                                 last_scale_time: null,
                                 observed_generation: null)
```


