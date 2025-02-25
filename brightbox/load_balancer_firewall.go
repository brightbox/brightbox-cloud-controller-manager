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
	"bytes"
	"context"
	"strconv"

	brightbox "github.com/brightbox/gobrightbox/v2"
	"github.com/brightbox/gobrightbox/v2/enums/listenerprotocol"
	"github.com/brightbox/k8ssdk/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

var defaultRegionCidr = "10.0.0.0/8"
var defaultIPv6RegionCidr = "2a02:1348:0140::/42"
var defaultRuleProtocol = listenerprotocol.Tcp.String()

// The approach is to create a separate server group, firewall policy
// and firewall rule for each loadbalancer primarily to avoid any
// potential race conditions in the driver.
// It also allows k8s to select subsets of nodes for each loadbalancer
// created if it wants to.
func (c *cloud) ensureFirewallOpenForService(ctx context.Context, name string, apiservice *v1.Service, nodes []*v1.Node) error {
	klog.V(4).Infof("ensureFireWallOpen(%v)", name)
	if len(apiservice.Spec.Ports) <= 0 {
		klog.V(4).Infof("no ports to open")
		return nil
	}
	serverGroup, err := c.ensureServerGroup(ctx, name, nodes)
	if err != nil {
		return err
	}
	firewallPolicy, err := c.ensureFirewallPolicy(ctx, serverGroup)
	if err != nil {
		return err
	}
	return c.ensureFirewallRules(ctx, apiservice, firewallPolicy)
}

func (c *cloud) ensureServerGroup(ctx context.Context, name string, nodes []*v1.Node) (*brightbox.ServerGroup, error) {
	klog.V(4).Infof("ensureServerGroup(%v)", name)
	group, err := c.GetServerGroupByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if group == nil {
		group, err = c.CreateServerGroup(ctx, name)
	}
	if err != nil {
		return nil, err
	}
	group, err = c.SyncServerGroup(ctx, group, mapNodesToServerIDs(nodes))
	if err == nil {
		return group, nil
	}
	return nil, err
}

func (c *cloud) ensureFirewallPolicy(ctx context.Context, group *brightbox.ServerGroup) (*brightbox.FirewallPolicy, error) {
	klog.V(4).Infof("ensureFireWallPolicy (%q)", group.Name)
	fp, err := c.GetFirewallPolicyByName(ctx, group.Name)
	if err != nil {
		return nil, err
	}
	if fp == nil {
		return c.CreateFirewallPolicy(ctx, *group)
	}
	return fp, nil
}

func (c *cloud) ensureFirewallRules(ctx context.Context, apiservice *v1.Service, fp *brightbox.FirewallPolicy) error {
	klog.V(4).Infof("ensureFireWallRules (%q)", fp.Name)
	portListStr := createPortListString(apiservice)
	newRule := brightbox.FirewallRuleOptions{
		FirewallPolicy:  fp.ID,
		Protocol:        &defaultRuleProtocol,
		Source:          &defaultRegionCidr,
		DestinationPort: &portListStr,
		Description:     &fp.Name,
	}
	if len(fp.Rules) == 0 {
		_, err := c.CreateFirewallRule(ctx, newRule)
		return err
	} else if isUpdateFirewallRuleRequired(fp.Rules[0], newRule) {
		newRule.ID = fp.Rules[0].ID
		_, err := c.UpdateFirewallRule(ctx, newRule)
		return err
	}
	klog.V(4).Infof("No rule update required for %q, skipping", fp.Rules[0].ID)
	return nil
}

func isUpdateFirewallRuleRequired(old brightbox.FirewallRule, new brightbox.FirewallRuleOptions) bool {
	return (new.Protocol != nil && *new.Protocol != old.Protocol) ||
		(new.Source != nil && *new.Source != old.Source) ||
		(new.DestinationPort != nil && *new.DestinationPort != old.DestinationPort) ||
		(new.Description != nil && *new.Description != old.Description)
}

func createPortListString(apiservice *v1.Service) string {
	var buffer bytes.Buffer
	ports := apiservice.Spec.Ports
	buffer.WriteString(strconv.Itoa(int(ports[0].NodePort)))
	for i := range ports[1:] {
		buffer.WriteString(",")
		buffer.WriteString(strconv.Itoa(int(ports[i+1].NodePort)))
	}
	if apiservice.Spec.HealthCheckNodePort != 0 {
		buffer.WriteString(",")
		buffer.WriteString(strconv.Itoa(int(apiservice.Spec.HealthCheckNodePort)))
	}
	return buffer.String()
}

func mapNodesToServerIDs(nodes []*v1.Node) []string {
	result := make([]string, 0, len(nodes))
	for i := range nodes {
		if nodes[i].Spec.ProviderID == "" {
			klog.Warningf("node %q did not have providerID set", nodes[i].Name)
			continue
		}
		result = append(result, k8ssdk.MapProviderIDToServerID(nodes[i].Spec.ProviderID))
	}
	return result
}

// Take all the servers out of the server group and remove it
func (c *cloud) ensureServerGroupDeleted(ctx context.Context, name string) error {
	klog.V(4).Infof("ensureServerGroupDeleted (%q)", name)
	group, err := c.GetServerGroupByName(ctx, name)
	if err != nil {
		klog.V(4).Infof("Error looking for Server Group for %q", name)
		return err
	}
	if group == nil {
		return nil
	}
	group, err = c.SyncServerGroup(ctx, group, nil)
	if err != nil {
		klog.V(4).Infof("Error removing servers from %q", group.ID)
		return err
	}
	if err := c.DestroyServerGroup(ctx, group.ID); err != nil {
		klog.V(4).Infof("Error destroying Server Group %q", group.ID)
		return err
	}
	return nil
}

// Remove the firewall policy
func (c *cloud) ensureFirewallClosed(ctx context.Context, name string) error {
	klog.V(4).Infof("ensureFirewallClosed (%q)", name)
	fp, err := c.GetFirewallPolicyByName(ctx, name)
	if err != nil {
		klog.V(4).Infof("Error looking for Firewall Policy %q", name)
		return err
	}
	if fp == nil {
		return nil
	}
	if err := c.DestroyFirewallPolicy(ctx, fp.ID); err != nil {
		klog.V(4).Infof("Error destroying Firewall Policy %q", fp.ID)
		return err
	}
	return nil
}
