# Kubernetes::OrgProjectcalicoCrdV1IPAMHandleList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**api_version** | **String** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**items** | **Array&lt;Object&gt;** | List of ipamhandles. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md | 
**kind** | **String** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**metadata** | [**OrgProjectcalicoCrdV1BGPConfigurationListMetadata**](OrgProjectcalicoCrdV1BGPConfigurationListMetadata.md) |  | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::OrgProjectcalicoCrdV1IPAMHandleList.new(api_version: null,
                                 items: null,
                                 kind: null,
                                 metadata: null)
```


