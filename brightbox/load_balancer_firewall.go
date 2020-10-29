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
	"strconv"

	brightbox "github.com/brightbox/gobrightbox"
	"github.com/brightbox/k8ssdk"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

var defaultRegionCidr = "10.0.0.0/8"
var defaultRuleProtocol = loadBalancerTCPProtocol

// The approach is to create a separate server group, firewall policy
// and firewall rule for each loadbalancer primarily to avoid any
// potential race conditions in the driver.
// It also allows k8s to select subsets of nodes for each loadbalancer
// created if it wants to.
func (c *cloud) ensureFirewallOpenForService(name string, apiservice *v1.Service, nodes []*v1.Node) error {
	klog.V(4).Infof("ensureFireWallOpen(%v)", name)
	if len(apiservice.Spec.Ports) <= 0 {
		klog.V(4).Infof("no ports to open")
		return nil
	}
	serverGroup, err := c.ensureServerGroup(name, nodes)
	if err != nil {
		return err
	}
	firewallPolicy, err := c.ensureFirewallPolicy(serverGroup)
	if err != nil {
		return err
	}
	return c.ensureFirewallRules(apiservice, firewallPolicy)
}

func (c *cloud) ensureServerGroup(name string, nodes []*v1.Node) (*brightbox.ServerGroup, error) {
	klog.V(4).Infof("ensureServerGroup(%v)", name)
	group, err := c.GetServerGroupByName(name)
	if err != nil {
		return nil, err
	}
	if group == nil {
		group, err = c.CreateServerGroup(name)
	}
	if err != nil {
		return nil, err
	}
	group, err = c.SyncServerGroup(group, mapNodesToServerIDs(nodes))
	if err == nil {
		return group, nil
	}
	return nil, err
}

func (c *cloud) ensureFirewallPolicy(group *brightbox.ServerGroup) (*brightbox.FirewallPolicy, error) {
	klog.V(4).Infof("ensureFireWallPolicy (%q)", group.Name)
	fp, err := c.GetFirewallPolicyByName(group.Name)
	if err != nil {
		return nil, err
	}
	if fp == nil {
		return c.CreateFirewallPolicy(group)
	}
	return fp, nil
}

func (c *cloud) ensureFirewallRules(apiservice *v1.Service, fp *brightbox.FirewallPolicy) error {
	klog.V(4).Infof("ensureFireWallRules (%q)", fp.Name)
	portListStr := createPortListString(apiservice)
	newRule := brightbox.FirewallRuleOptions{
		FirewallPolicy:  fp.Id,
		Protocol:        &defaultRuleProtocol,
		Source:          &defaultRegionCidr,
		DestinationPort: &portListStr,
		Description:     &fp.Name,
	}
	if len(fp.Rules) == 0 {
		_, err := c.CreateFirewallRule(&newRule)
		return err
	} else if isUpdateFirewallRuleRequired(fp.Rules[0], newRule) {
		newRule.Id = fp.Rules[0].Id
		_, err := c.UpdateFirewallRule(&newRule)
		return err
	}
	klog.V(4).Infof("No rule update required for %q, skipping", fp.Rules[0].Id)
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
