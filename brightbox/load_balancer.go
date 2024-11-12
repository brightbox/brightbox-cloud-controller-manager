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
	"fmt"
	"net"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	brightbox "github.com/brightbox/gobrightbox/v2"
	"github.com/brightbox/gobrightbox/v2/enums/balancingpolicy"
	"github.com/brightbox/gobrightbox/v2/enums/healthchecktype"
	"github.com/brightbox/gobrightbox/v2/enums/listenerprotocol"
	"github.com/brightbox/gobrightbox/v2/enums/proxyprotocol"
	"github.com/brightbox/k8ssdk/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/api/v1/service"
)

const (
	// Delete Backoff settings
	loadbalancerActiveInitDelay = 1 * time.Second
	loadbalancerActiveFactor    = 1.2
	loadbalancerActiveSteps     = 5

	// Listening protocols
	defaultLoadBalancerProtocol = listenerprotocol.Http

	// Default Proxy Protocol is none
	defaultProxyProtocol = 0

	standardSSLPort = 443

	// Healthcheck on http port if there are no endpoints for the loadbalancer
	defaultHealthCheckPort = 80

	// Maximum number of bits in unsigned integers specified in annotations.
	maxBits = 32

	// serviceAnnotationLoadBalancerBufferSize is the annotation used
	// on the server to specify the way balancing is done.
	// One of "least-connections", "round-robin" or "source-address"
	serviceAnnotationLoadBalancerPolicy = "service.beta.kubernetes.io/brightbox-load-balancer-policy"

	// ServiceAnnotationLoadBalancerListenerProtocol is the annotation used
	// on the service to specify the protocol spoken by the backend
	// (pod) behind a listener.
	// If `http` (default) or `http+ws`, an HTTP listener that terminates the
	// connection and parses headers is created.
	// If set to `TCP`, a "raw" listener is used.
	// The 'ws' extensions add support for Websockets to the listener.
	serviceAnnotationLoadBalancerListenerProtocol = "service.beta.kubernetes.io/brightbox-load-balancer-listener-protocol"

	// ServiceAnnotationLoadBalancerConnectionIdleTimeout is the
	// annotation used on the service to specify the idle connection
	// timeout.
	serviceAnnotationLoadBalancerListenerIdleTimeout = "service.beta.kubernetes.io/brightbox-load-balancer-listener-idle-timeout"

	// ServiceAnnotationLoadBalancerListenerProxyProtocol is the
	// annotation used on the service to activate the PROXY protocol to the backend
	// and specify the type of information that should be contained within it
	serviceAnnotationLoadBalancerListenerProxyProtocol = "service.beta.kubernetes.io/brightbox-load-balancer-listener-proxy-protocol"

	// ServiceAnnotationLoadBalancerSSLPorts is the annotation used on the service
	// to specify a comma-separated list of ports that will use SSL/HTTPS
	// listeners rather than plain 'http' listeners. Defaults to '443'.
	serviceAnnotationLoadBalancerSSLPorts = "service.beta.kubernetes.io/brightbox-load-balancer-ssl-ports"

	// ServiceAnnotationLoadBalancerSslDomains is the annotation used
	// on the service to specify the list of additional domains to add to the
	// Let's Encrypt SSL certificate used by the https listener.
	// The entry must be a comma separated list of DNS names that the
	// loadbalancer should accept as a target. These DNS names need to be
	// mapped externally onto the `Load Balancer Ingress` address
	// of the service, or via a CNAME onto the ingress address hostname
	serviceAnnotationLoadBalancerSslDomains = "service.beta.kubernetes.io/brightbox-load-balancer-ssl-domains"

	// ServiceAnnotationLoadBalancerCloudipAllocations is the
	// annotation used to specify the ID of the CloudIP that should
	// be mapped to the load balancer. It replaces the deprecated
	// `spec.loadBalancerIP` entry and should be in the form
	// `cip-xxxxx`. Only one cloudip can be specified.
	serviceAnnotationLoadBalancerCloudipAllocations = "service.beta.kubernetes.io/brightbox-load-balancer-cloudip-allocations"

	// ServiceAnnotationLoadBalancerHCHealthyThreshold is the
	// annotation used on the service to specify the number of successive
	// successful health checks required for a backend to be considered
	// healthy for traffic.
	serviceAnnotationLoadBalancerHCHealthyThreshold = "service.beta.kubernetes.io/brightbox-load-balancer-healthcheck-healthy-threshold"

	// ServiceAnnotationLoadBalancerHCUnhealthyThreshold is
	// the annotation used on the service to specify the number of
	// unsuccessful health checks required for a backend to be considered
	// unhealthy for traffic
	serviceAnnotationLoadBalancerHCUnhealthyThreshold = "service.beta.kubernetes.io/brightbox-load-balancer-healthcheck-unhealthy-threshold"

	// ServiceAnnotationLoadBalancerHeathcheckProtocol is the annotation used
	// on the service to specify the protocol used to do the healthcheck
	// Defaults to the same protocol as the listener
	serviceAnnotationLoadBalancerHCProtocol = "service.beta.kubernetes.io/brightbox-load-balancer-healthcheck-protocol"

	// ServiceAnnotationLoadBalancerHeathcheckRequest is the annotation used
	// on the service to specify the request path an http healthcheck should use to talk to the backend
	// Defaults to the kubernetes specified standard (currently '/healthz')
	serviceAnnotationLoadBalancerHCRequest = "service.beta.kubernetes.io/brightbox-load-balancer-healthcheck-request"

	// ServiceAnnotationLoadBalancerHCTimeout is the annotation used
	// on the service to specify, in seconds, how long to wait before
	// marking a health check as failed.
	serviceAnnotationLoadBalancerHCTimeout = "service.beta.kubernetes.io/brightbox-load-balancer-healthcheck-timeout"

	// ServiceAnnotationLoadBalancerHCInterval is the annotation used on the
	// service to specify, in seconds, the interval between health checks.
	serviceAnnotationLoadBalancerHCInterval = "service.beta.kubernetes.io/aws-load-balancer-healthcheck-interval"
)

