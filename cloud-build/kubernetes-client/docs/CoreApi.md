# Kubernetes::CoreApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**get_core_api_versions**](CoreApi.md#get_core_api_versions) | **GET** /api/ | 



## get_core_api_versions

> IoK8sApimachineryPkgApisMetaV1APIVersions get_core_api_versions



get available API versions

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

api_instance = Kubernetes::CoreApi.new

begin
  result = api_instance.get_core_api_versions
  p result
rescue Kubernetes::ApiError => e
  puts "Exception when calling CoreApi->get_core_api_versions: #{e}"
end
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**IoK8sApimachineryPkgApisMetaV1APIVersions**](IoK8sApimachineryPkgApisMetaV1APIVersions.md)

### Authorization

[BearerToken](../README.md#BearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf
