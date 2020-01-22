# Kubernetes::IoK8sApiCertificatesV1beta1CertificateSigningRequestStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**certificate** | **String** | If request was approved, the controller will place the issued certificate here. | [optional] 
**conditions** | [**Array&lt;IoK8sApiCertificatesV1beta1CertificateSigningRequestCondition&gt;**](IoK8sApiCertificatesV1beta1CertificateSigningRequestCondition.md) | Conditions applied to the request, such as approval or denial. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCertificatesV1beta1CertificateSigningRequestStatus.new(certificate: null,
                                 conditions: null)
```


