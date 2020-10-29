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
	"strconv"
	"strings"

	brightbox "github.com/brightbox/gobrightbox"
	"github.com/brightbox/k8ssdk"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/api/v1/service"
)

const (
	// Listening protocols
	loadBalancerTCPProtocol     = "tcp"
	loadBalancerHTTPProtocol    = "http"
	loadBalancerHTTPWSProtocol  = "http+ws"
	defaultLoadBalancerProtocol = loadBalancerHTTPProtocol

	// Proxy protocols
	loadBalancerProxyV1      = "v1"
	loadBalancerProxyV2      = "v2"
	loadBalancerProxyV2Ssl   = "v2-ssl"
	loadBalancerProxyV2SslCn = "v2-ssl-cn"

	standardSSLPort = 443

	// Healthcheck on http port if there are no endpoints for the loadbalancer
	defaultHealthCheckPort = 80

	// Maximum number of bits in unsigned integers specified in annotations.
	maxBits = 32

	// The minimum size of the buffer in a load balancer
	validMinimumBufferSize = 1024

	// The maximum size of the buffer in a load balancer
	validMaximumBufferSize = 16384

	// serviceAnnotationLoadBalancerBufferSize is the annotation used
	// on the server to specify the way balancing is done.
	// One of "least-connections", "round-robin" or "source-address"
	serviceAnnotationLoadBalancerPolicy = "service.beta.kubernetes.io/brightbox-load-balancer-policy"

	// serviceAnnotationLoadBalancerBufferSize is the annotation used
	// on the server to specify the size of the receive buffer for the service in bytes.
	// This is subject to a minimum size of 1024 bytes.
	serviceAnnotationLoadBalancerBufferSize = "service.beta.kubernetes.io/brightbox-load-balancer-buffer-size"

	// ServiceAnnotationLoadBalancerListenerProtocol is the annotation used
	// on the service to specify the protocol spoken by the backend
	// (pod) behind a listener.
	// If `http` (default) or `http+ws`, an HTTP listener that terminates the
	// connection and parses headers is created.
	// If set to `TCP`, a "raw" listener is used.
	// The 'ws' extensions add support for Websockets to the listener.
	serviceAnnotationLoadBalancerListenerProtocol = "service.beta.kubernetes.io/brightbox-load-balancer-listener-protocol"

	// ServiceAnnotationLoadBalancerSSLPorts is the annotation used on the service
	// to specify a comma-separated list of ports that will use SSL/HTTPS
	// listeners rather than plain 'http' listeners. Defaults to '443'.
	serviceAnnotationLoadBalancerSSLPorts = "service.beta.kubernetes.io/brightbox-load-balancer-ssl-ports"

	// ServiceAnnotationLoadBalancerConnectionIdleTimeout is the
	// annotation used on the service to specify the idle connection
	// timeout.
	serviceAnnotationLoadBalancerListenerIdleTimeout = "service.beta.kubernetes.io/brightbox-load-balancer-listener-idle-timeout"

	// ServiceAnnotationLoadBalancerListenerProxyProtocol is the
	// annotation used on the service to activate the PROXY protocol to the backend
	// and specify the type of information that should be contained within it
	serviceAnnotationLoadBalancerListenerProxyProtocol = "service.beta.kubernetes.io/brightbox-load-balancer-listener-proxy-protocol"

	// ServiceAnnotationLoadBalancerSslDomains is the annotation used
	// on the service to specify the list of additional domains to add to the
	// Let's Encrypt SSL certificate used by the https listener.
	// The entry must be a comma separated list of DNS names that the
	// loadbalancer should accept as a target. These DNS names need to be
	// mapped externally onto the `Load Balancer Ingress` address
	// of the service, or via a CNAME onto the ingress address hostname
	serviceAnnotationLoadBalancerSslDomains = "service.beta.kubernetes.io/brightbox-load-balancer-ssl-domains"

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
	validLoadBalancerPolicies = map[string]bool{
		"least-connections": true,
		"round-robin":       true,
		"source-address":    true,
	}
	validListenerProtocols = map[string]bool{
		loadBalancerHTTPProtocol:   true,
		loadBalancerTCPProtocol:    true,
		loadBalancerHTTPWSProtocol: true,
	}
	validHealthCheckProtocols = map[string]bool{
		loadBalancerHTTPProtocol: true,
		loadBalancerTCPProtocol:  true,
	}
	sslUpgradeProtocol = map[string]string{
		loadBalancerTCPProtocol:    loadBalancerTCPProtocol,
		loadBalancerHTTPProtocol:   "https",
		loadBalancerHTTPWSProtocol: "https+wss",
	}
	validListenerProxyProtocols = map[string]bool{
		loadBalancerProxyV1:      true,
		loadBalancerProxyV2:      true,
		loadBalancerProxyV2Ssl:   true,
		loadBalancerProxyV2SslCn: true,
	}
	truevar  = true
	falsevar = false
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
	klog.V(4).Infof("GetLoadBalancer(%v)", name)
	lb, err := c.GetLoadBalancerByName(name)
	return toLoadBalancerStatus(lb), err == nil && lb != nil, err
}

