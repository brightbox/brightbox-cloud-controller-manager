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
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/brightbox/gobrightbox"
)

// EC2Metadata is an abstraction over the AWS metadata service.
type EC2Metadata interface {
	// Query the EC2 metadata service (used to discover instance-id etc)
	GetMetadata(path string) (string, error)
}

// CloudAccess is an abstraction over the Brightbox API to allow testing
type CloudAccess interface {
	//Fetch a server
	Server(identifier string) (*brightbox.Server, error)
}

type cloud struct {
	client              CloudAccess
	metadataClientCache EC2Metadata
}

// Obtain a metadata client
func (c *cloud) metadataClient() (EC2Metadata, error) {
	if c.metadataClientCache == nil {
		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			return nil, err
		}
		c.metadataClientCache = ec2metadata.New(cfg)
	}

	return c.metadataClientCache, nil
}

// Obtain a cloud client
func (c *cloud) cloudClient() (CloudAccess, error) {
	if c.client == nil {
		client, err := brightbox.NewClient("", "", nil)
		if err != nil {
			return nil, err
		}
		c.client = client
	}

	return c.client, nil
}
