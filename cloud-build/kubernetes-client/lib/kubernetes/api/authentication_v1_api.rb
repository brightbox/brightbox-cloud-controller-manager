=begin
#Kubernetes

#No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

The version of the OpenAPI document: v1.17.1

Generated by: https://openapi-generator.tech
OpenAPI Generator version: 4.2.3-SNAPSHOT

=end

require 'cgi'

module Kubernetes
  class AuthenticationV1Api
    attr_accessor :api_client

    def initialize(api_client = ApiClient.default)
      @api_client = api_client
    end
    # create a TokenReview
    # @param body [IoK8sApiAuthenticationV1TokenReview] 
    # @param [Hash] opts the optional parameters
    # @option opts [String] :dry_run When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed
    # @option opts [String] :field_manager fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint.
    # @option opts [String] :pretty If &#39;true&#39;, then the output is pretty printed.
    # @return [IoK8sApiAuthenticationV1TokenReview]
    def create_authentication_v1_token_review(body, opts = {})
      data, _status_code, _headers = create_authentication_v1_token_review_with_http_info(body, opts)
      data
    end

    # create a TokenReview
    # @param body [IoK8sApiAuthenticationV1TokenReview] 
    # @param [Hash] opts the optional parameters
    # @option opts [String] :dry_run When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed
    # @option opts [String] :field_manager fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint.
    # @option opts [String] :pretty If &#39;true&#39;, then the output is pretty printed.
    # @return [Array<(IoK8sApiAuthenticationV1TokenReview, Integer, Hash)>] IoK8sApiAuthenticationV1TokenReview data, response status code and response headers
    def create_authentication_v1_token_review_with_http_info(body, opts = {})
      if @api_client.config.debugging
        @api_client.config.logger.debug 'Calling API: AuthenticationV1Api.create_authentication_v1_token_review ...'
      end
      # verify the required parameter 'body' is set
      if @api_client.config.client_side_validation && body.nil?
        fail ArgumentError, "Missing the required parameter 'body' when calling AuthenticationV1Api.create_authentication_v1_token_review"
      end
      # resource path
      local_var_path = '/apis/authentication.k8s.io/v1/tokenreviews'

      # query parameters
      query_params = opts[:query_params] || {}
      query_params[:'dryRun'] = opts[:'dry_run'] if !opts[:'dry_run'].nil?
      query_params[:'fieldManager'] = opts[:'field_manager'] if !opts[:'field_manager'].nil?
      query_params[:'pretty'] = opts[:'pretty'] if !opts[:'pretty'].nil?

      # header parameters
      header_params = opts[:header_params] || {}
      # HTTP header 'Accept' (if needed)
      header_params['Accept'] = @api_client.select_header_accept(['application/json', 'application/yaml', 'application/vnd.kubernetes.protobuf'])

      # form parameters
      form_params = opts[:form_params] || {}

      # http body (model)
      post_body = opts[:body] || @api_client.object_to_http_body(body) 

      # return_type
      return_type = opts[:return_type] || 'IoK8sApiAuthenticationV1TokenReview' 

      # auth_names
      auth_names = opts[:auth_names] || ['BearerToken']

      new_options = opts.merge(
        :header_params => header_params,
        :query_params => query_params,
        :form_params => form_params,
        :body => post_body,
        :auth_names => auth_names,
        :return_type => return_type
      )

      data, status_code, headers = @api_client.call_api(:POST, local_var_path, new_options)
      if @api_client.config.debugging
        @api_client.config.logger.debug "API called: AuthenticationV1Api#create_authentication_v1_token_review\nData: #{data.inspect}\nStatus code: #{status_code}\nHeaders: #{headers}"
      end
      return data, status_code, headers
    end

    # get available resources
    # @param [Hash] opts the optional parameters
    # @return [IoK8sApimachineryPkgApisMetaV1APIResourceList]
    def get_authentication_v1_api_resources(opts = {})
      data, _status_code, _headers = get_authentication_v1_api_resources_with_http_info(opts)
      data
    end

    # get available resources
    # @param [Hash] opts the optional parameters
    # @return [Array<(IoK8sApimachineryPkgApisMetaV1APIResourceList, Integer, Hash)>] IoK8sApimachineryPkgApisMetaV1APIResourceList data, response status code and response headers
    def get_authentication_v1_api_resources_with_http_info(opts = {})
      if @api_client.config.debugging
        @api_client.config.logger.debug 'Calling API: AuthenticationV1Api.get_authentication_v1_api_resources ...'
      end
      # resource path
      local_var_path = '/apis/authentication.k8s.io/v1/'

      # query parameters
      query_params = opts[:query_params] || {}

      # header parameters
      header_params = opts[:header_params] || {}
      # HTTP header 'Accept' (if needed)
      header_params['Accept'] = @api_client.select_header_accept(['application/json', 'application/yaml', 'application/vnd.kubernetes.protobuf'])

      # form parameters
      form_params = opts[:form_params] || {}

      # http body (model)
      post_body = opts[:body] 

      # return_type
      return_type = opts[:return_type] || 'IoK8sApimachineryPkgApisMetaV1APIResourceList' 

      # auth_names
      auth_names = opts[:auth_names] || ['BearerToken']

      new_options = opts.merge(
        :header_params => header_params,
        :query_params => query_params,
        :form_params => form_params,
        :body => post_body,
        :auth_names => auth_names,
        :return_type => return_type
      )

      data, status_code, headers = @api_client.call_api(:GET, local_var_path, new_options)
      if @api_client.config.debugging
        @api_client.config.logger.debug "API called: AuthenticationV1Api#get_authentication_v1_api_resources\nData: #{data.inspect}\nStatus code: #{status_code}\nHeaders: #{headers}"
      end
      return data, status_code, headers
    end
  end
end