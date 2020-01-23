# Kubernetes::VersionApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**get_code_version**](VersionApi.md#get_code_version) | **GET** /version/ | 



## get_code_version

> IoK8sApimachineryPkgVersionInfo get_code_version



get the code version

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

api_instance = Kubernetes::VersionApi.new

begin
  result = api_instance.get_code_version
  p result
rescue Kubernetes::ApiError => e
  puts "Exception when calling VersionApi->get_code_version: #{e}"
end
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**IoK8sApimachineryPkgVersionInfo**](IoK8sApimachineryPkgVersionInfo.md)

### Authorization

[BearerToken](../README.md#BearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json
