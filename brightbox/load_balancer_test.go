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
	"testing"

	"github.com/brightbox/gobrightbox"
	"github.com/go-test/deep"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	publicIP    = "180.180.180.180"
	fqdn        = "cip-180-180-180-180.gb1.brightbox.com"
	publicIP2   = "190.190.190.190"
	fqdn2       = "cip-190-190-190-190.gb1.brightbox.com"
	reverseDNS  = "k8s-lb.example.com"
	foundLba    = "lba-found"
	errorLba    = "lba-error"
	newUID      = "9d85099c-227c-46c0-a373-e954ec8eee2e"
	clusterName = "test-cluster-name"
)

//Constant variables you can take the address of!
var (
	newlbname  string    = "a9d85099c227c46c0a373e954ec8eee2.default." + clusterName
	lbuid      types.UID = "9bde5f33-1379-4b8c-877a-777f5da4d766"
	lbname     string    = "a9bde5f3313794b8c877a777f5da4d76.default." + clusterName
	lberror    string    = "888888f3313794b8c877a777f5da4d76.default." + clusterName
	testPolicy string    = "round-robin"
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
						PublicIP:   publicIP,
						ReverseDns: reverseDNS,
						Fqdn:       fqdn,
					},
				},
			},
			status: &v1.LoadBalancerStatus{
				Ingress: []v1.LoadBalancerIngress{
					v1.LoadBalancerIngress{
						IP:       publicIP,
						Hostname: reverseDNS,
					},
				},
			},
		},
		"two-cloudips": {
			lb: &brightbox.LoadBalancer{
				CloudIPs: []brightbox.CloudIP{
					brightbox.CloudIP{
						PublicIP:   publicIP2,
						ReverseDns: "",
						Fqdn:       fqdn2,
					},
					brightbox.CloudIP{
						PublicIP:   publicIP,
						ReverseDns: reverseDNS,
						Fqdn:       fqdn,
					},
				},
			},
			status: &v1.LoadBalancerStatus{
				Ingress: []v1.LoadBalancerIngress{
					v1.LoadBalancerIngress{
						IP:       publicIP2,
						Hostname: fqdn2,
					},
					v1.LoadBalancerIngress{
						IP:       publicIP,
						Hostname: reverseDNS,
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
		"invalid-listener-protocol": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: newUID,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerListenerProtocol: "https",
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
			status: "Invalid Load Balancer Listener Protocol \"https\"",
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
						serviceAnnotationLoadBalancerHCProtocol: "https",
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
			status: "Invalid Load Balancer Healthcheck Protocol \"https\"",
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
						serviceAnnotationLoadBalancerHCTimeout: "100000",
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
			status: "\"" + serviceAnnotationLoadBalancerHCTimeout + "\" needs to be a positive number (strconv.ParseUint: parsing \"100000\": value out of range)",
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
			if err == nil {
				t.Errorf("Expected error %q got nil", tc.status)
			} else if err.Error() != tc.status {
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
					Ports:           nil,
					SessionAffinity: "None",
					LoadBalancerIP:  "",
				},
			},
			lbstatus: nil,
			exists:   false,
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
			lbstatus: &v1.LoadBalancerStatus{
				Ingress: []v1.LoadBalancerIngress{
					v1.LoadBalancerIngress{
						IP:       publicIP,
						Hostname: reverseDNS,
					},
				},
			},
			exists: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &cloud{
				client: fakeInstanceCloudClient(context.TODO()),
			}

			lb, exists, err := client.GetLoadBalancer(
				context.TODO(),
				clusterName,
				tc.service,
			)
			if err != nil {
				t.Errorf("Error when none expected")
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
	groklbname := grokLoadBalancerName(lbname)
	bufferSize := 16384
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
						serviceAnnotationLoadBalancerBufferSize:          "16384",
						serviceAnnotationLoadBalancerListenerIdleTimeout: "6000",
						serviceAnnotationLoadBalancerPolicy:              testPolicy,
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
				Nodes: &[]brightbox.LoadBalancerNode{
					{
						Node: "srv-gdqms",
					},
					{
						Node: "srv-230b7",
					},
				},
				Listeners: &[]brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerTcpProtocol,
						In:       443,
						Out:      31347,
						Timeout:  6000,
					},
					{
						Protocol: loadBalancerTcpProtocol,
						Timeout:  6000,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerTcpProtocol,
					Port:    31347,
					Request: "/",
				},
				BufferSize: &bufferSize,
				Policy:     &testPolicy,
			},
		},
		"overrideToTcpHealthcheck": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerHCProtocol:           loadBalancerHttpProtocol,
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
				Nodes: &[]brightbox.LoadBalancerNode{
					{
						Node: "srv-gdqms",
					},
					{
						Node: "srv-230b7",
					},
				},
				Listeners: &[]brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerTcpProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerTcpProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:          loadBalancerHttpProtocol,
					Port:          31347,
					Request:       "/",
					Timeout:       6000,
					Interval:      4000,
					ThresholdUp:   4,
					ThresholdDown: 5,
				},
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
				Nodes: &[]brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: &[]brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerTcpProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerTcpProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHttpProtocol,
					Port:    8080,
					Request: "/healthz",
				},
			},
		},
		"overrideToHttpHealthcheck": {
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					UID: lbuid,
					Annotations: map[string]string{
						serviceAnnotationLoadBalancerHCProtocol: loadBalancerTcpProtocol,
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
				Nodes: &[]brightbox.LoadBalancerNode{
					{
						Node: "srv-230b7",
					},
				},
				Listeners: &[]brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerTcpProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerTcpProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerTcpProtocol,
					Port:    8080,
					Request: "/",
				},
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
					Type:    loadBalancerTcpProtocol,
					Port:    80,
					Request: "/",
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &cloud{
				client: fakeInstanceCloudClient(context.TODO()),
			}
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
					v1.LoadBalancerIngress{
						IP:       publicIP,
						Hostname: reverseDNS,
					},
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &cloud{
				client: fakeInstanceCloudClient(context.TODO()),
			}

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
			lbopts: &brightbox.LoadBalancer{
				Id:     foundLba,
				Name:   grokLoadBalancerName(lbname),
				Status: lbActive,
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
						Protocol: loadBalancerTcpProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerTcpProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerTcpProtocol,
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
			lbopts: &brightbox.LoadBalancer{
				Name:   grokLoadBalancerName(newlbname),
				Status: lbActive,
				Nodes: []brightbox.Server{
					{
						Id: "srv-230b7",
					},
				},
				Listeners: []brightbox.LoadBalancerListener{
					{
						Protocol: loadBalancerTcpProtocol,
						In:       443,
						Out:      31347,
					},
					{
						Protocol: loadBalancerTcpProtocol,
						In:       80,
						Out:      31348,
					},
				},
				Healthcheck: brightbox.LoadBalancerHealthcheck{
					Type:    loadBalancerHttpProtocol,
					Port:    8080,
					Request: "/different/path",
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &cloud{
				client: fakeInstanceCloudClient(context.TODO()),
			}

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

func TestEnsureMappedCip(t *testing.T) {
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
			client := &cloud{
				client: fakeInstanceCloudClient(context.TODO()),
			}

			err := client.ensureMappedCip(tc.lb, tc.cip)
			if err != nil && !tc.err {
				t.Errorf("Error when not expected: %q", err.Error())
			}
		})
	}
}

func TestEnsureAllocatedCip(t *testing.T) {
	testCases := map[string]struct {
		service *v1.Service
		cip     *brightbox.CloudIP
	}{
		"LBIP_found": {
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
				Id:       "cip-12345",
				PublicIP: publicIP,
			},
		},
		"LBIP_notfound": {
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
		"name_found": {
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
				Id:       "cip-found",
				Name:     lbname,
				PublicIP: "240.240.240.240",
			},
		},
		"new_allocation": {
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
				Id:       "cip-67890",
				Name:     newlbname,
				PublicIP: publicIP2,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &cloud{
				client: fakeInstanceCloudClient(context.TODO()),
			}

			desc := client.GetLoadBalancerName(context.TODO(), clusterName, tc.service)
			cip, err := client.ensureAllocatedCip(desc, tc.service)
			if err != nil && tc.cip != nil {
				t.Errorf("Error when not expected %q", err.Error())
			} else if diff := deep.Equal(cip, tc.cip); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestDeletionFunctions(t *testing.T) {
	testCases := []string{
		lbname,
		"not-found",
		lberror,
	}
	for _, name := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &cloud{
				client: fakeInstanceCloudClient(context.TODO()),
			}
			client.ensureServerGroupDeleted(name)
			client.ensureFirewallClosed(name)
			client.ensureLoadBalancerDeletedByName(name)
			client.ensureCloudIPsDeleted(name)
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
			err: fmt.Errorf("CloudIps still mapped to load balancer %q", foundLba),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &cloud{
				client: fakeInstanceCloudClient(context.TODO()),
			}

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
			Id:       "cip-12345",
			PublicIP: publicIP,
		},
		{
			Id:       "cip-found",
			Name:     lbname,
			PublicIP: "240.240.240.240",
		},
		{
			Id:       "cip-error",
			Name:     lberror,
			PublicIP: "255.255.255.255",
		},
	}, nil
}

func (f *fakeInstanceCloud) CreateCloudIP(newCloudIP *brightbox.CloudIPOptions) (*brightbox.CloudIP, error) {
	cip := &brightbox.CloudIP{
		Id:       "cip-67890",
		PublicIP: publicIP2,
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
	server_list := make([]brightbox.Server, len(*newLB.Nodes))
	for i, v := range *newLB.Nodes {
		server_list[i].Id = v.Node
	}
	return &brightbox.LoadBalancer{
		Id:          newLB.Id,
		Name:        *newLB.Name,
		Status:      lbActive,
		Nodes:       server_list,
		Listeners:   *newLB.Listeners,
		Healthcheck: *newLB.Healthcheck,
	}, nil
}

func (f *fakeInstanceCloud) LoadBalancers() ([]brightbox.LoadBalancer, error) {
	return []brightbox.LoadBalancer{
		{
			Id:       "lba-test1",
			Name:     grokLoadBalancerName(lbname),
			Status:   "deleted",
			CloudIPs: nil,
		},
		{
			Id:     foundLba,
			Name:   grokLoadBalancerName(lbname),
			Status: lbActive,
			CloudIPs: []brightbox.CloudIP{
				brightbox.CloudIP{
					PublicIP:   publicIP,
					ReverseDns: reverseDNS,
					Fqdn:       fqdn,
				},
			},
		},
		{
			Id:     "lba-test3",
			Name:   grokLoadBalancerName("abob"),
			Status: lbActive,
		},
		{
			Id:     errorLba,
			Name:   grokLoadBalancerName(lberror),
			Status: lbActive,
			CloudIPs: []brightbox.CloudIP{
				brightbox.CloudIP{
					PublicIP:   publicIP,
					ReverseDns: reverseDNS,
					Fqdn:       fqdn,
				},
			},
		},
	}, nil
}

func (f *fakeInstanceCloud) AddServersToServerGroup(identifier string, serverIds []string) (*brightbox.ServerGroup, error) {
	switch identifier {
	case "grp-found":
		groups, _ := f.ServerGroups()
		return &groups[0], nil
	default:
		result := &brightbox.ServerGroup{
			Id:      identifier,
			Name:    "Fake Name After AddServers",
			Servers: mapServerIdsToServers(serverIds),
		}
		return result, nil
	}
}

func (f *fakeInstanceCloud) RemoveServersFromServerGroup(identifier string, serverIds []string) (*brightbox.ServerGroup, error) {
	switch identifier {
	case "grp-found":
		groups, _ := f.ServerGroups()
		return &groups[0], nil
	default:
		result := &brightbox.ServerGroup{
			Id:      identifier,
			Name:    "Fake Name After RemoveServers",
			Servers: mapServerIdsToServers(serverIds),
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

func mapServerIdsToServers(serverIds []string) []brightbox.Server {
	result := make([]brightbox.Server, len(serverIds))
	for i := range serverIds {
		result[i].Id = serverIds[i]
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
	case "cip-found":
		return nil
	case "cip-error":
		return fmt.Errorf("Raising error in DestroyCloudIP")
	default:
		return fmt.Errorf("unexpected identifier %q sent to DestroyCloudIP", identifier)
	}
}

func (f *fakeInstanceCloud) CloudIP(identifier string) (*brightbox.CloudIP, error) {
	var lbId string
	switch identifier {
	case "cip-testy":
		lbId = "lba-testy"
	case "cip-12345":
		lbId = foundLba
	}
	result := &brightbox.CloudIP{
		Id: identifier,
	}
	if lbId != "" {
		result.LoadBalancer = &brightbox.LoadBalancer{Id: lbId}
	}
	return result, nil
}
