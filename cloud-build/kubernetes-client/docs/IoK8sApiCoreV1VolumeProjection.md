# Kubernetes::IoK8sApiCoreV1VolumeProjection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**config_map** | [**IoK8sApiCoreV1ConfigMapProjection**](IoK8sApiCoreV1ConfigMapProjection.md) |  | [optional] 
**downward_api** | [**IoK8sApiCoreV1DownwardAPIProjection**](IoK8sApiCoreV1DownwardAPIProjection.md) |  | [optional] 
**secret** | [**IoK8sApiCoreV1SecretProjection**](IoK8sApiCoreV1SecretProjection.md) |  | [optional] 
**service_account_token** | [**IoK8sApiCoreV1ServiceAccountTokenProjection**](IoK8sApiCoreV1ServiceAccountTokenProjection.md) |  | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiCoreV1VolumeProjection.new(config_map: null,
                                 downward_api: null,
                                 secret: null,
                                 service_account_token: null)
```


