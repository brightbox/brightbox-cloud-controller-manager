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

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// NodeAddresses returns the addresses of the specified instance.
// TODO(roberthbailey): This currently is only used in such a way that it
// returns the address of the calling instance. We should do a rename to
// make this clearer.
func (c *cloud) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	return nil, cloudprovider.NotImplemented
}

// NodeAddressesByProviderID returns the addresses of the specified instance.
// The instance is specified using the providerID of the node. The
// ProviderID is a unique identifier of the node. This will not be called
// from the node whose nodeaddresses are being queried. i.e. local metadata
// services cannot be used in this method to obtain nodeaddresses
func (c *cloud) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	return nil, cloudprovider.NotImplemented
	return nil, nil
}

// InstanceID returns the cloud provider ID of the node with the specified NodeName.
func (c *cloud) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	glog.V(4).Infof("InstanceID called for '%q'", nodeName)
	srv, err := c.getServer(ctx, mapNodeNameToServerID(nodeName))
	if err != nil {
		return "", cloudprovider.InstanceNotFound
	}
	return srv.Id, nil
}

// DEPRECATED: ExternalID returns the cloud provider ID of the node with
// the specified NodeName.
// Note that if the instance does not exist or is no longer running,
// we must return ("", cloudprovider.InstanceNotFound)
func (c *cloud) ExternalID(ctx context.Context, nodeName types.NodeName) (string, error) {
	glog.V(4).Infof("ExternalID called for '%q'", nodeName)
	return c.InstanceID(ctx, nodeName)
}

// InstanceType returns the type of the specified instance.
func (c *cloud) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	glog.V(4).Infof("InstanceType called for '%q'", name)
	srv, err := c.getServer(ctx, mapNodeNameToServerID(name))
	if err != nil {
		return "", err
	}
	return srv.Zone.Handle, nil
}

// InstanceTypeByProviderID returns the type of the specified instance.
func (c *cloud) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	glog.V(4).Infof("InstanceTypeByProviderID called for '%q'", providerID)
	srv, err := c.getServer(ctx, mapProviderIDToServerID(providerID))
	if err != nil {
		return "", err
	}
	return srv.Zone.Handle, nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (c *cloud) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	glog.V(4).Infof("AddSSHKey for '%q' called", user)
	return cloudprovider.NotImplemented
}

// CurrentNodeName returns the name of the node we are currently running on
// On most clouds (e.g. GCE) this is the hostname, so we provide the hostname
func (c *cloud) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	glog.V(4).Infof("CurrentNodeName(%q) called", hostname)
	return types.NodeName(hostname), nil
}

// InstanceExistsByProviderID returns true if the instance for the given provider id still is running.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
func (c *cloud) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	glog.V(4).Infof("InstanceExistsByProviderID called for '%q'", providerID)
	srv, err := c.getServer(ctx, mapProviderIDToServerID(providerID))
	if err != nil {
		if err == cloudprovider.InstanceNotFound {
			return false, nil
		}
		return false, err
	}
	if srv.Status != "active" {
		glog.Warningf("the instance %s is not active", srv.Id)
		return false, nil
	}
	return true, nil
}

// InstanceShutdownByProviderID returns true if the instance is shutdown in cloudprovider
func (c *cloud) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	glog.V(4).Infof("InstanceShutdownByProviderID called for '%q'", providerID)
	srv, err := c.getServer(ctx, mapProviderIDToServerID(providerID))
	if err != nil {
		if err == cloudprovider.InstanceNotFound {
			return false, nil
		}
		return false, err
	}
	if srv.Status != "active" {
		glog.Warningf("the instance %s is not active", srv.Id)
		return true, nil
	}
	return false, nil
}

