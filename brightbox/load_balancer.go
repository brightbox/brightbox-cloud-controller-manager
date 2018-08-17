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

	"github.com/brightbox/gobrightbox"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/api/v1/service"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

const (
	loadBalancerTcpProtocol  = "tcp"
	loadBalancerHttpProtocol = "http"
	//Healthcheck on http port if there are no endpoints for the loadbalancer
	defaultTcpHealthCheckPort = 80
)

// Ignoring clusterName completely as it doesn't appear to be used anywhere else at the moment
func (c *cloud) GetLoadBalancer(ctx context.Context, clusterName string, apiservice *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	glog.V(4).Infof("GetLoadBalancer(%v, %v)", clusterName, apiservice.UID)
	lb, err := c.getLoadBalancerFromService(apiservice)
	return toLoadBalancerStatus(lb), err == nil && lb != nil, err
}

// Make sure we have a cloud ip before asking for a load balancer. Try
// to get one matching the LoadBalancerIP spec in the service, and error
// if that isn't in the cloudip list.
func (c *cloud) EnsureLoadBalancer(ctx context.Context, clusterName string, apiservice *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	glog.V(4).Infof("EnsureLoadBalancer(%v, %v, %v, %v, %v, %v)",
		clusterName, apiservice.UID, apiservice.Namespace, apiservice.Spec.LoadBalancerIP, apiservice.Spec.Ports, apiservice.Annotations)
	if err := validateServiceSpec(apiservice); err != nil {
		return nil, err
	}
	cip, err := c.ensureAllocatedCip(apiservice)
	if err != nil {
		return nil, err
	}
	lb, err := c.ensureLoadBalancerFromService(apiservice, nodes)
	if err != nil {
		return nil, err
	}
	err = c.ensureMappedCip(lb, cip)
	if err != nil {
		return nil, err
	}
	lb, err = c.getLoadBalancerFromService(apiservice)
	if err != nil {
		return nil, err
	}
	return toLoadBalancerStatus(lb), errorIfNotComplete(lb, string(apiservice.UID))
}

func (c *cloud) UpdateLoadBalancer(ctx context.Context, clusterName string, apiservice *v1.Service, nodes []*v1.Node) error {
	glog.V(4).Infof("UpdateLoadBalancer called - delegating")
	_, err := c.EnsureLoadBalancer(ctx, clusterName, apiservice, nodes)
	return err
}

func (c *cloud) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, apiservice *v1.Service) error {
	name := cloudprovider.GetLoadBalancerName(apiservice)
	glog.V(4).Infof("EnsureLoadBalancerDeleted(%v, %v, %v, %v)", clusterName, name, apiservice.Namespace, apiservice.Spec.LoadBalancerIP)
	if err := c.ensureServerGroupDeleted(name); err != nil {
		return err
	}
	if err := c.ensureFirewallClosed(name); err != nil {
		return err
	}
	if err := c.ensureLoadBalancerDeletedByName(name); err != nil {
		return err
	}
	if err := c.ensureCloudIPsDeleted(name); err != nil {
		return err
	}
	lb, err := c.getLoadBalancerByName(name)
	if err != nil {
		return err
	}
	return errorIfNotErased(lb)
}

//Take all the servers out of the server group and remove it
func (c *cloud) ensureServerGroupDeleted(name string) error {
	glog.V(4).Infof("ensureServerGroupDeleted (%q)", name)
	group, err := c.getServerGroupByName(name)
	if err != nil {
		glog.V(4).Infof("Error looking for Server Group for %q", name)
		return err
	}
	if group == nil {
		return nil
	}
	group, err = c.syncServerGroup(group, nil)
	if err != nil {
		glog.V(4).Infof("Error removing servers from %q", group.Id)
		return err
	}
	if err := c.destroyServerGroup(group.Id); err != nil {
		glog.V(4).Infof("Error destroying Server Group %q", group.Id)
		return err
	}
	return nil
}

