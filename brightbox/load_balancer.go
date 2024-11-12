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

	brightbox "github.com/brightbox/gobrightbox/v2"
	"github.com/brightbox/gobrightbox/v2/enums/balancingpolicy"
	"github.com/brightbox/gobrightbox/v2/enums/listenerprotocol"
	"github.com/brightbox/k8ssdk/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

var (
	truevar  = true
	falsevar = false
)

// Remove load balancer by name
func (c *cloud) ensureLoadBalancerDeletedByName(ctx context.Context, name string) (*brightbox.LoadBalancer, error) {
	lb, err := c.GetLoadBalancerByName(ctx, name)
	if err != nil {
		klog.V(4).Infof("Error looking for Load Balancer %q", name)
		return nil, err
	}
	if lb != nil {
		if err = c.DestroyLoadBalancer(ctx, lb.ID); err != nil {
			klog.V(4).Infof("Error destroying Load Balancer %q", lb.ID)
			return nil, err
		}
	}
	return lb, nil
}

func buildLoadBalancerOptions(name string, domains []string, apiservice *v1.Service, nodes []*v1.Node) *brightbox.LoadBalancerOptions {
	klog.V(4).Infof("buildLoadBalancerOptions(%v)", name)
	result := &brightbox.LoadBalancerOptions{
		Name:        &name,
		Nodes:       buildLoadBalancerNodes(nodes),
		Listeners:   buildLoadBalancerListeners(apiservice),
		Healthcheck: buildLoadBalancerHealthCheck(apiservice),
	}
	if policy, ok := apiservice.Annotations[serviceAnnotationLoadBalancerPolicy]; ok {
		policyEnum, err := balancingpolicy.ParseEnum(policy)
		if err == nil {
			result.Policy = policyEnum
		} else {
			klog.V(4).Infof("Unexpected balancing policy %q", policy)
		}
	}
	for _, listener := range result.Listeners {
		if listener.Protocol == listenerprotocol.Https {
			result.Domains = &domains
			result.HTTPSRedirect = &truevar
			return result
		}
	}
	result.HTTPSRedirect = &falsevar
	return result
}

func buildLoadBalancerNodes(nodes []*v1.Node) []brightbox.LoadBalancerNode {
	if len(nodes) <= 0 {
		return nil
	}
	result := make([]brightbox.LoadBalancerNode, 0, len(nodes))
	for i := range nodes {
		if nodes[i].Spec.ProviderID == "" {
			klog.Warningf("node %q did not have providerID set", nodes[i].Name)
			continue
		}
		result = append(result, brightbox.LoadBalancerNode{Node: k8ssdk.MapProviderIDToServerID(nodes[i].Spec.ProviderID)})
	}
	return result
}

func buildLoadBalancerListeners(apiservice *v1.Service) []brightbox.LoadBalancerListener {
	if len(apiservice.Spec.Ports) <= 0 {
		return nil
	}
	sslPortSet := getPortSets(apiservice.Annotations[serviceAnnotationLoadBalancerSSLPorts])
	result := make([]brightbox.LoadBalancerListener, len(apiservice.Spec.Ports))
	for i := range apiservice.Spec.Ports {
		result[i].Protocol = getListenerProtocol(apiservice)
		result[i].ProxyProtocol = getListenerProxyProtocol(apiservice)
		if result[i].Protocol != listenerprotocol.Tcp && isSSLPort(&apiservice.Spec.Ports[i], sslPortSet) {
			result[i].Protocol = listenerprotocol.Https
		}
		result[i].In = uint16(apiservice.Spec.Ports[i].Port)
		result[i].Out = uint16(apiservice.Spec.Ports[i].NodePort)
		result[i].Timeout, _ = parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerListenerIdleTimeout)
	}
	return result
}

func buildLoadBalancerHealthCheck(apiservice *v1.Service) *brightbox.LoadBalancerHealthcheck {
	path, healthCheckNodePort := getServiceHealthCheckPathPort(apiservice)
	protocol := getHealthCheckProtocol(apiservice, path)
	//Validate has already checked all these so there should be no errors!
	interval, _ := parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerHCInterval)
	timeout, _ := parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerHCTimeout)
	thresholdUp, _ := parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerHCHealthyThreshold)
	thresholdDown, _ := parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerHCUnhealthyThreshold)
	return &brightbox.LoadBalancerHealthcheck{
		Type:          protocol,
		Port:          getHealthCheckPort(apiservice, uint16(healthCheckNodePort)),
		Request:       getHealthCheckPath(apiservice, protocol, path),
		Interval:      interval,
		Timeout:       timeout,
		ThresholdUp:   thresholdUp,
		ThresholdDown: thresholdDown,
	}
}

func toLoadBalancerStatus(lb *brightbox.LoadBalancer) *v1.LoadBalancerStatus {
	status := v1.LoadBalancerStatus{}
	if lb == nil {
		return &status
	}
	if len(lb.CloudIPs) > 0 {
		status.Ingress = make([]v1.LoadBalancerIngress, 0, len(lb.CloudIPs)*4)
		for _, v := range lb.CloudIPs {
			/*
				status.Ingress = append(status.Ingress,
					v1.LoadBalancerIngress{
						IP: v.PublicIPv4,
					},
					v1.LoadBalancerIngress{
						IP: v.PublicIPv6,
					},
				)
			*/
			if v.ReverseDNS != "" {
				status.Ingress = append(status.Ingress,
					v1.LoadBalancerIngress{
						Hostname: v.ReverseDNS,
					},
				)
			}
			if v.Fqdn != "" {
				status.Ingress = append(status.Ingress,
					v1.LoadBalancerIngress{
						Hostname: v.Fqdn,
					},
				)
			}
		}
	}
	return &status
}

func (c *cloud) ensureLoadBalancerFromService(ctx context.Context, name string, domains []string, apiservice *v1.Service, nodes []*v1.Node) (*brightbox.LoadBalancer, error) {
	klog.V(4).Infof("ensureLoadBalancerFromService(%v)", name)
	currentLb, err := c.GetLoadBalancerByName(ctx, name)
	if err != nil {
		return nil, err
	}
	err = c.ensureFirewallOpenForService(ctx, name, apiservice, nodes)
	if err != nil {
		return nil, err
	}
	newLB := buildLoadBalancerOptions(name, domains, apiservice, nodes)
	if currentLb == nil {
		return c.Cloud.CreateLoadBalancer(ctx, *newLB)
	} else if k8ssdk.IsUpdateLoadBalancerRequired(currentLb, *newLB) {
		newLB.ID = currentLb.ID
		return c.Cloud.UpdateLoadBalancer(ctx, *newLB)
	}
	klog.V(4).Infof("No Load Balancer update required for %q, skipping", currentLb.ID)
	return currentLb, nil
}
