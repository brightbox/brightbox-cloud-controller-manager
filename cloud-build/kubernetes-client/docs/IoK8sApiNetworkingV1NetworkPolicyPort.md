# Kubernetes::IoK8sApiNetworkingV1NetworkPolicyPort

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**port** | **String** | IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number. | [optional] 
**protocol** | **String** | The protocol (TCP, UDP, or SCTP) which traffic must match. If not specified, this field defaults to TCP. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiNetworkingV1NetworkPolicyPort.new(port: null,
                                 protocol: null)
```