//Remove the firewall policy
func (c *cloud) ensureFirewallClosed(name string) error {
	glog.V(4).Infof("ensureFirewallClosed (%q)", name)
	fp, err := c.getFirewallPolicyByName(name)
	if err != nil {
		glog.V(4).Infof("Error looking for Firewall Policy %q", name)
		return err
	}
	if fp == nil {
		return nil
	}
	if err := c.destroyFirewallPolicy(fp.Id); err != nil {
		glog.V(4).Infof("Error destroying Firewall Policy %q", fp.Id)
		return err
	}
	return nil
}

//Try to remove the loadbalancer
func (c *cloud) ensureLoadBalancerDeletedByName(name string) error {
	glog.V(4).Infof("ensureLoadBalancerDeletedByName (%q)", name)
	lb, err := c.getLoadBalancerByName(name)
	if err != nil {
		glog.V(4).Infof("Error looking for Load Balancer %q", name)
		return err
	}
	if lb == nil {
		return nil
	}
	if err = c.destroyLoadBalancer(lb.Id); err != nil {
		glog.V(4).Infof("Error destroying Load Balancer %q", lb.Id)
		return err
	}
	return nil
}

//Try to remove CloudIPs matching `name` from the list of cloudIPs
func (c *cloud) ensureCloudIPsDeleted(name string) error {
	glog.V(4).Infof("ensureCloudIPsDeleted (%q)", name)
	cloudIpList, err := c.getCloudIPs()
	if err != nil {
		glog.V(4).Infof("Error retrieving list of CloudIPs")
		return err
	}
	for i := range cloudIpList {
		if cloudIpList[i].Name == name {
			err := c.destroyCloudIP(cloudIpList[i].Id)
			if err != nil {
				glog.V(4).Infof("Error destroying CloudIP %q", cloudIpList[i].Id)
				return err
			}
		}
	}
	return nil
}

func toLoadBalancerStatus(lb *brightbox.LoadBalancer) *v1.LoadBalancerStatus {
	if lb == nil {
		return nil
	}
	status := v1.LoadBalancerStatus{}
	if len(lb.CloudIPs) > 0 {
		status.Ingress = make([]v1.LoadBalancerIngress, len(lb.CloudIPs))
		for i, v := range lb.CloudIPs {
			status.Ingress[i] = v1.LoadBalancerIngress{
				Hostname: selectHostname(&v),
				IP:       v.PublicIP,
			}
		}
	}
	return &status
}

func validateServiceSpec(apiservice *v1.Service) error {
	if apiservice.Spec.SessionAffinity != v1.ServiceAffinityNone {
		return fmt.Errorf("unsupported load balancer affinity: %v", apiservice.Spec.SessionAffinity)
	}
	if len(apiservice.Spec.Ports) == 0 {
		return fmt.Errorf("requested load balancer with no ports")
	}
	for _, port := range apiservice.Spec.Ports {
		if port.Protocol != v1.ProtocolTCP {
			return fmt.Errorf("UDP nodeports are not supported")
		}
	}
	return nil
}

// nil if no loadbalancer
func (c *cloud) getLoadBalancerFromService(apiservice *v1.Service) (*brightbox.LoadBalancer, error) {
	name := cloudprovider.GetLoadBalancerName(apiservice)
	glog.V(4).Infof("getLoadBalancerFromService(%v)", name)
	return c.getLoadBalancerByName(name)
}

func (c *cloud) ensureAllocatedCip(apiservice *v1.Service) (*brightbox.CloudIP, error) {
	name := cloudprovider.GetLoadBalancerName(apiservice)
	ip := apiservice.Spec.LoadBalancerIP
	glog.V(4).Infof("ensureAllocatedCip (%q, %q)", name, ip)
	cloudIpList, err := c.getCloudIPs()
	if err != nil {
		return nil, err
	}
	for i := range cloudIpList {
		if cloudIpList[i].PublicIP == ip || cloudIpList[i].Name == name {
			return &cloudIpList[i], nil
		}
	}
	if ip == "" {
		return c.allocateCip(name)
	} else {
		return nil, fmt.Errorf("Could not find allocated Cloud IP with address %q", ip)
	}
}

