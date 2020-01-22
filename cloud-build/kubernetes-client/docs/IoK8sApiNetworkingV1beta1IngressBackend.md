# Kubernetes::IoK8sApiNetworkingV1beta1IngressBackend

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**service_name** | **String** | Specifies the name of the referenced service. | 
**service_port** | **String** | IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiNetworkingV1beta1IngressBackend.new(service_name: null,
                                 service_port: null)
```


