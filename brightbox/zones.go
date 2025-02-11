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
	"context"

	"github.com/brightbox/k8ssdk/v2"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

var (
	emptyZone = cloudprovider.Zone{}
)

// GetZone returns the Zone containing the current failure zone
// and locality region that the program is running in. In most cases,
// this method is called from the kubelet querying a local metadata
// service to acquire its zone.  For the case of external cloud
// providers, use GetZoneByProviderID or GetZoneByNodeName since
// GetZone can no longer be called from the kubelets.
func (c *cloud) GetZone(ctx context.Context) (cloudprovider.Zone, error) {
	if err := logAction(ctx, "GetZone"); err != nil {
		return emptyZone, err
	}
	client, err := c.MetadataClient()
	if err != nil {
		return emptyZone, err
	}
	resp, err := client.GetMetadata("placement/availability-zone")
	if err != nil {
		return emptyZone, err
	}
	return createZone(resp)
}

// Create a Zone object from a zone name string
func createZone(zoneName string) (cloudprovider.Zone, error) {
	respRegion, err := k8ssdk.MapZoneHandleToRegion(zoneName)
	if err != nil {
		return emptyZone, err
	}

	return cloudprovider.Zone{
		FailureDomain: zoneName,
		Region:        respRegion,
	}, err
}

// GetZoneByProviderID returns the Zone containing the current zone
// and locality region of the node specified by providerId This method is
// particularly used in the context of external cloud providers where node
// initialization must be down outside the kubelets.
func (c *cloud) GetZoneByProviderID(ctx context.Context, providerID string) (cloudprovider.Zone, error) {
	if err := logAction(ctx, "GetZoneByProviderID %s", providerID); err != nil {
		return emptyZone, err
	}
	serverID := k8ssdk.MapProviderIDToServerID(providerID)
	return c.getZoneByServerID(ctx, serverID)
}

// GetZoneByNodeName returns the Zone containing the current zone
// and locality region of the node specified by node name This method is
// particularly used in the context of external cloud providers where node
// initialization must be down outside the kubelets.
func (c *cloud) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	serverID := mapNodeNameToServerID(nodeName)
	if err := logAction(ctx, "GetZoneByNodeName %s", serverID); err != nil {
		return emptyZone, err
	}
	return c.getZoneByServerID(ctx, serverID)
}

// Common function that gets the zone via a standard Brightbox serverid
func (c *cloud) getZoneByServerID(ctx context.Context, identifier string) (cloudprovider.Zone, error) {
	client, err := c.CloudClient()
	if err != nil {
		return emptyZone, err
	}
	server, err := client.Server(ctx, identifier)
	if err != nil {
		return emptyZone, err
	}
	zoneName := server.Zone.Handle
	return createZone(zoneName)
}