func (c *cloud) ensureLoadBalancerFromService(apiservice *v1.Service, nodes []*v1.Node) (*brightbox.LoadBalancer, error) {
	glog.V(4).Infof("ensureLoadBalancerFromService(%v)", apiservice.UID)
	current_lb, err := c.getLoadBalancerFromService(apiservice)
	if err != nil {
		return nil, err
	}
	err = c.ensureFirewallOpenForService(apiservice, nodes)
	if err != nil {
		return nil, err
	}
	newLB := buildLoadBalancerOptions(apiservice, nodes)
	if current_lb == nil {
		return c.createLoadBalancer(newLB)
	} else if isUpdateLoadBalancerRequired(current_lb, *newLB) {
		newLB.Id = current_lb.Id
		return c.updateLoadBalancer(newLB)
	}
	glog.V(4).Infof("No Load Balancer update required for %q, skipping", current_lb.Id)
	return current_lb, nil
}

func buildLoadBalancerOptions(apiservice *v1.Service, nodes []*v1.Node) *brightbox.LoadBalancerOptions {
	name := cloudprovider.GetLoadBalancerName(apiservice)
	glog.V(4).Infof("buildLoadBalancerOptions(%v)", name)
	return &brightbox.LoadBalancerOptions{
		Name:        &name,
		Nodes:       buildLoadBalancerNodes(nodes),
		Listeners:   buildLoadBalancerListeners(apiservice),
		Healthcheck: buildLoadBalancerHealthCheck(apiservice),
	}
}

func buildLoadBalancerNodes(nodes []*v1.Node) *[]brightbox.LoadBalancerNode {
	if len(nodes) <= 0 {
		return nil
	}
	result := make([]brightbox.LoadBalancerNode, 0, len(nodes))
	for i := range nodes {
		if nodes[i].Spec.ProviderID == "" {
			glog.Warningf("node %q did not have providerID set", nodes[i].Name)
			continue
		}
		result = append(result, brightbox.LoadBalancerNode{Node: mapProviderIDToServerID(nodes[i].Spec.ProviderID)})
	}
	return &result
}

func buildLoadBalancerListeners(apiservice *v1.Service) *[]brightbox.LoadBalancerListener {
	if len(apiservice.Spec.Ports) <= 0 {
		return nil
	}
	result := make([]brightbox.LoadBalancerListener, len(apiservice.Spec.Ports))
	for i := range apiservice.Spec.Ports {
		result[i].Protocol = loadBalancerTcpProtocol
		result[i].In = int(apiservice.Spec.Ports[i].Port)
		result[i].Out = int(apiservice.Spec.Ports[i].NodePort)
	}
	return &result
}

func buildLoadBalancerHealthCheck(apiservice *v1.Service) *brightbox.LoadBalancerHealthcheck {
	if path, healthCheckNodePort := service.GetServiceHealthCheckPathPort(apiservice); path != "" {
		return &brightbox.LoadBalancerHealthcheck{
			Type:    loadBalancerHttpProtocol,
			Port:    int(healthCheckNodePort),
			Request: path,
		}
	} else {
		return &brightbox.LoadBalancerHealthcheck{
			Type:    loadBalancerTcpProtocol,
			Port:    getTcpHealthCheckPort(apiservice),
			Request: "/",
		}
	}
}

func getTcpHealthCheckPort(apiservice *v1.Service) int {
	for i := range apiservice.Spec.Ports {
		return int(apiservice.Spec.Ports[i].NodePort)
	}
	return defaultTcpHealthCheckPort
}