// Make sure we have a cloud ip before asking for a load balancer. Try
// to get one matching the LoadBalancerIP spec in the service, and error
// if that isn't in the cloudip list.
func (c *cloud) EnsureLoadBalancer(ctx context.Context, clusterName string, apiservice *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	name := c.GetLoadBalancerName(ctx, clusterName, apiservice)
	klog.V(4).Infof("EnsureLoadBalancer(%v, %v, %v, %v)", name, apiservice.Spec.LoadBalancerIP, apiservice.Spec.Ports, apiservice.Annotations)
	if err := validateServiceSpec(apiservice); err != nil {
		return nil, err
	}
	cip, err := c.ensureAllocatedCloudIP(name, apiservice)
	if err != nil {
		return nil, err
	}
	if err := validateContextualAnnotations(apiservice.Annotations, cip); err != nil {
		return nil, err
	}
	lb, err := c.ensureLoadBalancerFromService(name, apiservice, nodes)
	if err != nil {
		return nil, err
	}
	err = c.EnsureMappedCloudIP(lb, cip)
	if err != nil {
		return nil, err
	}
	if apiservice.Spec.LoadBalancerIP != "" {
		err = c.EnsureOldCloudIPsDeposed(lb.CloudIPs, cip.Id, name)
		if err != nil {
			return nil, err
		}
		if err := c.ensureCloudIPsDeleted(name); err != nil {
			return nil, err
		}
	}
	lb, err = c.GetLoadBalancerByID(lb.Id)
	if err != nil {
		return nil, err
	}
	return toLoadBalancerStatus(lb), k8ssdk.ErrorIfNotComplete(lb, cip.Id, name)
}

func (c *cloud) UpdateLoadBalancer(ctx context.Context, clusterName string, apiservice *v1.Service, nodes []*v1.Node) error {
	klog.V(4).Infof("UpdateLoadBalancer called - delegating")
	_, err := c.EnsureLoadBalancer(ctx, clusterName, apiservice, nodes)
	return err
}

func (c *cloud) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, apiservice *v1.Service) error {
	name := c.GetLoadBalancerName(ctx, clusterName, apiservice)
	klog.V(4).Infof("EnsureLoadBalancerDeleted(%v, %v)", name, apiservice.Spec.LoadBalancerIP)
	if err := c.ensureServerGroupDeleted(name); err != nil {
		return err
	}
	if err := c.ensureFirewallClosed(name); err != nil {
		return err
	}
	lb, err := c.ensureLoadBalancerDeletedByName(name)
	if err != nil {
		return err
	}
	if err := c.ensureCloudIPsDeleted(name); err != nil {
		return err
	}
	if lb != nil {
		lb, err = c.GetLoadBalancerByID(lb.Id)
		if err != nil {
			return err
		}
	}
	return k8ssdk.ErrorIfNotErased(lb)
}

//Take all the servers out of the server group and remove it
func (c *cloud) ensureServerGroupDeleted(name string) error {
	klog.V(4).Infof("ensureServerGroupDeleted (%q)", name)
	group, err := c.GetServerGroupByName(name)
	if err != nil {
		klog.V(4).Infof("Error looking for Server Group for %q", name)
		return err
	}
	if group == nil {
		return nil
	}
	group, err = c.SyncServerGroup(group, nil)
	if err != nil {
		klog.V(4).Infof("Error removing servers from %q", group.Id)
		return err
	}
	if err := c.DestroyServerGroup(group.Id); err != nil {
		klog.V(4).Infof("Error destroying Server Group %q", group.Id)
		return err
	}
	return nil
}

