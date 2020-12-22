// Copyright 2020 Brightbox Systems Ltd
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
	"context"

	"github.com/brightbox/k8ssdk"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

// InstanceExists returns true if the instance for the given node exists according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (c *cloud) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	klog.V(4).Infof("InstanceExists (%q)", node.Spec.ProviderID)
	return c.InstanceExistsByProviderID(ctx, mapNodeToProviderID(node))
}

// InstanceShutdown returns true if the instance is shutdown according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (c *cloud) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	klog.V(4).Infof("InstanceShutdown (%q)", node.Spec.ProviderID)
	return c.InstanceShutdownByProviderID(ctx, mapNodeToProviderID(node))
}

// InstanceMetadata returns the instance's metadata. The values returned in InstanceMetadata are
// translated into specific fields and labels in the Node object on registration.
// Implementations should always check node.spec.providerID first when trying to discover the instance
// for a given node. In cases where node.spec.providerID is empty, implementations can use other
// properties of the node like its name, labels and annotations.
func (c *cloud) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {
	klog.V(4).Infof("InstanceMetadata (%q)", node.Spec.ProviderID)
	srv, err := c.GetServer(ctx, mapNodeToServerID(node), cloudprovider.InstanceNotFound)
	if err != nil {
		return nil, err
	}
	addresses, err := nodeAddressesFromServer(srv)
	if err != nil {
		return nil, err
	}
	region, err := k8ssdk.MapZoneHandleToRegion(srv.Zone.Handle)
	if err != nil {
		return nil, err
	}
	return &cloudprovider.InstanceMetadata{
		ProviderID:    k8ssdk.MapServerIDToProviderID(srv.Id),
		InstanceType:  srv.ServerType.Handle,
		NodeAddresses: addresses,
		Zone:          srv.Zone.Handle,
		Region:        region,
	}, nil
}
