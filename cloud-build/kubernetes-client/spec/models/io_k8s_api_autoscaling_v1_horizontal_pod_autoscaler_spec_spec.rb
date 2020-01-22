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

# Unit tests for Kubernetes::IoK8sApiAutoscalingV1HorizontalPodAutoscalerSpec
# Automatically generated by openapi-generator (https://openapi-generator.tech)
# Please update as you see appropriate
describe 'IoK8sApiAutoscalingV1HorizontalPodAutoscalerSpec' do
  before do
    # run before each test
    @instance = Kubernetes::IoK8sApiAutoscalingV1HorizontalPodAutoscalerSpec.new
  end

  after do
    # run after each test
  end

  describe 'test an instance of IoK8sApiAutoscalingV1HorizontalPodAutoscalerSpec' do
    it 'should create an instance of IoK8sApiAutoscalingV1HorizontalPodAutoscalerSpec' do
      expect(@instance).to be_instance_of(Kubernetes::IoK8sApiAutoscalingV1HorizontalPodAutoscalerSpec)
    end
  end
  describe 'test attribute "max_replicas"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

  describe 'test attribute "min_replicas"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

  describe 'test attribute "scale_target_ref"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

  describe 'test attribute "target_cpu_utilization_percentage"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

end
