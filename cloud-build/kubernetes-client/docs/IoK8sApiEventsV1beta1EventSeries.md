# Kubernetes::IoK8sApiEventsV1beta1EventSeries

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**count** | **Integer** | Number of occurrences in this series up to the last heartbeat time | 
**last_observed_time** | **DateTime** | MicroTime is version of Time with microsecond level precision. | 
**state** | **String** | Information whether this series is ongoing or finished. Deprecated. Planned removal for 1.18 | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiEventsV1beta1EventSeries.new(count: null,
                                 last_observed_time: null,
                                 state: null)
```


