# Kubernetes::IoK8sApiCoreV1EndpointPort

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | The name of this port.  This must match the &#39;name&#39; field in the corresponding ServicePort. Must be a DNS_LABEL. Optional only if one port is defined. | [optional] 
**port** | **Integer** | The port number of the endpoint. | 
**protocol** | **String** | The IP protocol for this port. Must be UDP, TCP, or SCTP. Default is TCP. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1EndpointPort.new(name: null,
                                 port: null,
                                 protocol: null)
```


