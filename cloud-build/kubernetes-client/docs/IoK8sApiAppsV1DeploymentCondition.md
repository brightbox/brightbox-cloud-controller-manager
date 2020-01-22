# Kubernetes::IoK8sApiAppsV1DeploymentCondition

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**last_transition_time** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**last_update_time** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**message** | **String** | A human readable message indicating details about the transition. | [optional] 
**reason** | **String** | The reason for the condition&#39;s last transition. | [optional] 
**status** | **String** | Status of the condition, one of True, False, Unknown. | 
**type** | **String** | Type of deployment condition. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAppsV1DeploymentCondition.new(last_transition_time: null,
                                 last_update_time: null,
                                 message: null,
                                 reason: null,
                                 status: null,
                                 type: null)
```


