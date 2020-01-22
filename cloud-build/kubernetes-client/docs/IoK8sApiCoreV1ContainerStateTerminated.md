# Kubernetes::IoK8sApiCoreV1ContainerStateTerminated

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**container_id** | **String** | Container&#39;s ID in the format &#39;docker://&lt;container_id&gt;&#39; | [optional] 
**exit_code** | **Integer** | Exit status from the last termination of the container | 
**finished_at** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**message** | **String** | Message regarding the last termination of the container | [optional] 
**reason** | **String** | (brief) reason from the last termination of the container | [optional] 
**signal** | **Integer** | Signal from the last termination of the container | [optional] 
**started_at** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1ContainerStateTerminated.new(container_id: null,
                                 exit_code: null,
                                 finished_at: null,
                                 message: null,
                                 reason: null,
                                 signal: null,
                                 started_at: null)
```


