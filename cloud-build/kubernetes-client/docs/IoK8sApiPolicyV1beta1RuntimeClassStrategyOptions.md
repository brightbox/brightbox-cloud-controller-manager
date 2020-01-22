# Kubernetes::IoK8sApiPolicyV1beta1RuntimeClassStrategyOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**allowed_runtime_class_names** | **Array&lt;String&gt;** | allowedRuntimeClassNames is a whitelist of RuntimeClass names that may be specified on a pod. A value of \&quot;*\&quot; means that any RuntimeClass name is allowed, and must be the only item in the list. An empty list requires the RuntimeClassName field to be unset. | 
**default_runtime_class_name** | **String** | defaultRuntimeClassName is the default RuntimeClassName to set on the pod. The default MUST be allowed by the allowedRuntimeClassNames list. A value of nil does not mutate the Pod. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiPolicyV1beta1RuntimeClassStrategyOptions.new(allowed_runtime_class_names: null,
                                 default_runtime_class_name: null)
```


