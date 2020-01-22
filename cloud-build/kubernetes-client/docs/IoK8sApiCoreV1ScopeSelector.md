# Kubernetes::IoK8sApiCoreV1ScopeSelector

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**match_expressions** | [**Array&lt;IoK8sApiCoreV1ScopedResourceSelectorRequirement&gt;**](IoK8sApiCoreV1ScopedResourceSelectorRequirement.md) | A list of scope selector requirements by scope of the resources. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1ScopeSelector.new(match_expressions: null)
```


