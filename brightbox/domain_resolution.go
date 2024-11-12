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
	"fmt"
	"net"
	"slices"

	brightbox "github.com/brightbox/gobrightbox/v2"
)

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
