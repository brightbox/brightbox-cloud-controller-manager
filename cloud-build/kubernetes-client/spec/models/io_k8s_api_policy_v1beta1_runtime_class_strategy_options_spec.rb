=begin
#Kubernetes

#No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

The version of the OpenAPI document: v1.17.1

Generated by: https://openapi-generator.tech
OpenAPI Generator version: 4.2.3-SNAPSHOT

=end

require 'spec_helper'
require 'json'
require 'date'

# Unit tests for Kubernetes::IoK8sApiPolicyV1beta1RuntimeClassStrategyOptions
# Automatically generated by openapi-generator (https://openapi-generator.tech)
# Please update as you see appropriate
describe 'IoK8sApiPolicyV1beta1RuntimeClassStrategyOptions' do
  before do
    # run before each test
    @instance = Kubernetes::IoK8sApiPolicyV1beta1RuntimeClassStrategyOptions.new
  end

  after do
    # run after each test
  end

  describe 'test an instance of IoK8sApiPolicyV1beta1RuntimeClassStrategyOptions' do
    it 'should create an instance of IoK8sApiPolicyV1beta1RuntimeClassStrategyOptions' do
      expect(@instance).to be_instance_of(Kubernetes::IoK8sApiPolicyV1beta1RuntimeClassStrategyOptions)
    end
  end
  describe 'test attribute "allowed_runtime_class_names"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

  describe 'test attribute "default_runtime_class_name"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

end
