# Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1CustomResourceConversion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**strategy** | **String** | strategy specifies how custom resources are converted between versions. Allowed values are: - &#x60;None&#x60;: The converter only change the apiVersion and would not touch any other field in the custom resource. - &#x60;Webhook&#x60;: API Server will call to an external webhook to do the conversion. Additional information   is needed for this option. This requires spec.preserveUnknownFields to be false, and spec.conversion.webhook to be set. | 
**webhook** | [**IoK8sApiextensionsApiserverPkgApisApiextensionsV1WebhookConversion**](IoK8sApiextensionsApiserverPkgApisApiextensionsV1WebhookConversion.md) |  | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1CustomResourceConversion.new(strategy: null,
                                 webhook: null)
```


