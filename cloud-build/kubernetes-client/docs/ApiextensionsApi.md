# Kubernetes::ApiextensionsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**get_apiextensions_api_group**](ApiextensionsApi.md#get_apiextensions_api_group) | **GET** /apis/apiextensions.k8s.io/ | 



## get_apiextensions_api_group

> IoK8sApimachineryPkgApisMetaV1APIGroup get_apiextensions_api_group



get information of a group

### Example

```ruby
# load the gem
require 'kubernetes'
# setup authorization
Kubernetes.configure do |config|
  # Configure API key authorization: BearerToken
  config.api_key['authorization'] = 'YOUR API KEY'
  # Uncomment the following line to set a prefix for the API key, e.g. 'Bearer' (defaults to nil)
  #config.api_key_prefix['authorization'] = 'Bearer'
end

api_instance = Kubernetes::ApiextensionsApi.new

begin
  result = api_instance.get_apiextensions_api_group
  p result
rescue Kubernetes::ApiError => e
  puts "Exception when calling ApiextensionsApi->get_apiextensions_api_group: #{e}"
end
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**IoK8sApimachineryPkgApisMetaV1APIGroup**](IoK8sApimachineryPkgApisMetaV1APIGroup.md)

### Authorization

[BearerToken](../README.md#BearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf

