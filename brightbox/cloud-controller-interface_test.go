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
	"strings"
	"testing"
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

func TestInterfaceAdaption(t *testing.T) {
	var config io.Reader = strings.NewReader(config_const)
	var interface_tests = []struct {
		name string
		fn   func(cloudprovider.Interface) func(*testing.T)
	}{
		{"ProviderName", interfaceProviderName},
		{"HasClusterID", interfaceHasClusterID},
		{"Routes", interfaceRoutes},
		{"Clusters", interfaceClusters},
		{"Zones", interfaceZones},
	}

	cloud, err := cloudprovider.GetCloudProvider(provider, config)
	if cloud == nil {
		t.Fatalf("Failed to initialise %s provider", provider)
	} else if err != nil {
		t.Fatalf("Failed to obtain cloud structure: %v", err)
	}
	for _, example := range interface_tests {
		t.Run(example.name, example.fn(cloud))
	}
}
