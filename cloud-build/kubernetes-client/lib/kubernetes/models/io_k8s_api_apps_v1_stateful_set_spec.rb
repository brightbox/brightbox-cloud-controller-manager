=begin
#Kubernetes

#No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

The version of the OpenAPI document: v1.17.1

Generated by: https://openapi-generator.tech
OpenAPI Generator version: 4.2.3-SNAPSHOT

=end

require 'date'

module Kubernetes
  # A StatefulSetSpec is the specification of a StatefulSet.
  class IoK8sApiAppsV1StatefulSetSpec
    # podManagementPolicy controls how pods are created during initial scale up, when replacing pods on nodes, or when scaling down. The default policy is `OrderedReady`, where pods are created in increasing order (pod-0, then pod-1, etc) and the controller will wait until each pod is ready before continuing. When scaling down, the pods are removed in the opposite order. The alternative policy is `Parallel` which will create pods in parallel to match the desired scale without waiting, and on scale down will delete all pods at once.
    attr_accessor :pod_management_policy

    # replicas is the desired number of replicas of the given Template. These are replicas in the sense that they are instantiations of the same Template, but individual replicas also have a consistent identity. If unspecified, defaults to 1.
    attr_accessor :replicas

    # revisionHistoryLimit is the maximum number of revisions that will be maintained in the StatefulSet's revision history. The revision history consists of all revisions not represented by a currently applied StatefulSetSpec version. The default value is 10.
    attr_accessor :revision_history_limit

    attr_accessor :selector

    # serviceName is the name of the service that governs this StatefulSet. This service must exist before the StatefulSet, and is responsible for the network identity of the set. Pods get DNS/hostnames that follow the pattern: pod-specific-string.serviceName.default.svc.cluster.local where \"pod-specific-string\" is managed by the StatefulSet controller.
    attr_accessor :service_name

    attr_accessor :template

    attr_accessor :update_strategy

    # volumeClaimTemplates is a list of claims that pods are allowed to reference. The StatefulSet controller is responsible for mapping network identities to claims in a way that maintains the identity of a pod. Every claim in this list must have at least one matching (by name) volumeMount in one container in the template. A claim in this list takes precedence over any volumes in the template, with the same name.
    attr_accessor :volume_claim_templates

    # Attribute mapping from ruby-style variable name to JSON key.
    def self.attribute_map
      {
        :'pod_management_policy' => :'podManagementPolicy',
        :'replicas' => :'replicas',
        :'revision_history_limit' => :'revisionHistoryLimit',
        :'selector' => :'selector',
        :'service_name' => :'serviceName',
        :'template' => :'template',
        :'update_strategy' => :'updateStrategy',
        :'volume_claim_templates' => :'volumeClaimTemplates'
      }
    end

    # Attribute type mapping.
    def self.openapi_types
      {
        :'pod_management_policy' => :'String',
        :'replicas' => :'Integer',
        :'revision_history_limit' => :'Integer',
        :'selector' => :'IoK8sApimachineryPkgApisMetaV1LabelSelector',
        :'service_name' => :'String',
        :'template' => :'IoK8sApiCoreV1PodTemplateSpec',
        :'update_strategy' => :'IoK8sApiAppsV1StatefulSetUpdateStrategy',
        :'volume_claim_templates' => :'Array<IoK8sApiCoreV1PersistentVolumeClaim>'
      }
    end

    # List of attributes with nullable: true
    def self.openapi_nullable
      Set.new([
      ])
    end

    # Initializes the object
    # @param [Hash] attributes Model attributes in the form of hash
    def initialize(attributes = {})
      if (!attributes.is_a?(Hash))
        fail ArgumentError, "The input argument (attributes) must be a hash in `Kubernetes::IoK8sApiAppsV1StatefulSetSpec` initialize method"
      end

      # check to see if the attribute exists and convert string to symbol for hash key
      attributes = attributes.each_with_object({}) { |(k, v), h|
        if (!self.class.attribute_map.key?(k.to_sym))
          fail ArgumentError, "`#{k}` is not a valid attribute in `Kubernetes::IoK8sApiAppsV1StatefulSetSpec`. Please check the name to make sure it's valid. List of attributes: " + self.class.attribute_map.keys.inspect
        end
        h[k.to_sym] = v
      }

      if attributes.key?(:'pod_management_policy')
        self.pod_management_policy = attributes[:'pod_management_policy']
      end

      if attributes.key?(:'replicas')
        self.replicas = attributes[:'replicas']
      end

      if attributes.key?(:'revision_history_limit')
        self.revision_history_limit = attributes[:'revision_history_limit']
      end

      if attributes.key?(:'selector')
        self.selector = attributes[:'selector']
      end

      if attributes.key?(:'service_name')
        self.service_name = attributes[:'service_name']
      end

      if attributes.key?(:'template')
        self.template = attributes[:'template']
      end

      if attributes.key?(:'update_strategy')
        self.update_strategy = attributes[:'update_strategy']
      end

      if attributes.key?(:'volume_claim_templates')
        if (value = attributes[:'volume_claim_templates']).is_a?(Array)
          self.volume_claim_templates = value
        end
      end
    end

    # Show invalid properties with the reasons. Usually used together with valid?
    # @return Array for valid properties with the reasons
    def list_invalid_properties
      invalid_properties = Array.new
      if @selector.nil?
        invalid_properties.push('invalid value for "selector", selector cannot be nil.')
      end

      if @service_name.nil?
        invalid_properties.push('invalid value for "service_name", service_name cannot be nil.')
      end

      if @template.nil?
        invalid_properties.push('invalid value for "template", template cannot be nil.')
      end

      invalid_properties
    end

    # Check to see if the all the properties in the model are valid
    # @return true if the model is valid
    def valid?
      return false if @selector.nil?
      return false if @service_name.nil?
      return false if @template.nil?
      true
    end

    # Checks equality by comparing each attribute.
    # @param [Object] Object to be compared
    def ==(o)
      return true if self.equal?(o)
      self.class == o.class &&
          pod_management_policy == o.pod_management_policy &&
          replicas == o.replicas &&
          revision_history_limit == o.revision_history_limit &&
          selector == o.selector &&
          service_name == o.service_name &&
          template == o.template &&
          update_strategy == o.update_strategy &&
          volume_claim_templates == o.volume_claim_templates
    end

    # @see the `==` method
    # @param [Object] Object to be compared
    def eql?(o)
      self == o
    end

    # Calculates hash code according to all attributes.
    # @return [Integer] Hash code
    def hash
      [pod_management_policy, replicas, revision_history_limit, selector, service_name, template, update_strategy, volume_claim_templates].hash
    end

    # Builds the object from hash
    # @param [Hash] attributes Model attributes in the form of hash
    # @return [Object] Returns the model itself
    def self.build_from_hash(attributes)
      new.build_from_hash(attributes)
    end

    # Builds the object from hash
    # @param [Hash] attributes Model attributes in the form of hash
    # @return [Object] Returns the model itself
    def build_from_hash(attributes)
      return nil unless attributes.is_a?(Hash)
      self.class.openapi_types.each_pair do |key, type|
        if type =~ /\AArray<(.*)>/i
          # check to ensure the input is an array given that the attribute
          # is documented as an array but the input is not
          if attributes[self.class.attribute_map[key]].is_a?(Array)
            self.send("#{key}=", attributes[self.class.attribute_map[key]].map { |v| _deserialize($1, v) })
          end
        elsif !attributes[self.class.attribute_map[key]].nil?
          self.send("#{key}=", _deserialize(type, attributes[self.class.attribute_map[key]]))
        end # or else data not found in attributes(hash), not an issue as the data can be optional
      end

      self
    end

    # Deserializes the data based on type
    # @param string type Data type
    # @param string value Value to be deserialized
    # @return [Object] Deserialized data
    def _deserialize(type, value)
      case type.to_sym
      when :DateTime
        DateTime.parse(value)
      when :Date
        Date.parse(value)
      when :String
        value.to_s
      when :Integer
        value.to_i
      when :Float
        value.to_f
      when :Boolean
        if value.to_s =~ /\A(true|t|yes|y|1)\z/i
          true
        else
          false
        end
      when :Object
        # generic object (usually a Hash), return directly
        value
      when /\AArray<(?<inner_type>.+)>\z/
        inner_type = Regexp.last_match[:inner_type]
        value.map { |v| _deserialize(inner_type, v) }
      when /\AHash<(?<k_type>.+?), (?<v_type>.+)>\z/
        k_type = Regexp.last_match[:k_type]
        v_type = Regexp.last_match[:v_type]
        {}.tap do |hash|
          value.each do |k, v|
            hash[_deserialize(k_type, k)] = _deserialize(v_type, v)
          end
        end
      else # model
        Kubernetes.const_get(type).build_from_hash(value)
      end
    end

    # Returns the string representation of the object
    # @return [String] String presentation of the object
    def to_s
      to_hash.to_s
    end

    # to_body is an alias to to_hash (backward compatibility)
    # @return [Hash] Returns the object in the form of hash
    def to_body
      to_hash
    end

    # Returns the object in the form of hash
    # @return [Hash] Returns the object in the form of hash
    def to_hash
      hash = {}
      self.class.attribute_map.each_pair do |attr, param|
        value = self.send(attr)
        if value.nil?
          is_nullable = self.class.openapi_nullable.include?(attr)
          next if !is_nullable || (is_nullable && !instance_variable_defined?(:"@#{attr}"))
        end
        
        hash[param] = _to_hash(value)
      end
      hash
    end

    # Outputs non-array value in the form of hash
    # For object, use to_hash. Otherwise, just return the value
    # @param [Object] value Any valid value
    # @return [Hash] Returns the value in the form of hash
    def _to_hash(value)
      if value.is_a?(Array)
        value.compact.map { |v| _to_hash(v) }
      elsif value.is_a?(Hash)
        {}.tap do |hash|
          value.each { |k, v| hash[k] = _to_hash(v) }
        end
      elsif value.respond_to? :to_hash
        value.to_hash
      else
        value
      end
    end
  end
end