//Remove the firewall policy
func (c *cloud) ensureFirewallClosed(name string) error {
	klog.V(4).Infof("ensureFirewallClosed (%q)", name)
	fp, err := c.GetFirewallPolicyByName(name)
	if err != nil {
		klog.V(4).Infof("Error looking for Firewall Policy %q", name)
		return err
	}
	if fp == nil {
		return nil
	}
	if err := c.DestroyFirewallPolicy(fp.Id); err != nil {
		klog.V(4).Infof("Error destroying Firewall Policy %q", fp.Id)
		return err
	}
	return nil
}

//Remove load balancer by name
func (c *cloud) ensureLoadBalancerDeletedByName(name string) (*brightbox.LoadBalancer, error) {
	lb, err := c.GetLoadBalancerByName(name)
	if err != nil {
		klog.V(4).Infof("Error looking for Load Balancer %q", name)
		return nil, err
	}
	if lb != nil {
		if err = c.DestroyLoadBalancer(lb.Id); err != nil {
			klog.V(4).Infof("Error destroying Load Balancer %q", lb.Id)
			return nil, err
		}
	}
	return lb, nil
}

//Try to remove CloudIPs matching `name` from the list of cloudIPs
func (c *cloud) ensureCloudIPsDeleted(name string) error {
	klog.V(4).Infof("ensureCloudIPsDeleted (%q)", name)
	cloudIPList, err := c.GetCloudIPs()
	if err != nil {
		klog.V(4).Infof("Error retrieving list of CloudIPs")
		return err
	}
	return c.DestroyCloudIPs(cloudIPList, name)
}

