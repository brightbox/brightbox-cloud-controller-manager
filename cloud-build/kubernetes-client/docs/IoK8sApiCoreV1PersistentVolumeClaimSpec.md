# Kubernetes::IoK8sApiCoreV1PersistentVolumeClaimSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**access_modes** | **Array&lt;String&gt;** | AccessModes contains the desired access modes the volume should have. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1 | [optional] 
**data_source** | [**IoK8sApiCoreV1TypedLocalObjectReference**](IoK8sApiCoreV1TypedLocalObjectReference.md) |  | [optional] 
**resources** | [**IoK8sApiCoreV1ResourceRequirements**](IoK8sApiCoreV1ResourceRequirements.md) |  | [optional] 
**selector** | [**IoK8sApimachineryPkgApisMetaV1LabelSelector**](IoK8sApimachineryPkgApisMetaV1LabelSelector.md) |  | [optional] 
**storage_class_name** | **String** | Name of the StorageClass required by the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1 | [optional] 
**volume_mode** | **String** | volumeMode defines what type of volume is required by the claim. Value of Filesystem is implied when not included in claim spec. This is a beta feature. | [optional] 
**volume_name** | **String** | VolumeName is the binding reference to the PersistentVolume backing this claim. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1PersistentVolumeClaimSpec.new(access_modes: null,
                                 data_source: null,
                                 resources: null,
                                 selector: null,
                                 storage_class_name: null,
                                 volume_mode: null,
                                 volume_name: null)
```

