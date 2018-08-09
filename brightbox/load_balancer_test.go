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
	"testing"

	"github.com/brightbox/gobrightbox"
	"github.com/go-test/deep"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

const (
	publicIP   = "180.180.180.180"
	fqdn       = "cip-180-180-180-180.gb1.brightbox.com"
	publicIP2  = "190.190.190.190"
	fqdn2      = "cip-190-190-190-190.gb1.brightbox.com"
	reverseDNS = "k8s-lb.example.com"
)

//Constant variables you can take the address of!
var (
	lbuid  types.UID = "9bde5f33-1379-4b8c-877a-777f5da4d766"
	lbname string    = "a9bde5f3313794b8c877a777f5da4d76"
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
							Protocol:   "tcp",
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   "tcp",
							Port:       80,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "dns",
							Protocol:   "udp",
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
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateServiceSpec(tc.service)
			if err == nil {
				t.Errorf("Expected error got nil")
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
					UID: "9d85099c-227c-46c0-a373-e954ec8eee2e",
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
				"dummy_cluster",
				tc.service,
			)
			if err != nil {
				t.Errorf("Error when none expected")
			} else if tc.exists != exists {
				t.Errorf("Exists status wrong, got %v, expected %v for %v", exists, tc.exists, cloudprovider.GetLoadBalancerName(tc.service))
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
				},
				Spec: v1.ServiceSpec{
					Type: v1.ServiceTypeLoadBalancer,
					Ports: []v1.ServicePort{
						{
							Name:       "https",
							Protocol:   "tcp",
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   "tcp",
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
				Name: &lbname,
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
					Type: loadBalancerTcpProtocol,
					Port: 31347,
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
							Protocol:   "tcp",
							Port:       443,
							TargetPort: intstr.FromInt(8080),
							NodePort:   31347,
						},
						{
							Name:       "http",
							Protocol:   "tcp",
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
				Name: &lbname,
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
				Name: &lbname,
				Healthcheck: &brightbox.LoadBalancerHealthcheck{
					Type: loadBalancerTcpProtocol,
					Port: 80,
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			lbopts := buildLoadBalancerOptions(tc.service, tc.nodes)
			if diff := deep.Equal(lbopts, tc.lbopts); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func (f *fakeInstanceCloud) CreateLoadBalancer(newLB *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	return nil, nil
}

func (f *fakeInstanceCloud) UpdateLoadBalancer(newLB *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	return nil, nil
}

func (f *fakeInstanceCloud) LoadBalancers() ([]brightbox.LoadBalancer, error) {
	return []brightbox.LoadBalancer{
		{
			Id:       "lba-test1",
			Name:     lbname,
			Status:   "Deleted",
			CloudIPs: nil,
		},
		{
			Id:     "lba-test2",
			Name:   lbname,
			Status: "Active",
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
			Name:   "abob",
			Status: "Active",
		},
	}, nil
}
