# Kubernetes::IoK8sApiAuthorizationV1beta1SelfSubjectAccessReviewSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**non_resource_attributes** | [**IoK8sApiAuthorizationV1beta1NonResourceAttributes**](IoK8sApiAuthorizationV1beta1NonResourceAttributes.md) |  | [optional] 
**resource_attributes** | [**IoK8sApiAuthorizationV1beta1ResourceAttributes**](IoK8sApiAuthorizationV1beta1ResourceAttributes.md) |  | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAuthorizationV1beta1SelfSubjectAccessReviewSpec.new(non_resource_attributes: null,
                                 resource_attributes: null)
```


