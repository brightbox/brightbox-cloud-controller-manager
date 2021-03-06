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

	brightbox "github.com/brightbox/gobrightbox"
	"github.com/brightbox/gobrightbox/status"
	"github.com/brightbox/k8ssdk"
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

//Constant variables you can take the address of!
var (
	newlbname     string            = "a9d85099c227c46c0a373e954ec8eee2.default." + clusterName
	lbuid         types.UID         = "9bde5f33-1379-4b8c-877a-777f5da4d766"
	lbname        string            = "a9bde5f3313794b8c877a777f5da4d76.default." + clusterName
	lberror       string            = "888888f3313794b8c877a777f5da4d76.default." + clusterName
	testPolicy    string            = "round-robin"
	trueVar       bool              = true
	falseVar      bool              = false
	groklbname    string            = lbname
	groknewlbname string            = newlbname
	bufferSize    int               = 16384
	resolvCip     brightbox.CloudIP = brightbox.CloudIP{
		Id:         "cip-vsalc",
		PublicIP:   "109.107.39.92",
		PublicIPv4: "109.107.39.92",
		PublicIPv6: "2a02:1348:ffff:ffff::6d6b:275c",
		Fqdn:       resolvedDomain,
		ReverseDns: "cip-109-107-39-92.gb1s.brightbox.com",
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
						Id:         publicCipID,
						PublicIPv4: publicIP,
						PublicIPv6: publicIPv6,
						ReverseDns: reverseDNS,
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
						Id:         publicCipID2,
						PublicIPv4: publicIP2,
						PublicIPv6: publicIPv62,
						ReverseDns: "",
						Fqdn:       fqdn2,
						Name:       "manually allocated",
					},
					brightbox.CloudIP{
						Id:         publicCipID,
						PublicIPv4: publicIP,
						PublicIPv6: publicIPv6,
						ReverseDns: reverseDNS,
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
			status: "Invalid Load Balancer Policy \"magic-routing\"",
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
			status: "Invalid Load Balancer Listener Proxy Protocol \"v1-ssl\"",
		},
		"Domains with TCP": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: loadBalancerTCPProtocol,
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
						serviceAnnotationLoadBalancerListenerProtocol: loadBalancerTCPProtocol,
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
						serviceAnnotationLoadBalancerListenerProxyProtocol: loadBalancerProxyV2SslCn,
						serviceAnnotationLoadBalancerListenerProtocol:      loadBalancerTCPProtocol,
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
						serviceAnnotationLoadBalancerListenerProxyProtocol: loadBalancerProxyV2Ssl,
						serviceAnnotationLoadBalancerListenerProtocol:      loadBalancerHTTPProtocol,
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
						serviceAnnotationLoadBalancerListenerProtocol: loadBalancerHTTPProtocol,
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
						serviceAnnotationLoadBalancerListenerProtocol: loadBalancerHTTPWSProtocol,
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
						serviceAnnotationLoadBalancerListenerProtocol: loadBalancerHTTPProtocol,
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
						serviceAnnotationLoadBalancerListenerProtocol: loadBalancerHTTPProtocol,
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
			status: "Invalid Load Balancer Listener Protocol \"gopher\"",
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
						serviceAnnotationLoadBalancerHCProtocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
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
			status: "Invalid Load Balancer Healthcheck Protocol \"" + sslUpgradeProtocol[loadBalancerHTTPProtocol] + "\"",
		},
		"https without domains": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: loadBalancerHTTPProtocol,
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
		"invalid-value-for-buffer-size": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerBufferSize: "buffer",
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
			status: "\"" + serviceAnnotationLoadBalancerBufferSize + "\" needs to be a positive number (strconv.ParseUint: parsing \"buffer\": invalid syntax)",
		},
		"invalid-small-buffer-size": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerBufferSize: "1023",
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
			status: "\"" + serviceAnnotationLoadBalancerBufferSize + "\" needs to be no less than 1024",
		},
		"invalid-big-buffer-size": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerBufferSize: "16385",
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
			status: "\"" + serviceAnnotationLoadBalancerBufferSize + "\" needs to be no more than 16384",
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
						serviceAnnotationLoadBalancerBufferSize:          strconv.Itoa(bufferSize),
						serviceAnnotationLoadBalancerListenerIdleTimeout: strconv.Itoa(testTimeout),
						serviceAnnotationLoadBalancerPolicy:              testPolicy,
						serviceAnnotationLoadBalancerSslDomains:          resolvedDomain + "," + fqdn,
						serviceAnnotationLoadBalancerListenerProtocol:    loadBalancerHTTPProtocol,
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
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
						Timeout:  testTimeout,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						Timeout:  testTimeout,
						In:       80,
						Out:      31348,
					},
				},
				Domains: &[]string{resolvedDomain, fqdn},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    31347,
					Request: "/healthz",
				},
				BufferSize:    &bufferSize,
				Policy:        &testPolicy,
				HttpsRedirect: &trueVar,
			},
		},
		"standard_proxy_protocol": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerBufferSize:            strconv.Itoa(bufferSize),
						serviceAnnotationLoadBalancerListenerIdleTimeout:   strconv.Itoa(testTimeout),
						serviceAnnotationLoadBalancerPolicy:                testPolicy,
						serviceAnnotationLoadBalancerSslDomains:            resolvedDomain + "," + fqdn,
						serviceAnnotationLoadBalancerListenerProtocol:      loadBalancerHTTPProtocol,
						serviceAnnotationLoadBalancerListenerProxyProtocol: loadBalancerProxyV2SslCn,
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
						Protocol:      sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:            443,
						Out:           31347,
						Timeout:       testTimeout,
						ProxyProtocol: loadBalancerProxyV2SslCn,
					},
					{
						Protocol:      loadBalancerHTTPProtocol,
						Timeout:       testTimeout,
						In:            80,
						Out:           31348,
						ProxyProtocol: loadBalancerProxyV2SslCn,
					},
				},
				Domains: &[]string{resolvedDomain, fqdn},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    31347,
					Request: "/healthz",
				},
				BufferSize:    &bufferSize,
				Policy:        &testPolicy,
				HttpsRedirect: &trueVar,
			},
		},
		"websocket": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerBufferSize:          strconv.Itoa(bufferSize),
						serviceAnnotationLoadBalancerListenerIdleTimeout: strconv.Itoa(testTimeout),
						serviceAnnotationLoadBalancerPolicy:              testPolicy,
						serviceAnnotationLoadBalancerSslDomains:          resolvedDomain + "," + fqdn,
						serviceAnnotationLoadBalancerListenerProtocol:    loadBalancerHTTPWSProtocol,
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
						Protocol: sslUpgradeProtocol[loadBalancerHTTPWSProtocol],
						In:       443,
						Out:      31347,
						Timeout:  testTimeout,
					},
					{
						Protocol: loadBalancerHTTPWSProtocol,
						Timeout:  testTimeout,
						In:       80,
						Out:      31348,
					},
				},
				Domains: &[]string{resolvedDomain, fqdn},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    31347,
					Request: "/healthz",
				},
				BufferSize:    &bufferSize,
				Policy:        &testPolicy,
				HttpsRedirect: &trueVar,
			},
		},
		"extraSSLports": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerBufferSize:          strconv.Itoa(bufferSize),
						serviceAnnotationLoadBalancerListenerIdleTimeout: strconv.Itoa(testTimeout),
						serviceAnnotationLoadBalancerPolicy:              testPolicy,
						serviceAnnotationLoadBalancerSslDomains:          resolvedDomain + "," + fqdn,
						serviceAnnotationLoadBalancerSSLPorts:            "fancy,3030",
						serviceAnnotationLoadBalancerListenerProtocol:    loadBalancerHTTPProtocol,
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
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
						Timeout:  testTimeout,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						Timeout:  testTimeout,
						In:       80,
						Out:      31348,
					},
					{
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						Timeout:  testTimeout,
						In:       5050,
						Out:      31348,
					},
					{
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						Timeout:  testTimeout,
						In:       3030,
						Out:      31348,
					},
				},
				Domains: &[]string{resolvedDomain, fqdn},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    31347,
					Request: "/healthz",
				},
				BufferSize:    &bufferSize,
				Policy:        &testPolicy,
				HttpsRedirect: &trueVar,
			},
		},
		"OverrideToTcpListener": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol:    loadBalancerTCPProtocol,
						serviceAnnotationLoadBalancerBufferSize:          strconv.Itoa(bufferSize),
						serviceAnnotationLoadBalancerListenerIdleTimeout: strconv.Itoa(testTimeout),
						serviceAnnotationLoadBalancerPolicy:              testPolicy,
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
						Protocol: loadBalancerTCPProtocol,
						Timeout:  testTimeout,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerTCPProtocol,
					Port:    31348,
					Request: "/",
				},
				BufferSize:    &bufferSize,
				Policy:        &testPolicy,
				HttpsRedirect: &falseVar,
			},
		},
		"overrideToHttpHealthcheck": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol:     loadBalancerTCPProtocol,
						serviceAnnotationLoadBalancerHCProtocol:           loadBalancerHTTPProtocol,
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
						Protocol: loadBalancerTCPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerTCPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:          loadBalancerHTTPProtocol,
					Port:          31347,
					Request:       "/healthz",
					Timeout:       6000,
					Interval:      4000,
					ThresholdUp:   4,
					ThresholdDown: 5,
				},
				HttpsRedirect: &falsevar,
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
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
		},
		"overrideToTcpHealthcheck": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerHCProtocol: loadBalancerTCPProtocol,
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
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerTCPProtocol,
					Port:    8080,
					Request: "/",
				},
				HttpsRedirect: &falsevar,
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
					Type:    loadBalancerHTTPProtocol,
					Port:    80,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
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
						serviceAnnotationLoadBalancerListenerProtocol: loadBalancerTCPProtocol,
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
				Id:     foundLba,
				Name:   lbname,
				Status: status.Active,
				Nodes: []brightbox.Server{
					{
						Id: "srv-gdqms",
					},
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerTCPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerTCPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerTCPProtocol,
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
						serviceAnnotationLoadBalancerListenerProtocol: loadBalancerTCPProtocol,
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
				Id:     foundLba,
				Name:   lbname,
				Status: status.Active,
				Nodes: []brightbox.Server{
					{
						Id: "srv-gdqms",
					},
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerTCPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerTCPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerTCPProtocol,
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
				Status: status.Active,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/different/path",
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()

			desc := client.GetLoadBalancerName(context.TODO(), clusterName, tc.service)
			lbopts, err := client.ensureLoadBalancerFromService(desc, tc.service, tc.nodes)
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
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
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
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Domains: &[]string{resolvedDomain, fqdn},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: false,
		},
		"swap domains": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
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
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Domains: &[]string{fqdn, resolvedDomain},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: false,
		},
		"add_domain": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
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
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Domains: &[]string{resolvedDomain, fqdn, reverseDNS},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"remove_domain": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
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
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Domains: &[]string{fqdn},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"change domain": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
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
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Domains: &[]string{reverseDNS, fqdn},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: sslUpgradeProtocol[loadBalancerHTTPProtocol],
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"No change": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: false,
		},
		"add listener": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
					{
						Protocol: loadBalancerTCPProtocol,
						In:       25,
						Out:      32456,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"remove listener": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"change_listener": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31350,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"add_proxy_protocol": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol:      loadBalancerHTTPProtocol,
						In:            443,
						Out:           31347,
						ProxyProtocol: "v2",
					},
					{
						Protocol:      loadBalancerHTTPProtocol,
						In:            80,
						Out:           31348,
						ProxyProtocol: "v2",
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"remove_proxy_protocol": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol:      loadBalancerHTTPProtocol,
						In:            443,
						Out:           31347,
						ProxyProtocol: "v2",
					},
					{
						Protocol:      loadBalancerHTTPProtocol,
						In:            80,
						Out:           31348,
						ProxyProtocol: "v2",
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"change_proxy_protocol": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol:      loadBalancerHTTPProtocol,
						In:            443,
						Out:           31347,
						ProxyProtocol: "v2",
					},
					{
						Protocol:      loadBalancerHTTPProtocol,
						In:            80,
						Out:           31348,
						ProxyProtocol: "v2",
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol:      loadBalancerHTTPProtocol,
						In:            443,
						Out:           31347,
						ProxyProtocol: "v2-ssl",
					},
					{
						Protocol:      loadBalancerHTTPProtocol,
						In:            80,
						Out:           31348,
						ProxyProtocol: "v2-ssl",
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"change node": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-newon",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"remove node": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:    foundLba,
				Name:  &groklbname,
				Nodes: []brightbox.LoadBalancerNode{},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"name_change": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groknewlbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"add node": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
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
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
				HttpsRedirect: &falsevar,
			},
			expected: true,
		},
		"Healthcheck change": {
			lb: &brightbox.LoadBalancer{
				Id:   foundLba,
				Name: groklbname,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
			lbopts: brightbox.LoadBalancerOptions{
				Id:   foundLba,
				Name: &groklbname,
				Nodes: []brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerHTTPProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHTTPProtocol,
					Port:    8080,
					Request: "/check",
				},
				HttpsRedirect: &falsevar,
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
				Id: "lba-testy",
				CloudIPs: []brightbox.CloudIP{
					{
						Id: "cip-testy",
					},
				},
			},
			cip: &brightbox.CloudIP{
				Id: "cip-testy",
				LoadBalancer: &brightbox.LoadBalancer{
					Id: "lba-testy",
				},
				Status: "mapped",
			},
			err: false,
		},
		"badmap": {
			lb: &brightbox.LoadBalancer{
				Id:       "lba-testy",
				CloudIPs: []brightbox.CloudIP{},
			},
			cip: &brightbox.CloudIP{
				Id:     "cip-testy",
				Status: "mapped",
			},
			err: true,
		},
		"unmapped": {
			lb: &brightbox.LoadBalancer{
				Id:       "lba-testy",
				CloudIPs: []brightbox.CloudIP{},
			},
			cip: &brightbox.CloudIP{
				Id:     "cip-testy",
				Status: "unmapped",
			},
			err: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := makeFakeInstanceCloudClient()

			err := client.EnsureMappedCloudIP(tc.lb, tc.cip)
			if err != nil && !tc.err {
				t.Errorf("Error when not expected: %q", err.Error())
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
				Id:         "cip-12345",
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
				Id:         publicCipID,
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
				Id:         "cip-67890",
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
				Id:         "cip-12345",
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

			desc := client.GetLoadBalancerName(context.TODO(), clusterName, tc.service)
			cip, err := client.ensureAllocatedCloudIP(desc, tc.service)
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
				Id: "lba-oldip",
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						Id:         publicCipID,
						PublicIPv4: publicIP,
						PublicIPv6: publicIPv6,
						ReverseDns: reverseDNS,
						Fqdn:       fqdn,
						Name:       "test",
					},
				},
			},
			cip: &brightbox.CloudIP{
				Id:         publicCipID,
				PublicIPv4: publicIP,
				PublicIPv6: publicIPv6,
				ReverseDns: reverseDNS,
				Fqdn:       fqdn,
				Name:       "test",
			},
			name: "test",
		},
		"no_change_manual": {
			lb: &brightbox.LoadBalancer{
				Id: "lba-oldip",
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						Id:         publicCipID2,
						PublicIPv4: publicIP2,
						PublicIPv6: publicIPv62,
						Fqdn:       fqdn2,
						Name:       "manually allocated",
					},
				},
			},
			cip: &brightbox.CloudIP{
				Id:         publicCipID2,
				PublicIPv4: publicIP2,
				PublicIPv6: publicIPv62,
				Fqdn:       fqdn2,
				Name:       "manually allocated",
			},
			name: "test",
		},
		"changed_delete": {
			lb: &brightbox.LoadBalancer{
				Id: "lba-oldip",
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						Id:         publicCipID2,
						PublicIPv4: publicIP2,
						PublicIPv6: publicIPv62,
						Fqdn:       fqdn2,
						Name:       "manually allocated",
					},
					brightbox.CloudIP{
						Id:         publicCipID,
						PublicIPv4: publicIP,
						PublicIPv6: publicIPv6,
						ReverseDns: reverseDNS,
						Fqdn:       fqdn,
						Name:       "test",
					},
				},
			},
			cip: &brightbox.CloudIP{
				Id:         publicCipID2,
				PublicIPv4: publicIP2,
				PublicIPv6: publicIPv62,
				Fqdn:       fqdn2,
				Name:       "manually allocated",
			},
			name: "test",
		},
		"changed_unmap": {
			lb: &brightbox.LoadBalancer{
				Id: "lba-oldip",
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						Id:         publicCipID2,
						PublicIPv4: publicIP2,
						PublicIPv6: publicIPv62,
						Fqdn:       fqdn2,
						Name:       "manually allocated",
					},
					brightbox.CloudIP{
						Id:         publicCipID,
						PublicIPv4: publicIP,
						PublicIPv6: publicIPv6,
						ReverseDns: reverseDNS,
						Fqdn:       fqdn,
						Name:       "test",
					},
				},
			},
			cip: &brightbox.CloudIP{
				Id:         publicCipID,
				PublicIPv4: publicIP,
				PublicIPv6: publicIPv6,
				ReverseDns: reverseDNS,
				Fqdn:       fqdn,
				Name:       "test",
			},
			name: "test",
		},
		"already unmapped": {
			lb: &brightbox.LoadBalancer{
				Id:       "lba-oldip",
				CloudIPs: []brightbox.CloudIP{},
			},
			cip: &brightbox.CloudIP{
				Id:         publicCipID,
				PublicIPv4: publicIP,
				PublicIPv6: publicIPv6,
				ReverseDns: reverseDNS,
				Fqdn:       fqdn,
				Name:       "test",
			},
			name: "test",
		},
		"already unmapped_manual": {
			lb: &brightbox.LoadBalancer{
				Id:       "lba-oldip",
				CloudIPs: []brightbox.CloudIP{},
			},
			cip: &brightbox.CloudIP{
				Id:         publicCipID2,
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

			err := client.EnsureOldCloudIPsDeposed(tc.lb.CloudIPs, tc.cip.Id, tc.name)
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

			client.ensureServerGroupDeleted(name)
			client.ensureLoadBalancerDeletedByName(name)
			client.ensureFirewallClosed(name)
			client.ensureCloudIPsDeleted(name)
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

func (f *fakeInstanceCloud) MapCloudIP(identifier string, destination string) error {
	return nil
}

func (f *fakeInstanceCloud) CloudIPs() ([]brightbox.CloudIP, error) {
	return []brightbox.CloudIP{
		{
			Id:         "cip-12345",
			PublicIPv4: publicIP,
			PublicIPv6: publicIPv6,
		},
		{
			Id:         publicCipID,
			Name:       lbname,
			PublicIPv4: "240.240.240.240",
		},
		{
			Id:         errorCipID,
			Name:       lberror,
			PublicIPv4: "255.255.255.255",
		},
	}, nil
}

func (f *fakeInstanceCloud) CreateCloudIP(newCloudIP *brightbox.CloudIPOptions) (*brightbox.CloudIP, error) {
	cip := &brightbox.CloudIP{
		Id:         "cip-67890",
		PublicIPv4: publicIP2,
		PublicIPv6: publicIPv62,
	}
	if newCloudIP.Name != nil {
		cip.Name = *newCloudIP.Name
	}
	return cip, nil
}

func (f *fakeInstanceCloud) CreateLoadBalancer(newLB *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	return f.UpdateLoadBalancer(newLB)
}

func (f *fakeInstanceCloud) UpdateLoadBalancer(newLB *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	server_list := make([]brightbox.Server, len(newLB.Nodes))
	for i, v := range newLB.Nodes {
		server_list[i].Id = v.Node
	}
	return &brightbox.LoadBalancer{
		Id:          newLB.Id,
		Name:        *newLB.Name,
		Status:      status.Active,
		Nodes:       server_list,
		Listeners:   newLB.Listeners,
		Healthcheck: *newLB.Healthcheck,
	}, nil
}

func (f *fakeInstanceCloud) LoadBalancer(id string) (*brightbox.LoadBalancer, error) {
	list, _ := f.LoadBalancers()
	for _, balancer := range list {
		if balancer.Id == id {
			return &balancer, nil
		}
	}
	return nil, fmt.Errorf("unexpected identifier %q sent to LoadBalancer", id)
}

func (f *fakeInstanceCloud) LoadBalancers() ([]brightbox.LoadBalancer, error) {
	return []brightbox.LoadBalancer{
		{
			Id:       "lba-test1",
			Name:     lbname,
			Status:   "deleted",
			CloudIPs: nil,
		},
		{
			Id:     foundLba,
			Name:   lbname,
			Status: status.Active,
			CloudIPs: []brightbox.CloudIP{
				brightbox.CloudIP{
					Id:         "cip-12345",
					PublicIPv4: publicIP,
					PublicIPv6: publicIPv6,
					ReverseDns: reverseDNS,
					Fqdn:       fqdn,
				},
			},
		},
		{
			Id:     "lba-test3",
			Name:   "abob",
			Status: status.Active,
		},
		{
			Id:     errorLba,
			Name:   lberror,
			Status: status.Active,
			CloudIPs: []brightbox.CloudIP{
				brightbox.CloudIP{
					Id:         publicCipID,
					PublicIPv4: publicIP,
					PublicIPv6: publicIPv6,
					ReverseDns: reverseDNS,
					Fqdn:       fqdn,
				},
			},
		},
	}, nil
}

func (f *fakeInstanceCloud) AddServersToServerGroup(identifier string, serverIDs []string) (*brightbox.ServerGroup, error) {
	switch identifier {
	case "grp-found":
		groups, _ := f.ServerGroups()
		return &groups[0], nil
	default:
		result := &brightbox.ServerGroup{
			Id:      identifier,
			Name:    "Fake Name After AddServers",
			Servers: mapServerIDsToServers(serverIDs),
		}
		return result, nil
	}
}

func (f *fakeInstanceCloud) RemoveServersFromServerGroup(identifier string, serverIDs []string) (*brightbox.ServerGroup, error) {
	switch identifier {
	case "grp-found":
		groups, _ := f.ServerGroups()
		return &groups[0], nil
	default:
		result := &brightbox.ServerGroup{
			Id:      identifier,
			Name:    "Fake Name After RemoveServers",
			Servers: mapServerIDsToServers(serverIDs),
		}
		return result, nil
	}
}

func (f *fakeInstanceCloud) ServerGroups() ([]brightbox.ServerGroup, error) {
	result := []brightbox.ServerGroup{
		brightbox.ServerGroup{
			Id:   "grp-found",
			Name: lbname,
			Servers: []brightbox.Server{
				{
					Id: "srv-gdqms",
				},
				{
					Id: "srv-230b7",
				},
			},
			FirewallPolicy: &brightbox.FirewallPolicy{
				Id:   "fwp-found",
				Name: lbname,
				Rules: []brightbox.FirewallRule{
					{
						Id:          "fwr-found",
						Description: lbname,
					},
				},
			},
		},
		brightbox.ServerGroup{
			Id:   "grp-error",
			Name: lberror,
			Servers: []brightbox.Server{
				{
					Id: "srv-gdqms",
				},
				{
					Id: "srv-230b7",
				},
			},
			FirewallPolicy: &brightbox.FirewallPolicy{
				Id:   "fwp-error",
				Name: lberror,
				Rules: []brightbox.FirewallRule{
					{
						Id:          "fwr-found",
						Description: lberror,
					},
				},
			},
		},
	}
	return result, nil
}

func (f *fakeInstanceCloud) CreateServerGroup(newServerGroup *brightbox.ServerGroupOptions) (*brightbox.ServerGroup, error) {
	result := &brightbox.ServerGroup{
		Id: "grp-testy",
	}
	if newServerGroup.Name != nil {
		result.Name = *newServerGroup.Name
	}
	if newServerGroup.Description != nil {
		result.Description = *newServerGroup.Description
	}
	return result, nil
}

func (f *fakeInstanceCloud) CreateFirewallPolicy(policyOptions *brightbox.FirewallPolicyOptions) (*brightbox.FirewallPolicy, error) {
	result := &brightbox.FirewallPolicy{
		Id:   "fwp-testy",
		Name: *policyOptions.Name,
		ServerGroup: &brightbox.ServerGroup{
			Id:   *policyOptions.ServerGroup,
			Name: *policyOptions.Name,
		},
	}
	return result, nil
}

func (f *fakeInstanceCloud) CreateFirewallRule(ruleOptions *brightbox.FirewallRuleOptions) (*brightbox.FirewallRule, error) {
	result := &brightbox.FirewallRule{
		Id:          "fwr-testy",
		Description: "After Create Firewll Rule",
		FirewallPolicy: brightbox.FirewallPolicy{
			Id:   ruleOptions.FirewallPolicy,
			Name: "After Create Firewall Rule",
		},
	}
	return result, nil
}

func (f *fakeInstanceCloud) UpdateFirewallRule(ruleOptions *brightbox.FirewallRuleOptions) (*brightbox.FirewallRule, error) {
	result := &brightbox.FirewallRule{
		Id:          ruleOptions.Id,
		Description: *ruleOptions.Description,
		FirewallPolicy: brightbox.FirewallPolicy{
			Id:   ruleOptions.FirewallPolicy,
			Name: *ruleOptions.Description,
		},
	}
	return result, nil
}

func mapServerIDsToServers(serverIDs []string) []brightbox.Server {
	result := make([]brightbox.Server, len(serverIDs))
	for i := range serverIDs {
		result[i].Id = serverIDs[i]
	}
	return result
}

func (f *fakeInstanceCloud) FirewallPolicies() ([]brightbox.FirewallPolicy, error) {
	result := []brightbox.FirewallPolicy{
		brightbox.FirewallPolicy{
			Id:   "fwp-found",
			Name: lbname,
			Rules: []brightbox.FirewallRule{
				{
					Id:          "fwr-found",
					Description: lbname,
				},
			},
		},
		brightbox.FirewallPolicy{
			Id:   "fwp-error",
			Name: lberror,
			Rules: []brightbox.FirewallRule{
				{
					Id:          "fwr-error",
					Description: lberror,
				},
			},
			ServerGroup: &brightbox.ServerGroup{
				Id:   "grp-error",
				Name: lberror,
				Servers: []brightbox.Server{
					{
						Id: "srv-gdqms",
					},
					{
						Id: "srv-230b7",
					},
				},
			},
		},
	}
	return result, nil
}

func (f *fakeInstanceCloud) DestroyServerGroup(identifier string) error {
	switch identifier {
	case "grp-found":
		return nil
	case "grp-error":
		return fmt.Errorf("Raising error in DestroyServerGroup")
	default:
		return fmt.Errorf("unexpected identifier %q sent to DestroyServerGroup", identifier)
	}
}

func (f *fakeInstanceCloud) DestroyFirewallPolicy(identifier string) error {
	switch identifier {
	case "fwp-found":
		return nil
	case "fwp-error":
		return fmt.Errorf("Raising error in DestroyFirewallPolicy")
	default:
		return fmt.Errorf("unexpected identifier %q sent to DestroyFirewallPolicy", identifier)
	}
}

func (f *fakeInstanceCloud) DestroyLoadBalancer(identifier string) error {
	switch identifier {
	case foundLba:
		return nil
	case errorLba:
		return fmt.Errorf("Raising error in DestroyLoadBalancer")
	default:
		return fmt.Errorf("unexpected identifier %q sent to DestroyLoadBalancer", identifier)
	}
}

func (f *fakeInstanceCloud) DestroyCloudIP(identifier string) error {
	switch identifier {
	case publicCipID:
		return nil
	case errorCipID:
		return fmt.Errorf("Raising error in DestroyCloudIP")
	default:
		return fmt.Errorf("unexpected identifier %q sent to DestroyCloudIP", identifier)
	}
}

func (f *fakeInstanceCloud) UnMapCloudIP(identifier string) error {
	switch identifier {
	case publicCipID, publicCipID2:
		return nil
	case errorCipID:
		return fmt.Errorf("Raising error in UnMapCloudIP")
	default:
		return fmt.Errorf("unexpected identifier %q sent to UnMapCloudIP", identifier)
	}
}

func (f *fakeInstanceCloud) CloudIP(identifier string) (*brightbox.CloudIP, error) {
	result := &brightbox.CloudIP{
		Id: identifier,
	}
	switch identifier {
	case "cip-testy":
		result.LoadBalancer = &brightbox.LoadBalancer{Id: "lba-testy"}
	case "cip-12345":
		result.PublicIPv4 = publicIP
		result.PublicIPv6 = publicIPv6
		result.LoadBalancer = &brightbox.LoadBalancer{Id: foundLba}
	}
	return result, nil
}
