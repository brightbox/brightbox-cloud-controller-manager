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
	"testing"

	"github.com/brightbox/k8ssdk"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

func TestInstanceExists(t *testing.T) {
	client := makeFakeInstanceCloudClient()

	exists, err := client.InstanceExists(context.TODO(), mapServerIDToNode(serverExist))
	if err != nil {
		t.Errorf(err.Error())
	} else if !exists {
		t.Errorf("Active: expected Instance to exist")
	}
	exists, err = client.InstanceExists(context.TODO(), mapServerIDToNodeProviderID(serverExist))
	if err != nil {
		t.Errorf(err.Error())
	} else if !exists {
		t.Errorf("Active: expected Instance to exist")
	}
	exists, err = client.InstanceExists(context.TODO(), mapServerIDToNode(serverDeleted))
	if err != nil {
		t.Errorf(err.Error())
	} else if exists {
		t.Errorf("Deleted: expected Instance to not exist")
	}
	exists, err = client.InstanceExists(context.TODO(), mapServerIDToNodeProviderID(serverDeleted))
	if err != nil {
		t.Errorf(err.Error())
	} else if exists {
		t.Errorf("Deleted: expected Instance to not exist")
	}
	exists, err = client.InstanceExists(context.TODO(), mapServerIDToNode(serverMissing))
	if err != nil {
		t.Errorf(err.Error())
	} else if exists {
		t.Errorf("Missing: expected Instance to not exist")
	}
	exists, err = client.InstanceExists(context.TODO(), mapServerIDToNodeProviderID(serverMissing))
	if err != nil {
		t.Errorf(err.Error())
	} else if exists {
		t.Errorf("Missing: expected Instance to not exist")
	}
	exists, err = client.InstanceExists(context.TODO(), mapServerIDToNode(serverShutdown))
	if err != nil {
		t.Errorf(err.Error())
	} else if !exists {
		t.Errorf("Inactive: expected Instance to exist")
	}
	exists, err = client.InstanceExists(context.TODO(), mapServerIDToNodeProviderID(serverShutdown))
	if err != nil {
		t.Errorf(err.Error())
	} else if !exists {
		t.Errorf("Inactive: expected Instance to exist")
	}
	exists, err = client.InstanceExists(context.TODO(), mapServerIDToNode(serverBust))
	if err == nil {
		t.Errorf("expected Instance to fail")
	} else if err == cloudprovider.InstanceNotFound {
		t.Errorf("Got Instance not found error rather than failure")
	}
	exists, err = client.InstanceExists(context.TODO(), mapServerIDToNodeProviderID(serverBust))
	if err == nil {
		t.Errorf("expected Instance to fail")
	} else if err == cloudprovider.InstanceNotFound {
		t.Errorf("Got Instance not found error rather than failure")
	}
}

func TestInstanceShutdown(t *testing.T) {
	client := makeFakeInstanceCloudClient()

	down, err := client.InstanceShutdown(context.TODO(), mapServerIDToNode(serverExist))
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("Active: expected Instance to be not shutdown")
	}
	down, err = client.InstanceShutdown(context.TODO(), mapServerIDToNodeProviderID(serverExist))
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("Active: expected Instance to be not shutdown")
	}
	down, err = client.InstanceShutdown(context.TODO(), mapServerIDToNode(serverDeleted))
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("Deleted: expected Instance to be not shutdown")
	}
	down, err = client.InstanceShutdown(context.TODO(), mapServerIDToNodeProviderID(serverDeleted))
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("Deleted: expected Instance to be not shutdown")
	}
	down, err = client.InstanceShutdown(context.TODO(), mapServerIDToNode(serverMissing))
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("Missing: expected Instance to be not shutdown")
	}
	down, err = client.InstanceShutdown(context.TODO(), mapServerIDToNodeProviderID(serverMissing))
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("Missing: expected Instance to be not shutdown")
	}
	down, err = client.InstanceShutdown(context.TODO(), mapServerIDToNode(serverShutdown))
	if err != nil {
		t.Errorf(err.Error())
	} else if !down {
		t.Errorf("Inactive: expected Instance to be shutdown")
	}
	down, err = client.InstanceShutdown(context.TODO(), mapServerIDToNodeProviderID(serverShutdown))
	if err != nil {
		t.Errorf(err.Error())
	} else if !down {
		t.Errorf("Inactive: expected Instance to be shutdown")
	}
}

func TestInstanceMetadataFailures(t *testing.T) {
	client := makeFakeInstanceCloudClient()

	metadata, err := client.InstanceMetadata(context.TODO(), mapServerIDToNodeProviderID(serverBust))
	if err == nil {
		t.Errorf("Expected error, got %+v", metadata)
	}
}

func TestInstanceMetadata(t *testing.T) {
	client := makeFakeInstanceCloudClient()

	var instanceTests = []struct {
		server                string
		serverType            string
		expectedNodeAddresses []v1.NodeAddress
		zone                  string
		region                string
	}{
		{
			server:                serverExist,
			serverType:            typeHandle,
			expectedNodeAddresses: expectedExistNodeAddresses,
			zone:                  zoneHandle,
			region:                region,
		},
		{
			server:                serverShutdown,
			serverType:            typeHandle2,
			expectedNodeAddresses: expectedShutdownNodeAddresses,
			zone:                  zoneHandle2,
			region:                region,
		},
	}
	for _, example := range instanceTests {
		t.Run(
			example.server,
			func(t *testing.T) {
				metadata, err := client.InstanceMetadata(context.TODO(),
					mapServerIDToNodeProviderID(example.server))
				if err != nil {
					t.Fatalf(err.Error())
				}
				providerID := k8ssdk.MapServerIDToProviderID(example.server)
				if metadata.ProviderID != providerID {
					t.Errorf("Expected Provider ID %s, got %s", providerID, metadata.ProviderID)
				}
				if metadata.InstanceType != example.serverType {
					t.Errorf("Expected Instance Type %s, got %s", example.serverType, metadata.InstanceType)
				}
				if metadata.Zone != example.zone {
					t.Errorf("Expected Zone %s, got %s", example.zone, metadata.Zone)
				}
				if metadata.Region != example.region {
					t.Errorf("Expected Region %s, got %s", example.region, metadata.Region)
				}
				addresses := metadata.NodeAddresses
				lenExpected := len(example.expectedNodeAddresses)
				lenAddresses := len(addresses)
				if lenAddresses != lenExpected {
					t.Errorf("Expected %d items, got %d", lenExpected, lenAddresses)
				}
				for _, expected := range example.expectedNodeAddresses {
					if !containsNodeAddress(addresses, expected) {
						t.Errorf("Expected node is missing: %+v, got %+v", expected, addresses)
					}
				}
			},
		)
	}
}
