# Kubernetes::IoK8sApiCoreV1NodeSelectorTerm

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**match_expressions** | [**Array&lt;IoK8sApiCoreV1NodeSelectorRequirement&gt;**](IoK8sApiCoreV1NodeSelectorRequirement.md) | A list of node selector requirements by node&#39;s labels. | [optional] 
**match_fields** | [**Array&lt;IoK8sApiCoreV1NodeSelectorRequirement&gt;**](IoK8sApiCoreV1NodeSelectorRequirement.md) | A list of node selector requirements by node&#39;s fields. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1NodeSelectorTerm.new(match_expressions: null,
                                 match_fields: null)
```


