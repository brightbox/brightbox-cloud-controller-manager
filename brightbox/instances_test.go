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
	"testing"

	brightbox "github.com/brightbox/gobrightbox/v2"
	"github.com/brightbox/gobrightbox/v2/enums/serverstatus"
	"github.com/brightbox/k8ssdk/v2"
	"github.com/brightbox/k8ssdk/v2/mocks"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

const (
	serverDeleted              = "srv-goner"
	serverExist                = "srv-exist"
	serverMissing              = "srv-missy"
	serverShutdown             = "srv-downy"
	serverBust                 = "srv-busty"
	region                     = "gb1s"
	zoneHandle                 = region + "-a"
	zoneHandle2                = region + "-b"
	typeHandle                 = "1gb.ssd"
	typeHandle2                = "2gb.ssd"
	regionRoot                 = ".brightbox.com"
	serverExistIP              = "81.15.16.17"
	serverExistIPv6            = "64:ff9b::510f:1011"
	serverShutdownIP           = "81.15.16.21"
	serverShutdownIPv6         = "64:ff9b::510f:1015"
	serverShutdownExternalIP   = "109.107.50.0"
	serverShutdownExternalName = "cip-k4a25"
	serverDodgy4               = "srv-dodg4"
	serverDodgy6               = "srv-dodg6"
	serverDodgyCIP             = "srv-dodgc"
	serverDodgyIPv6            = "bust::edfe"
	serverDodgyIPv4            = "::ffff:256.156.256.256"
	serverDodgyCIPv4           = "300.30.300.30"
)

var (
	domain                     = currentDomain()
	expectedExistNodeAddresses = []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: serverExist,
		},
		{
			Type:    v1.NodeExternalDNS,
			Address: serverExist + "." + domain,
		},
		{
			Type:    v1.NodeExternalIP,
			Address: serverExistIP,
		},
		{
			Type:    v1.NodeExternalIP,
			Address: serverExistIPv6,
		},
	}
	expectedShutdownNodeAddresses = []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: serverShutdown,
		},
		{
			Type:    v1.NodeExternalDNS,
			Address: serverShutdown + "." + domain,
		},
		{
			Type:    v1.NodeExternalIP,
			Address: serverShutdownIP,
		},
		{
			Type:    v1.NodeExternalIP,
			Address: serverShutdownExternalIP,
		},
		{
			Type:    v1.NodeExternalDNS,
			Address: serverShutdownExternalName + "." + domain,
		},
		{
			Type:    v1.NodeExternalIP,
			Address: serverShutdownIPv6,
		},
	}
)

func TestCurrentNodeName(t *testing.T) {
	const server = "srv-testy"
	client := makeFakeInstanceCloudClient()
	nodeName, err := client.CurrentNodeName(context.TODO(), server)
	if err != nil {
		t.Errorf(err.Error())
	}
	if nodeName != types.NodeName(server) {
		t.Errorf("Nodename does not match %q", server)
	}
}

func TestAddSSHKey(t *testing.T) {
	client := makeFakeInstanceCloudClient()

	err := client.AddSSHKeyToAllInstances(context.TODO(), "fred", []byte("hello"))
	if err == nil {
		t.Errorf("Add SSH should be unimplemented")
	}
}

func TestNodeNameChecks(t *testing.T) {
	client := makeFakeInstanceCloudClient()

	var instanceTests = []struct {
		name string
		fn   func(*testing.T)
	}{
		{
			"ExternalID",
			nodeNameTestFactory(client.ExternalID,
				serverExist,
				serverMissing,
				serverExist),
		},
		{
			"InstanceID",
			nodeNameTestFactory(client.InstanceID,
				serverExist,
				serverMissing,
				serverExist),
		},
		{
			"InstanceType",
			nodeNameTestFactory(client.InstanceType,
				serverExist,
				serverMissing,
				typeHandle),
		},
		{
			"InstanceTypeByProviderID",
			providerIdTestFactory(client.InstanceTypeByProviderID,
				k8ssdk.MapServerIDToProviderID(serverExist),
				k8ssdk.MapServerIDToProviderID(serverMissing),
				typeHandle),
		},
	}
	for _, example := range instanceTests {
		t.Run(example.name, example.fn)
	}
}

