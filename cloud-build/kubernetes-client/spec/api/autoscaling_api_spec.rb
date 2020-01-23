=begin
#Kubernetes

#No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

The version of the OpenAPI document: v1.17.1

Generated by: https://openapi-generator.tech
OpenAPI Generator version: 4.2.3-SNAPSHOT

=end

require 'spec_helper'
require 'json'

# Unit tests for Kubernetes::AutoscalingApi
# Automatically generated by openapi-generator (https://openapi-generator.tech)
# Please update as you see appropriate
describe 'AutoscalingApi' do
  before do
    # run before each test
    @api_instance = Kubernetes::AutoscalingApi.new
  end

  after do
    # run after each test
  end

  describe 'test an instance of AutoscalingApi' do
    it 'should create an instance of AutoscalingApi' do
      expect(@api_instance).to be_instance_of(Kubernetes::AutoscalingApi)
    end
  end

  # unit tests for get_autoscaling_api_group
  # get information of a group
  # @param [Hash] opts the optional parameters
  # @return [IoK8sApimachineryPkgApisMetaV1APIGroup]
  describe 'get_autoscaling_api_group test' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

end