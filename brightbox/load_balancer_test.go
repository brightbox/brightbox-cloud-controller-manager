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
	"strconv"
	"strings"
	"testing"

	brightbox "github.com/brightbox/gobrightbox/v2"
	"github.com/brightbox/gobrightbox/v2/enums/balancingpolicy"
	"github.com/brightbox/gobrightbox/v2/enums/cloudipstatus"
	"github.com/brightbox/gobrightbox/v2/enums/healthchecktype"
	"github.com/brightbox/gobrightbox/v2/enums/listenerprotocol"
	"github.com/brightbox/gobrightbox/v2/enums/loadbalancerstatus"
	"github.com/brightbox/gobrightbox/v2/enums/proxyprotocol"
	"github.com/brightbox/k8ssdk/v2"
	"github.com/go-test/deep"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	publicCipID    = "cip-found"
	errorCipID     = "cip-error"
	publicIP       = "180.180.180.180"
	publicIPv6     = "2a02:1348:ffff:ffff::6d6b:275c"
	publicIPv62    = "2a02:1348:ffff:ffff::6d6b:375c"
	fqdn           = "cip-180-180-180-180.gb1.brightbox.com"
	publicCipID2   = "cip-manul"
	publicIP2      = "190.190.190.190"
	fqdn2          = "cip-190-190-190-190.gb1.brightbox.com"
	reverseDNS     = "k8s-lb.example.com"
	foundLba       = "lba-found"
	errorLba       = "lba-error"
	newUID         = "9d85099c-227c-46c0-a373-e954ec8eee2e"
	clusterName    = "test-cluster-name"
	missingDomain  = "probablynotthere.co"
	resolvedDomain = "cip-vsalc.gb1s.brightbox.com"
	tooBigInt      = 1 << maxBits
	testTimeout    = tooBigInt - 1
)

// Constant variables you can take the address of!
var (
	newlbname     string               = "a9d85099c227c46c0a373e954ec8eee2.default." + clusterName
	lbuid         types.UID            = "9bde5f33-1379-4b8c-877a-777f5da4d766"
	lbname        string               = "a9bde5f3313794b8c877a777f5da4d76.default." + clusterName
	lberror       string               = "888888f3313794b8c877a777f5da4d76.default." + clusterName
	testPolicy    balancingpolicy.Enum = balancingpolicy.RoundRobin
	trueVar       bool                 = true
	falseVar      bool                 = false
	groklbname    string               = lbname
	groknewlbname string               = newlbname
	resolvCip     brightbox.CloudIP    = brightbox.CloudIP{
		ID:         "cip-vsalc",
		PublicIP:   "109.107.39.92",
		PublicIPv4: "109.107.39.92",
		PublicIPv6: "2a02:1348:ffff:ffff::6d6b:275c",
		Fqdn:       resolvedDomain,
		ReverseDNS: "cip-109-107-39-92.gb1s.brightbox.com",
	}
)

func TestLoadBalancerStatus(t *testing.T) {
	testCases := map[string]struct {
		lb     *brightbox.LoadBalancer
		status *v1.LoadBalancerStatus
	}{
		"no-cloudip": {
			lb:     &brightbox.LoadBalancer{},
			status: &v1.LoadBalancerStatus{},
		},
		"one-cloudip": {
			lb: &brightbox.LoadBalancer{
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						ID:         publicCipID,
						PublicIPv4: publicIP,
						PublicIPv6: publicIPv6,
						ReverseDNS: reverseDNS,
						Fqdn:       fqdn,
					},
				},
			},
			status: &v1.LoadBalancerStatus{
				Ingress: []v1.LoadBalancerIngress{
					/*
						v1.LoadBalancerIngress{
							IP: publicIP,
						},
						v1.LoadBalancerIngress{
							IP: publicIPv6,
						},
					*/
					v1.LoadBalancerIngress{
						Hostname: reverseDNS,
					},
					v1.LoadBalancerIngress{
						Hostname: fqdn,
					},
				},
			},
		},
		"two-cloudips": {
			lb: &brightbox.LoadBalancer{
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						ID:         publicCipID2,
						PublicIPv4: publicIP2,
						PublicIPv6: publicIPv62,
						ReverseDNS: "",
						Fqdn:       fqdn2,
						Name:       "manually allocated",
					},
					brightbox.CloudIP{
						ID:         publicCipID,
						PublicIPv4: publicIP,
						PublicIPv6: publicIPv6,
						ReverseDNS: reverseDNS,
						Fqdn:       fqdn,
					},
				},
			},
			status: &v1.LoadBalancerStatus{
				Ingress: []v1.LoadBalancerIngress{
					/*
						v1.LoadBalancerIngress{
							IP: publicIP2,
						},
						v1.LoadBalancerIngress{
							IP: publicIPv62,
						},
					*/
					v1.LoadBalancerIngress{
						Hostname: fqdn2,
					},
					/*
						v1.LoadBalancerIngress{
							IP: publicIP,
						},
						v1.LoadBalancerIngress{
							IP: publicIPv6,
						},
					*/
					v1.LoadBalancerIngress{
						Hostname: reverseDNS,
					},
					v1.LoadBalancerIngress{
						Hostname: fqdn,
					},
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := toLoadBalancerStatus(tc.lb)
			if diff := deep.Equal(result, tc.status); diff != nil {
				t.Error(diff)
			}
		})
	}

}

func TestErrorIfAcmeNotComplete(t *testing.T) {
	testCases := map[string]struct {
		acme   *brightbox.LoadBalancerAcme
		status string
	}{
		"domain_invalid": {
			acme: &brightbox.LoadBalancerAcme{
				Domains: []brightbox.LoadBalancerAcmeDomain{
					{
						Identifier:  missingDomain,
						Status:      "invalid",
						LastMessage: "failed to resolve",
					},
				},
			},
			status: "Domain \"" + missingDomain + "\" has not yet been validated for SSL use (\"invalid\":\"failed to resolve\")",
		},
		"just_one_domain_invalid": {
			acme: &brightbox.LoadBalancerAcme{
				Domains: []brightbox.LoadBalancerAcmeDomain{
					{
						Identifier: resolvedDomain,
						Status:     k8ssdk.ValidAcmeDomainStatus,
					},
					{
						Identifier:  missingDomain,
						Status:      "invalid",
						LastMessage: "failed to resolve",
					},
				},
			},
			status: "Domain \"" + missingDomain + "\" has not yet been validated for SSL use (\"invalid\":\"failed to resolve\")",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := k8ssdk.ErrorIfAcmeNotComplete(tc.acme)
			if err == nil {
				t.Errorf("Expected error %q got nil", tc.status)
			} else if err.Error() != tc.status {
				t.Errorf("Expected %q, got %q", tc.status, err.Error())
			}
		})
	}
}