func TestNodeAddresses(t *testing.T) {
	client := makeFakeInstanceCloudClient()

	var instanceTests = []struct {
		server                types.NodeName
		expectedNodeAddresses []v1.NodeAddress
	}{
		{
			server:                serverExist,
			expectedNodeAddresses: expectedExistNodeAddresses,
		},
		{
			server:                serverShutdown,
			expectedNodeAddresses: expectedShutdownNodeAddresses,
		},
	}
	for _, example := range instanceTests {
		t.Run(
			mapNodeNameToServerID(example.server),
			func(t *testing.T) {
				addresses, err := client.NodeAddresses(context.TODO(), example.server)
				if err != nil {
					t.Fatalf(err.Error())
				}
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

func TestNodeAddressesByProviderID(t *testing.T) {
	client := makeFakeInstanceCloudClient()

	addresses, err := client.NodeAddressesByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverBust))
	if err == nil {
		t.Errorf("Expected error, got %+v", addresses)
	}
	addresses, err = client.NodeAddressesByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverDodgy4))
	if err == nil {
		t.Errorf("Expected error, got %+v", addresses)
	}
	addresses, err = client.NodeAddressesByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverDodgy6))
	if err == nil {
		t.Errorf("Expected error, got %+v", addresses)
	}
	addresses, err = client.NodeAddressesByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverDodgyCIP))
	if err == nil {
		t.Errorf("Expected error, got %+v", addresses)
	}
}

func TestInstanceExistsByProviderID(t *testing.T) {
	client := makeFakeInstanceCloudClient()

	exists, err := client.InstanceExistsByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverExist))
	if err != nil {
		t.Errorf(err.Error())
	} else if !exists {
		t.Errorf("Active: expected Instance to exist")
	}
	exists, err = client.InstanceExistsByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverDeleted))
	if err != nil {
		t.Errorf(err.Error())
	} else if exists {
		t.Errorf("Deleted: expected Instance to not exist")
	}
	exists, err = client.InstanceExistsByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverMissing))
	if err != nil {
		t.Errorf(err.Error())
	} else if exists {
		t.Errorf("Missing: expected Instance to not exist")
	}
	exists, err = client.InstanceExistsByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverShutdown))
	if err != nil {
		t.Errorf(err.Error())
	} else if !exists {
		t.Errorf("Inactive: expected Instance to exist")
	}
	exists, err = client.InstanceExistsByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverBust))
	if err == nil {
		t.Errorf("expected Instance to fail")
	} else if err == cloudprovider.InstanceNotFound {
		t.Errorf("Got Instance not found error rather than failure")
	}
}

func TestInstanceShutdownByProviderID(t *testing.T) {
	client := makeFakeInstanceCloudClient()

	down, err := client.InstanceShutdownByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverExist))
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("Active: expected Instance to be not shutdown")
	}
	down, err = client.InstanceShutdownByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverDeleted))
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("Deleted: expected Instance to be not shutdown")
	}
	down, err = client.InstanceShutdownByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverMissing))
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("Missing: expected Instance to be not shutdown")
	}
	down, err = client.InstanceShutdownByProviderID(context.TODO(), k8ssdk.MapServerIDToProviderID(serverShutdown))
	if err != nil {
		t.Errorf(err.Error())
	} else if !down {
		t.Errorf("Inactive: expected Instance to be shutdown")
	}
}

func TestGetInstanceCloudClientFailure(t *testing.T) {
	k8ssdk.ResetAuthEnvironment()
	defer k8ssdk.ResetAuthEnvironment()
	client := makeFakeCloudClient()
	instance, err := client.InstanceID(context.TODO(), types.NodeName("srv-duffy"))
	if err == nil {
		t.Errorf("Expected error")
	} else if instance != "" {
		t.Errorf("Expected empty instance, got %+v", instance)
	}
}

