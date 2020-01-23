# Kubernetes::IoK8sApiPolicyV1beta1RunAsUserStrategyOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ranges** | [**Array&lt;IoK8sApiPolicyV1beta1IDRange&gt;**](IoK8sApiPolicyV1beta1IDRange.md) | ranges are the allowed ranges of uids that may be used. If you would like to force a single uid then supply a single range with the same start and end. Required for MustRunAs. | [optional] 
**rule** | **String** | rule is the strategy that will dictate the allowable RunAsUser values that may be set. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiPolicyV1beta1RunAsUserStrategyOptions.new(ranges: null,
                                 rule: null)
```

