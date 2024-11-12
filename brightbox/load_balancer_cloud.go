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

	"github.com/brightbox/k8ssdk/v2"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
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