func TestValidateService(t *testing.T) {
	testCases := map[string]struct {
		service *v1.Service
		status  string
	}{
		"session affinity": {
			service: &v1.Service{
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						v1.ServicePort{},
					},
					SessionAffinity: v1.ServiceAffinityClientIP,
				},
			},
			status: "unsupported load balancer affinity: ClientIP",
		},
		"empty ports": {
			service: &v1.Service{
				Spec: v1.ServiceSpec{
					Type:            v1.ServiceTypeLoadBalancer,
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "requested load balancer with no ports",
		},
		"udp ports": {
			service: &v1.Service{
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "dns",
							Protocol:   v1.ProtocolUDP,
							Port:       53,
							TargetPort: intstr.FromInt(1024),
							NodePort:   31348,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "UDP nodeports are not supported",
		},
		"invalid-policy": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerPolicy: "magic-routing",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "Invalid Load Balancer Policy \"magic-routing\": magic-routing is not a valid balancingpolicy.Enum",
		},
		"invalid-proxy-protocol": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProxyProtocol: "v1-ssl",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "Invalid Load Balancer Listener Proxy Protocol \"v1-ssl\": v1-ssl is not a valid proxyprotocol.Enum",
		},
		"Domains with TCP": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: listenerprotocol.Tcp.String(),
						serviceAnnotationLoadBalancerSslDomains:       resolvedDomain,
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "SSL Domains are not supported with the tcp protocol",
		},
		"Ports with TCP": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: listenerprotocol.Tcp.String(),
						serviceAnnotationLoadBalancerSSLPorts:         "443",
						serviceAnnotationLoadBalancerSslDomains:       resolvedDomain,
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "SSL Ports are not supported with the tcp protocol",
		},
		"valid-proxy-protocol": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProxyProtocol: proxyprotocol.V2SslCn.String(),
						serviceAnnotationLoadBalancerListenerProtocol:      listenerprotocol.Tcp.String(),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "",
		},
		"valid-http-proxy-protocol": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProxyProtocol: proxyprotocol.V2Ssl.String(),
						serviceAnnotationLoadBalancerListenerProtocol:      listenerprotocol.Http.String(),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "",
		},
		"valid-listener-protocol": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: listenerprotocol.Http.String(),
						serviceAnnotationLoadBalancerSslDomains:       resolvedDomain,
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "",
		},
		"valid-websocket-listener-protocol": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: listenerprotocol.Https.String(),
						serviceAnnotationLoadBalancerSslDomains:       resolvedDomain,
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "",
		},
		"domains without port 443": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: listenerprotocol.Http.String(),
						serviceAnnotationLoadBalancerSslDomains:       resolvedDomain,
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "SSL support requires a Port definition for 443",
		},
		"Ports without port 443": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: listenerprotocol.Http.String(),
						serviceAnnotationLoadBalancerSSLPorts:         "http",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "SSL support requires a Port definition for 443",
		},
		"invalid-listener-protocol": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: "gopher",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "Invalid Load Balancer Listener Protocol \"gopher\": gopher is not a valid listenerprotocol.Enum",
		},
		"invalid-healthcheck-request": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerHCRequest: "fred",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "\"" + serviceAnnotationLoadBalancerHCRequest + "\" needs to be a valid Url request path",
		},
		"invalid-healthcheck-protocol": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerHCProtocol: listenerprotocol.Https.String(),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "Invalid Load Balancer Healthcheck Protocol \"" + listenerprotocol.Https.String() + "\": " +
				listenerprotocol.Https.String() + " is not a valid healthchecktype.Enum",
		},
		"https without domains": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: listenerprotocol.Http.String(),
						serviceAnnotationLoadBalancerSSLPorts:         "443",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "SSL needs a list of domains to certify. Add the \"" + serviceAnnotationLoadBalancerSslDomains + "\" annotation",
		},
		"invalid-uint-negative": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerHCInterval: "-1",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "\"" + serviceAnnotationLoadBalancerHCInterval + "\" needs to be a positive number (strconv.ParseUint: parsing \"-1\": invalid syntax)",
		},
		"invalid-uint-alpha": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerHCInterval: "0x56",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "\"" + serviceAnnotationLoadBalancerHCInterval + "\" needs to be a positive number (strconv.ParseUint: parsing \"0x56\": invalid syntax)",
		},
		"invalid-uint-too-big": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerHCTimeout: strconv.Itoa(tooBigInt),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
					},
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			status: "\"" + serviceAnnotationLoadBalancerHCTimeout + "\" needs to be a positive number (strconv.ParseUint: parsing \"" +
				strconv.Itoa(tooBigInt) + "\": value out of range)",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateServiceSpec(tc.service)
			got := ""
			if err != nil {
				got = err.Error()
			}
			if tc.status != got {
				t.Errorf("Expected %q, got %q", tc.status, got)
			}
		})
	}
}

func TestValidateDomains(t *testing.T) {
	testCases := map[string]struct {
		annotations map[string]string
		cloudIp     *brightbox.CloudIP
		status      string
	}{
		"missing domain": {
			annotations: map[string]string{
				serviceAnnotationLoadBalancerSslDomains: missingDomain,
			},
			cloudIp: &resolvCip,
			status:  "Failed to resolve \"" + missingDomain + "\" to load balancer address (" + resolvCip.PublicIPv4 + "," + resolvCip.PublicIPv6 + "):",
		},
		"missing domain in list": {
			annotations: map[string]string{
				serviceAnnotationLoadBalancerSslDomains: resolvedDomain + "," + missingDomain,
			},
			cloudIp: &resolvCip,
			status:  "Failed to resolve \"" + missingDomain + "\" to load balancer address (" + resolvCip.PublicIPv4 + "," + resolvCip.PublicIPv6 + "):",
		},
		"other addresses": {
			annotations: map[string]string{
				serviceAnnotationLoadBalancerSslDomains: resolvedDomain + ",archive.ubuntu.com",
			},
			cloudIp: &resolvCip,
			status:  "Failed to resolve \"archive.ubuntu.com\" to load balancer address (" + resolvCip.PublicIPv4 + "," + resolvCip.PublicIPv6 + ")",
		},
		"dodgy cloudip": {
			annotations: map[string]string{
				serviceAnnotationLoadBalancerSslDomains: resolvedDomain,
			},
			cloudIp: &brightbox.CloudIP{},
			status:  "Cloud IP \"\" failed to parse IP addresses",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateContextualAnnotations(tc.annotations, tc.cloudIp)
			if err == nil {
				t.Errorf("Expected error %q got nil", tc.status)
			} else if !strings.HasPrefix(err.Error(), tc.status) {
				t.Errorf("Expected %q, got %q", tc.status, err.Error())
			}
		})
	}
}

