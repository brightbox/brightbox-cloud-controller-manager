# Kubernetes::AuthorizationV1Api

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**create_authorization_v1_namespaced_local_subject_access_review**](AuthorizationV1Api.md#create_authorization_v1_namespaced_local_subject_access_review) | **POST** /apis/authorization.k8s.io/v1/namespaces/{namespace}/localsubjectaccessreviews | 
[**create_authorization_v1_self_subject_access_review**](AuthorizationV1Api.md#create_authorization_v1_self_subject_access_review) | **POST** /apis/authorization.k8s.io/v1/selfsubjectaccessreviews | 
[**create_authorization_v1_self_subject_rules_review**](AuthorizationV1Api.md#create_authorization_v1_self_subject_rules_review) | **POST** /apis/authorization.k8s.io/v1/selfsubjectrulesreviews | 
[**create_authorization_v1_subject_access_review**](AuthorizationV1Api.md#create_authorization_v1_subject_access_review) | **POST** /apis/authorization.k8s.io/v1/subjectaccessreviews | 
[**get_authorization_v1_api_resources**](AuthorizationV1Api.md#get_authorization_v1_api_resources) | **GET** /apis/authorization.k8s.io/v1/ | 



## create_authorization_v1_namespaced_local_subject_access_review

> IoK8sApiAuthorizationV1LocalSubjectAccessReview create_authorization_v1_namespaced_local_subject_access_review(namespace, body, opts)



create a LocalSubjectAccessReview

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

api_instance = Kubernetes::AuthorizationV1Api.new
namespace = 'namespace_example' # String | object name and auth scope, such as for teams and projects
body = Kubernetes::IoK8sApiAuthorizationV1LocalSubjectAccessReview.new # IoK8sApiAuthorizationV1LocalSubjectAccessReview | 
opts = {
  dry_run: 'dry_run_example', # String | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed
  field_manager: 'field_manager_example', # String | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint.
  pretty: 'pretty_example' # String | If 'true', then the output is pretty printed.
}

begin
  result = api_instance.create_authorization_v1_namespaced_local_subject_access_review(namespace, body, opts)
  p result
rescue Kubernetes::ApiError => e
  puts "Exception when calling AuthorizationV1Api->create_authorization_v1_namespaced_local_subject_access_review: #{e}"
end
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **namespace** | **String**| object name and auth scope, such as for teams and projects | 
 **body** | [**IoK8sApiAuthorizationV1LocalSubjectAccessReview**](IoK8sApiAuthorizationV1LocalSubjectAccessReview.md)|  | 
 **dry_run** | **String**| When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | [optional] 
 **field_manager** | **String**| fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | [optional] 
 **pretty** | **String**| If &#39;true&#39;, then the output is pretty printed. | [optional] 

### Return type

[**IoK8sApiAuthorizationV1LocalSubjectAccessReview**](IoK8sApiAuthorizationV1LocalSubjectAccessReview.md)

### Authorization

[BearerToken](../README.md#BearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf


## create_authorization_v1_self_subject_access_review

> IoK8sApiAuthorizationV1SelfSubjectAccessReview create_authorization_v1_self_subject_access_review(body, opts)



create a SelfSubjectAccessReview

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

api_instance = Kubernetes::AuthorizationV1Api.new
body = Kubernetes::IoK8sApiAuthorizationV1SelfSubjectAccessReview.new # IoK8sApiAuthorizationV1SelfSubjectAccessReview | 
opts = {
  dry_run: 'dry_run_example', # String | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed
  field_manager: 'field_manager_example', # String | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint.
  pretty: 'pretty_example' # String | If 'true', then the output is pretty printed.
}

begin
  result = api_instance.create_authorization_v1_self_subject_access_review(body, opts)
  p result
rescue Kubernetes::ApiError => e
  puts "Exception when calling AuthorizationV1Api->create_authorization_v1_self_subject_access_review: #{e}"
end
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**IoK8sApiAuthorizationV1SelfSubjectAccessReview**](IoK8sApiAuthorizationV1SelfSubjectAccessReview.md)|  | 
 **dry_run** | **String**| When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | [optional] 
 **field_manager** | **String**| fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | [optional] 
 **pretty** | **String**| If &#39;true&#39;, then the output is pretty printed. | [optional] 

### Return type

[**IoK8sApiAuthorizationV1SelfSubjectAccessReview**](IoK8sApiAuthorizationV1SelfSubjectAccessReview.md)

### Authorization

[BearerToken](../README.md#BearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf


## create_authorization_v1_self_subject_rules_review

> IoK8sApiAuthorizationV1SelfSubjectRulesReview create_authorization_v1_self_subject_rules_review(body, opts)



create a SelfSubjectRulesReview

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

api_instance = Kubernetes::AuthorizationV1Api.new
body = Kubernetes::IoK8sApiAuthorizationV1SelfSubjectRulesReview.new # IoK8sApiAuthorizationV1SelfSubjectRulesReview | 
opts = {
  dry_run: 'dry_run_example', # String | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed
  field_manager: 'field_manager_example', # String | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint.
  pretty: 'pretty_example' # String | If 'true', then the output is pretty printed.
}

begin
  result = api_instance.create_authorization_v1_self_subject_rules_review(body, opts)
  p result
rescue Kubernetes::ApiError => e
  puts "Exception when calling AuthorizationV1Api->create_authorization_v1_self_subject_rules_review: #{e}"
end
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**IoK8sApiAuthorizationV1SelfSubjectRulesReview**](IoK8sApiAuthorizationV1SelfSubjectRulesReview.md)|  | 
 **dry_run** | **String**| When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | [optional] 
 **field_manager** | **String**| fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | [optional] 
 **pretty** | **String**| If &#39;true&#39;, then the output is pretty printed. | [optional] 

### Return type

[**IoK8sApiAuthorizationV1SelfSubjectRulesReview**](IoK8sApiAuthorizationV1SelfSubjectRulesReview.md)

### Authorization

[BearerToken](../README.md#BearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf


## create_authorization_v1_subject_access_review

> IoK8sApiAuthorizationV1SubjectAccessReview create_authorization_v1_subject_access_review(body, opts)



create a SubjectAccessReview

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

api_instance = Kubernetes::AuthorizationV1Api.new
body = Kubernetes::IoK8sApiAuthorizationV1SubjectAccessReview.new # IoK8sApiAuthorizationV1SubjectAccessReview | 
opts = {
  dry_run: 'dry_run_example', # String | When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed
  field_manager: 'field_manager_example', # String | fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint.
  pretty: 'pretty_example' # String | If 'true', then the output is pretty printed.
}

begin
  result = api_instance.create_authorization_v1_subject_access_review(body, opts)
  p result
rescue Kubernetes::ApiError => e
  puts "Exception when calling AuthorizationV1Api->create_authorization_v1_subject_access_review: #{e}"
end
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**IoK8sApiAuthorizationV1SubjectAccessReview**](IoK8sApiAuthorizationV1SubjectAccessReview.md)|  | 
 **dry_run** | **String**| When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed | [optional] 
 **field_manager** | **String**| fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint. | [optional] 
 **pretty** | **String**| If &#39;true&#39;, then the output is pretty printed. | [optional] 

### Return type

[**IoK8sApiAuthorizationV1SubjectAccessReview**](IoK8sApiAuthorizationV1SubjectAccessReview.md)

### Authorization

[BearerToken](../README.md#BearerToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/yaml, application/vnd.kubernetes.protobuf


## get_authorization_v1_api_resources

> IoK8sApimachineryPkgApisMetaV1APIResourceList get_authorization_v1_api_resources



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

api_instance = Kubernetes::AuthorizationV1Api.new

begin
  result = api_instance.get_authorization_v1_api_resources
  p result
rescue Kubernetes::ApiError => e
  puts "Exception when calling AuthorizationV1Api->get_authorization_v1_api_resources: #{e}"
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
