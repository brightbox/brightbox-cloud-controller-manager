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

	brightbox "github.com/brightbox/gobrightbox/v2"
	"github.com/brightbox/k8ssdk/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

func logAction(ctx context.Context, action string, args ...interface{}) error {
	// Log the action with the provided arguments
	klog.V(4).Infof(action, args...)

	// Check if context is canceled and log the cause if it exists
	if cause := context.Cause(ctx); cause != nil {
		klog.Warningf("unexpected context - %q abandoned due to: %v", fmt.Sprintf(action, args...), cause)
		return cause
	}

	return nil
}

// mapNodeNameToServerID maps a k8s NodeName to a Brightbox Server ID
// This is a simple string cast.
func mapNodeNameToServerID(nodeName types.NodeName) string {
	return string(nodeName)
}

// mapServerIDToNodeName maps a Brightbox Server ID to a nodename
// Again a simple string cast
func mapServerIDToNodeName(name string) types.NodeName {
	return types.NodeName(name)
}

func mapProviderIDToNodeName(providerID string) types.NodeName {
	return mapServerIDToNodeName(k8ssdk.MapProviderIDToServerID(providerID))
}

// func mapNodeNameToProviderID(nodeName types.NodeName) string {
// 	return k8ssdk.MapServerIDToProviderID(mapNodeNameToServerID(nodeName))
// }

func mapNodeToProviderID(node *v1.Node) string {
	if node.Spec.ProviderID == "" {
		return k8ssdk.MapServerIDToProviderID(node.Name)
	}
	return node.Spec.ProviderID
}

func mapNodeToServerID(node *v1.Node) string {
	if node.Spec.ProviderID == "" {
		return node.Name
	}
	return k8ssdk.MapProviderIDToServerID(node.Spec.ProviderID)
}

func mapServerIDToNode(name string) *v1.Node {
	result := &v1.Node{}
	result.Name = name
	return result
}

func mapServerIDToNodeProviderID(name string) *v1.Node {
	result := &v1.Node{}
	result.Spec.ProviderID = k8ssdk.MapServerIDToProviderID(name)
	return result
}

func nodeAddressesFromServer(srv *brightbox.Server) ([]v1.NodeAddress, error) {
	addresses := []v1.NodeAddress{
		{Type: v1.NodeHostName, Address: srv.Hostname},
		{Type: v1.NodeExternalDNS, Address: srv.Fqdn},
	}
	for _, iface := range srv.Interfaces {
		ipv4Node, err := parseIPString(iface.IPv4Address, "IPv4", srv.ID, "Server", v1.NodeExternalIP)
		if err != nil {
			return nil, err
		}
		ipv6Node, err := parseIPString(iface.IPv6Address, "IPv6", srv.ID, "Server", v1.NodeExternalIP)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, *ipv4Node, *ipv6Node)
	}
	for _, cip := range srv.CloudIPs {
		ipv4Node, err := parseIPString(cip.PublicIP, "IPv4", cip.ID, "Cloud IP", v1.NodeExternalIP)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, *ipv4Node, v1.NodeAddress{Type: v1.NodeExternalDNS, Address: cip.Fqdn})
	}
	return addresses, nil
}
