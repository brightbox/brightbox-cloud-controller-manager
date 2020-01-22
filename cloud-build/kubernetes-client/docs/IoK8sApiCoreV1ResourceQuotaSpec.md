# Kubernetes::IoK8sApiCoreV1ResourceQuotaSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**hard** | **Hash&lt;String, String&gt;** | hard is the set of desired hard limits for each named resource. More info: https://kubernetes.io/docs/concepts/policy/resource-quotas/ | [optional] 
**scope_selector** | [**IoK8sApiCoreV1ScopeSelector**](IoK8sApiCoreV1ScopeSelector.md) |  | [optional] 
**scopes** | **Array&lt;String&gt;** | A collection of filters that must match each object tracked by a quota. If not specified, the quota matches all objects. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1ResourceQuotaSpec.new(hard: null,
                                 scope_selector: null,
                                 scopes: null)
```


