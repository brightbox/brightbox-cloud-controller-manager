# Kubernetes::IoK8sApiCoreV1NodeSelector

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**node_selector_terms** | [**Array&lt;IoK8sApiCoreV1NodeSelectorTerm&gt;**](IoK8sApiCoreV1NodeSelectorTerm.md) | Required. A list of node selector terms. The terms are ORed. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1NodeSelector.new(node_selector_terms: null)
```


