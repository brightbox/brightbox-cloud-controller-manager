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

	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

// Initialize provides the cloud with a kubernetes client builder and
// may spawn goroutines to perform housekeeping activities within the
// cloud provider.
func (c *cloud) Initialize(clientBuilder controller.ControllerClientBuilder) {
}

// LoadBalancer returns a balancer interface. Also returns true if the
// interface is supported, false otherwise.
func (c *cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

// Instances returns an instances interface. Also returns true if the
// interface is supported, false otherwise.
func (c *cloud) Instances() (cloudprovider.Instances, bool) {
	return c, true
}

// Zones returns a zones interface. Also returns true if the interface
// is supported, false otherwise.
func (c *cloud) Zones() (cloudprovider.Zones, bool) {
	return c, true
}

// Clusters returns a clusters interface.  Also returns true if the
// interface is supported, false otherwise.
func (c *cloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

// Routes returns a routes interface along with whether the interface
// is supported.
func (c *cloud) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

// ProviderName returns the cloud provider ID.
func (c *cloud) ProviderName() string {
	return providerName
}

// HasClusterID returns true if a ClusterID is required and set
func (c *cloud) HasClusterID() bool {
	return true
}

// Register this provider's creation function with the manager
func init() {
	cloudprovider.RegisterCloudProvider(providerName, newCloudConnection)
}

//
func newCloudConnection(config io.Reader) (cloudprovider.Interface, error) {
	return &cloud{}, nil
}
