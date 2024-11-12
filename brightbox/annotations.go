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

const (
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
