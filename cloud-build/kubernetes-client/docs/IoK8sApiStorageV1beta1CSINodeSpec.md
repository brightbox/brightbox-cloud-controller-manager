# Kubernetes::IoK8sApiStorageV1beta1CSINodeSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**drivers** | [**Array&lt;IoK8sApiStorageV1beta1CSINodeDriver&gt;**](IoK8sApiStorageV1beta1CSINodeDriver.md) | drivers is a list of information of all CSI Drivers existing on a node. If all drivers in the list are uninstalled, this can become empty. | 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiStorageV1beta1CSINodeSpec.new(drivers: null)
```


