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

	"github.com/brightbox/gobrightbox"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

const (
	serverExist = "srv-exist"
	serverMissing = "srv-missy"
	serverShutdown = "srv-downy"
	zoneHandle = "gb1s-a"

)

func TestCurrentNodeName(t *testing.T) {
	const server = "srv-testy"
	client := &cloud{
		client: fakeInstanceCloudClient(context.TODO()),
	}
	nodeName, err := client.CurrentNodeName(context.TODO(), server)
	if err != nil {
		t.Errorf(err.Error())
	}
	if nodeName != types.NodeName(server) {
		t.Errorf("Nodename does not match %q", server)
	}
}

func TestAddSSHKey(t *testing.T) {
	client := &cloud{
		client: fakeInstanceCloudClient(context.TODO()),
	}
	err := client.AddSSHKeyToAllInstances(context.TODO(), "fred", []byte("hello"))
	if err == nil {
		t.Errorf("Add SSH should be unimplemented")
	}
}

func TestNodeNameChecks(t *testing.T) {
	client := &cloud{
		client: fakeInstanceCloudClient(context.TODO()),
	}
	var instance_tests = []struct{
		name string
		fn func(* testing.T)
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
				zoneHandle),
			},
		{
			"InstanceTypeByProviderID",
			providerIdTestFactory(client.InstanceTypeByProviderID,
				providerPrefix + serverExist,
				providerPrefix + serverMissing,
				zoneHandle),
			},
		}
	for _, example := range instance_tests {
		t.Run(example.name, example.fn)
	}
}

func TestInstanceExistsByProviderID(t *testing.T) {
	client := &cloud{
		client: fakeInstanceCloudClient(context.TODO()),
	}
	exists, err := client.InstanceExistsByProviderID(context.TODO(), providerPrefix + serverExist)
	if err != nil {
		t.Errorf(err.Error())
	} else if !exists {
		t.Errorf("expected Instance to exist")
	}
	exists, err = client.InstanceExistsByProviderID(context.TODO(), providerPrefix + serverMissing)
	if err != nil {
		t.Errorf(err.Error())
	} else if exists {
		t.Errorf("expected Instance to be missing")
	}
	exists, err = client.InstanceExistsByProviderID(context.TODO(), providerPrefix + serverShutdown)
	if err != nil {
		t.Errorf(err.Error())
	} else if exists {
		t.Errorf("expected Instance to be missing")
	}
}

func TestInstanceShutdownByProviderID(t *testing.T) {
	client := &cloud{
		client: fakeInstanceCloudClient(context.TODO()),
	}
	down, err := client.InstanceShutdownByProviderID(context.TODO(), providerPrefix + serverExist)
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("expected Instance to be active not down")
	}
	down, err = client.InstanceShutdownByProviderID(context.TODO(), providerPrefix + serverMissing)
	if err != nil {
		t.Errorf(err.Error())
	} else if down {
		t.Errorf("expected Instance to be missing not down")
	}
	down, err = client.InstanceShutdownByProviderID(context.TODO(), providerPrefix + serverShutdown)
	if err != nil {
		t.Errorf(err.Error())
	} else if !down {
		t.Errorf("expected Instance to be down")
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

type fakeInstanceCloud struct {
}

func fakeInstanceCloudClient(ctx context.Context) *fakeInstanceCloud {
	return &fakeInstanceCloud{}
}

func (f *fakeInstanceCloud) Server(identifier string) (*brightbox.Server, error) {
	switch identifier {
	case serverExist:
		return &brightbox.Server{
			Id: identifier,
			Status: "active",
			Zone: brightbox.Zone{
				Id: "zon-testy",
				Handle: zoneHandle,
			},
		}, nil
	case serverShutdown:
		return &brightbox.Server{
			Id: identifier,
			Status: "inactive",
			Zone: brightbox.Zone{
				Id: "zon-testy",
				Handle: zoneHandle,
			},
		}, nil
	default:
		return nil, brightbox.ApiError{
			StatusCode: 404,
			Status:     "404 Not Found",
		}
	}
}