func TestGetLoadBalancer(t *testing.T) {
	testCases := map[string]struct {
		service  *v1.Service
		lbstatus *v1.LoadBalancerStatus
		exists   bool
		err      bool
	}{
		"missing": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
				},
				Spec: v1.ServiceSpec{
					Type:            v1.ServiceTypeLoadBalancer,
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			lbstatus: &v1.LoadBalancerStatus{},
			exists:   false,
		},
		"found": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
				},
				Spec: v1.ServiceSpec{
					Type:            v1.ServiceTypeLoadBalancer,
					SessionAffinity: v1.ServiceAffinityNone,
				},
			},
			lbstatus: &v1.LoadBalancerStatus{
				Ingress: []v1.LoadBalancerIngress{
					/*
						v1.LoadBalancerIngress{
							IP: publicIP,
						},
						v1.LoadBalancerIngress{
							IP: publicIPv6,
						},
					*/
					v1.LoadBalancerIngress{
						Hostname: reverseDNS,
					},
					v1.LoadBalancerIngress{
						Hostname: fqdn,
					},
				},
			},
			exists: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()

			lb, exists, err := client.GetLoadBalancer(
				context.TODO(),
				clusterName,
				tc.service,
			)
			if err != nil {
				t.Errorf("Error when not expected: %q", err.Error())
			} else if tc.exists != exists {
				t.Errorf("Exists status wrong, got %v, expected %v for %v", exists, tc.exists,
					client.GetLoadBalancerName(context.TODO(), clusterName, tc.service))
			} else if diff := deep.Equal(lb, tc.lbstatus); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestBuildLoadBalancerOptions(t *testing.T) {
	testCases := map[string]struct {
		service *v1.Service
		nodes   []*v1.Node
		lbopts  *brightbox.LoadBalancerOptions
	}{
		"standard": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerIdleTimeout: strconv.Itoa(testTimeout),
						serviceAnnotationLoadBalancerPolicy:              testPolicy.String(),
						serviceAnnotationLoadBalancerSslDomains:          resolvedDomain + "," + fqdn,
						serviceAnnotationLoadBalancerListenerProtocol:    listenerprotocol.Http.String(),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-gdqms",
					},
				},
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			lbopts: &brightbox.LoadBalancerOptions{
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-gdqms",
					},
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
						Timeout:  testTimeout,
					},
					{
						Protocol: listenerprotocol.Http,
						Timeout:  testTimeout,
						In:       80,
						Out:      31348,
					},
				},
				Domains: &[]string{resolvedDomain, fqdn},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    31347,
					Request: "/healthz",
				},
				Policy:        testPolicy,
				HTTPSRedirect: &trueVar,
			},
		},
		"standard_proxy_protocol": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerIdleTimeout:   strconv.Itoa(testTimeout),
						serviceAnnotationLoadBalancerPolicy:                testPolicy.String(),
						serviceAnnotationLoadBalancerSslDomains:            resolvedDomain + "," + fqdn,
						serviceAnnotationLoadBalancerListenerProtocol:      listenerprotocol.Http.String(),
						serviceAnnotationLoadBalancerListenerProxyProtocol: proxyprotocol.V2SslCn.String(),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-gdqms",
					},
				},
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			lbopts: &brightbox.LoadBalancerOptions{
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-gdqms",
					},
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol:      listenerprotocol.Https,
						In:            443,
						Out:           31347,
						Timeout:       testTimeout,
						ProxyProtocol: proxyprotocol.V2SslCn,
					},
					{
						Protocol:      listenerprotocol.Http,
						Timeout:       testTimeout,
						In:            80,
						Out:           31348,
						ProxyProtocol: proxyprotocol.V2SslCn,
					},
				},
				Domains: &[]string{resolvedDomain, fqdn},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    31347,
					Request: "/healthz",
				},
				Policy:        testPolicy,
				HTTPSRedirect: &trueVar,
			},
		},
		"extraSSLports": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerIdleTimeout: strconv.Itoa(testTimeout),
						serviceAnnotationLoadBalancerPolicy:              testPolicy.String(),
						serviceAnnotationLoadBalancerSslDomains:          resolvedDomain + "," + fqdn,
						serviceAnnotationLoadBalancerSSLPorts:            "fancy,3030",
						serviceAnnotationLoadBalancerListenerProtocol:    listenerprotocol.Http.String(),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
						{
							Name:       "fancy",
							Protocol:   v1.ProtocolTCP,
							Port:       5050,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
						{
							Name:       "charm",
							Protocol:   v1.ProtocolTCP,
							Port:       3030,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-gdqms",
					},
				},
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			lbopts: &brightbox.LoadBalancerOptions{
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-gdqms",
					},
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
						Timeout:  testTimeout,
					},
					{
						Protocol: listenerprotocol.Http,
						Timeout:  testTimeout,
						In:       80,
						Out:      31348,
					},
					{
						Protocol: listenerprotocol.Https,
						Timeout:  testTimeout,
						In:       5050,
						Out:      31348,
					},
					{
						Protocol: listenerprotocol.Https,
						Timeout:  testTimeout,
						In:       3030,
						Out:      31348,
					},
				},
				Domains: &[]string{resolvedDomain, fqdn},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    31347,
					Request: "/healthz",
				},
				Policy:        testPolicy,
				HTTPSRedirect: &trueVar,
			},
		},
		"OverrideToTcpListener": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol:    listenerprotocol.Tcp.String(),
						serviceAnnotationLoadBalancerListenerIdleTimeout: strconv.Itoa(testTimeout),
						serviceAnnotationLoadBalancerPolicy:              testPolicy.String(),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-gdqms",
					},
				},
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			lbopts: &brightbox.LoadBalancerOptions{
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-gdqms",
					},
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Tcp,
						Timeout:  testTimeout,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Tcp,
					Port:    31348,
					Request: "/",
				},
				Policy:        testPolicy,
				HTTPSRedirect: &falseVar,
			},
		},
		"overrideToHttpHealthcheck": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol:     listenerprotocol.Tcp.String(),
						serviceAnnotationLoadBalancerHCProtocol:           healthchecktype.Http.String(),
						serviceAnnotationLoadBalancerHCInterval:           "4000",
						serviceAnnotationLoadBalancerHCTimeout:            "6000",
						serviceAnnotationLoadBalancerHCHealthyThreshold:   "4",
						serviceAnnotationLoadBalancerHCUnhealthyThreshold: "5",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-gdqms",
					},
				},
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			lbopts: &brightbox.LoadBalancerOptions{
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-gdqms",
					},
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Tcp,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Tcp,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:          healthchecktype.Http,
					Port:          31347,
					Request:       "/healthz",
					Timeout:       6000,
					Interval:      4000,
					ThresholdUp:   4,
					ThresholdDown: 5,
				},
				HTTPSRedirect: &falsevar,
			},
		},
		"httphealthcheck": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeLocal,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "srv-gdprt",
					},
					Spec: v1.NodeSpec{},
				},
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "srv-230b7",
					},
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			lbopts: &brightbox.LoadBalancerOptions{
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
		},
		"overrideToTcpHealthcheck": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerHCProtocol: healthchecktype.Tcp.String(),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeLocal,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "srv-gdprt",
					},
					Spec: v1.NodeSpec{},
				},
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "srv-230b7",
					},
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			lbopts: &brightbox.LoadBalancerOptions{
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Tcp,
					Port:    8080,
					Request: "/",
				},
				HTTPSRedirect: &falsevar,
			},
		},
		"empty": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
				},
				Spec: v1.ServiceSpec{
					Type:                  v1.ServiceTypeLoadBalancer,
					Ports:                 []v1.ServicePort{},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{},
			lbopts: &brightbox.LoadBalancerOptions{
				Name: &groklbname,
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    80,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()
			desc := client.GetLoadBalancerName(context.TODO(), clusterName, tc.service)

			lbopts := buildLoadBalancerOptions(desc, tc.service, tc.nodes)
			if diff := deep.Equal(lbopts, tc.lbopts); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestEnsureAndUpdateLoadBalancer(t *testing.T) {
	testCases := map[string]struct {
		service *v1.Service
		nodes   []*v1.Node
		status  *v1.LoadBalancerStatus
	}{
		"found": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-gdqms",
					},
				},
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			status: &v1.LoadBalancerStatus{
				Ingress: []v1.LoadBalancerIngress{
					/*
						v1.LoadBalancerIngress{
							IP: publicIP,
						},
						v1.LoadBalancerIngress{
							IP: publicIPv6,
						},
					*/
					v1.LoadBalancerIngress{
						Hostname: reverseDNS,
					},
					v1.LoadBalancerIngress{
						Hostname: fqdn,
					},
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()

			lbstatus, err := client.EnsureLoadBalancer(context.TODO(), clusterName, tc.service, tc.nodes)
			if err != nil {
				t.Errorf("Error when not expected: %q", err.Error())
			} else if diff := deep.Equal(lbstatus, tc.status); diff != nil {
				t.Error(diff)
			}
			err = client.UpdateLoadBalancer(context.TODO(), clusterName, tc.service, tc.nodes)
			if err != nil {
				t.Errorf("Error when not expected: %q", err.Error())
			}
		})
	}
}

