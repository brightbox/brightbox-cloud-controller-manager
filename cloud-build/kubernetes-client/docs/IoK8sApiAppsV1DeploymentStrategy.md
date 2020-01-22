# Kubernetes::IoK8sApiAppsV1DeploymentStrategy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**rolling_update** | [**IoK8sApiAppsV1RollingUpdateDeployment**](IoK8sApiAppsV1RollingUpdateDeployment.md) |  | [optional] 
**type** | **String** | Type of deployment. Can be \&quot;Recreate\&quot; or \&quot;RollingUpdate\&quot;. Default is RollingUpdate. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAppsV1DeploymentStrategy.new(rolling_update: null,
                                 type: null)
```