func nodeNameTestFactory(testFunction func(context.Context, types.NodeName) (string, error), sourceExist types.NodeName, sourceMissing types.NodeName, expected string) func(*testing.T) {
	return func(t *testing.T) {
		id, err := testFunction(context.TODO(), sourceExist)
		if err != nil {
			t.Errorf(err.Error())
		} else if id != expected {
			t.Errorf("expected %q, got %q", expected, id)
		}
		id, err = testFunction(context.TODO(), sourceMissing)
		if err == nil {
			t.Errorf("expected not found, got %q", id)
		} else if err != cloudprovider.InstanceNotFound {
			t.Errorf("expected not found, got %+v", err)
		}
	}
}

func providerIdTestFactory(testFunction func(context.Context, string) (string, error), sourceExist string, sourceMissing string, expected string) func(*testing.T) {
	return func(t *testing.T) {
		id, err := testFunction(context.TODO(), sourceExist)
		if err != nil {
			t.Errorf(err.Error())
		} else if id != expected {
			t.Errorf("expected %q, got %q", expected, id)
		}
		id, err = testFunction(context.TODO(), sourceMissing)
		if err == nil {
			t.Errorf("expected not found, got %q", id)
		} else if err != cloudprovider.InstanceNotFound {
			t.Errorf("expected not found, got %+v", err)
		}
	}
}

func currentDomain() string {
	region, err := k8ssdk.MapZoneHandleToRegion(zoneHandle)
	if err != nil {
		return ""
	}
	return region + regionRoot
}

type fakeInstanceCloud struct {
	mocks.CloudAccess
}

func fakeInstanceCloudClient(ctx context.Context) *fakeInstanceCloud {
	return &fakeInstanceCloud{}
}

