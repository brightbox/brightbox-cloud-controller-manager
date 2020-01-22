# Kubernetes::IoK8sApiCoreV1NodeSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**config_source** | [**IoK8sApiCoreV1NodeConfigSource**](IoK8sApiCoreV1NodeConfigSource.md) |  | [optional] 
**external_id** | **String** | Deprecated. Not all kubelets will set this field. Remove field after 1.13. see: https://issues.k8s.io/61966 | [optional] 
**pod_cidr** | **String** | PodCIDR represents the pod IP range assigned to the node. | [optional] 
**pod_cid_rs** | **Array&lt;String&gt;** | podCIDRs represents the IP ranges assigned to the node for usage by Pods on that node. If this field is specified, the 0th entry must match the podCIDR field. It may contain at most 1 value for each of IPv4 and IPv6. | [optional] 
**provider_id** | **String** | ID of the node assigned by the cloud provider in the format: &lt;ProviderName&gt;://&lt;ProviderSpecificNodeID&gt; | [optional] 
**taints** | [**Array&lt;IoK8sApiCoreV1Taint&gt;**](IoK8sApiCoreV1Taint.md) | If specified, the node&#39;s taints. | [optional] 
**unschedulable** | **Boolean** | Unschedulable controls node schedulability of new pods. By default, node is schedulable. More info: https://kubernetes.io/docs/concepts/nodes/node/#manual-node-administration | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1NodeSpec.new(config_source: null,
                                 external_id: null,
                                 pod_cidr: null,
                                 pod_cid_rs: null,
                                 provider_id: null,
                                 taints: null,
                                 unschedulable: null)
```


