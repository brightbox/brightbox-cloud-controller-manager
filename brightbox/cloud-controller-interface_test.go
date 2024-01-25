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
	"strings"
	"testing"

	"github.com/brightbox/k8ssdk"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/controller-manager/pkg/clientbuilder"
)

const (
	config_const = "dummy"
	provider     = "brightbox"
)

func interfaceProviderName(impl cloudprovider.Interface) func(*testing.T) {
	return func(t *testing.T) {
		if impl.ProviderName() != provider {
			t.Errorf("ProviderName should be %s", provider)
		}
	}
}

func interfaceHasClusterID(impl cloudprovider.Interface) func(*testing.T) {
	return func(t *testing.T) {
		if !impl.HasClusterID() {
			t.Errorf("HasClusterID should return true")
		}
	}
}

func interfaceRoutes(impl cloudprovider.Interface) func(*testing.T) {
	return func(t *testing.T) {
		_, supported := impl.Routes()
		if supported {
			t.Errorf("Routes should return false")
		}
	}
}

func interfaceLoadBalancer(impl cloudprovider.Interface) func(*testing.T) {
	return func(t *testing.T) {
		client, supported := impl.LoadBalancer()
		if !supported {
			t.Errorf("Instances should return true")
		}
		switch client.(type) {
		case (*cloud):
		default:
			t.Errorf("Instances returned incorrect client interface")
		}
	}
}

func interfaceClusters(impl cloudprovider.Interface) func(*testing.T) {
	return func(t *testing.T) {
		_, supported := impl.Clusters()
		if supported {
			t.Errorf("Clusters should return false")
		}
	}
}

func interfaceZones(impl cloudprovider.Interface) func(*testing.T) {
	return func(t *testing.T) {
		client, supported := impl.Zones()
		if !supported {
			t.Errorf("Zones should return true")
		}
		switch client.(type) {
		case (*cloud):
		default:
			t.Errorf("Zones returned incorrect client interface")
		}
	}
}

func interfaceInstances(impl cloudprovider.Interface) func(*testing.T) {
	return func(t *testing.T) {
		client, supported := impl.Instances()
		if !supported {
			t.Errorf("Instances should return true")
		}
		switch client.(type) {
		case (*cloud):
		default:
			t.Errorf("Instances returned incorrect client interface")
		}
	}
}

func interfaceInstancesV2(impl cloudprovider.Interface) func(*testing.T) {
	return func(t *testing.T) {
		client, supported := impl.InstancesV2()
		if !supported {
			t.Errorf("Instances should return true")
		}
		switch client.(type) {
		case (*cloud):
		default:
			t.Errorf("Instances returned incorrect client interface")
		}
	}
}

func TestInterfaceAdaption(t *testing.T) {
	var config io.Reader = strings.NewReader(config_const)
	var interfaceTests = []struct {
		name string
		fn   func(cloudprovider.Interface) func(*testing.T)
	}{
		{"ProviderName", interfaceProviderName},
		{"HasClusterID", interfaceHasClusterID},
		{"Routes", interfaceRoutes},
		{"Clusters", interfaceClusters},
		{"Zones", interfaceZones},
		{"Instances", interfaceInstances},
		{"InstancesV2", interfaceInstancesV2},
		{"LoadBalancer", interfaceLoadBalancer},
	}

	ts := k8ssdk.GetAuthEnvTokenHandler(t)
	defer k8ssdk.ResetAuthEnvironment()
	defer ts.Close()
	cloud, err := cloudprovider.GetCloudProvider(provider, config)
	if err != nil {
		t.Fatalf("Failed to obtain cloud structure: %v", err)
	} else if cloud == nil {
		t.Fatalf("Failed to initialise %s provider", provider)
	}
	cloud.Initialize(clientbuilder.SimpleControllerClientBuilder{}, make(chan struct{}))
	for _, example := range interfaceTests {
		t.Run(example.name, example.fn(cloud))
	}
}

func TestGetCloudProviderFailure(t *testing.T) {
	var config io.Reader = strings.NewReader(config_const)
	k8ssdk.ResetAuthEnvironment()
	defer k8ssdk.ResetAuthEnvironment()
	cloud, err := cloudprovider.GetCloudProvider(provider, config)
	if err == nil {
		t.Errorf("Expected error, didn't get one")
	} else if cloud != nil {
		t.Errorf("Expected nil cloud provider, got %+v", cloud)
	}
}
