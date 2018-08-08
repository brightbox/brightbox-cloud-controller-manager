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
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// Ignoring clusterName completely as it doesn't appear to be used anywhere else at the moment

func (c *cloud) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	name := cloudprovider.GetLoadBalancerName(service)
	glog.V(2).Infof("GetLoadBalancer(%v, %v)", clusterName, name)
	lb, err := c.getLoadBalancerByName(name)
	if lb == nil || err != nil {
		return nil, false, err
	}
	return toLoadBalancerStatus(lb), true, nil
}

func (c *cloud) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	name := cloudprovider.GetLoadBalancerName(service)
	annotations := service.Annotations
	glog.V(2).Infof("EnsureLoadBalancer(%v, %v, %v, %v, %v, %v)",
		clusterName, name, service.Namespace, service.Spec.LoadBalancerIP, service.Spec.Ports, annotations)
	if err := validateServiceSpec(service); err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *cloud) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	glog.V(2).Infof("UpdateLoadBalancer called - delegating")
	_, err := c.EnsureLoadBalancer(ctx, clusterName, service, nodes)
	return err
}

func (c *cloud) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	name := cloudprovider.GetLoadBalancerName(service)
	glog.V(2).Infof("EnsureLoadBalancerDeleted(%v, %v, %v, %v)", clusterName, name, service.Namespace, service.Spec.LoadBalancerIP)
	return nil
}

// lb is expected to be not nil
func toLoadBalancerStatus(lb *brightbox.LoadBalancer) *v1.LoadBalancerStatus {
	status := &v1.LoadBalancerStatus{}
	if len(lb.CloudIPs) > 0 {
		status.Ingress = make([]v1.LoadBalancerIngress, len(lb.CloudIPs))
		for i, v := range lb.CloudIPs {
			status.Ingress[i] = v1.LoadBalancerIngress{
				Hostname: selectHostname(&v),
				IP:       v.PublicIP,
			}
		}
	}
	return status
}

func selectHostname(ip *brightbox.CloudIP) string {
	if ip.ReverseDns != "" {
		return ip.ReverseDns
	} else {
		return ip.Fqdn
	}
}

func validateServiceSpec(service *v1.Service) error {
	if service.Spec.SessionAffinity != v1.ServiceAffinityNone {
		//  supports sticky sessions, but only when configured for HTTP/HTTPS
		return fmt.Errorf("unsupported load balancer affinity: %v", service.Spec.SessionAffinity)
	}
	if len(service.Spec.Ports) == 0 {
		return fmt.Errorf("requested load balancer with no ports")
	}
	return nil
}