func TestBuildEnsureLoadBalancer(t *testing.T) {
	testCases := map[string]struct {
		service *v1.Service
		nodes   []*v1.Node
		lbopts  *brightbox.LoadBalancer
	}{
		"found_no_change": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: listenerprotocol.Tcp.String(),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-gdqms",
					},
				},
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			lbopts: &brightbox.LoadBalancer{
				ID:     foundLba,
				Name:   lbname,
				Status: loadbalancerstatus.Active,
				Nodes: []brightbox.Server{
					{
						ID: "srv-gdqms",
					},
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Tcp,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Tcp,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Tcp,
					Port:    31347,
					Request: "/",
				},
			},
		},
		"found_updated": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: listenerprotocol.Tcp.String(),
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-gdqms",
					},
				},
				&v1.Node{
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			lbopts: &brightbox.LoadBalancer{
				ID:     foundLba,
				Name:   lbname,
				Status: loadbalancerstatus.Active,
				Nodes: []brightbox.Server{
					{
						ID: "srv-gdqms",
					},
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Tcp,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Tcp,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Tcp,
					Port:    31347,
					Request: "/",
				},
			},
		},
		"notfound": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerHCRequest: "/different/path",
					},
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeLocal,
					HealthCheckNodePort:   8080,
				},
			},
			nodes: []*v1.Node{
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "srv-gdprt",
					},
					Spec: v1.NodeSpec{},
				},
				&v1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "srv-230b7",
					},
					Spec: v1.NodeSpec{
						ProviderID: "brightbox://srv-230b7",
					},
				},
			},
			lbopts: &brightbox.LoadBalancer{
				Name:   newlbname,
				Status: loadbalancerstatus.Active,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/different/path",
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()

			ctx := context.Background()
			desc := client.GetLoadBalancerName(ctx, clusterName, tc.service)
			lbopts, err := client.ensureLoadBalancerFromService(ctx, desc, tc.service, tc.nodes)
			if err != nil {
				t.Errorf("Error when not expected")
			} else if diff := deep.Equal(lbopts, tc.lbopts); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestUpdateLoadBalancerCheck(t *testing.T) {
	testCases := map[string]struct {
		lb       *brightbox.LoadBalancer
		lbopts   brightbox.LoadBalancerOptions
		expected bool
	}{
		"No change domains": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Acme: &brightbox.LoadBalancerAcme{
					Domains: []brightbox.LoadBalancerAcmeDomain{
						{
							Identifier: resolvedDomain,
							Status:     k8ssdk.ValidAcmeDomainStatus,
						},
						{
							Identifier: fqdn,
							Status:     k8ssdk.ValidAcmeDomainStatus,
						},
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Domains: &[]string{resolvedDomain, fqdn},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: false,
		},
		"swap domains": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Acme: &brightbox.LoadBalancerAcme{
					Domains: []brightbox.LoadBalancerAcmeDomain{
						{
							Identifier: resolvedDomain,
							Status:     k8ssdk.ValidAcmeDomainStatus,
						},
						{
							Identifier: fqdn,
							Status:     k8ssdk.ValidAcmeDomainStatus,
						},
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Domains: &[]string{fqdn, resolvedDomain},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: false,
		},
		"add_domain": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Acme: &brightbox.LoadBalancerAcme{
					Domains: []brightbox.LoadBalancerAcmeDomain{
						{
							Identifier: resolvedDomain,
							Status:     k8ssdk.ValidAcmeDomainStatus,
						},
						{
							Identifier: fqdn,
							Status:     k8ssdk.ValidAcmeDomainStatus,
						},
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Domains: &[]string{resolvedDomain, fqdn, reverseDNS},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"remove_domain": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Acme: &brightbox.LoadBalancerAcme{
					Domains: []brightbox.LoadBalancerAcmeDomain{
						{
							Identifier: resolvedDomain,
							Status:     k8ssdk.ValidAcmeDomainStatus,
						},
						{
							Identifier: fqdn,
							Status:     k8ssdk.ValidAcmeDomainStatus,
						},
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Domains: &[]string{fqdn},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"change domain": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Acme: &brightbox.LoadBalancerAcme{
					Domains: []brightbox.LoadBalancerAcmeDomain{
						{
							Identifier: resolvedDomain,
							Status:     k8ssdk.ValidAcmeDomainStatus,
						},
						{
							Identifier: fqdn,
							Status:     k8ssdk.ValidAcmeDomainStatus,
						},
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Domains: &[]string{reverseDNS, fqdn},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Https,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"No change": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: false,
		},
		"add listener": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
					{
						Protocol: listenerprotocol.Tcp,
						In:       25,
						Out:      32456,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"remove listener": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"change_listener": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31350,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"add_proxy_protocol": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol:      listenerprotocol.Http,
						In:            443,
						Out:           31347,
						ProxyProtocol: proxyprotocol.V2,
					},
					{
						Protocol:      listenerprotocol.Http,
						In:            80,
						Out:           31348,
						ProxyProtocol: proxyprotocol.V2,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"remove_proxy_protocol": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol:      listenerprotocol.Http,
						In:            443,
						Out:           31347,
						ProxyProtocol: proxyprotocol.V2,
					},
					{
						Protocol:      listenerprotocol.Http,
						In:            80,
						Out:           31348,
						ProxyProtocol: proxyprotocol.V2,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"change_proxy_protocol": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol:      listenerprotocol.Http,
						In:            443,
						Out:           31347,
						ProxyProtocol: proxyprotocol.V2,
					},
					{
						Protocol:      listenerprotocol.Http,
						In:            80,
						Out:           31348,
						ProxyProtocol: proxyprotocol.V2,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol:      listenerprotocol.Http,
						In:            443,
						Out:           31347,
						ProxyProtocol: proxyprotocol.V2Ssl,
					},
					{
						Protocol:      listenerprotocol.Http,
						In:            80,
						Out:           31348,
						ProxyProtocol: proxyprotocol.V2Ssl,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"change node": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-newon",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"remove node": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:    foundLba,
				Name:  &groklbname,
				Nodes: []brightbox.LoadBalancerNode{},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"name_change": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groknewlbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"add node": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
					{
						Node: "srv-newon",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
		"Healthcheck change": {
			lb: &brightbox.LoadBalancer{
				ID:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						ID: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				ID:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: listenerprotocol.Http,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: listenerprotocol.Http,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    healthchecktype.Http,
					Port:    8080,
					Request: "/check",
				},
				HTTPSRedirect: &falsevar,
			},
			expected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if change := k8ssdk.IsUpdateLoadBalancerRequired(tc.lb, tc.lbopts); change != tc.expected {
				t.Errorf("Expected %v got %v", tc.expected, change)
			}
		})
	}
}

func TestEnsureMappedCloudIP(t *testing.T) {
	testCases := map[string]struct {
		lb  *brightbox.LoadBalancer
		cip *brightbox.CloudIP
		err bool
	}{
		"mapped": {
			lb: &brightbox.LoadBalancer{
				ID: "lba-testy",
				CloudIPs: []brightbox.CloudIP{
					{
						ID: "cip-testy",
					},
				},
			},
			cip: &brightbox.CloudIP{
				ID: "cip-testy",
				LoadBalancer: &brightbox.LoadBalancer{
					ID: "lba-testy",
				},
				Status: cloudipstatus.Mapped,
			},
			err: false,
		},
		"mapped_elsewhere": {
			lb: &brightbox.LoadBalancer{
				ID:       "lba-testy",
				CloudIPs: []brightbox.CloudIP{},
			},
			cip: &brightbox.CloudIP{
				ID:     "cip-testy",
				Status: cloudipstatus.Mapped,
			},
			err: false,
		},
		"unmapped": {
			lb: &brightbox.LoadBalancer{
				ID:       "lba-testy",
				CloudIPs: []brightbox.CloudIP{},
			},
			cip: &brightbox.CloudIP{
				ID:     "cip-testy",
				Status: cloudipstatus.Unmapped,
			},
			err: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()

			err := client.EnsureMappedCloudIP(context.Background(), tc.lb, tc.cip)
			if tc.err && err == nil {
				t.Errorf("Expected error and none returned")
			} else if !tc.err && err != nil {
				t.Errorf("Error return when not expected: %q", err.Error())
			}
		})
	}
}

func TestEnsureAllocatedCloudIP(t *testing.T) {
	testCases := map[string]struct {
		service *v1.Service
		cip     *brightbox.CloudIP
	}{
		"LBIP_invalid": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        "fred",
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			cip: nil,
		},
		"LBIP_found_no_name": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			cip: &brightbox.CloudIP{
				ID:         "cip-12345",
				PublicIPv4: publicIP,
				PublicIPv6: publicIPv6,
			},
		},
		"LBIP_notfound_no_name": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					LoadBalancerIP:        publicIP2,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			cip: nil,
		},
		"name_found_noLBIP": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			cip: &brightbox.CloudIP{
				ID:         publicCipID,
				Name:       lbname,
				PublicIPv4: "240.240.240.240",
			},
		},
		"new_allocation_noLBIP": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
				},
			},
			cip: &brightbox.CloudIP{
				ID:         "cip-67890",
				Name:       newlbname,
				PublicIPv4: publicIP2,
				PublicIPv6: publicIPv62,
			},
		},
		"LBIP_specified_with_name_found": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
					LoadBalancerIP:        publicIP,
				},
			},
			cip: &brightbox.CloudIP{
				ID:         "cip-12345",
				PublicIPv4: publicIP,
				PublicIPv6: publicIPv6,
			},
		},
		"LBIP_not_found_with_name_found": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
					HealthCheckNodePort:   8080,
					LoadBalancerIP:        publicIP2,
				},
			},
			cip: nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()
			ctx := context.Background()

			desc := client.GetLoadBalancerName(ctx, clusterName, tc.service)
			cip, err := client.ensureAllocatedCloudIP(ctx, desc, tc.service)
			if err != nil && tc.cip != nil {
				t.Errorf("Error when not expected %q", err.Error())
			} else if diff := deep.Equal(cip, tc.cip); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestDeposeCloudIPFunctions(t *testing.T) {
	testCases := map[string]struct {
		lb   *brightbox.LoadBalancer
		cip  *brightbox.CloudIP
		name string
	}{
		"no_change": {
			lb: &brightbox.LoadBalancer{
				ID: "lba-oldip",
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						ID:         publicCipID,
						PublicIPv4: publicIP,
						PublicIPv6: publicIPv6,
						ReverseDNS: reverseDNS,
						Fqdn:       fqdn,
						Name:       "test",
					},
				},
			},
			cip: &brightbox.CloudIP{
				ID:         publicCipID,
				PublicIPv4: publicIP,
				PublicIPv6: publicIPv6,
				ReverseDNS: reverseDNS,
				Fqdn:       fqdn,
				Name:       "test",
			},
			name: "test",
		},
		"no_change_manual": {
			lb: &brightbox.LoadBalancer{
				ID: "lba-oldip",
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						ID:         publicCipID2,
						PublicIPv4: publicIP2,
						PublicIPv6: publicIPv62,
						Fqdn:       fqdn2,
						Name:       "manually allocated",
					},
				},
			},
			cip: &brightbox.CloudIP{
				ID:         publicCipID2,
				PublicIPv4: publicIP2,
				PublicIPv6: publicIPv62,
				Fqdn:       fqdn2,
				Name:       "manually allocated",
			},
			name: "test",
		},
		"changed_delete": {
			lb: &brightbox.LoadBalancer{
				ID: "lba-oldip",
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						ID:         publicCipID2,
						PublicIPv4: publicIP2,
						PublicIPv6: publicIPv62,
						Fqdn:       fqdn2,
						Name:       "manually allocated",
					},
					brightbox.CloudIP{
						ID:         publicCipID,
						PublicIPv4: publicIP,
						PublicIPv6: publicIPv6,
						ReverseDNS: reverseDNS,
						Fqdn:       fqdn,
						Name:       "test",
					},
				},
			},
			cip: &brightbox.CloudIP{
				ID:         publicCipID2,
				PublicIPv4: publicIP2,
				PublicIPv6: publicIPv62,
				Fqdn:       fqdn2,
				Name:       "manually allocated",
			},
			name: "test",
		},
		"changed_unmap": {
			lb: &brightbox.LoadBalancer{
				ID: "lba-oldip",
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						ID:         publicCipID2,
						PublicIPv4: publicIP2,
						PublicIPv6: publicIPv62,
						Fqdn:       fqdn2,
						Name:       "manually allocated",
					},
					brightbox.CloudIP{
						ID:         publicCipID,
						PublicIPv4: publicIP,
						PublicIPv6: publicIPv6,
						ReverseDNS: reverseDNS,
						Fqdn:       fqdn,
						Name:       "test",
					},
				},
			},
			cip: &brightbox.CloudIP{
				ID:         publicCipID,
				PublicIPv4: publicIP,
				PublicIPv6: publicIPv6,
				ReverseDNS: reverseDNS,
				Fqdn:       fqdn,
				Name:       "test",
			},
			name: "test",
		},
		"already unmapped": {
			lb: &brightbox.LoadBalancer{
				ID:       "lba-oldip",
				CloudIPs: []brightbox.CloudIP{},
			},
			cip: &brightbox.CloudIP{
				ID:         publicCipID,
				PublicIPv4: publicIP,
				PublicIPv6: publicIPv6,
				ReverseDNS: reverseDNS,
				Fqdn:       fqdn,
				Name:       "test",
			},
			name: "test",
		},
		"already unmapped_manual": {
			lb: &brightbox.LoadBalancer{
				ID:       "lba-oldip",
				CloudIPs: []brightbox.CloudIP{},
			},
			cip: &brightbox.CloudIP{
				ID:         publicCipID2,
				PublicIPv4: publicIP2,
				PublicIPv6: publicIPv62,
				Fqdn:       fqdn2,
				Name:       "manually allocated",
			},
			name: "test",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()

			err := client.EnsureOldCloudIPsDeposed(context.Background(), tc.lb.CloudIPs, tc.cip.ID, tc.name)
			if err != nil {
				t.Errorf("Error when not expected: %q", err.Error())
			}
		})
	}
}

