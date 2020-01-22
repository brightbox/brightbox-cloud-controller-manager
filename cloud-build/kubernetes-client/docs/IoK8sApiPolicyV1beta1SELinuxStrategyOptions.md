# Kubernetes::IoK8sApiPolicyV1beta1SELinuxStrategyOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**rule** | **String** | rule is the strategy that will dictate the allowable labels that may be set. | 
**se_linux_options** | [**IoK8sApiCoreV1SELinuxOptions**](IoK8sApiCoreV1SELinuxOptions.md) |  | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiPolicyV1beta1SELinuxStrategyOptions.new(rule: null,
                                 se_linux_options: null)
```