var (
	truevar        = true
	falsevar       = false
	cloudIPPattern = regexp.MustCompile(`^cip-[0-9a-z]{5,}$`)
)

// Return a name that is 'name'.'namespace'.'clusterName'
// Use the default name derived from the UID if no name field is set
func (c *cloud) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	namespace := service.Namespace
	if namespace == "" {
		namespace = "default"
	}
	name := service.Name
	if name == "" {
		name = cloudprovider.DefaultLoadBalancerName(service)
	}
	var buffer bytes.Buffer
	buffer.WriteString(name)
	buffer.WriteString(".")
	buffer.WriteString(namespace)
	buffer.WriteString(".")
	buffer.WriteString(clusterName)
	return buffer.String()
}

func (c *cloud) GetLoadBalancer(ctx context.Context, clusterName string, apiservice *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	name := c.GetLoadBalancerName(ctx, clusterName, apiservice)
	if err := logAction(ctx, "GetLoadBalancer(%v)", name); err != nil {
		return nil, false, err
	}
	lb, err := c.GetLoadBalancerByName(ctx, name)
	return toLoadBalancerStatus(lb), err == nil && lb != nil, err
}

// Make sure we have a cloud ip before asking for a load balancer. Try
// to get one matching the LoadBalancerIP spec in the service, and error
// if that isn't in the cloudip list.
func (c *cloud) EnsureLoadBalancer(ctx context.Context, clusterName string, apiservice *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	name := c.GetLoadBalancerName(ctx, clusterName, apiservice)
	if err := logAction(ctx, "EnsureLoadBalancer(%v, %v, %v, %v)", name, apiservice.Spec.LoadBalancerIP, apiservice.Spec.Ports, apiservice.Annotations); err != nil {
		return nil, err
	}
	if err := validateServiceSpec(apiservice); err != nil {
		return nil, err
	}
	cip, err := c.ensureAllocatedCloudIP(ctx, name, apiservice)
	if err != nil {
		return nil, err
	}
	domains, err := ensureLoadBalancerDomainResolution(apiservice.Annotations, cip)
	if err != nil {
		return nil, err
	}
	lb, err := c.ensureLoadBalancerFromService(ctx, name, domains, apiservice, nodes)
	if err != nil {
		return nil, err
	}
	err = c.EnsureMappedCloudIP(ctx, lb, cip)
	if err != nil {
		return nil, err
	}
	err = c.EnsureOldCloudIPsDeposed(ctx, lb.CloudIPs, cip.ID)
	if err != nil {
		return nil, err
	}
	if err := c.ensureCloudIPsDeleted(ctx, cip.ID, name); err != nil {
		return nil, err
	}
	lb, err = c.GetLoadBalancerByID(ctx, lb.ID)
	if err != nil {
		return nil, err
	}
	return toLoadBalancerStatus(lb), k8ssdk.ErrorIfNotComplete(lb, cip.ID, name)
}

func (c *cloud) UpdateLoadBalancer(ctx context.Context, clusterName string, apiservice *v1.Service, nodes []*v1.Node) error {
	if err := logAction(ctx, "UpdateLoadBalancer called - delegating"); err != nil {
		return err
	}
	_, err := c.EnsureLoadBalancer(ctx, clusterName, apiservice, nodes)
	return err
}

