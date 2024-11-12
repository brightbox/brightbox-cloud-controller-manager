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
	"net/url"
	"regexp"
	"strconv"

	"github.com/brightbox/gobrightbox/v2/enums/balancingpolicy"
	"github.com/brightbox/gobrightbox/v2/enums/healthchecktype"
	"github.com/brightbox/gobrightbox/v2/enums/listenerprotocol"
	"github.com/brightbox/gobrightbox/v2/enums/proxyprotocol"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

const (
	// Maximum number of bits in unsigned integers specified in annotations.
	maxBits = 32
)

var cloudIPPattern = regexp.MustCompile(`^cip-[0-9a-z]{5,}$`)

// If annotation is missing returns zero value
func parseUintAnnotation(annotationList map[string]string, annotation string) (uint, error) {
	klog.V(6).Infof("parseUintAnnotation(%+v, %+v)", annotationList, annotation)
	strValue, ok := annotationList[annotation]
	if !ok {
		return 0, nil
	}
	val, err := strconv.ParseUint(strValue, 10, maxBits)
	klog.V(6).Infof("Value Converted from %+v to %+v", strValue, val)
	return uint(val), err
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
	if !sslPortFound && protocol == listenerprotocol.Http {
		_, ports := apiservice.Annotations[serviceAnnotationLoadBalancerSSLPorts]
		_, domains := apiservice.Annotations[serviceAnnotationLoadBalancerSslDomains]
		if ports || domains {
			return fmt.Errorf("SSL support requires a Port definition for %d", standardSSLPort)
		}
	}
	// CloudIP allocation annotation and spec.loadBalancerIP conflict
	if apiservice.Spec.LoadBalancerIP != "" {
		if _, ok := apiservice.Annotations[serviceAnnotationLoadBalancerCloudipAllocations]; ok {
			return fmt.Errorf("Remove obsolete field: spec.loadBalancerIP")
		}
	}
	return validateAnnotations(apiservice.Annotations)
}

func validateAnnotations(annotationList map[string]string) error {
	for annotation, value := range annotationList {
		switch annotation {
		case serviceAnnotationLoadBalancerPolicy:
			if _, err := balancingpolicy.ParseEnum(value); err != nil {
				return fmt.Errorf("Invalid Load Balancer Policy %q: %w", value, err)
			}
		case serviceAnnotationLoadBalancerListenerProtocol:
			valueEnum, err := listenerprotocol.ParseEnum(value)
			if err != nil {
				return fmt.Errorf("Invalid Load Balancer Listener Protocol %q: %w", value, err)
			}
			if valueEnum == listenerprotocol.Tcp {
				if _, ok := annotationList[serviceAnnotationLoadBalancerSSLPorts]; ok {
					return fmt.Errorf("SSL Ports are not supported with the %s protocol", valueEnum)
				}
				if _, ok := annotationList[serviceAnnotationLoadBalancerSslDomains]; ok {
					return fmt.Errorf("SSL Domains are not supported with the %s protocol", valueEnum)
				}
			}
		case serviceAnnotationLoadBalancerListenerProxyProtocol:
			if _, err := proxyprotocol.ParseEnum(value); err != nil {
				return fmt.Errorf("Invalid Load Balancer Listener Proxy Protocol %q: %w", value, err)
			}
		case serviceAnnotationLoadBalancerSSLPorts:
			if _, ok := annotationList[serviceAnnotationLoadBalancerSslDomains]; !ok {
				return fmt.Errorf("SSL needs a list of domains to certify. Add the %q annotation", serviceAnnotationLoadBalancerSslDomains)
			}
		case serviceAnnotationLoadBalancerHCProtocol:
			if _, err := healthchecktype.ParseEnum(value); err != nil {
				return fmt.Errorf("Invalid Load Balancer Healthcheck Protocol %q: %w", value, err)
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
		case serviceAnnotationLoadBalancerHCRequest:
			testURL := "http://example.com:6443" + value
			u, err := url.Parse(testURL)
			if err != nil || u.Path != value {
				return fmt.Errorf("%q needs to be a valid Url request path", annotation)
			}
		case serviceAnnotationLoadBalancerCloudipAllocations:
			if !cloudIPPattern.MatchString(value) {
				return fmt.Errorf("%q needs to match the pattern %q", annotation, cloudIPPattern)
			}
		}
	}
	return nil
}
