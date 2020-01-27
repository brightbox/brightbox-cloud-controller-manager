# Kubernetes::IoK8sApiSchedulingV1beta1PriorityClass

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**api_version** | **String** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**description** | **String** | description is an arbitrary string that usually provides guidelines on when this priority class should be used. | [optional] 
**global_default** | **Boolean** | globalDefault specifies whether this PriorityClass should be considered as the default priority for pods that do not have any priority class. Only one PriorityClass can be marked as &#x60;globalDefault&#x60;. However, if more than one PriorityClasses exists with their &#x60;globalDefault&#x60; field set to true, the smallest value of such global default PriorityClasses will be used as the default priority. | [optional] 
**kind** | **String** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**metadata** | [**IoK8sApimachineryPkgApisMetaV1ObjectMeta**](IoK8sApimachineryPkgApisMetaV1ObjectMeta.md) |  | [optional] 
**preemption_policy** | **String** | PreemptionPolicy is the Policy for preempting pods with lower priority. One of Never, PreemptLowerPriority. Defaults to PreemptLowerPriority if unset. This field is alpha-level and is only honored by servers that enable the NonPreemptingPriority feature. | [optional] 
**value** | **Integer** | The value of this priority class. This is the actual priority that pods receive when they have the name of this class in their pod spec. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiSchedulingV1beta1PriorityClass.new(api_version: null,
                                 description: null,
                                 global_default: null,
                                 kind: null,
                                 metadata: null,
                                 preemption_policy: null,
                                 value: null)
```


