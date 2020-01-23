# Kubernetes::IoK8sApiCoreV1PodSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**active_deadline_seconds** | **Integer** | Optional duration in seconds the pod may be active on the node relative to StartTime before the system will actively try to mark it failed and kill associated containers. Value must be a positive integer. | [optional] 
**affinity** | [**IoK8sApiCoreV1Affinity**](IoK8sApiCoreV1Affinity.md) |  | [optional] 
**automount_service_account_token** | **Boolean** | AutomountServiceAccountToken indicates whether a service account token should be automatically mounted. | [optional] 
**containers** | [**Array&lt;IoK8sApiCoreV1Container&gt;**](IoK8sApiCoreV1Container.md) | List of containers belonging to the pod. Containers cannot currently be added or removed. There must be at least one container in a Pod. Cannot be updated. | 
**dns_config** | [**IoK8sApiCoreV1PodDNSConfig**](IoK8sApiCoreV1PodDNSConfig.md) |  | [optional] 
**dns_policy** | **String** | Set DNS policy for the pod. Defaults to \&quot;ClusterFirst\&quot;. Valid values are &#39;ClusterFirstWithHostNet&#39;, &#39;ClusterFirst&#39;, &#39;Default&#39; or &#39;None&#39;. DNS parameters given in DNSConfig will be merged with the policy selected with DNSPolicy. To have DNS options set along with hostNetwork, you have to specify DNS policy explicitly to &#39;ClusterFirstWithHostNet&#39;. | [optional] 
**enable_service_links** | **Boolean** | EnableServiceLinks indicates whether information about services should be injected into pod&#39;s environment variables, matching the syntax of Docker links. Optional: Defaults to true. | [optional] 
**ephemeral_containers** | [**Array&lt;IoK8sApiCoreV1EphemeralContainer&gt;**](IoK8sApiCoreV1EphemeralContainer.md) | List of ephemeral containers run in this pod. Ephemeral containers may be run in an existing pod to perform user-initiated actions such as debugging. This list cannot be specified when creating a pod, and it cannot be modified by updating the pod spec. In order to add an ephemeral container to an existing pod, use the pod&#39;s ephemeralcontainers subresource. This field is alpha-level and is only honored by servers that enable the EphemeralContainers feature. | [optional] 
**host_aliases** | [**Array&lt;IoK8sApiCoreV1HostAlias&gt;**](IoK8sApiCoreV1HostAlias.md) | HostAliases is an optional list of hosts and IPs that will be injected into the pod&#39;s hosts file if specified. This is only valid for non-hostNetwork pods. | [optional] 
**host_ipc** | **Boolean** | Use the host&#39;s ipc namespace. Optional: Default to false. | [optional] 
**host_network** | **Boolean** | Host networking requested for this pod. Use the host&#39;s network namespace. If this option is set, the ports that will be used must be specified. Default to false. | [optional] 
**host_pid** | **Boolean** | Use the host&#39;s pid namespace. Optional: Default to false. | [optional] 
**hostname** | **String** | Specifies the hostname of the Pod If not specified, the pod&#39;s hostname will be set to a system-defined value. | [optional] 
**image_pull_secrets** | [**Array&lt;IoK8sApiCoreV1LocalObjectReference&gt;**](IoK8sApiCoreV1LocalObjectReference.md) | ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec. If specified, these secrets will be passed to individual puller implementations for them to use. For example, in the case of docker, only DockerConfig type secrets are honored. More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod | [optional] 
**init_containers** | [**Array&lt;IoK8sApiCoreV1Container&gt;**](IoK8sApiCoreV1Container.md) | List of initialization containers belonging to the pod. Init containers are executed in order prior to containers being started. If any init container fails, the pod is considered to have failed and is handled according to its restartPolicy. The name for an init container or normal container must be unique among all containers. Init containers may not have Lifecycle actions, Readiness probes, Liveness probes, or Startup probes. The resourceRequirements of an init container are taken into account during scheduling by finding the highest request/limit for each resource type, and then using the max of of that value or the sum of the normal containers. Limits are applied to init containers in a similar fashion. Init containers cannot currently be added or removed. Cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/ | [optional] 
**node_name** | **String** | NodeName is a request to schedule this pod onto a specific node. If it is non-empty, the scheduler simply schedules this pod onto that node, assuming that it fits resource requirements. | [optional] 
**node_selector** | **Hash&lt;String, String&gt;** | NodeSelector is a selector which must be true for the pod to fit on a node. Selector which must match a node&#39;s labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/ | [optional] 
**overhead** | **Hash&lt;String, String&gt;** | Overhead represents the resource overhead associated with running a pod for a given RuntimeClass. This field will be autopopulated at admission time by the RuntimeClass admission controller. If the RuntimeClass admission controller is enabled, overhead must not be set in Pod create requests. The RuntimeClass admission controller will reject Pod create requests which have the overhead already set. If RuntimeClass is configured and selected in the PodSpec, Overhead will be set to the value defined in the corresponding RuntimeClass, otherwise it will remain unset and treated as zero. More info: https://git.k8s.io/enhancements/keps/sig-node/20190226-pod-overhead.md This field is alpha-level as of Kubernetes v1.16, and is only honored by servers that enable the PodOverhead feature. | [optional] 
**preemption_policy** | **String** | PreemptionPolicy is the Policy for preempting pods with lower priority. One of Never, PreemptLowerPriority. Defaults to PreemptLowerPriority if unset. This field is alpha-level and is only honored by servers that enable the NonPreemptingPriority feature. | [optional] 
**priority** | **Integer** | The priority value. Various system components use this field to find the priority of the pod. When Priority Admission Controller is enabled, it prevents users from setting this field. The admission controller populates this field from PriorityClassName. The higher the value, the higher the priority. | [optional] 
**priority_class_name** | **String** | If specified, indicates the pod&#39;s priority. \&quot;system-node-critical\&quot; and \&quot;system-cluster-critical\&quot; are two special keywords which indicate the highest priorities with the former being the highest priority. Any other name must be defined by creating a PriorityClass object with that name. If not specified, the pod priority will be default or zero if there is no default. | [optional] 
**readiness_gates** | [**Array&lt;IoK8sApiCoreV1PodReadinessGate&gt;**](IoK8sApiCoreV1PodReadinessGate.md) | If specified, all readiness gates will be evaluated for pod readiness. A pod is ready when all its containers are ready AND all conditions specified in the readiness gates have status equal to \&quot;True\&quot; More info: https://git.k8s.io/enhancements/keps/sig-network/0007-pod-ready%2B%2B.md | [optional] 
**restart_policy** | **String** | Restart policy for all containers within the pod. One of Always, OnFailure, Never. Default to Always. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy | [optional] 
**runtime_class_name** | **String** | RuntimeClassName refers to a RuntimeClass object in the node.k8s.io group, which should be used to run this pod.  If no RuntimeClass resource matches the named class, the pod will not be run. If unset or empty, the \&quot;legacy\&quot; RuntimeClass will be used, which is an implicit class with an empty definition that uses the default runtime handler. More info: https://git.k8s.io/enhancements/keps/sig-node/runtime-class.md This is a beta feature as of Kubernetes v1.14. | [optional] 
**scheduler_name** | **String** | If specified, the pod will be dispatched by specified scheduler. If not specified, the pod will be dispatched by default scheduler. | [optional] 
**security_context** | [**IoK8sApiCoreV1PodSecurityContext**](IoK8sApiCoreV1PodSecurityContext.md) |  | [optional] 
**service_account** | **String** | DeprecatedServiceAccount is a depreciated alias for ServiceAccountName. Deprecated: Use serviceAccountName instead. | [optional] 
**service_account_name** | **String** | ServiceAccountName is the name of the ServiceAccount to use to run this pod. More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/ | [optional] 
**share_process_namespace** | **Boolean** | Share a single process namespace between all of the containers in a pod. When this is set containers will be able to view and signal processes from other containers in the same pod, and the first process in each container will not be assigned PID 1. HostPID and ShareProcessNamespace cannot both be set. Optional: Default to false. | [optional] 
**subdomain** | **String** | If specified, the fully qualified Pod hostname will be \&quot;&lt;hostname&gt;.&lt;subdomain&gt;.&lt;pod namespace&gt;.svc.&lt;cluster domain&gt;\&quot;. If not specified, the pod will not have a domainname at all. | [optional] 
**termination_grace_period_seconds** | **Integer** | Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period will be used instead. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. Defaults to 30 seconds. | [optional] 
**tolerations** | [**Array&lt;IoK8sApiCoreV1Toleration&gt;**](IoK8sApiCoreV1Toleration.md) | If specified, the pod&#39;s tolerations. | [optional] 
**topology_spread_constraints** | [**Array&lt;IoK8sApiCoreV1TopologySpreadConstraint&gt;**](IoK8sApiCoreV1TopologySpreadConstraint.md) | TopologySpreadConstraints describes how a group of pods ought to spread across topology domains. Scheduler will schedule pods in a way which abides by the constraints. This field is alpha-level and is only honored by clusters that enables the EvenPodsSpread feature. All topologySpreadConstraints are ANDed. | [optional] 
**volumes** | [**Array&lt;IoK8sApiCoreV1Volume&gt;**](IoK8sApiCoreV1Volume.md) | List of volumes that can be mounted by containers belonging to the pod. More info: https://kubernetes.io/docs/concepts/storage/volumes | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1PodSpec.new(active_deadline_seconds: null,
                                 affinity: null,
                                 automount_service_account_token: null,
                                 containers: null,
                                 dns_config: null,
                                 dns_policy: null,
                                 enable_service_links: null,
                                 ephemeral_containers: null,
                                 host_aliases: null,
                                 host_ipc: null,
                                 host_network: null,
                                 host_pid: null,
                                 hostname: null,
                                 image_pull_secrets: null,
                                 init_containers: null,
                                 node_name: null,
                                 node_selector: null,
                                 overhead: null,
                                 preemption_policy: null,
                                 priority: null,
                                 priority_class_name: null,
                                 readiness_gates: null,
                                 restart_policy: null,
                                 runtime_class_name: null,
                                 scheduler_name: null,
                                 security_context: null,
                                 service_account: null,
                                 service_account_name: null,
                                 share_process_namespace: null,
                                 subdomain: null,
                                 termination_grace_period_seconds: null,
                                 tolerations: null,
                                 topology_spread_constraints: null,
                                 volumes: null)
```

