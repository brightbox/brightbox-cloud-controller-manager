// Copyright 2018 Brightbox Systems Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package brightbox

import (
	"io"

	"github.com/brightbox/k8ssdk"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"
)

type cloud struct {
	*k8ssdk.Cloud
}

// Initialize provides the cloud with a kubernetes client builder and
// may spawn goroutines to perform housekeeping activities within the
// cloud provider.
func (c *cloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
	klog.V(4).Infof("Initialise called with %+v", clientBuilder)
}

// LoadBalancer returns a balancer interface. Also returns true if the
// interface is supported, false otherwise.
func (c *cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	klog.V(4).Infof("LoadBalancer called")
	return c, true
}

// Instances returns an instances interface. Also returns true if the
// interface is supported, false otherwise.
func (c *cloud) Instances() (cloudprovider.Instances, bool) {
	klog.V(4).Infof("Instances called")
	return c, true
}

// Zones returns a zones interface. Also returns true if the interface
// is supported, false otherwise.
func (c *cloud) Zones() (cloudprovider.Zones, bool) {
	klog.V(4).Infof("Zones called")
	return c, true
}

// Clusters returns a clusters interface.  Also returns true if the
// interface is supported, false otherwise.
func (c *cloud) Clusters() (cloudprovider.Clusters, bool) {
	klog.V(4).Infof("Clusters called")
	return nil, false
}

// Routes returns a routes interface along with whether the interface
// is supported.
func (c *cloud) Routes() (cloudprovider.Routes, bool) {
	klog.V(4).Infof("Routes called")
	return nil, false
}

// ProviderName returns the cloud provider ID.
func (c *cloud) ProviderName() string {
	klog.V(4).Infof("ProviderName called")
	return k8ssdk.ProviderName
}

// HasClusterID returns true if a ClusterID is required and set
func (c *cloud) HasClusterID() bool {
	klog.V(4).Infof("HasClusterID called")
	return true
}

// Register this provider's creation function with the manager
func init() {
	cloudprovider.RegisterCloudProvider(k8ssdk.ProviderName, newCloudConnection)
}

// Read a config and generate a cloud structure
// Open a cloud connection early in this version to validate environment
// settings.
// TODO: Look at whether open on demand works better
func newCloudConnection(config io.Reader) (cloudprovider.Interface, error) {
	klog.V(4).Infof("newCloudConnection called with %+v", config)
	if config != nil {
		klog.Warningf("supplied config is not read by this version. Using environment")
	}
	newCloud := &cloud{
		&k8ssdk.Cloud{},
	}
	_, err := newCloud.CloudClient()
	if err != nil {
		return nil, err
	}
	return newCloud, nil
}
