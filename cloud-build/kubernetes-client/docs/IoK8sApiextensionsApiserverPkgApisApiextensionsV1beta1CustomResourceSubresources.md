# Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceSubresources

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**scale** | [**IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceSubresourceScale**](IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceSubresourceScale.md) |  | [optional] 
**status** | [**Object**](.md) | CustomResourceSubresourceStatus defines how to serve the status subresource for CustomResources. Status is represented by the &#x60;.status&#x60; JSON path inside of a CustomResource. When set, * exposes a /status subresource for the custom resource * PUT requests to the /status subresource take a custom resource object, and ignore changes to anything except the status stanza * PUT/POST/PATCH requests to the custom resource ignore changes to the status stanza | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceSubresources.new(scale: null,
                                 status: null)
```