func toLoadBalancerStatus(lb *brightbox.LoadBalancer) *v1.LoadBalancerStatus {
	if lb == nil {
		return nil
	}
	status := v1.LoadBalancerStatus{}
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
			if v.ReverseDns != "" {
				status.Ingress = append(status.Ingress,
					v1.LoadBalancerIngress{
						Hostname: v.ReverseDns,
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
	if !sslPortFound && protocol == loadBalancerHTTPProtocol {
		_, ports := apiservice.Annotations[serviceAnnotationLoadBalancerSSLPorts]
		_, domains := apiservice.Annotations[serviceAnnotationLoadBalancerSslDomains]
		if ports || domains {
			return fmt.Errorf("SSL support requires a Port definition for %d", standardSSLPort)
		}
	}
	return validateAnnotations(apiservice.Annotations)
}

func validateAnnotations(annotationList map[string]string) error {
	for annotation, value := range annotationList {
		switch annotation {
		case serviceAnnotationLoadBalancerPolicy:
			if !validLoadBalancerPolicies[value] {
				return fmt.Errorf("Invalid Load Balancer Policy %q", value)
			}
		case serviceAnnotationLoadBalancerListenerProtocol:
			if !validListenerProtocols[value] {
				return fmt.Errorf("Invalid Load Balancer Listener Protocol %q", value)
			}
			if value == loadBalancerTCPProtocol {
				if _, ok := annotationList[serviceAnnotationLoadBalancerSSLPorts]; ok {
					return fmt.Errorf("SSL Ports are not supported with the %s protocol", loadBalancerTCPProtocol)
				}
				if _, ok := annotationList[serviceAnnotationLoadBalancerSslDomains]; ok {
					return fmt.Errorf("SSL Domains are not supported with the %s protocol", loadBalancerTCPProtocol)
				}
			}
		case serviceAnnotationLoadBalancerListenerProxyProtocol:
			if !validListenerProxyProtocols[value] {
				return fmt.Errorf("Invalid Load Balancer Listener Proxy Protocol %q", value)
			}
		case serviceAnnotationLoadBalancerSSLPorts:
			if _, ok := annotationList[serviceAnnotationLoadBalancerSslDomains]; !ok {
				return fmt.Errorf("SSL needs a list of domains to certify. Add the %q annotation", serviceAnnotationLoadBalancerSslDomains)
			}
		case serviceAnnotationLoadBalancerHCProtocol:
			if !validHealthCheckProtocols[value] {
				return fmt.Errorf("Invalid Load Balancer Healthcheck Protocol %q", value)
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
		case serviceAnnotationLoadBalancerBufferSize:
			val, err := parseUintAnnotation(annotationList, annotation)
			if err != nil {
				return fmt.Errorf("%q needs to be a positive number (%v)", annotation, err)
			}
			if val < validMinimumBufferSize {
				return fmt.Errorf("%q needs to be no less than %d", annotation, validMinimumBufferSize)
			}
			if val > validMaximumBufferSize {
				return fmt.Errorf("%q needs to be no more than %d", annotation, validMaximumBufferSize)
			}
		case serviceAnnotationLoadBalancerHCRequest:
			testURL := "http://example.com:6443" + value
			u, err := url.Parse(testURL)
			if err != nil || u.Path != value {
				return fmt.Errorf("%q needs to be a valid Url request path", annotation)
			}
		}
	}
	return nil
}

func validateContextualAnnotations(annotationList map[string]string, cloudIP *brightbox.CloudIP) error {
	domains := buildLoadBalancerDomains(annotationList)
	if domains != nil {
		cloudIPList, err := toIPList(cloudIP)
		if err != nil {
			return err
		}
		for _, domain := range domains {
			resolvedAddresses, err := net.LookupIP(domain)
			if err != nil {
				return fmt.Errorf("Failed to resolve %q to load balancer address (%s,%s): %v", domain, cloudIP.PublicIPv4, cloudIP.PublicIPv6, err.Error())
			}
			if !anyAddressMatch(cloudIPList, resolvedAddresses) {
				return fmt.Errorf("Failed to resolve %q to load balancer address (%s,%s)", domain, cloudIP.PublicIPv4, cloudIP.PublicIPv6)
			}
		}
	}
	return nil
}

func toIPList(cloudIP *brightbox.CloudIP) ([]net.IP, error) {
	result := append([]net.IP{}, net.ParseIP(cloudIP.PublicIPv4), net.ParseIP(cloudIP.PublicIPv6))
	if result[0] == nil || result[1] == nil {
		return nil, fmt.Errorf("Cloud IP %q failed to parse IP addresses", cloudIP.Id)
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

func (c *cloud) ensureAllocatedCloudIP(name string, apiservice *v1.Service) (*brightbox.CloudIP, error) {
	ip := apiservice.Spec.LoadBalancerIP
	klog.V(4).Infof("ensureAllocatedCloudIP (%q, %q)", name, ip)
	var compareFunc func(cip *brightbox.CloudIP) bool
	switch ip {
	case "":
		compareFunc = func(cip *brightbox.CloudIP) bool {
			return cip.Name == name
		}
	default:
		ipval := net.ParseIP(ip)
		if ipval == nil {
			return nil, fmt.Errorf("Invalid LoadBalancerIP: %q", ip)
		}
		compareFunc = func(cip *brightbox.CloudIP) bool {
			return ipval.Equal(net.ParseIP(cip.PublicIPv4)) || ipval.Equal(net.ParseIP(cip.PublicIPv6))
		}
	}
	cloudIPList, err := c.GetCloudIPs()
	if err != nil {
		return nil, err
	}
	for i := range cloudIPList {
		if compareFunc(&cloudIPList[i]) {
			return &cloudIPList[i], nil
		}
	}
	if ip == "" {
		return c.AllocateCloudIP(name)
	}
	return nil, fmt.Errorf("Could not find allocated Cloud IP with address %q", ip)
}

func (c *cloud) ensureLoadBalancerFromService(name string, apiservice *v1.Service, nodes []*v1.Node) (*brightbox.LoadBalancer, error) {
	klog.V(4).Infof("ensureLoadBalancerFromService(%v)", name)
	currentLb, err := c.GetLoadBalancerByName(name)
	if err != nil {
		return nil, err
	}
	err = c.ensureFirewallOpenForService(name, apiservice, nodes)
	if err != nil {
		return nil, err
	}
	newLB := buildLoadBalancerOptions(name, apiservice, nodes)
	if currentLb == nil {
		return c.Cloud.CreateLoadBalancer(newLB)
	} else if k8ssdk.IsUpdateLoadBalancerRequired(currentLb, *newLB) {
		newLB.Id = currentLb.Id
		return c.Cloud.UpdateLoadBalancer(newLB)
	}
	klog.V(4).Infof("No Load Balancer update required for %q, skipping", currentLb.Id)
	return currentLb, nil
}

func buildLoadBalancerOptions(name string, apiservice *v1.Service, nodes []*v1.Node) *brightbox.LoadBalancerOptions {
	klog.V(4).Infof("buildLoadBalancerOptions(%v)", name)
	result := &brightbox.LoadBalancerOptions{
		Name:        &name,
		Nodes:       buildLoadBalancerNodes(nodes),
		Listeners:   buildLoadBalancerListeners(apiservice),
		Healthcheck: buildLoadBalancerHealthCheck(apiservice),
		Domains:     buildLoadBalancerDomains(apiservice.Annotations),
	}
	bufferSize, _ := parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerBufferSize)
	if bufferSize != 0 {
		result.BufferSize = &bufferSize
	}
	if policy, ok := apiservice.Annotations[serviceAnnotationLoadBalancerPolicy]; ok {
		result.Policy = &policy
	}
	if result.Domains != nil {
		result.HttpsRedirect = &truevar
	} else {
		result.HttpsRedirect = &falsevar
	}
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
		if result[i].Protocol != loadBalancerTCPProtocol && isSSLPort(&apiservice.Spec.Ports[i], sslPortSet) {
			result[i].Protocol = sslUpgradeProtocol[result[i].Protocol]
		}
		result[i].In = int(apiservice.Spec.Ports[i].Port)
		result[i].Out = int(apiservice.Spec.Ports[i].NodePort)
		result[i].Timeout, _ = parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerListenerIdleTimeout)
	}
	return result
}

func buildLoadBalancerDomains(annotations map[string]string) []string {
	if domains, ok := annotations[serviceAnnotationLoadBalancerSslDomains]; ok {
		return strings.Split(domains, ",")
	}
	return nil
}

func getListenerProtocol(apiservice *v1.Service) string {
	if protocol, ok := apiservice.Annotations[serviceAnnotationLoadBalancerListenerProtocol]; ok {
		return protocol
	}
	return defaultLoadBalancerProtocol
}

func getListenerProxyProtocol(apiservice *v1.Service) string {
	return apiservice.Annotations[serviceAnnotationLoadBalancerListenerProxyProtocol]
}

func isSSLPort(port *v1.ServicePort, sslPorts *portSets) bool {
	return port.Port == standardSSLPort ||
		sslPorts != nil && (sslPorts.numbers.Has(int64(port.Port)) || sslPorts.names.Has(port.Name))
}

func buildLoadBalancerHealthCheck(apiservice *v1.Service) *brightbox.LoadBalancerHealthcheck {
	path, healthCheckNodePort := service.GetServiceHealthCheckPathPort(apiservice)
	protocol := getHealthCheckProtocol(apiservice, path)
	//Validate has already checked all these so there should be no errors!
	interval, _ := parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerHCInterval)
	timeout, _ := parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerHCTimeout)
	thresholdUp, _ := parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerHCHealthyThreshold)
	thresholdDown, _ := parseUintAnnotation(apiservice.Annotations, serviceAnnotationLoadBalancerHCUnhealthyThreshold)
	return &brightbox.LoadBalancerHealthcheck{
		Type:          protocol,
		Port:          getHealthCheckPort(apiservice, int(healthCheckNodePort)),
		Request:       getHealthCheckPath(apiservice, protocol, path),
		Interval:      interval,
		Timeout:       timeout,
		ThresholdUp:   thresholdUp,
		ThresholdDown: thresholdDown,
	}
}

