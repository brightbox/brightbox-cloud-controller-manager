# Kubernetes::IoK8sApiCoreV1PreferredSchedulingTerm

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**preference** | [**IoK8sApiCoreV1NodeSelectorTerm**](IoK8sApiCoreV1NodeSelectorTerm.md) |  | 
**weight** | **Integer** | Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1PreferredSchedulingTerm.new(preference: null,
                                 weight: null)
```


