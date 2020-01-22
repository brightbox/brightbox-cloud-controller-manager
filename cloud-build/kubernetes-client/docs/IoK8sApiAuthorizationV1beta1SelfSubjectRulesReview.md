# Kubernetes::IoK8sApiAuthorizationV1beta1SelfSubjectRulesReview

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**api_version** | **String** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**kind** | **String** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**metadata** | [**IoK8sApimachineryPkgApisMetaV1ObjectMeta**](IoK8sApimachineryPkgApisMetaV1ObjectMeta.md) |  | [optional] 
**spec** | [**IoK8sApiAuthorizationV1beta1SelfSubjectRulesReviewSpec**](IoK8sApiAuthorizationV1beta1SelfSubjectRulesReviewSpec.md) |  | 
**status** | [**IoK8sApiAuthorizationV1beta1SubjectRulesReviewStatus**](IoK8sApiAuthorizationV1beta1SubjectRulesReviewStatus.md) |  | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAuthorizationV1beta1SelfSubjectRulesReview.new(api_version: null,
                                 kind: null,
                                 metadata: null,
                                 spec: null,
                                 status: null)
```


