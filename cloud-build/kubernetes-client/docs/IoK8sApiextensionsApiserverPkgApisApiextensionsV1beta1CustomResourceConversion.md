# Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceConversion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**conversion_review_versions** | **Array&lt;String&gt;** | conversionReviewVersions is an ordered list of preferred &#x60;ConversionReview&#x60; versions the Webhook expects. The API server will use the first version in the list which it supports. If none of the versions specified in this list are supported by API server, conversion will fail for the custom resource. If a persisted Webhook configuration specifies allowed versions and does not include any versions known to the API Server, calls to the webhook will fail. Defaults to &#x60;[\&quot;v1beta1\&quot;]&#x60;. | [optional] 
**strategy** | **String** | strategy specifies how custom resources are converted between versions. Allowed values are: - &#x60;None&#x60;: The converter only change the apiVersion and would not touch any other field in the custom resource. - &#x60;Webhook&#x60;: API Server will call to an external webhook to do the conversion. Additional information   is needed for this option. This requires spec.preserveUnknownFields to be false, and spec.conversion.webhookClientConfig to be set. | 
**webhook_client_config** | [**IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1WebhookClientConfig**](IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1WebhookClientConfig.md) |  | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceConversion.new(conversion_review_versions: null,
                                 strategy: null,
                                 webhook_client_config: null)
```