func TestDeletionByNameFunctions(t *testing.T) {
	testCases := []string{
		lbname,
		"not-found",
		lberror,
	}
	for _, name := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()
			ctx := context.Background()

			client.ensureServerGroupDeleted(ctx, name)
			client.ensureLoadBalancerDeletedByName(ctx, name)
			client.ensureFirewallClosed(ctx, name)
			client.ensureCloudIPsDeleted(ctx, "", name)
		})
	}
}

func TestPortListString(t *testing.T) {
	testCases := map[string]struct {
		service    *v1.Service
		portstring string
	}{
		"with healthcheck port": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeLocal,
					HealthCheckNodePort:   8080,
				},
			},
			portstring: "31347,31348,8080",
		},
		"withoutHealthcheckPort": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   v1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   v1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31348,
						},
					},
					SessionAffinity:       v1.ServiceAffinityNone,
					ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
				},
			},
			portstring: "31347,31348",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			portstring := createPortListString(tc.service)
			if diff := deep.Equal(portstring, tc.portstring); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestEnsureLoadBalancerDeleted(t *testing.T) {
	testCases := map[string]struct {
		service *v1.Service
		err     error
	}{
		"missing": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
				},
				Spec: v1.ServiceSpec{
					Type:            v1.ServiceTypeLoadBalancer,
					Ports:           nil,
					SessionAffinity: "None",
					LoadBalancerIP:  "",
				},
			},
			err: nil,
		},
		"found": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
				},
				Spec: v1.ServiceSpec{
					Type:            v1.ServiceTypeLoadBalancer,
					Ports:           nil,
					SessionAffinity: "None",
					LoadBalancerIP:  "",
				},
			},
			err: fmt.Errorf("CloudIPs still mapped to load balancer %q", foundLba),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()

			err := client.EnsureLoadBalancerDeleted(
				context.TODO(),
				clusterName,
				tc.service,
			)
			if diff := deep.Equal(err, tc.err); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func (f *fakeInstanceCloud) MapCloudIP(_ context.Context, identifier string, destination brightbox.CloudIPAttachment) (*brightbox.CloudIP, error) {
	return nil, nil
}

