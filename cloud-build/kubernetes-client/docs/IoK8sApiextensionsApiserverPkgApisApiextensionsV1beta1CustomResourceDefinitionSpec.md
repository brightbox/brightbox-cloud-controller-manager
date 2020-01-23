# Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceDefinitionSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**additional_printer_columns** | [**Array&lt;IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceColumnDefinition&gt;**](IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceColumnDefinition.md) | additionalPrinterColumns specifies additional columns returned in Table output. See https://kubernetes.io/docs/reference/using-api/api-concepts/#receiving-resources-as-tables for details. If present, this field configures columns for all versions. Top-level and per-version columns are mutually exclusive. If no top-level or per-version columns are specified, a single column displaying the age of the custom resource is used. | [optional] 
**conversion** | [**IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceConversion**](IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceConversion.md) |  | [optional] 
**group** | **String** | group is the API group of the defined custom resource. The custom resources are served under &#x60;/apis/&lt;group&gt;/...&#x60;. Must match the name of the CustomResourceDefinition (in the form &#x60;&lt;names.plural&gt;.&lt;group&gt;&#x60;). | 
**names** | [**IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceDefinitionNames**](IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceDefinitionNames.md) |  | 
**preserve_unknown_fields** | **Boolean** | preserveUnknownFields indicates that object fields which are not specified in the OpenAPI schema should be preserved when persisting to storage. apiVersion, kind, metadata and known fields inside metadata are always preserved. If false, schemas must be defined for all versions. Defaults to true in v1beta for backwards compatibility. Deprecated: will be required to be false in v1. Preservation of unknown fields can be specified in the validation schema using the &#x60;x-kubernetes-preserve-unknown-fields: true&#x60; extension. See https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#pruning-versus-preserving-unknown-fields for details. | [optional] 
**scope** | **String** | scope indicates whether the defined custom resource is cluster- or namespace-scoped. Allowed values are &#x60;Cluster&#x60; and &#x60;Namespaced&#x60;. Default is &#x60;Namespaced&#x60;. | 
**subresources** | [**IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceSubresources**](IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceSubresources.md) |  | [optional] 
**validation** | [**IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceValidation**](IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceValidation.md) |  | [optional] 
**version** | **String** | version is the API version of the defined custom resource. The custom resources are served under &#x60;/apis/&lt;group&gt;/&lt;version&gt;/...&#x60;. Must match the name of the first item in the &#x60;versions&#x60; list if &#x60;version&#x60; and &#x60;versions&#x60; are both specified. Optional if &#x60;versions&#x60; is specified. Deprecated: use &#x60;versions&#x60; instead. | [optional] 
**versions** | [**Array&lt;IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceDefinitionVersion&gt;**](IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceDefinitionVersion.md) | versions is the list of all API versions of the defined custom resource. Optional if &#x60;version&#x60; is specified. The name of the first item in the &#x60;versions&#x60; list must match the &#x60;version&#x60; field if &#x60;version&#x60; and &#x60;versions&#x60; are both specified. Version names are used to compute the order in which served versions are listed in API discovery. If the version string is \&quot;kube-like\&quot;, it will sort above non \&quot;kube-like\&quot; version strings, which are ordered lexicographically. \&quot;Kube-like\&quot; versions start with a \&quot;v\&quot;, then are followed by a number (the major version), then optionally the string \&quot;alpha\&quot; or \&quot;beta\&quot; and another number (the minor version). These are sorted first by GA &gt; beta &gt; alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing major version, then minor version. An example sorted list of versions: v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiextensionsApiserverPkgApisApiextensionsV1beta1CustomResourceDefinitionSpec.new(additional_printer_columns: null,
                                 conversion: null,
                                 group: null,
                                 names: null,
                                 preserve_unknown_fields: null,
                                 scope: null,
                                 subresources: null,
                                 validation: null,
                                 version: null,
                                 versions: null)
```

