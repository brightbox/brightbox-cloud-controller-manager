# Kubernetes::IoK8sApiCoreV1Affinity

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**node_affinity** | [**IoK8sApiCoreV1NodeAffinity**](IoK8sApiCoreV1NodeAffinity.md) |  | [optional] 
**pod_affinity** | [**IoK8sApiCoreV1PodAffinity**](IoK8sApiCoreV1PodAffinity.md) |  | [optional] 
**pod_anti_affinity** | [**IoK8sApiCoreV1PodAntiAffinity**](IoK8sApiCoreV1PodAntiAffinity.md) |  | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1Affinity.new(node_affinity: null,
                                 pod_affinity: null,
                                 pod_anti_affinity: null)
```