func (f *fakeInstanceCloud) CloudIPs(context.Context) ([]brightbox.CloudIP, error) {
	return []brightbox.CloudIP{
		{
			ID:         "cip-12345",
			PublicIPv4: publicIP,
			PublicIPv6: publicIPv6,
		},
		{
			ID:         publicCipID,
			Name:       lbname,
			PublicIPv4: "240.240.240.240",
		},
		{
			ID:         errorCipID,
			Name:       lberror,
			PublicIPv4: "255.255.255.255",
		},
	}, nil
}

func (f *fakeInstanceCloud) CreateCloudIP(_ context.Context, newCloudIP brightbox.CloudIPOptions) (*brightbox.CloudIP, error) {
	cip := &brightbox.CloudIP{
		ID:         "cip-67890",
		PublicIPv4: publicIP2,
		PublicIPv6: publicIPv62,
	}
	if newCloudIP.Name != nil {
		cip.Name = *newCloudIP.Name
	}
	return cip, nil
}

func (f *fakeInstanceCloud) CreateLoadBalancer(ctx context.Context, newLB brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	return f.UpdateLoadBalancer(ctx, newLB)
}

func (f *fakeInstanceCloud) UpdateLoadBalancer(_ context.Context, newLB brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	server_list := make([]brightbox.Server, len(newLB.Nodes))
	for i, v := range newLB.Nodes {
		server_list[i].ID = v.Node
	}
	return &brightbox.LoadBalancer{
		ID:          newLB.ID,
		Name:        *newLB.Name,
		Status:      loadbalancerstatus.Active,
		Nodes:       server_list,
		Listeners:   newLB.Listeners,
		Healthcheck: *newLB.Healthcheck,
	}, nil
}