func (f *fakeInstanceCloud) Server(_ context.Context, identifier string) (*brightbox.Server, error) {
	region, err := k8ssdk.MapZoneHandleToRegion(zoneHandle)
	if err != nil {
		return nil, err
	}
	domain := region + ".brightbox.com"
	switch identifier {
	case serverExist:
		return &brightbox.Server{
			ID:       identifier,
			Status:   serverstatus.Active,
			Hostname: serverExist,
			Fqdn:     serverExist + "." + domain,
			Zone: &brightbox.Zone{
				ID:     "zon-testy",
				Handle: zoneHandle,
			},
			ServerType: &brightbox.ServerType{
				ID:     "typ-8985i",
				Handle: typeHandle,
			},
			Interfaces: []brightbox.Interface{
				{
					ID:          "int-ds42k",
					MacAddress:  "02:24:19:00:00:ee",
					IPv4Address: serverExistIP,
					IPv6Address: serverExistIPv6,
				},
			},
			CloudIPs: []brightbox.CloudIP{},
		}, nil
	case serverShutdown:
		return &brightbox.Server{
			ID:       identifier,
			Status:   serverstatus.Inactive,
			Hostname: serverShutdown,
			Fqdn:     serverShutdown + "." + domain,
			Zone: &brightbox.Zone{
				ID:     "zon-testy",
				Handle: zoneHandle2,
			},
			ServerType: &brightbox.ServerType{
				ID:     "typ-wusvn",
				Handle: typeHandle2,
			},
			Interfaces: []brightbox.Interface{
				{
					ID:          "int-ds42l",
					MacAddress:  "02:24:19:00:00:ef",
					IPv4Address: serverShutdownIP,
					IPv6Address: serverShutdownIPv6,
				},
			},
			CloudIPs: []brightbox.CloudIP{
				{
					ID:         serverShutdownExternalName,
					PublicIP:   serverShutdownExternalIP,
					Fqdn:       serverShutdownExternalName + "." + domain,
					ReverseDNS: "",
				},
			},
		}, nil
	case serverDeleted:
		return &brightbox.Server{
			ID:       identifier,
			Status:   serverstatus.Deleted,
			Hostname: serverDeleted,
			Fqdn:     serverDeleted + "." + domain,
			Zone: &brightbox.Zone{
				ID:     "zon-testy",
				Handle: zoneHandle,
			},
			ServerType: &brightbox.ServerType{
				ID:     "typ-wusvn",
				Handle: typeHandle2,
			},
			Interfaces: []brightbox.Interface{
				{
					ID:          "int-ds42l",
					MacAddress:  "02:24:19:00:00:ef",
					IPv4Address: serverShutdownIP,
					IPv6Address: serverShutdownIPv6,
				},
			},
			CloudIPs: []brightbox.CloudIP{
				{
					ID:         serverShutdownExternalName,
					PublicIP:   serverShutdownExternalIP,
					Fqdn:       serverShutdownExternalName + "." + domain,
					ReverseDNS: "",
				},
			},
		}, nil
	case serverBust:
		return nil, &brightbox.APIError{
			StatusCode: 500,
			Status:     "Internal Server Error",
		}
	case serverDodgy4:
		return &brightbox.Server{
			ID:       identifier,
			Status:   serverstatus.Active,
			Hostname: serverDodgy4,
			Fqdn:     serverDodgy4 + "." + domain,
			Zone: &brightbox.Zone{
				ID:     "zon-testy",
				Handle: zoneHandle,
			},
			ServerType: &brightbox.ServerType{
				ID:     "typ-wusvn",
				Handle: typeHandle2,
			},
			Interfaces: []brightbox.Interface{
				{
					ID:          "int-ds42k",
					MacAddress:  "02:24:19:00:00:ee",
					IPv4Address: serverDodgyIPv4,
					IPv6Address: serverExistIPv6,
				},
			},
			CloudIPs: []brightbox.CloudIP{},
		}, nil
	case serverDodgy6:
		return &brightbox.Server{
			ID:       identifier,
			Status:   serverstatus.Active,
			Hostname: serverDodgy6,
			Fqdn:     serverDodgy6 + "." + domain,
			Zone: &brightbox.Zone{
				ID:     "zon-testy",
				Handle: zoneHandle,
			},
			ServerType: &brightbox.ServerType{
				ID:     "typ-wusvn",
				Handle: typeHandle2,
			},
			Interfaces: []brightbox.Interface{
				{
					ID:          "int-ds42k",
					MacAddress:  "02:24:19:00:00:ee",
					IPv4Address: serverExistIP,
					IPv6Address: serverDodgyIPv6,
				},
			},
			CloudIPs: []brightbox.CloudIP{},
		}, nil

	case serverDodgyCIP:
		return &brightbox.Server{
			ID:       identifier,
			Status:   serverstatus.Inactive,
			Hostname: serverShutdown,
			Fqdn:     serverShutdown + "." + domain,
			Zone: &brightbox.Zone{
				ID:     "zon-testy",
				Handle: zoneHandle,
			},
			ServerType: &brightbox.ServerType{
				ID:     "typ-wusvn",
				Handle: typeHandle2,
			},
			Interfaces: []brightbox.Interface{
				{
					ID:          "int-ds42l",
					MacAddress:  "02:24:19:00:00:ef",
					IPv4Address: serverShutdownIP,
					IPv6Address: serverShutdownIPv6,
				},
			},
			CloudIPs: []brightbox.CloudIP{
				{
					ID:         serverShutdownExternalName,
					PublicIP:   serverDodgyCIPv4,
					Fqdn:       serverShutdownExternalName + "." + domain,
					ReverseDNS: "",
				},
			},
		}, nil
	default:
		return nil, &brightbox.APIError{
			StatusCode: 404,
			Status:     "404 Not Found",
		}
	}
}

func containsNodeAddress(list []v1.NodeAddress, item v1.NodeAddress) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func makeFakeInstanceCloudClient() *cloud {
	return &cloud{
		k8ssdk.MakeTestClient(
			fakeInstanceCloudClient(context.TODO()),
			nil,
		),
	}
}

func makeFakeCloudClient() *cloud {
	return &cloud{
		k8ssdk.MakeTestClient(nil, nil),
	}
}
