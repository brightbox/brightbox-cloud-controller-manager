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
	"errors"
	"testing"

	brightbox "github.com/brightbox/gobrightbox/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// TestLogAction verifies logging behavior, including handling a canceled context
func TestLogAction(t *testing.T) {
	// Test normal logging without context cancellation
	ctx := context.Background()
	err := logAction(ctx, "Test action: %s", "running")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test logging with canceled context
	myError := errors.New("test cancellation")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(myError) // immediately cancel the context
	err = logAction(ctx, "Test cancellation")
	if err == nil {
		t.Errorf("Expected error due to context cancellation, got nil")
	} else if err != myError {
		t.Errorf("Expected error %q, got %q", myError, err)
	}
}

// TestMapNodeNameToServerID verifies that a NodeName is mapped to the correct ServerID
func TestMapNodeNameToServerID(t *testing.T) {
	nodeName := types.NodeName("node123")
	serverID := mapNodeNameToServerID(nodeName)
	if serverID != "node123" {
		t.Errorf("Expected serverID to be 'node123', got '%s'", serverID)
	}
}

// TestMapServerIDToNodeName verifies that a ServerID is mapped to the correct NodeName
func TestMapServerIDToNodeName(t *testing.T) {
	serverID := "server123"
	nodeName := mapServerIDToNodeName(serverID)
	if nodeName != types.NodeName("server123") {
		t.Errorf("Expected NodeName to be 'server123', got '%s'", nodeName)
	}
}

// TestMapProviderIDToNodeName verifies that a ProviderID is correctly mapped to a NodeName
func TestMapProviderIDToNodeName(t *testing.T) {
	providerID := "brightbox://server123"
	nodeName := mapProviderIDToNodeName(providerID)
	expected := types.NodeName("server123")
	if nodeName != expected {
		t.Errorf("Expected NodeName to be '%s', got '%s'", expected, nodeName)
	}
}

// TestMapNodeToProviderID verifies that a Node is correctly mapped to its ProviderID
func TestMapNodeToProviderID(t *testing.T) {
	node := &v1.Node{Spec: v1.NodeSpec{ProviderID: "brightbox://server456"}}
	providerID := mapNodeToProviderID(node)
	expected := "brightbox://server456"
	if providerID != expected {
		t.Errorf("Expected ProviderID to be '%s', got '%s'", expected, providerID)
	}
}

// TestMapNodeToServerID verifies that a Node is correctly mapped to its ServerID
func TestMapNodeToServerID(t *testing.T) {
	node := &v1.Node{Spec: v1.NodeSpec{ProviderID: "brightbox://server789"}}
	serverID := mapNodeToServerID(node)
	expected := "server789"
	if serverID != expected {
		t.Errorf("Expected ServerID to be '%s', got '%s'", expected, serverID)
	}
}

// TestMapServerIDToNode verifies that a ServerID is correctly mapped to a Node with that ID as name
func TestMapServerIDToNode(t *testing.T) {
	serverID := "server101"
	node := mapServerIDToNode(serverID)
	if node.Name != "server101" {
		t.Errorf("Expected Node name to be 'server101', got '%s'", node.Name)
	}
}

// TestParseIPString verifies that valid and invalid IPs are handled correctly
func TestParseIPString(t *testing.T) {
	validIPv4 := "192.168.0.1"
	validIPv6 := "2001:db8::ff00:42:8329"
	invalidIP := "invalid_ip"

	_, err := parseIPString(validIPv4, "IPv4", "obj123", "Server", v1.NodeExternalIP)
	if err != nil {
		t.Errorf("Expected no error for valid IPv4, got %v", err)
	}

	_, err = parseIPString(validIPv6, "IPv6", "obj123", "Server", v1.NodeExternalIP)
	if err != nil {
		t.Errorf("Expected no error for valid IPv6, got %v", err)
	}

	_, err = parseIPString(invalidIP, "IPv4", "obj123", "Server", v1.NodeExternalIP)
	if err == nil {
		t.Errorf("Expected error for invalid IP, got nil")
	}
}

// TestNodeAddressesFromServer verifies that node addresses are correctly extracted from a Server
func TestNodeAddressesFromServer(t *testing.T) {
	server := &brightbox.Server{
		Hostname: "host123",
		Fqdn:     "host123.example.com",
		Interfaces: []brightbox.Interface{
			{IPv4Address: "192.168.0.1", IPv6Address: "2001:db8::ff00:42:8329"},
		},
		CloudIPs: []brightbox.CloudIP{
			{PublicIP: "198.51.100.1", Fqdn: "cloudip.example.com"},
		},
	}

	addresses, err := nodeAddressesFromServer(server)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedAddresses := []struct {
		Type    v1.NodeAddressType
		Address string
	}{
		{v1.NodeHostName, "host123"},
		{v1.NodeExternalDNS, "host123.example.com"},
		{v1.NodeExternalIP, "192.168.0.1"},
		{v1.NodeExternalIP, "2001:db8::ff00:42:8329"},
		{v1.NodeExternalIP, "198.51.100.1"},
		{v1.NodeExternalDNS, "cloudip.example.com"},
	}

	if len(addresses) != len(expectedAddresses) {
		t.Fatalf("Expected %d addresses, got %d", len(expectedAddresses), len(addresses))
	}

	for i, expected := range expectedAddresses {
		if addresses[i].Type != expected.Type || addresses[i].Address != expected.Address {
			t.Errorf("Expected address %d to be %v, got %v", i, expected, addresses[i])
		}
	}
}
