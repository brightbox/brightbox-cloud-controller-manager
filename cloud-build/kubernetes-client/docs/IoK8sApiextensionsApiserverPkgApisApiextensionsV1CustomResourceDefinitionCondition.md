# Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1CustomResourceDefinitionCondition

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**last_transition_time** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**message** | **String** | message is a human-readable message indicating details about last transition. | [optional] 
**reason** | **String** | reason is a unique, one-word, CamelCase reason for the condition&#39;s last transition. | [optional] 
**status** | **String** | status is the status of the condition. Can be True, False, Unknown. | 
**type** | **String** | type is the type of the condition. Types include Established, NamesAccepted and Terminating. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1CustomResourceDefinitionCondition.new(last_transition_time: null,
                                 message: null,
                                 reason: null,
                                 status: null,
                                 type: null)
```


