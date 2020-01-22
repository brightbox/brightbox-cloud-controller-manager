=begin
#Kubernetes

#No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

The version of the OpenAPI document: v1.17.1

Generated by: https://openapi-generator.tech
OpenAPI Generator version: 4.2.3-SNAPSHOT

=end

require 'date'

module Kubernetes
  # ServiceSpec describes the attributes that a user creates on a service.
  class IoK8sApiCoreV1ServiceSpec
    # clusterIP is the IP address of the service and is usually assigned randomly by the master. If an address is specified manually and is not in use by others, it will be allocated to the service; otherwise, creation of the service will fail. This field can not be changed through updates. Valid values are \"None\", empty string (\"\"), or a valid IP address. \"None\" can be specified for headless services when proxying is not required. Only applies to types ClusterIP, NodePort, and LoadBalancer. Ignored if type is ExternalName. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
    attr_accessor :cluster_ip

    # externalIPs is a list of IP addresses for which nodes in the cluster will also accept traffic for this service.  These IPs are not managed by Kubernetes.  The user is responsible for ensuring that traffic arrives at a node with this IP.  A common example is external load-balancers that are not part of the Kubernetes system.
    attr_accessor :external_i_ps

    # externalName is the external reference that kubedns or equivalent will return as a CNAME record for this service. No proxying will be involved. Must be a valid RFC-1123 hostname (https://tools.ietf.org/html/rfc1123) and requires Type to be ExternalName.
    attr_accessor :external_name

    # externalTrafficPolicy denotes if this Service desires to route external traffic to node-local or cluster-wide endpoints. \"Local\" preserves the client source IP and avoids a second hop for LoadBalancer and Nodeport type services, but risks potentially imbalanced traffic spreading. \"Cluster\" obscures the client source IP and may cause a second hop to another node, but should have good overall load-spreading.
    attr_accessor :external_traffic_policy

    # healthCheckNodePort specifies the healthcheck nodePort for the service. If not specified, HealthCheckNodePort is created by the service api backend with the allocated nodePort. Will use user-specified nodePort value if specified by the client. Only effects when Type is set to LoadBalancer and ExternalTrafficPolicy is set to Local.
    attr_accessor :health_check_node_port

    # ipFamily specifies whether this Service has a preference for a particular IP family (e.g. IPv4 vs. IPv6).  If a specific IP family is requested, the clusterIP field will be allocated from that family, if it is available in the cluster.  If no IP family is requested, the cluster's primary IP family will be used. Other IP fields (loadBalancerIP, loadBalancerSourceRanges, externalIPs) and controllers which allocate external load-balancers should use the same IP family.  Endpoints for this Service will be of this family.  This field is immutable after creation. Assigning a ServiceIPFamily not available in the cluster (e.g. IPv6 in IPv4 only cluster) is an error condition and will fail during clusterIP assignment.
    attr_accessor :ip_family

    # Only applies to Service Type: LoadBalancer LoadBalancer will get created with the IP specified in this field. This feature depends on whether the underlying cloud-provider supports specifying the loadBalancerIP when a load balancer is created. This field will be ignored if the cloud-provider does not support the feature.
    attr_accessor :load_balancer_ip

    # If specified and supported by the platform, this will restrict traffic through the cloud-provider load-balancer will be restricted to the specified client IPs. This field will be ignored if the cloud-provider does not support the feature.\" More info: https://kubernetes.io/docs/tasks/access-application-cluster/configure-cloud-provider-firewall/
    attr_accessor :load_balancer_source_ranges

    # The list of ports that are exposed by this service. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
    attr_accessor :ports

    # publishNotReadyAddresses, when set to true, indicates that DNS implementations must publish the notReadyAddresses of subsets for the Endpoints associated with the Service. The default value is false. The primary use case for setting this field is to use a StatefulSet's Headless Service to propagate SRV records for its Pods without respect to their readiness for purpose of peer discovery.
    attr_accessor :publish_not_ready_addresses

    # Route service traffic to pods with label keys and values matching this selector. If empty or not present, the service is assumed to have an external process managing its endpoints, which Kubernetes will not modify. Only applies to types ClusterIP, NodePort, and LoadBalancer. Ignored if type is ExternalName. More info: https://kubernetes.io/docs/concepts/services-networking/service/
    attr_accessor :selector

    # Supports \"ClientIP\" and \"None\". Used to maintain session affinity. Enable client IP based session affinity. Must be ClientIP or None. Defaults to None. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
    attr_accessor :session_affinity

    attr_accessor :session_affinity_config

    # topologyKeys is a preference-order list of topology keys which implementations of services should use to preferentially sort endpoints when accessing this Service, it can not be used at the same time as externalTrafficPolicy=Local. Topology keys must be valid label keys and at most 16 keys may be specified. Endpoints are chosen based on the first topology key with available backends. If this field is specified and all entries have no backends that match the topology of the client, the service has no backends for that client and connections should fail. The special value \"*\" may be used to mean \"any topology\". This catch-all value, if used, only makes sense as the last value in the list. If this is not specified or empty, no topology constraints will be applied.
    attr_accessor :topology_keys

    # type determines how the Service is exposed. Defaults to ClusterIP. Valid options are ExternalName, ClusterIP, NodePort, and LoadBalancer. \"ExternalName\" maps to the specified externalName. \"ClusterIP\" allocates a cluster-internal IP address for load-balancing to endpoints. Endpoints are determined by the selector or if that is not specified, by manual construction of an Endpoints object. If clusterIP is \"None\", no virtual IP is allocated and the endpoints are published as a set of endpoints rather than a stable IP. \"NodePort\" builds on ClusterIP and allocates a port on every node which routes to the clusterIP. \"LoadBalancer\" builds on NodePort and creates an external load-balancer (if supported in the current cloud) which routes to the clusterIP. More info: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
    attr_accessor :type

    # Attribute mapping from ruby-style variable name to JSON key.
    def self.attribute_map
      {
        :'cluster_ip' => :'clusterIP',
        :'external_i_ps' => :'externalIPs',
        :'external_name' => :'externalName',
        :'external_traffic_policy' => :'externalTrafficPolicy',
        :'health_check_node_port' => :'healthCheckNodePort',
        :'ip_family' => :'ipFamily',
        :'load_balancer_ip' => :'loadBalancerIP',
        :'load_balancer_source_ranges' => :'loadBalancerSourceRanges',
        :'ports' => :'ports',
        :'publish_not_ready_addresses' => :'publishNotReadyAddresses',
        :'selector' => :'selector',
        :'session_affinity' => :'sessionAffinity',
        :'session_affinity_config' => :'sessionAffinityConfig',
        :'topology_keys' => :'topologyKeys',
        :'type' => :'type'
      }
    end

    # Attribute type mapping.
    def self.openapi_types
      {
        :'cluster_ip' => :'String',
        :'external_i_ps' => :'Array<String>',
        :'external_name' => :'String',
        :'external_traffic_policy' => :'String',
        :'health_check_node_port' => :'Integer',
        :'ip_family' => :'String',
        :'load_balancer_ip' => :'String',
        :'load_balancer_source_ranges' => :'Array<String>',
        :'ports' => :'Array<IoK8sApiCoreV1ServicePort>',
        :'publish_not_ready_addresses' => :'Boolean',
        :'selector' => :'Hash<String, String>',
        :'session_affinity' => :'String',
        :'session_affinity_config' => :'IoK8sApiCoreV1SessionAffinityConfig',
        :'topology_keys' => :'Array<String>',
        :'type' => :'String'
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
        fail ArgumentError, "The input argument (attributes) must be a hash in `Kubernetes::IoK8sApiCoreV1ServiceSpec` initialize method"
      end

      # check to see if the attribute exists and convert string to symbol for hash key
      attributes = attributes.each_with_object({}) { |(k, v), h|
        if (!self.class.attribute_map.key?(k.to_sym))
          fail ArgumentError, "`#{k}` is not a valid attribute in `Kubernetes::IoK8sApiCoreV1ServiceSpec`. Please check the name to make sure it's valid. List of attributes: " + self.class.attribute_map.keys.inspect
        end
        h[k.to_sym] = v
      }

      if attributes.key?(:'cluster_ip')
        self.cluster_ip = attributes[:'cluster_ip']
      end

      if attributes.key?(:'external_i_ps')
        if (value = attributes[:'external_i_ps']).is_a?(Array)
          self.external_i_ps = value
        end
      end

      if attributes.key?(:'external_name')
        self.external_name = attributes[:'external_name']
      end

      if attributes.key?(:'external_traffic_policy')
        self.external_traffic_policy = attributes[:'external_traffic_policy']
      end

      if attributes.key?(:'health_check_node_port')
        self.health_check_node_port = attributes[:'health_check_node_port']
      end

      if attributes.key?(:'ip_family')
        self.ip_family = attributes[:'ip_family']
      end

      if attributes.key?(:'load_balancer_ip')
        self.load_balancer_ip = attributes[:'load_balancer_ip']
      end

      if attributes.key?(:'load_balancer_source_ranges')
        if (value = attributes[:'load_balancer_source_ranges']).is_a?(Array)
          self.load_balancer_source_ranges = value
        end
      end

      if attributes.key?(:'ports')
        if (value = attributes[:'ports']).is_a?(Array)
          self.ports = value
        end
      end

      if attributes.key?(:'publish_not_ready_addresses')
        self.publish_not_ready_addresses = attributes[:'publish_not_ready_addresses']
      end

      if attributes.key?(:'selector')
        if (value = attributes[:'selector']).is_a?(Hash)
          self.selector = value
        end
      end

      if attributes.key?(:'session_affinity')
        self.session_affinity = attributes[:'session_affinity']
      end

      if attributes.key?(:'session_affinity_config')
        self.session_affinity_config = attributes[:'session_affinity_config']
      end

      if attributes.key?(:'topology_keys')
        if (value = attributes[:'topology_keys']).is_a?(Array)
          self.topology_keys = value
        end
      end

      if attributes.key?(:'type')
        self.type = attributes[:'type']
      end
    end

    # Show invalid properties with the reasons. Usually used together with valid?
    # @return Array for valid properties with the reasons
    def list_invalid_properties
      invalid_properties = Array.new
      invalid_properties
    end

    # Check to see if the all the properties in the model are valid
    # @return true if the model is valid
    def valid?
      true
    end

    # Checks equality by comparing each attribute.
    # @param [Object] Object to be compared
    def ==(o)
      return true if self.equal?(o)
      self.class == o.class &&
          cluster_ip == o.cluster_ip &&
          external_i_ps == o.external_i_ps &&
          external_name == o.external_name &&
          external_traffic_policy == o.external_traffic_policy &&
          health_check_node_port == o.health_check_node_port &&
          ip_family == o.ip_family &&
          load_balancer_ip == o.load_balancer_ip &&
          load_balancer_source_ranges == o.load_balancer_source_ranges &&
          ports == o.ports &&
          publish_not_ready_addresses == o.publish_not_ready_addresses &&
          selector == o.selector &&
          session_affinity == o.session_affinity &&
          session_affinity_config == o.session_affinity_config &&
          topology_keys == o.topology_keys &&
          type == o.type
    end

    # @see the `==` method
    # @param [Object] Object to be compared
    def eql?(o)
      self == o
    end

    # Calculates hash code according to all attributes.
    # @return [Integer] Hash code
    def hash
      [cluster_ip, external_i_ps, external_name, external_traffic_policy, health_check_node_port, ip_family, load_balancer_ip, load_balancer_source_ranges, ports, publish_not_ready_addresses, selector, session_affinity, session_affinity_config, topology_keys, type].hash
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
