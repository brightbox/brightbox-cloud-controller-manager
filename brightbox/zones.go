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

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// GetZone returns the Zone containing the current failure zone
// and locality region that the program is running in In most cases,
// this method is called from the kubelet querying a local metadata
// service to acquire its zone.  For the case of external cloud
// providers, use GetZoneByProviderID or GetZoneByNodeName since
// GetZone can no longer be called from the kubelets.
func (c *cloud) GetZone(ctx context.Context) (cloudprovider.Zone, error) {
	client, err := c.metadataClient(ctx)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	resp, err := client.GetMetadata("placement/availability-zone")
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return cloudprovider.Zone{
		FailureDomain: resp,
		Region:        mapZoneHandleToRegion(resp),
	}, err
}

// GetZoneByProviderID returns the Zone containing the current zone
// and locality region of the node specified by providerId This method is
// particularly used in the context of external cloud providers where node
// initialization must be down outside the kubelets.
func (c *cloud) GetZoneByProviderID(ctx context.Context, providerID string) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, nil
}

// GetZoneByNodeName returns the Zone containing the current zone
// and locality region of the node specified by node name This method is
// particularly used in the context of external cloud providers where node
// initialization must be down outside the kubelets.
func (c *cloud) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{}, nil
}

// Obtain a metadata client
func (c *cloud) metadataClient(ctx context.Context) (EC2Metadata, error) {
	if c.metadataClientCache == nil {
		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			return nil, err
		}
		c.metadataClientCache = ec2metadata.New(cfg)
	}

	return c.metadataClientCache, nil
}
