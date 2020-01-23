=begin
#Kubernetes

#No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

The version of the OpenAPI document: v1.17.1

Generated by: https://openapi-generator.tech
OpenAPI Generator version: 4.2.3-SNAPSHOT

=end

require 'date'

module Kubernetes
  # Info contains versioning information. how we'll want to distribute that information.
  class IoK8sApimachineryPkgVersionInfo
    attr_accessor :build_date

    attr_accessor :compiler

    attr_accessor :git_commit

    attr_accessor :git_tree_state

    attr_accessor :git_version

    attr_accessor :go_version

    attr_accessor :major

    attr_accessor :minor

    attr_accessor :platform

    # Attribute mapping from ruby-style variable name to JSON key.
    def self.attribute_map
      {
        :'build_date' => :'buildDate',
        :'compiler' => :'compiler',
        :'git_commit' => :'gitCommit',
        :'git_tree_state' => :'gitTreeState',
        :'git_version' => :'gitVersion',
        :'go_version' => :'goVersion',
        :'major' => :'major',
        :'minor' => :'minor',
        :'platform' => :'platform'
      }
    end

    # Attribute type mapping.
    def self.openapi_types
      {
        :'build_date' => :'String',
        :'compiler' => :'String',
        :'git_commit' => :'String',
        :'git_tree_state' => :'String',
        :'git_version' => :'String',
        :'go_version' => :'String',
        :'major' => :'String',
        :'minor' => :'String',
        :'platform' => :'String'
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
        fail ArgumentError, "The input argument (attributes) must be a hash in `Kubernetes::IoK8sApimachineryPkgVersionInfo` initialize method"
      end

      # check to see if the attribute exists and convert string to symbol for hash key
      attributes = attributes.each_with_object({}) { |(k, v), h|
        if (!self.class.attribute_map.key?(k.to_sym))
          fail ArgumentError, "`#{k}` is not a valid attribute in `Kubernetes::IoK8sApimachineryPkgVersionInfo`. Please check the name to make sure it's valid. List of attributes: " + self.class.attribute_map.keys.inspect
        end
        h[k.to_sym] = v
      }

      if attributes.key?(:'build_date')
        self.build_date = attributes[:'build_date']
      end

      if attributes.key?(:'compiler')
        self.compiler = attributes[:'compiler']
      end

      if attributes.key?(:'git_commit')
        self.git_commit = attributes[:'git_commit']
      end

      if attributes.key?(:'git_tree_state')
        self.git_tree_state = attributes[:'git_tree_state']
      end

      if attributes.key?(:'git_version')
        self.git_version = attributes[:'git_version']
      end

      if attributes.key?(:'go_version')
        self.go_version = attributes[:'go_version']
      end

      if attributes.key?(:'major')
        self.major = attributes[:'major']
      end

      if attributes.key?(:'minor')
        self.minor = attributes[:'minor']
      end

      if attributes.key?(:'platform')
        self.platform = attributes[:'platform']
      end
    end

    # Show invalid properties with the reasons. Usually used together with valid?
    # @return Array for valid properties with the reasons
    def list_invalid_properties
      invalid_properties = Array.new
      if @build_date.nil?
        invalid_properties.push('invalid value for "build_date", build_date cannot be nil.')
      end

      if @compiler.nil?
        invalid_properties.push('invalid value for "compiler", compiler cannot be nil.')
      end

      if @git_commit.nil?
        invalid_properties.push('invalid value for "git_commit", git_commit cannot be nil.')
      end

      if @git_tree_state.nil?
        invalid_properties.push('invalid value for "git_tree_state", git_tree_state cannot be nil.')
      end

      if @git_version.nil?
        invalid_properties.push('invalid value for "git_version", git_version cannot be nil.')
      end

      if @go_version.nil?
        invalid_properties.push('invalid value for "go_version", go_version cannot be nil.')
      end

      if @major.nil?
        invalid_properties.push('invalid value for "major", major cannot be nil.')
      end

      if @minor.nil?
        invalid_properties.push('invalid value for "minor", minor cannot be nil.')
      end

      if @platform.nil?
        invalid_properties.push('invalid value for "platform", platform cannot be nil.')
      end

      invalid_properties
    end

    # Check to see if the all the properties in the model are valid
    # @return true if the model is valid
    def valid?
      return false if @build_date.nil?
      return false if @compiler.nil?
      return false if @git_commit.nil?
      return false if @git_tree_state.nil?
      return false if @git_version.nil?
      return false if @go_version.nil?
      return false if @major.nil?
      return false if @minor.nil?
      return false if @platform.nil?
      true
    end

    # Checks equality by comparing each attribute.
    # @param [Object] Object to be compared
    def ==(o)
      return true if self.equal?(o)
      self.class == o.class &&
          build_date == o.build_date &&
          compiler == o.compiler &&
          git_commit == o.git_commit &&
          git_tree_state == o.git_tree_state &&
          git_version == o.git_version &&
          go_version == o.go_version &&
          major == o.major &&
          minor == o.minor &&
          platform == o.platform
    end

    # @see the `==` method
    # @param [Object] Object to be compared
    def eql?(o)
      self == o
    end

    # Calculates hash code according to all attributes.
    # @return [Integer] Hash code
    def hash
      [build_date, compiler, git_commit, git_tree_state, git_version, go_version, major, minor, platform].hash
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