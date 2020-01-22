# Kubernetes::IoK8sApiRbacV1beta1ClusterRole

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**aggregation_rule** | [**IoK8sApiRbacV1beta1AggregationRule**](IoK8sApiRbacV1beta1AggregationRule.md) |  | [optional] 
**api_version** | **String** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**kind** | **String** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**metadata** | [**IoK8sApimachineryPkgApisMetaV1ObjectMeta**](IoK8sApimachineryPkgApisMetaV1ObjectMeta.md) |  | [optional] 
**rules** | [**Array&lt;IoK8sApiRbacV1beta1PolicyRule&gt;**](IoK8sApiRbacV1beta1PolicyRule.md) | Rules holds all the PolicyRules for this ClusterRole | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiRbacV1beta1ClusterRole.new(aggregation_rule: null,
                                 api_version: null,
                                 kind: null,
                                 metadata: null,
                                 rules: null)
```


