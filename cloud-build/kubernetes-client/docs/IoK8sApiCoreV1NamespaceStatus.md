# Kubernetes::IoK8sApiCoreV1NamespaceStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**conditions** | [**Array&lt;IoK8sApiCoreV1NamespaceCondition&gt;**](IoK8sApiCoreV1NamespaceCondition.md) | Represents the latest available observations of a namespace&#39;s current state. | [optional] 
**phase** | **String** | Phase is the current lifecycle phase of the namespace. More info: https://kubernetes.io/docs/tasks/administer-cluster/namespaces/ | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1NamespaceStatus.new(conditions: null,
                                 phase: null)
```


