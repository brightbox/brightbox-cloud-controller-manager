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

	"github.com/brightbox/gobrightbox"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

var defaultRegionCidr = "10.0.0.0/8"
var defaultRuleProtocol = loadBalancerTcpProtocol

// The approach is to create a separate server group, firewall policy
// and firewall rule for each loadbalancer primarily to avoid any
// potential race conditions in the driver.
// It also allows k8s to select subsets of nodes for each loadbalancer
// created if it wants to.
func (c *cloud) ensureFirewallOpenForService(apiservice *v1.Service, nodes []*v1.Node) error {
	name := cloudprovider.GetLoadBalancerName(apiservice)
	glog.V(4).Infof("ensureFireWallOpen(%v)", name)
	if len(apiservice.Spec.Ports) <= 0 {
		glog.V(4).Infof("no ports to open")
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
	glog.V(4).Infof("ensureServerGroup(%v)", name)
	group, err := c.getServerGroupByName(name)
	if err != nil {
		return nil, err
	}
	if group == nil {
		group, err = c.createServerGroup(name)
	}
	if err != nil {
		return nil, err
	}
	group, err = c.syncServerGroup(group, mapNodesToServerIDs(nodes))
	if err == nil {
		return group, nil
	}
	return nil, err
}

func (c *cloud) ensureFirewallPolicy(group *brightbox.ServerGroup) (*brightbox.FirewallPolicy, error) {
	glog.V(4).Infof("ensureFireWallPolicy (%q)", group.Name)
	fp, err := c.getFirewallPolicyByName(group.Name)
	if err != nil {
		return nil, err
	}
	if fp == nil {
		return c.createFirewallPolicy(group)
	}
	return fp, nil
}

func (c *cloud) ensureFirewallRules(apiservice *v1.Service, fp *brightbox.FirewallPolicy) error {
	glog.V(4).Infof("ensureFireWallRules (%q)", fp.Name)
	portListStr := createPortListString(apiservice.Spec.Ports)
	newRule := brightbox.FirewallRuleOptions{
		FirewallPolicy:  fp.Id,
		Protocol:        &defaultRuleProtocol,
		Source:          &defaultRegionCidr,
		DestinationPort: &portListStr,
		Description:     &fp.Name,
	}
	if len(fp.Rules) == 0 {
		_, err := c.createFirewallRule(&newRule)
		return err
	} else if isUpdateFirewallRuleRequired(fp.Rules[0], newRule) {
		newRule.Id = fp.Rules[0].Id
		_, err := c.updateFirewallRule(&newRule)
		return err
	}
	glog.V(4).Infof("No rule update required for %q, skipping", fp.Rules[0].Id)
	return nil
}

func isUpdateFirewallRuleRequired(old brightbox.FirewallRule, new brightbox.FirewallRuleOptions) bool {
	return (new.Protocol != nil && *new.Protocol != old.Protocol) ||
		(new.Source != nil && *new.Source != old.Source) ||
		(new.DestinationPort != nil && *new.DestinationPort != old.DestinationPort) ||
		(new.Description != nil && *new.Description != old.Description)
}

func createPortListString(ports []v1.ServicePort) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(ports[0].NodePort)))
	for i := range ports[1:] {
		buffer.WriteString(",")
		buffer.WriteString(strconv.Itoa(int(ports[i].NodePort)))
	}
	return buffer.String()
}

func mapNodesToServerIDs(nodes []*v1.Node) []string {
	result := make([]string, 0, len(nodes))
	for i := range nodes {
		if nodes[i].Spec.ProviderID == "" {
			glog.Warningf("node %q did not have providerID set", nodes[i].Name)
			continue
		}
		result = append(result, mapProviderIDToServerID(nodes[i].Spec.ProviderID))
	}
	return result
}
