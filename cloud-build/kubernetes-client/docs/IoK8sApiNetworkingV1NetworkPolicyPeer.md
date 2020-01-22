# Kubernetes::IoK8sApiNetworkingV1NetworkPolicyPeer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ip_block** | [**IoK8sApiNetworkingV1IPBlock**](IoK8sApiNetworkingV1IPBlock.md) |  | [optional] 
**namespace_selector** | [**IoK8sApimachineryPkgApisMetaV1LabelSelector**](IoK8sApimachineryPkgApisMetaV1LabelSelector.md) |  | [optional] 
**pod_selector** | [**IoK8sApimachineryPkgApisMetaV1LabelSelector**](IoK8sApimachineryPkgApisMetaV1LabelSelector.md) |  | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiNetworkingV1NetworkPolicyPeer.new(ip_block: null,
                                 namespace_selector: null,
                                 pod_selector: null)
```


