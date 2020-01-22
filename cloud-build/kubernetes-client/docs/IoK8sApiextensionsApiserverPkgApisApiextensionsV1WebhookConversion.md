# Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1WebhookConversion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**client_config** | [**IoK8sApiextensionsApiserverPkgApisApiextensionsV1WebhookClientConfig**](IoK8sApiextensionsApiserverPkgApisApiextensionsV1WebhookClientConfig.md) |  | [optional] 
**conversion_review_versions** | **Array&lt;String&gt;** | conversionReviewVersions is an ordered list of preferred &#x60;ConversionReview&#x60; versions the Webhook expects. The API server will use the first version in the list which it supports. If none of the versions specified in this list are supported by API server, conversion will fail for the custom resource. If a persisted Webhook configuration specifies allowed versions and does not include any versions known to the API Server, calls to the webhook will fail. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1WebhookConversion.new(client_config: null,
                                 conversion_review_versions: null)
```


