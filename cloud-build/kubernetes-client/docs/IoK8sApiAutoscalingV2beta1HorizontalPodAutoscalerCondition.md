# Kubernetes::IoK8sApiAutoscalingV2beta1HorizontalPodAutoscalerCondition

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**last_transition_time** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**message** | **String** | message is a human-readable explanation containing details about the transition | [optional] 
**reason** | **String** | reason is the reason for the condition&#39;s last transition. | [optional] 
**status** | **String** | status is the status of the condition (True, False, Unknown) | 
**type** | **String** | type describes the current condition | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAutoscalingV2beta1HorizontalPodAutoscalerCondition.new(last_transition_time: null,
                                 message: null,
                                 reason: null,
                                 status: null,
                                 type: null)
```


