# Kubernetes::IoK8sApiCoreV1SecretEnvSource

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names | [optional] 
**optional** | **Boolean** | Specify whether the Secret must be defined | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1SecretEnvSource.new(name: null,
                                 optional: null)
```