func (f *fakeInstanceCloud) LoadBalancer(ctx context.Context, id string) (*brightbox.LoadBalancer, error) {
	list, _ := f.LoadBalancers(ctx)
	for _, balancer := range list {
		if balancer.ID == id {
			return &balancer, nil
		}
	}
	return nil, fmt.Errorf("unexpected identifier %q sent to LoadBalancer", id)
}

func (f *fakeInstanceCloud) LoadBalancers(context.Context) ([]brightbox.LoadBalancer, error) {
	return []brightbox.LoadBalancer{
		{
			ID:       "lba-test1",
			Name:     lbname,
			Status:   loadbalancerstatus.Deleted,
			CloudIPs: nil,
		},
		{
			ID:     foundLba,
			Name:   lbname,
			Status: loadbalancerstatus.Active,
			CloudIPs: []brightbox.CloudIP{
				brightbox.CloudIP{
					ID:         "cip-12345",
					PublicIPv4: publicIP,
					PublicIPv6: publicIPv6,
					ReverseDNS: reverseDNS,
					Fqdn:       fqdn,
				},
			},
		},
		{
			ID:     "lba-test3",
			Name:   "abob",
			Status: loadbalancerstatus.Active,
		},
		{
			ID:     errorLba,
			Name:   lberror,
			Status: loadbalancerstatus.Active,
			CloudIPs: []brightbox.CloudIP{
				brightbox.CloudIP{
					ID:         publicCipID,
					PublicIPv4: publicIP,
					PublicIPv6: publicIPv6,
					ReverseDNS: reverseDNS,
					Fqdn:       fqdn,
				},
			},
		},
	}, nil
}

func (f *fakeInstanceCloud) AddServersToServerGroup(ctx context.Context, identifier string, serverIDs brightbox.ServerGroupMemberList) (*brightbox.ServerGroup, error) {
	switch identifier {
	case "grp-found":
		groups, _ := f.ServerGroups(ctx)
		return &groups[0], nil
	default:
		result := &brightbox.ServerGroup{
			ID:      identifier,
			Name:    "Fake Name After AddServers",
			Servers: mapServerIDsToServers(serverIDs),
		}
		return result, nil
	}
}

func (f *fakeInstanceCloud) RemoveServersFromServerGroup(ctx context.Context, identifier string, serverIDs brightbox.ServerGroupMemberList) (*brightbox.ServerGroup, error) {
	switch identifier {
	case "grp-found":
		groups, _ := f.ServerGroups(ctx)
		return &groups[0], nil
	default:
		result := &brightbox.ServerGroup{
			ID:      identifier,
			Name:    "Fake Name After RemoveServers",
			Servers: mapServerIDsToServers(serverIDs),
		}
		return result, nil
	}
}

func (f *fakeInstanceCloud) ServerGroups(context.Context) ([]brightbox.ServerGroup, error) {
	result := []brightbox.ServerGroup{
		brightbox.ServerGroup{
			ID:   "grp-found",
			Name: lbname,
			Servers: []brightbox.Server{
				{
					ID: "srv-gdqms",
				},
				{
					ID: "srv-230b7",
				},
			},
			FirewallPolicy: &brightbox.FirewallPolicy{
				ID:   "fwp-found",
				Name: lbname,
				Rules: []brightbox.FirewallRule{
					{
						ID:          "fwr-found",
						Description: lbname,
					},
				},
			},
		},
		brightbox.ServerGroup{
			ID:   "grp-error",
			Name: lberror,
			Servers: []brightbox.Server{
				{
					ID: "srv-gdqms",
				},
				{
					ID: "srv-230b7",
				},
			},
			FirewallPolicy: &brightbox.FirewallPolicy{
				ID:   "fwp-error",
				Name: lberror,
				Rules: []brightbox.FirewallRule{
					{
						ID:          "fwr-found",
						Description: lberror,
					},
				},
			},
		},
	}
	return result, nil
}

