# Kubernetes::IoK8sApiBatchV1JobStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**active** | **Integer** | The number of actively running pods. | [optional] 
**completion_time** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**conditions** | [**Array&lt;IoK8sApiBatchV1JobCondition&gt;**](IoK8sApiBatchV1JobCondition.md) | The latest available observations of an object&#39;s current state. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/ | [optional] 
**failed** | **Integer** | The number of pods which reached phase Failed. | [optional] 
**start_time** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**succeeded** | **Integer** | The number of pods which reached phase Succeeded. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiBatchV1JobStatus.new(active: null,
                                 completion_time: null,
                                 conditions: null,
                                 failed: null,
                                 start_time: null,
                                 succeeded: null)
```


