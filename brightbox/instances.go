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
	"net"

	"github.com/brightbox/gobrightbox/status"
	"github.com/brightbox/k8ssdk"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"
)

// NodeAddresses returns the addresses of the specified instance.
// TODO(roberthbailey): This currently is only used in such a way that it
// returns the address of the calling instance. We should do a rename to
// make this clearer.
func (c *cloud) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	klog.V(4).Infof("NodeAddresses (%q)", name)
	srv, err := c.GetServer(ctx, mapNodeNameToServerID(name), cloudprovider.InstanceNotFound)
	if err != nil {
		return nil, err
	}
	addresses := []v1.NodeAddress{
		{Type: v1.NodeHostName, Address: srv.Hostname},
		{Type: v1.NodeInternalDNS, Address: srv.Fqdn},
	}
	for _, iface := range srv.Interfaces {
		ipv4_node, err := parseIPString(iface.IPv4Address, "IPv4", srv.Id, "Server", v1.NodeInternalIP)
		if err != nil {
			return nil, err
		}
		ipv6_node, err := parseIPString(iface.IPv6Address, "IPv6", srv.Id, "Server", v1.NodeExternalIP)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, *ipv4_node, *ipv6_node)
	}
	for _, cip := range srv.CloudIPs {
		ipv4_node, err := parseIPString(cip.PublicIP, "IPv4", cip.Id, "Cloud IP", v1.NodeExternalIP)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, *ipv4_node, v1.NodeAddress{Type: v1.NodeExternalDNS, Address: cip.Fqdn})
	}
	return addresses, nil
}

// NodeAddressesByProviderID returns the addresses of the specified instance.
// The instance is specified using the providerID of the node. The
// ProviderID is a unique identifier of the node. This will not be called
// from the node whose nodeaddresses are being queried. i.e. local metadata
// services cannot be used in this method to obtain nodeaddresses
func (c *cloud) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	klog.V(4).Infof("NodeAddressesByProviderID (%q)", providerID)
	return c.NodeAddresses(ctx, mapProviderIDToNodeName(providerID))
}

// InstanceID returns the cloud provider ID of the node with the specified NodeName.
func (c *cloud) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	klog.V(4).Infof("InstanceID (%q)", nodeName)
	srv, err := c.GetServer(ctx, mapNodeNameToServerID(nodeName), cloudprovider.InstanceNotFound)
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
	klog.V(4).Infof("ExternalID (%q)", nodeName)
	return c.InstanceID(ctx, nodeName)
}

// InstanceType returns the type of the specified instance.
func (c *cloud) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	klog.V(4).Infof("InstanceType (%q)", name)
	srv, err := c.GetServer(ctx, mapNodeNameToServerID(name), cloudprovider.InstanceNotFound)
	if err != nil {
		return "", err
	}
	return srv.ServerType.Handle, nil
}

// InstanceTypeByProviderID returns the type of the specified instance.
func (c *cloud) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	klog.V(4).Infof("InstanceTypeByProviderID (%q)", providerID)
	return c.InstanceType(ctx, mapProviderIDToNodeName(providerID))
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (c *cloud) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	klog.V(4).Infof("AddSSHKey (%q)", user)
	return cloudprovider.NotImplemented
}

// CurrentNodeName returns the name of the node we are currently running on
// On most clouds (e.g. GCE) this is the hostname, so we provide the hostname
func (c *cloud) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	klog.V(4).Infof("CurrentNodeName (%q)", hostname)
	return mapServerIDToNodeName(hostname), nil
}

// InstanceExistsByProviderID returns true if the instance for the given provider exists.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
// This method should still return true for instances that exist but are stopped/sleeping.
func (c *cloud) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.V(4).Infof("InstanceExistsByProviderID (%q)", providerID)
	srv, err := c.GetServer(ctx, k8ssdk.MapProviderIDToServerID(providerID), cloudprovider.InstanceNotFound)
	if err != nil {
		if err == cloudprovider.InstanceNotFound {
			return false, nil
		}
		return false, err
	}
	switch srv.Status {
	case status.Active,
		status.Inactive,
		status.Deleting,
		status.Creating,
		status.Unavailable:
		klog.V(4).Infof("the instance %s exists", srv.Id)
		return true, nil
	case status.Deleted,
		status.Failed:
		klog.V(4).Infof("the instance %s does not exist", srv.Id)
		return false, nil
	default:
		return false, fmt.Errorf("Instance %s: Unrecognised status %q", srv.Id, srv.Status)
	}
}

// InstanceShutdownByProviderID returns true if the instance still exists and is shutdown in cloudprovider
func (c *cloud) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.V(4).Infof("InstanceShutdownByProviderID (%q)", providerID)
	srv, err := c.GetServer(ctx, k8ssdk.MapProviderIDToServerID(providerID), cloudprovider.InstanceNotFound)
	if err != nil {
		if err == cloudprovider.InstanceNotFound {
			return false, nil
		}
		return false, err
	}
	switch srv.Status {
	case status.Inactive,
		status.Unavailable:
		klog.V(4).Infof("the instance %s is shutdown", srv.Id)
		return true, nil
	case status.Active,
		status.Creating,
		status.Deleting,
		status.Deleted,
		status.Failed:
		klog.V(4).Infof("the instance %s is not shutdown", srv.Id)
		return false, nil
	default:
		return false, fmt.Errorf("Instance %s: Unrecognised status %q", srv.Id, srv.Status)
	}
}

func parseIPString(ipString string, ipType string, objectId string,
	objectType string, nodeType v1.NodeAddressType) (*v1.NodeAddress, error) {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("%s has invalid %s address: %s (%q)", objectType, ipType, objectId, ipString)
	}
	return &v1.NodeAddress{Type: nodeType, Address: ip.String()}, nil
}
