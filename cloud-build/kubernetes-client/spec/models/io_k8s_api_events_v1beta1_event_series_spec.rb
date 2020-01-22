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

# Unit tests for Kubernetes::IoK8sApiEventsV1beta1EventSeries
# Automatically generated by openapi-generator (https://openapi-generator.tech)
# Please update as you see appropriate
describe 'IoK8sApiEventsV1beta1EventSeries' do
  before do
    # run before each test
    @instance = Kubernetes::IoK8sApiEventsV1beta1EventSeries.new
  end

  after do
    # run after each test
  end

  describe 'test an instance of IoK8sApiEventsV1beta1EventSeries' do
    it 'should create an instance of IoK8sApiEventsV1beta1EventSeries' do
      expect(@instance).to be_instance_of(Kubernetes::IoK8sApiEventsV1beta1EventSeries)
    end
  end
  describe 'test attribute "count"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

  describe 'test attribute "last_observed_time"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

  describe 'test attribute "state"' do
    it 'should work' do
      # assertion here. ref: https://www.relishapp.com/rspec/rspec-expectations/docs/built-in-matchers
    end
  end

end