func getHealthCheckPath(apiservice *v1.Service, protocol string, path string) string {
	if protocol == loadBalancerTCPProtocol {
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

func getHealthCheckProtocol(apiservice *v1.Service, path string) string {
	if protocol, ok := apiservice.Annotations[serviceAnnotationLoadBalancerHCProtocol]; ok {
		return protocol
	}
	if getListenerProtocol(apiservice) == loadBalancerTCPProtocol && path == "" {
		return loadBalancerTCPProtocol
	}
	return loadBalancerHTTPProtocol
}

func getHealthCheckPort(apiservice *v1.Service, nodeport int) int {
	if nodeport != 0 {
		return nodeport
	}
	for i := range apiservice.Spec.Ports {
		return int(apiservice.Spec.Ports[i].NodePort)
	}
	return defaultHealthCheckPort
}

//If annotation is missing returns zero value
func parseUintAnnotation(annotationList map[string]string, annotation string) (int, error) {
	klog.V(6).Infof("parseUintAnnotation(%+v, %+v)", annotationList, annotation)
	strValue, ok := annotationList[annotation]
	if !ok {
		return 0, nil
	}
	val, err := strconv.ParseUint(strValue, 10, maxBits)
	klog.V(6).Infof("Value Converted from %+v to %+v", strValue, val)
	return int(val), err
}
