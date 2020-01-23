# Kubernetes::AuthenticationV1Api

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**create_authentication_v1_token_review**](AuthenticationV1Api.md#create_authentication_v1_token_review) | **POST** /apis/authentication.k8s.io/v1/tokenreviews | 
[**get_authentication_v1_api_resources**](AuthenticationV1Api.md#get_authentication_v1_api_resources) | **GET** /apis/authentication.k8s.io/v1/ | 



## create_authentication_v1_token_review

> IoK8sApiAuthenticationV1TokenReview create_authentication_v1_token_review(body, opts)



create a TokenReview

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

api_instance = Kubernetes::AuthenticationV1Api.new
body = Kubernetes::IoK8sApiAuthenticationV1TokenReview.new # IoK8sApiAuthenticationV1TokenReview | 
opts = {
  dry_run: 'dry_run_example', # String | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed
  field_manager: 'field_manager_example', # String | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint.
  pretty: 'pretty_example' # String | If 'true', then the output is pretty printed.
}

begin
  result = api_instance.create_authentication_v1_token_review(body, opts)
  p result
rescue Kubernetes::ApiError => e
  puts "Exception when calling AuthenticationV1Api->create_authentication_v1_token_review: #{e}"
end
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**IoK8sApiAuthenticationV1TokenReview**](IoK8sApiAuthenticationV1TokenReview.md)|  | 
 **dry_run** | **String**| When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | [optional] 
 **field_manager** | **String**| fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | [optional] 
 **pretty** | **String**| If &#39;true&#39;, then the output is pretty printed. | [optional] 

### Return type

[**IoK8sApiAuthenticationV1TokenReview**](IoK8sApiAuthenticationV1TokenReview.md)

### Authorization

[BearerToken](../README.md#BearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf


## get_authentication_v1_api_resources

> IoK8sApimachineryPkgApisMetaV1APIResourceList get_authentication_v1_api_resources



get available resources

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

api_instance = Kubernetes::AuthenticationV1Api.new

begin
  result = api_instance.get_authentication_v1_api_resources
  p result
rescue Kubernetes::ApiError => e
  puts "Exception when calling AuthenticationV1Api->get_authentication_v1_api_resources: #{e}"
end
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**IoK8sApimachineryPkgApisMetaV1APIResourceList**](IoK8sApimachineryPkgApisMetaV1APIResourceList.md)

### Authorization

[BearerToken](../README.md#BearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf
