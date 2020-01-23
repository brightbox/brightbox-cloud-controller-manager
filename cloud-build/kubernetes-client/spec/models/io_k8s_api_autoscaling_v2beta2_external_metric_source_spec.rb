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

# Unit tests for Kubernetes::IoK8sApiAutoscalingV2beta2ExternalMetricSource
# Automatically generated by openapi-generator (https://openapi-generator.tech)
# Please update as you see appropriate
describe 'IoK8sApiAutoscalingV2beta2ExternalMetricSource' do
  before do
    # run before each test
    @instance = Kubernetes::IoK8sApiAutoscalingV2beta2ExternalMetricSource.new
  end

  after do
    # run after each test
  end

  describe 'test an instance of IoK8sApiAutoscalingV2beta2ExternalMetricSource' do
    it 'should create an instance of IoK8sApiAutoscalingV2beta2ExternalMetricSource' do
      expect(@instance).to be_instance_of(Kubernetes::IoK8sApiAutoscalingV2beta2ExternalMetricSource)
    end
  end
  describe 'test attribute "metric"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

  describe 'test attribute "target"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

end