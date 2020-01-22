# Kubernetes::IoK8sApiCertificatesV1beta1CertificateSigningRequestCondition

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**last_update_time** | **DateTime** | Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers. | [optional] 
**message** | **String** | human readable message with details about the request state | [optional] 
**reason** | **String** | brief reason for the request state | [optional] 
**type** | **String** | request approval state, currently Approved or Denied. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCertificatesV1beta1CertificateSigningRequestCondition.new(last_update_time: null,
                                 message: null,
                                 reason: null,
                                 type: null)
```


