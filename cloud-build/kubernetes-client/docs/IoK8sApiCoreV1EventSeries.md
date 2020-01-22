# Kubernetes::IoK8sApiCoreV1EventSeries

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**count** | **Integer** | Number of occurrences in this series up to the last heartbeat time | [optional] 
**last_observed_time** | **DateTime** | MicroTime is version of Time with microsecond level precision. | [optional] 
**state** | **String** | State of this Series: Ongoing or Finished Deprecated. Planned removal for 1.18 | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1EventSeries.new(count: null,
                                 last_observed_time: null,
                                 state: null)
```


