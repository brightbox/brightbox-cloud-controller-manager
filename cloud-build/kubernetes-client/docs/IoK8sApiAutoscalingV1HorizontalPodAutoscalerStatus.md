# Kubernetes::IoK8sApiAutoscalingV1HorizontalPodAutoscalerStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**current_cpu_utilization_percentage** | **Integer** | current average CPU utilization over all pods, represented as a percentage of requested CPU, e.g. 70 means that an average pod is using now 70% of its requested CPU. | [optional] 
**current_replicas** | **Integer** | current number of replicas of pods managed by this autoscaler. | 
**desired_replicas** | **Integer** | desired number of replicas of pods managed by this autoscaler. | 
**last_scale_time** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**observed_generation** | **Integer** | most recent generation observed by this autoscaler. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAutoscalingV1HorizontalPodAutoscalerStatus.new(current_cpu_utilization_percentage: null,
                                 current_replicas: null,
                                 desired_replicas: null,
                                 last_scale_time: null,
                                 observed_generation: null)
```


