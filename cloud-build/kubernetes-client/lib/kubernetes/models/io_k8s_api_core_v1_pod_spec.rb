=begin
#Kubernetes

#No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

The version of the OpenAPI document: v1.17.1

Generated by: https://openapi-generator.tech
OpenAPI Generator version: 4.2.3-SNAPSHOT

=end

require 'date'

module Kubernetes
  # PodSpec is a description of a pod.
  class IoK8sApiCoreV1PodSpec
    # Optional duration in seconds the pod may be active on the node relative to StartTime before the system will actively try to mark it failed and kill associated containers. Value must be a positive integer.
    attr_accessor :active_deadline_seconds

    attr_accessor :affinity

    # AutomountServiceAccountToken indicates whether a service account token should be automatically mounted.
    attr_accessor :automount_service_account_token

    # List of containers belonging to the pod. Containers cannot currently be added or removed. There must be at least one container in a Pod. Cannot be updated.
    attr_accessor :containers

    attr_accessor :dns_config

    # Set DNS policy for the pod. Defaults to \"ClusterFirst\". Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'. DNS parameters given in DNSConfig will be merged with the policy selected with DNSPolicy. To have DNS options set along with hostNetwork, you have to specify DNS policy explicitly to 'ClusterFirstWithHostNet'.
    attr_accessor :dns_policy

    # EnableServiceLinks indicates whether information about services should be injected into pod's environment variables, matching the syntax of Docker links. Optional: Defaults to true.
    attr_accessor :enable_service_links

    # List of ephemeral containers run in this pod. Ephemeral containers may be run in an existing pod to perform user-initiated actions such as debugging. This list cannot be specified when creating a pod, and it cannot be modified by updating the pod spec. In order to add an ephemeral container to an existing pod, use the pod's ephemeralcontainers subresource. This field is alpha-level and is only honored by servers that enable the EphemeralContainers feature.
    attr_accessor :ephemeral_containers

    # HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts file if specified. This is only valid for non-hostNetwork pods.
    attr_accessor :host_aliases

    # Use the host's ipc namespace. Optional: Default to false.
    attr_accessor :host_ipc

    # Host networking requested for this pod. Use the host's network namespace. If this option is set, the ports that will be used must be specified. Default to false.
    attr_accessor :host_network

    # Use the host's pid namespace. Optional: Default to false.
    attr_accessor :host_pid

    # Specifies the hostname of the Pod If not specified, the pod's hostname will be set to a system-defined value.
    attr_accessor :hostname

    # ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec. If specified, these secrets will be passed to individual puller implementations for them to use. For example, in the case of docker, only DockerConfig type secrets are honored. More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
    attr_accessor :image_pull_secrets

    # List of initialization containers belonging to the pod. Init containers are executed in order prior to containers being started. If any init container fails, the pod is considered to have failed and is handled according to its restartPolicy. The name for an init container or normal container must be unique among all containers. Init containers may not have Lifecycle actions, Readiness probes, Liveness probes, or Startup probes. The resourceRequirements of an init container are taken into account during scheduling by finding the highest request/limit for each resource type, and then using the max of of that value or the sum of the normal containers. Limits are applied to init containers in a similar fashion. Init containers cannot currently be added or removed. Cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/
    attr_accessor :init_containers

    # NodeName is a request to schedule this pod onto a specific node. If it is non-empty, the scheduler simply schedules this pod onto that node, assuming that it fits resource requirements.
    attr_accessor :node_name

    # NodeSelector is a selector which must be true for the pod to fit on a node. Selector which must match a node's labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
    attr_accessor :node_selector

    # Overhead represents the resource overhead associated with running a pod for a given RuntimeClass. This field will be autopopulated at admission time by the RuntimeClass admission controller. If the RuntimeClass admission controller is enabled, overhead must not be set in Pod create requests. The RuntimeClass admission controller will reject Pod create requests which have the overhead already set. If RuntimeClass is configured and selected in the PodSpec, Overhead will be set to the value defined in the corresponding RuntimeClass, otherwise it will remain unset and treated as zero. More info: https://git.k8s.io/enhancements/keps/sig-node/20190226-pod-overhead.md This field is alpha-level as of Kubernetes v1.16, and is only honored by servers that enable the PodOverhead feature.
    attr_accessor :overhead

    # PreemptionPolicy is the Policy for preempting pods with lower priority. One of Never, PreemptLowerPriority. Defaults to PreemptLowerPriority if unset. This field is alpha-level and is only honored by servers that enable the NonPreemptingPriority feature.
    attr_accessor :preemption_policy

    # The priority value. Various system components use this field to find the priority of the pod. When Priority Admission Controller is enabled, it prevents users from setting this field. The admission controller populates this field from PriorityClassName. The higher the value, the higher the priority.
    attr_accessor :priority

    # If specified, indicates the pod's priority. \"system-node-critical\" and \"system-cluster-critical\" are two special keywords which indicate the highest priorities with the former being the highest priority. Any other name must be defined by creating a PriorityClass object with that name. If not specified, the pod priority will be default or zero if there is no default.
    attr_accessor :priority_class_name

    # If specified, all readiness gates will be evaluated for pod readiness. A pod is ready when all its containers are ready AND all conditions specified in the readiness gates have status equal to \"True\" More info: https://git.k8s.io/enhancements/keps/sig-network/0007-pod-ready%2B%2B.md
    attr_accessor :readiness_gates

    # Restart policy for all containers within the pod. One of Always, OnFailure, Never. Default to Always. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy
    attr_accessor :restart_policy

    # RuntimeClassName refers to a RuntimeClass object in the node.k8s.io group, which should be used to run this pod.  If no RuntimeClass resource matches the named class, the pod will not be run. If unset or empty, the \"legacy\" RuntimeClass will be used, which is an implicit class with an empty definition that uses the default runtime handler. More info: https://git.k8s.io/enhancements/keps/sig-node/runtime-class.md This is a beta feature as of Kubernetes v1.14.
    attr_accessor :runtime_class_name

    # If specified, the pod will be dispatched by specified scheduler. If not specified, the pod will be dispatched by default scheduler.
    attr_accessor :scheduler_name

    attr_accessor :security_context

    # DeprecatedServiceAccount is a depreciated alias for ServiceAccountName. Deprecated: Use serviceAccountName instead.
    attr_accessor :service_account

    # ServiceAccountName is the name of the ServiceAccount to use to run this pod. More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
    attr_accessor :service_account_name

    # Share a single process namespace between all of the containers in a pod. When this is set containers will be able to view and signal processes from other containers in the same pod, and the first process in each container will not be assigned PID 1. HostPID and ShareProcessNamespace cannot both be set. Optional: Default to false.
    attr_accessor :share_process_namespace

    # If specified, the fully qualified Pod hostname will be \"<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>\". If not specified, the pod will not have a domainname at all.
    attr_accessor :subdomain

    # Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period will be used instead. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. Defaults to 30 seconds.
    attr_accessor :termination_grace_period_seconds

    # If specified, the pod's tolerations.
    attr_accessor :tolerations

    # TopologySpreadConstraints describes how a group of pods ought to spread across topology domains. Scheduler will schedule pods in a way which abides by the constraints. This field is alpha-level and is only honored by clusters that enables the EvenPodsSpread feature. All topologySpreadConstraints are ANDed.
    attr_accessor :topology_spread_constraints

    # List of volumes that can be mounted by containers belonging to the pod. More info: https://kubernetes.io/docs/concepts/storage/volumes
    attr_accessor :volumes

    # Attribute mapping from ruby-style variable name to JSON key.
    def self.attribute_map
      {
        :'active_deadline_seconds' => :'activeDeadlineSeconds',
        :'affinity' => :'affinity',
        :'automount_service_account_token' => :'automountServiceAccountToken',
        :'containers' => :'containers',
        :'dns_config' => :'dnsConfig',
        :'dns_policy' => :'dnsPolicy',
        :'enable_service_links' => :'enableServiceLinks',
        :'ephemeral_containers' => :'ephemeralContainers',
        :'host_aliases' => :'hostAliases',
        :'host_ipc' => :'hostIPC',
        :'host_network' => :'hostNetwork',
        :'host_pid' => :'hostPID',
        :'hostname' => :'hostname',
        :'image_pull_secrets' => :'imagePullSecrets',
        :'init_containers' => :'initContainers',
        :'node_name' => :'nodeName',
        :'node_selector' => :'nodeSelector',
        :'overhead' => :'overhead',
        :'preemption_policy' => :'preemptionPolicy',
        :'priority' => :'priority',
        :'priority_class_name' => :'priorityClassName',
        :'readiness_gates' => :'readinessGates',
        :'restart_policy' => :'restartPolicy',
        :'runtime_class_name' => :'runtimeClassName',
        :'scheduler_name' => :'schedulerName',
        :'security_context' => :'securityContext',
        :'service_account' => :'serviceAccount',
        :'service_account_name' => :'serviceAccountName',
        :'share_process_namespace' => :'shareProcessNamespace',
        :'subdomain' => :'subdomain',
        :'termination_grace_period_seconds' => :'terminationGracePeriodSeconds',
        :'tolerations' => :'tolerations',
        :'topology_spread_constraints' => :'topologySpreadConstraints',
        :'volumes' => :'volumes'
      }
    end

    # Attribute type mapping.
    def self.openapi_types
      {
        :'active_deadline_seconds' => :'Integer',
        :'affinity' => :'IoK8sApiCoreV1Affinity',
        :'automount_service_account_token' => :'Boolean',
        :'containers' => :'Array<IoK8sApiCoreV1Container>',
        :'dns_config' => :'IoK8sApiCoreV1PodDNSConfig',
        :'dns_policy' => :'String',
        :'enable_service_links' => :'Boolean',
        :'ephemeral_containers' => :'Array<IoK8sApiCoreV1EphemeralContainer>',
        :'host_aliases' => :'Array<IoK8sApiCoreV1HostAlias>',
        :'host_ipc' => :'Boolean',
        :'host_network' => :'Boolean',
        :'host_pid' => :'Boolean',
        :'hostname' => :'String',
        :'image_pull_secrets' => :'Array<IoK8sApiCoreV1LocalObjectReference>',
        :'init_containers' => :'Array<IoK8sApiCoreV1Container>',
        :'node_name' => :'String',
        :'node_selector' => :'Hash<String, String>',
        :'overhead' => :'Hash<String, String>',
        :'preemption_policy' => :'String',
        :'priority' => :'Integer',
        :'priority_class_name' => :'String',
        :'readiness_gates' => :'Array<IoK8sApiCoreV1PodReadinessGate>',
        :'restart_policy' => :'String',
        :'runtime_class_name' => :'String',
        :'scheduler_name' => :'String',
        :'security_context' => :'IoK8sApiCoreV1PodSecurityContext',
        :'service_account' => :'String',
        :'service_account_name' => :'String',
        :'share_process_namespace' => :'Boolean',
        :'subdomain' => :'String',
        :'termination_grace_period_seconds' => :'Integer',
        :'tolerations' => :'Array<IoK8sApiCoreV1Toleration>',
        :'topology_spread_constraints' => :'Array<IoK8sApiCoreV1TopologySpreadConstraint>',
        :'volumes' => :'Array<IoK8sApiCoreV1Volume>'
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
        fail ArgumentError, "The input argument (attributes) must be a hash in `Kubernetes::IoK8sApiCoreV1PodSpec` initialize method"
      end

      # check to see if the attribute exists and convert string to symbol for hash key
      attributes = attributes.each_with_object({}) { |(k, v), h|
        if (!self.class.attribute_map.key?(k.to_sym))
          fail ArgumentError, "`#{k}` is not a valid attribute in `Kubernetes::IoK8sApiCoreV1PodSpec`. Please check the name to make sure it's valid. List of attributes: " + self.class.attribute_map.keys.inspect
        end
        h[k.to_sym] = v
      }

      if attributes.key?(:'active_deadline_seconds')
        self.active_deadline_seconds = attributes[:'active_deadline_seconds']
      end

      if attributes.key?(:'affinity')
        self.affinity = attributes[:'affinity']
      end

      if attributes.key?(:'automount_service_account_token')
        self.automount_service_account_token = attributes[:'automount_service_account_token']
      end

      if attributes.key?(:'containers')
        if (value = attributes[:'containers']).is_a?(Array)
          self.containers = value
        end
      end

      if attributes.key?(:'dns_config')
        self.dns_config = attributes[:'dns_config']
      end

      if attributes.key?(:'dns_policy')
        self.dns_policy = attributes[:'dns_policy']
      end

      if attributes.key?(:'enable_service_links')
        self.enable_service_links = attributes[:'enable_service_links']
      end

      if attributes.key?(:'ephemeral_containers')
        if (value = attributes[:'ephemeral_containers']).is_a?(Array)
          self.ephemeral_containers = value
        end
      end

      if attributes.key?(:'host_aliases')
        if (value = attributes[:'host_aliases']).is_a?(Array)
          self.host_aliases = value
        end
      end

      if attributes.key?(:'host_ipc')
        self.host_ipc = attributes[:'host_ipc']
      end

      if attributes.key?(:'host_network')
        self.host_network = attributes[:'host_network']
      end

      if attributes.key?(:'host_pid')
        self.host_pid = attributes[:'host_pid']
      end

      if attributes.key?(:'hostname')
        self.hostname = attributes[:'hostname']
      end

      if attributes.key?(:'image_pull_secrets')
        if (value = attributes[:'image_pull_secrets']).is_a?(Array)
          self.image_pull_secrets = value
        end
      end

      if attributes.key?(:'init_containers')
        if (value = attributes[:'init_containers']).is_a?(Array)
          self.init_containers = value
        end
      end

      if attributes.key?(:'node_name')
        self.node_name = attributes[:'node_name']
      end

      if attributes.key?(:'node_selector')
        if (value = attributes[:'node_selector']).is_a?(Hash)
          self.node_selector = value
        end
      end

      if attributes.key?(:'overhead')
        if (value = attributes[:'overhead']).is_a?(Hash)
          self.overhead = value
        end
      end

      if attributes.key?(:'preemption_policy')
        self.preemption_policy = attributes[:'preemption_policy']
      end

      if attributes.key?(:'priority')
        self.priority = attributes[:'priority']
      end

      if attributes.key?(:'priority_class_name')
        self.priority_class_name = attributes[:'priority_class_name']
      end

      if attributes.key?(:'readiness_gates')
        if (value = attributes[:'readiness_gates']).is_a?(Array)
          self.readiness_gates = value
        end
      end

      if attributes.key?(:'restart_policy')
        self.restart_policy = attributes[:'restart_policy']
      end

      if attributes.key?(:'runtime_class_name')
        self.runtime_class_name = attributes[:'runtime_class_name']
      end

      if attributes.key?(:'scheduler_name')
        self.scheduler_name = attributes[:'scheduler_name']
      end

      if attributes.key?(:'security_context')
        self.security_context = attributes[:'security_context']
      end

      if attributes.key?(:'service_account')
        self.service_account = attributes[:'service_account']
      end

      if attributes.key?(:'service_account_name')
        self.service_account_name = attributes[:'service_account_name']
      end

      if attributes.key?(:'share_process_namespace')
        self.share_process_namespace = attributes[:'share_process_namespace']
      end

      if attributes.key?(:'subdomain')
        self.subdomain = attributes[:'subdomain']
      end

      if attributes.key?(:'termination_grace_period_seconds')
        self.termination_grace_period_seconds = attributes[:'termination_grace_period_seconds']
      end

      if attributes.key?(:'tolerations')
        if (value = attributes[:'tolerations']).is_a?(Array)
          self.tolerations = value
        end
      end

      if attributes.key?(:'topology_spread_constraints')
        if (value = attributes[:'topology_spread_constraints']).is_a?(Array)
          self.topology_spread_constraints = value
        end
      end

      if attributes.key?(:'volumes')
        if (value = attributes[:'volumes']).is_a?(Array)
          self.volumes = value
        end
      end
    end

    # Show invalid properties with the reasons. Usually used together with valid?
    # @return Array for valid properties with the reasons
    def list_invalid_properties
      invalid_properties = Array.new
      if @containers.nil?
        invalid_properties.push('invalid value for "containers", containers cannot be nil.')
      end

      invalid_properties
    end

    # Check to see if the all the properties in the model are valid
    # @return true if the model is valid
    def valid?
      return false if @containers.nil?
      true
    end

    # Checks equality by comparing each attribute.
    # @param [Object] Object to be compared
    def ==(o)
      return true if self.equal?(o)
      self.class == o.class &&
          active_deadline_seconds == o.active_deadline_seconds &&
          affinity == o.affinity &&
          automount_service_account_token == o.automount_service_account_token &&
          containers == o.containers &&
          dns_config == o.dns_config &&
          dns_policy == o.dns_policy &&
          enable_service_links == o.enable_service_links &&
          ephemeral_containers == o.ephemeral_containers &&
          host_aliases == o.host_aliases &&
          host_ipc == o.host_ipc &&
          host_network == o.host_network &&
          host_pid == o.host_pid &&
          hostname == o.hostname &&
          image_pull_secrets == o.image_pull_secrets &&
          init_containers == o.init_containers &&
          node_name == o.node_name &&
          node_selector == o.node_selector &&
          overhead == o.overhead &&
          preemption_policy == o.preemption_policy &&
          priority == o.priority &&
          priority_class_name == o.priority_class_name &&
          readiness_gates == o.readiness_gates &&
          restart_policy == o.restart_policy &&
          runtime_class_name == o.runtime_class_name &&
          scheduler_name == o.scheduler_name &&
          security_context == o.security_context &&
          service_account == o.service_account &&
          service_account_name == o.service_account_name &&
          share_process_namespace == o.share_process_namespace &&
          subdomain == o.subdomain &&
          termination_grace_period_seconds == o.termination_grace_period_seconds &&
          tolerations == o.tolerations &&
          topology_spread_constraints == o.topology_spread_constraints &&
          volumes == o.volumes
    end

    # @see the `==` method
    # @param [Object] Object to be compared
    def eql?(o)
      self == o
    end

    # Calculates hash code according to all attributes.
    # @return [Integer] Hash code
    def hash
      [active_deadline_seconds, affinity, automount_service_account_token, containers, dns_config, dns_policy, enable_service_links, ephemeral_containers, host_aliases, host_ipc, host_network, host_pid, hostname, image_pull_secrets, init_containers, node_name, node_selector, overhead, preemption_policy, priority, priority_class_name, readiness_gates, restart_policy, runtime_class_name, scheduler_name, security_context, service_account, service_account_name, share_process_namespace, subdomain, termination_grace_period_seconds, tolerations, topology_spread_constraints, volumes].hash
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
