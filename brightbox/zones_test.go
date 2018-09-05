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
	"fmt"
	"testing"

	"github.com/brightbox/gobrightbox"
	"k8s.io/apimachinery/pkg/types"
)

const (
	dodgyServer = "srv-dodgy"
)

func interfaceGetZone(zoneName string) func(*testing.T) {
	return func(t *testing.T) {
		client := &cloud{
			metadataClientCache: fakeZoneMetadataClient(zoneName),
		}
		zone, err := client.GetZone(context.TODO())
		if err != nil {
			if zoneName != "" && zoneName != "dummy" {
				t.Errorf("Failed to obtain zone: %s", err.Error())
			}
			if zone != emptyZone {
				t.Errorf("Unexpected zone return on error, got %+v", zone)
			}
			// zoneName is blank triggers an expected metadata failure
			// zoneName of dummy triggers expected metadata failure
		} else {
			if zone.FailureDomain != zoneName {
				t.Errorf("Expected %v, got %v", zoneName, zone.FailureDomain)
			}
			testRegion, err := mapZoneHandleToRegion(zoneName)
			if err != nil {
				t.Errorf(err.Error())
			}
			if zone.Region != testRegion {
				t.Errorf("Expected %v, got %v", testRegion, zone.Region)
			}
		}
	}
}

func interfaceGetZoneByProviderID(ProviderName string, zoneName string) func(*testing.T) {
	return func(t *testing.T) {
		client := &cloud{
			client: fakeZoneCloudClient(context.TODO()),
		}
		zone, err := client.GetZoneByProviderID(context.TODO(), ProviderName)
		if err != nil {
			if ProviderName != providerPrefix+dodgyServer {
				t.Errorf("Failed to obtain zone: %s", err.Error())
			}
			// dodgy providerName should fail
		} else {
			if zone.FailureDomain != zoneName {
				t.Errorf("Expected %v, got %v", zoneName, zone.FailureDomain)
			}
			testRegion, err := mapZoneHandleToRegion(zoneName)
			if err != nil {
				t.Errorf(err.Error())
			}
			if zone.Region != testRegion {
				t.Errorf("Expected %v, got %v", testRegion, zone.Region)
			}
		}
	}
}

func interfaceGetZoneByNodeName(NodeName types.NodeName, zoneName string) func(*testing.T) {
	return func(t *testing.T) {
		client := &cloud{
			client: fakeZoneCloudClient(context.TODO()),
		}
		zone, err := client.GetZoneByNodeName(context.TODO(), NodeName)
		if err != nil {
			if NodeName != dodgyServer {
				t.Errorf("Failed to obtain zone: %s", err.Error())
			}
			// dodgy providerName should fail
		} else {
			if zone.FailureDomain != zoneName {
				t.Errorf("Expected %v, got %v", zoneName, zone.FailureDomain)
			}
			testRegion, err := mapZoneHandleToRegion(zoneName)
			if err != nil {
				t.Errorf(err.Error())
			}
			if zone.Region != testRegion {
				t.Errorf("Expected %v, got %v", testRegion, zone.Region)
			}
		}
	}
}

func TestGetZone(t *testing.T) {
	testCases := []string{"", "dummy", "gb1s-a", "gb1s-b", "gb1-a", "gb1-b"}
	for _, tc := range testCases {
		t.Run(tc, interfaceGetZone(tc))
	}
}

func TestGetZoneByProviderID(t *testing.T) {
	testCases := map[string]string{
		"brightbox://srv-testy":      "gb1-a",
		"brightbox://srv-teste":      "gb1-b",
		providerPrefix + dodgyServer: "",
	}
	for name, zone := range testCases {
		t.Run(name, interfaceGetZoneByProviderID(name, zone))
	}
}

func TestGetZoneByNodeName(t *testing.T) {
	testCases := map[string]string{
		"srv-testy": "gb1-a",
		"srv-teste": "gb1-b",
		dodgyServer: "",
	}
	for name, zone := range testCases {
		node := types.NodeName(name)
		t.Run(name, interfaceGetZoneByNodeName(node, zone))
	}
}

func TestGetZoneCloudClientFailure(t *testing.T) {
	resetAuthEnvironment()
	defer resetAuthEnvironment()
	client := &cloud{}
	zone, err := client.GetZoneByNodeName(context.TODO(), types.NodeName("srv-duffy"))
	if err == nil {
		t.Errorf("Expected error")
	} else if zone != emptyZone {
		t.Errorf("Expected empty zone, got %+v", zone)
	}
}

type fakeZoneMetadata struct {
	fail bool
	zone string
}

func fakeZoneMetadataClient(zoneName string) *fakeZoneMetadata {
	return &fakeZoneMetadata{
		fail: zoneName == "",
		zone: zoneName,
	}
}

func (f *fakeZoneMetadata) GetMetadata(target string) (string, error) {
	if f.fail {
		return "", fmt.Errorf("metadata deactivated")
	}
	if target != "placement/availability-zone" {
		return "", fmt.Errorf("Incorrect metadata requested")
	}
	return f.zone, nil
}

type fakeZoneCloud struct {
	CloudAccess
	serverzone map[string]string
}

func fakeZoneCloudClient(ctx context.Context) *fakeZoneCloud {
	return &fakeZoneCloud{
		serverzone: map[string]string{"srv-testy": "gb1-a",
			"srv-teste": "gb1-b",
			"srv-testa": "gb1s-a",
			"srv-testb": "gb1s-b"},
	}
}

func (f *fakeZoneCloud) Server(identifier string) (*brightbox.Server, error) {
	result := f.serverzone[identifier]
	if result == "" {
		return nil, brightbox.ApiError{
			StatusCode: 404,
			Status:     "404 Not Found",
		}
	} else {
		return &brightbox.Server{
			Id: identifier,
			Zone: brightbox.Zone{
				Id:     "typ-testy",
				Handle: result,
			},
		}, nil
	}
}