func (c *cloud) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, apiservice *v1.Service) error {
	name := c.GetLoadBalancerName(ctx, clusterName, apiservice)
	if err := logAction(ctx, "EnsureLoadBalancerDeleted(%v, %v)", name, apiservice.Spec.LoadBalancerIP); err != nil {
		return err
	}
	if err := c.ensureServerGroupDeleted(ctx, name); err != nil {
		return err
	}
	if err := c.ensureFirewallClosed(ctx, name); err != nil {
		return err
	}
	lb, err := c.ensureLoadBalancerDeletedByName(ctx, name)
	if err != nil {
		return err
	}
	if err := c.ensureCloudIPsDeleted(ctx, "", name); err != nil {
		return err
	}
	if lb != nil {
		lb, err = c.GetLoadBalancerByID(ctx, lb.ID)
		if err != nil {
			return err
		}
	}
	return k8ssdk.ErrorIfNotErased(lb)
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

// Try to remove CloudIPs matching `name` from the list of cloudIPs
func (c *cloud) ensureCloudIPsDeleted(ctx context.Context, currentID string, name string) error {
	klog.V(4).Infof("ensureCloudIPsDeleted (%q)", name)
	backoff := wait.Backoff{
		Duration: loadbalancerActiveInitDelay,
		Factor:   loadbalancerActiveFactor,
		Steps:    loadbalancerActiveSteps,
	}

	return wait.ExponentialBackoff(backoff, func() (bool, error) {
		cloudIPList, err := c.GetCloudIPs(ctx)
		if err != nil {
			klog.V(4).Info("Error retrieving list of CloudIPs")
			return false, err
		}
		if err := c.DestroyCloudIPs(ctx, cloudIPList, currentID, name); err != nil {
			klog.V(4).Info(err)
			return false, nil
		}
		return true, nil
	},
	)
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

func validateServiceSpec(apiservice *v1.Service) error {
	if apiservice.Spec.SessionAffinity != v1.ServiceAffinityNone {
		return fmt.Errorf("unsupported load balancer affinity: %v", apiservice.Spec.SessionAffinity)
	}
	if len(apiservice.Spec.Ports) == 0 {
		return fmt.Errorf("requested load balancer with no ports")
	}
	protocol := getListenerProtocol(apiservice)
	sslPortFound := false
	for _, port := range apiservice.Spec.Ports {
		if port.Protocol != v1.ProtocolTCP {
			return fmt.Errorf("UDP nodeports are not supported")
		}
		sslPortFound = sslPortFound || port.Port == standardSSLPort
	}
	if !sslPortFound && protocol == listenerprotocol.Http {
		_, ports := apiservice.Annotations[serviceAnnotationLoadBalancerSSLPorts]
		_, domains := apiservice.Annotations[serviceAnnotationLoadBalancerSslDomains]
		if ports || domains {
			return fmt.Errorf("SSL support requires a Port definition for %d", standardSSLPort)
		}
	}
	// CloudIP allocation annotation and spec.loadBalancerIP conflict
	if apiservice.Spec.LoadBalancerIP != "" {
		if _, ok := apiservice.Annotations[serviceAnnotationLoadBalancerCloudipAllocations]; ok {
			return fmt.Errorf("Remove obsolete field: spec.loadBalancerIP")
		}
	}
	return validateAnnotations(apiservice.Annotations)
}

func validateAnnotations(annotationList map[string]string) error {
	for annotation, value := range annotationList {
		switch annotation {
		case serviceAnnotationLoadBalancerPolicy:
			if _, err := balancingpolicy.ParseEnum(value); err != nil {
				return fmt.Errorf("Invalid Load Balancer Policy %q: %w", value, err)
			}
		case serviceAnnotationLoadBalancerListenerProtocol:
			valueEnum, err := listenerprotocol.ParseEnum(value)
			if err != nil {
				return fmt.Errorf("Invalid Load Balancer Listener Protocol %q: %w", value, err)
			}
			if valueEnum == listenerprotocol.Tcp {
				if _, ok := annotationList[serviceAnnotationLoadBalancerSSLPorts]; ok {
					return fmt.Errorf("SSL Ports are not supported with the %s protocol", valueEnum)
				}
				if _, ok := annotationList[serviceAnnotationLoadBalancerSslDomains]; ok {
					return fmt.Errorf("SSL Domains are not supported with the %s protocol", valueEnum)
				}
			}
		case serviceAnnotationLoadBalancerListenerProxyProtocol:
			if _, err := proxyprotocol.ParseEnum(value); err != nil {
				return fmt.Errorf("Invalid Load Balancer Listener Proxy Protocol %q: %w", value, err)
			}
		case serviceAnnotationLoadBalancerSSLPorts:
			if _, ok := annotationList[serviceAnnotationLoadBalancerSslDomains]; !ok {
				return fmt.Errorf("SSL needs a list of domains to certify. Add the %q annotation", serviceAnnotationLoadBalancerSslDomains)
			}
		case serviceAnnotationLoadBalancerHCProtocol:
			if _, err := healthchecktype.ParseEnum(value); err != nil {
				return fmt.Errorf("Invalid Load Balancer Healthcheck Protocol %q: %w", value, err)
			}
		case serviceAnnotationLoadBalancerHCInterval,
			serviceAnnotationLoadBalancerHCTimeout,
			serviceAnnotationLoadBalancerHCHealthyThreshold,
			serviceAnnotationLoadBalancerHCUnhealthyThreshold,
			serviceAnnotationLoadBalancerListenerIdleTimeout:
			_, err := parseUintAnnotation(annotationList, annotation)
			if err != nil {
				return fmt.Errorf("%q needs to be a positive number (%v)", annotation, err)
			}
		case serviceAnnotationLoadBalancerHCRequest:
			testURL := "http://example.com:6443" + value
			u, err := url.Parse(testURL)
			if err != nil || u.Path != value {
				return fmt.Errorf("%q needs to be a valid Url request path", annotation)
			}
		case serviceAnnotationLoadBalancerCloudipAllocations:
			if !cloudIPPattern.MatchString(value) {
				return fmt.Errorf("%q needs to match the pattern %q", annotation, cloudIPPattern)
			}
		}
	}
	return nil
}

func ensureLoadBalancerDomainResolution(annotationList map[string]string, cloudIP *brightbox.CloudIP) ([]string, error) {
	domains := append(extraLoadBalancerDomains(annotationList), cloudIP.Fqdn, cloudIP.ReverseDNS)
	slices.Sort(domains)
	domains = slices.Compact(domains)
	cloudIPList, err := toIPList(cloudIP)
	if err != nil {
		return nil, err
	}
	for _, domain := range domains {
		resolvedAddresses, err := net.LookupIP(domain)
		if err != nil {
			return nil, fmt.Errorf("Failed to resolve %q to load balancer address (%s,%s): %v", domain, cloudIP.PublicIPv4, cloudIP.PublicIPv6, err.Error())
		}
		if !anyAddressMatch(cloudIPList, resolvedAddresses) {
			return nil, fmt.Errorf("Failed to resolve %q to load balancer address (%s,%s)", domain, cloudIP.PublicIPv4, cloudIP.PublicIPv6)
		}
	}
	return domains, nil
}

func toIPList(cloudIP *brightbox.CloudIP) ([]net.IP, error) {
	result := append([]net.IP{}, net.ParseIP(cloudIP.PublicIPv4), net.ParseIP(cloudIP.PublicIPv6))
	if result[0] == nil || result[1] == nil {
		return nil, fmt.Errorf("Cloud IP %q failed to parse IP addresses", cloudIP.ID)
	}
	return result, nil
}

func anyAddressMatch(ipListA, ipListB []net.IP) bool {
	for a := range ipListA {
		for b := range ipListB {
			if ipListA[a].Equal(ipListB[b]) {
				return true
			}
		}
	}
	return false
}

func (c *cloud) ensureAllocatedCloudIP(ctx context.Context, name string, apiservice *v1.Service) (*brightbox.CloudIP, error) {
	klog.V(4).Info("ensureAllocatedCloudIP")
	if cipID, ok := apiservice.Annotations[serviceAnnotationLoadBalancerCloudipAllocations]; ok {
		return c.GetCloudIP(ctx, cipID)
	}
	if ip := apiservice.Spec.LoadBalancerIP; ip != "" {
		return lookupCloudIPByIP(ctx, c, ip)
	}
	return lookupCloudIPByName(ctx, c, name)
}

func lookupCloudIPByIP(ctx context.Context, c *cloud, ip string) (*brightbox.CloudIP, error) {
	ipval := net.ParseIP(ip)
	if ipval == nil {
		return nil, fmt.Errorf("Invalid LoadBalancerIP: %q", ip)
	}

	cloudIPList, err := c.GetCloudIPs(ctx)
	if err != nil {
		return nil, err
	}

	cip := findMatchingCloudIP(cloudIPList, func(cip *brightbox.CloudIP) bool {
		return ipval.Equal(net.ParseIP(cip.PublicIPv4)) || ipval.Equal(net.ParseIP(cip.PublicIPv6))
	})
	if cip == nil {
		return nil, fmt.Errorf("Could not find allocated Cloud IP with address %q", ip)
	}
	klog.Warningf("spec.loadBalancerIP is deprecated. Remove the entry and add the annotation: %s=%s", serviceAnnotationLoadBalancerCloudipAllocations, cip.ID)
	return cip, nil
}

func lookupCloudIPByName(ctx context.Context, c *cloud, name string) (*brightbox.CloudIP, error) {
	cloudIPList, err := c.GetCloudIPs(ctx)
	if err != nil {
		return nil, err
	}

	cip := findMatchingCloudIP(cloudIPList, func(cip *brightbox.CloudIP) bool {
		return cip.Name == name || (cip.LoadBalancer != nil && cip.LoadBalancer.Name == name)
	})

	if cip == nil {
		return c.AllocateCloudIP(ctx, name)
	}
	return cip, nil
}

func findMatchingCloudIP(cloudIPList []brightbox.CloudIP, matches func(*brightbox.CloudIP) bool) *brightbox.CloudIP {
	for i := range cloudIPList {
		if matches(&cloudIPList[i]) {
			return &cloudIPList[i]
		}
	}
	return nil
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

func extraLoadBalancerDomains(annotations map[string]string) []string {
	if domains, ok := annotations[serviceAnnotationLoadBalancerSslDomains]; ok {
		return strings.Split(domains, ",")
	}
	return nil
}

func getListenerProtocol(apiservice *v1.Service) listenerprotocol.Enum {
	if protocol, ok := apiservice.Annotations[serviceAnnotationLoadBalancerListenerProtocol]; ok {
		if protocolEnum, err := listenerprotocol.ParseEnum(protocol); err == nil {
			return protocolEnum
		}
	}
	return defaultLoadBalancerProtocol
}

func getListenerProxyProtocol(apiservice *v1.Service) proxyprotocol.Enum {
	if protocol, ok := apiservice.Annotations[serviceAnnotationLoadBalancerListenerProxyProtocol]; ok {
		if protocolEnum, err := proxyprotocol.ParseEnum(protocol); err == nil {
			return protocolEnum
		}
	}
	return defaultProxyProtocol
}

func isSSLPort(port *v1.ServicePort, sslPorts *portSets) bool {
	return port.Port == standardSSLPort ||
		sslPorts != nil && (sslPorts.numbers.Has(int64(port.Port)) || sslPorts.names.Has(port.Name))
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

func getHealthCheckPath(apiservice *v1.Service, protocol healthchecktype.Enum, path string) string {
	if protocol == healthchecktype.Tcp {
		return "/"
	}
	if request, ok := apiservice.Annotations[serviceAnnotationLoadBalancerHCRequest]; ok {
		return request
	}
	if path == "" {
		return "/healthz"
	}
	return path
}

func getServiceHealthCheckPathPort(apiservice *v1.Service) (string, int32) {
	if !service.NeedsHealthCheck(apiservice) {
		return "", 0
	}
	port := apiservice.Spec.HealthCheckNodePort
	if port == 0 {
		return "", 0
	}
	return "/healthz", port
}

func getHealthCheckProtocol(apiservice *v1.Service, path string) healthchecktype.Enum {
	if protocol, ok := apiservice.Annotations[serviceAnnotationLoadBalancerHCProtocol]; ok {
		if protocolEnum, err := healthchecktype.ParseEnum(protocol); err == nil {
			return protocolEnum
		}
	}
	if getListenerProtocol(apiservice) == listenerprotocol.Tcp && path == "" {
		return healthchecktype.Tcp
	}
	return healthchecktype.Http
}

func getHealthCheckPort(apiservice *v1.Service, nodeport uint16) uint16 {
	if nodeport != 0 {
		return nodeport
	}
	for i := range apiservice.Spec.Ports {
		return uint16(apiservice.Spec.Ports[i].NodePort)
	}
	return defaultHealthCheckPort
}

// If annotation is missing returns zero value
func parseUintAnnotation(annotationList map[string]string, annotation string) (uint, error) {
	klog.V(6).Infof("parseUintAnnotation(%+v, %+v)", annotationList, annotation)
	strValue, ok := annotationList[annotation]
	if !ok {
		return 0, nil
	}
	val, err := strconv.ParseUint(strValue, 10, maxBits)
	klog.V(6).Infof("Value Converted from %+v to %+v", strValue, val)
	return uint(val), err
}