func (f *fakeInstanceCloud) CreateServerGroup(_ context.Context, newServerGroup brightbox.ServerGroupOptions) (*brightbox.ServerGroup, error) {
	result := &brightbox.ServerGroup{
		ID: "grp-testy",
	}
	if newServerGroup.Name != nil {
		result.Name = *newServerGroup.Name
	}
	if newServerGroup.Description != nil {
		result.Description = *newServerGroup.Description
	}
	return result, nil
}

func (f *fakeInstanceCloud) CreateFirewallPolicy(_ context.Context, policyOptions brightbox.FirewallPolicyOptions) (*brightbox.FirewallPolicy, error) {
	result := &brightbox.FirewallPolicy{
		ID:   "fwp-testy",
		Name: *policyOptions.Name,
		ServerGroup: &brightbox.ServerGroup{
			ID:   policyOptions.ServerGroup,
			Name: *policyOptions.Name,
		},
	}
	return result, nil
}

func (f *fakeInstanceCloud) CreateFirewallRule(_ context.Context, ruleOptions brightbox.FirewallRuleOptions) (*brightbox.FirewallRule, error) {
	result := &brightbox.FirewallRule{
		ID:          "fwr-testy",
		Description: "After Create Firewll Rule",
		FirewallPolicy: &brightbox.FirewallPolicy{
			ID:   ruleOptions.FirewallPolicy,
			Name: "After Create Firewall Rule",
		},
	}
	return result, nil
}

func (f *fakeInstanceCloud) UpdateFirewallRule(_ context.Context, ruleOptions brightbox.FirewallRuleOptions) (*brightbox.FirewallRule, error) {
	result := &brightbox.FirewallRule{
		ID:          ruleOptions.ID,
		Description: *ruleOptions.Description,
		FirewallPolicy: &brightbox.FirewallPolicy{
			ID:   ruleOptions.FirewallPolicy,
			Name: *ruleOptions.Description,
		},
	}
	return result, nil
}

func mapServerIDsToServers(serverIDs brightbox.ServerGroupMemberList) []brightbox.Server {
	result := make([]brightbox.Server, len(serverIDs.Servers))
	for i, server := range serverIDs.Servers {
		result[i].ID = server.Server
	}
	return result
}

func (f *fakeInstanceCloud) FirewallPolicies(context.Context) ([]brightbox.FirewallPolicy, error) {
	result := []brightbox.FirewallPolicy{
		brightbox.FirewallPolicy{
			ID:   "fwp-found",
			Name: lbname,
			Rules: []brightbox.FirewallRule{
				{
					ID:          "fwr-found",
					Description: lbname,
				},
			},
		},
		brightbox.FirewallPolicy{
			ID:   "fwp-error",
			Name: lberror,
			Rules: []brightbox.FirewallRule{
				{
					ID:          "fwr-error",
					Description: lberror,
				},
			},
			ServerGroup: &brightbox.ServerGroup{
				ID:   "grp-error",
				Name: lberror,
				Servers: []brightbox.Server{
					{
						ID: "srv-gdqms",
					},
					{
						ID: "srv-230b7",
					},
				},
			},
		},
	}
	return result, nil
}

func (f *fakeInstanceCloud) DestroyServerGroup(_ context.Context, identifier string) (*brightbox.ServerGroup, error) {
	switch identifier {
	case "grp-found":
		return nil, nil
	case "grp-error":
		return nil, fmt.Errorf("Raising error in DestroyServerGroup")
	default:
		return nil, fmt.Errorf("unexpected identifier %q sent to DestroyServerGroup", identifier)
	}
}

func (f *fakeInstanceCloud) DestroyFirewallPolicy(_ context.Context, identifier string) (*brightbox.FirewallPolicy, error) {
	switch identifier {
	case "fwp-found":
		return nil, nil
	case "fwp-error":
		return nil, fmt.Errorf("Raising error in DestroyFirewallPolicy")
	default:
		return nil, fmt.Errorf("unexpected identifier %q sent to DestroyFirewallPolicy", identifier)
	}
}

func (f *fakeInstanceCloud) DestroyLoadBalancer(_ context.Context, identifier string) (*brightbox.LoadBalancer, error) {
	switch identifier {
	case foundLba:
		return nil, nil
	case errorLba:
		return nil, fmt.Errorf("Raising error in DestroyLoadBalancer")
	default:
		return nil, fmt.Errorf("unexpected identifier %q sent to DestroyLoadBalancer", identifier)
	}
}

func (f *fakeInstanceCloud) DestroyCloudIP(_ context.Context, identifier string) (*brightbox.CloudIP, error) {
	switch identifier {
	case publicCipID:
		return nil, nil
	case errorCipID:
		return nil, fmt.Errorf("Raising error in DestroyCloudIP")
	default:
		return nil, fmt.Errorf("unexpected identifier %q sent to DestroyCloudIP", identifier)
	}
}

func (f *fakeInstanceCloud) UnMapCloudIP(_ context.Context, identifier string) (*brightbox.CloudIP, error) {
	switch identifier {
	case publicCipID, publicCipID2:
		return nil, nil
	case errorCipID:
		return nil, fmt.Errorf("Raising error in UnMapCloudIP")
	default:
		return nil, fmt.Errorf("unexpected identifier %q sent to UnMapCloudIP", identifier)
	}
}

func (f *fakeInstanceCloud) CloudIP(_ context.Context, identifier string) (*brightbox.CloudIP, error) {
	result := &brightbox.CloudIP{
		ID: identifier,
	}
	switch identifier {
	case "cip-testy":
		result.LoadBalancer = &brightbox.LoadBalancer{ID: "lba-testy"}
	case "cip-12345":
		result.PublicIPv4 = publicIP
		result.PublicIPv6 = publicIPv6
		result.LoadBalancer = &brightbox.LoadBalancer{ID: foundLba}
	}
	return result, nil
}
