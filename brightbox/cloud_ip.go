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
	"net"
	"time"

	brightbox "github.com/brightbox/gobrightbox/v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
)

const (
	// Delete Backoff settings
	loadbalancerActiveInitDelay = 1 * time.Second
	loadbalancerActiveFactor    = 1.2
	loadbalancerActiveSteps     = 5
)

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
