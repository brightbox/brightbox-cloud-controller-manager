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
	"os"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/types"
)

const (
	providerName   = "brightbox"
	providerPrefix = providerName + "://"
)

// Parse the provider id string and return a string that should be a server id
// Should be no need for  error checking here, since the input string
// is constrained in format by the k8s process
func mapProviderIDToServerID(providerID string) string {
	if strings.HasPrefix(providerID, providerPrefix) {
		return strings.TrimPrefix(providerID, providerPrefix)
	}
	return providerID
}

// Parse the zone handle and return the embedded region id
// Zone names are of the form: ${region-name}-${ix}
// So we look for the last '-' and trim just before that
func mapZoneHandleToRegion(zoneHandle string) (string, error) {
	ix := strings.LastIndex(zoneHandle, "-")
	if ix == -1 {
		return "", fmt.Errorf("unexpected zone: %s", zoneHandle)
	}
	return zoneHandle[:ix], nil
}

// mapNodeNameToServerID maps a k8s NodeName to a Brightbox Server ID
// This is a simple string cast.
func mapNodeNameToServerID(nodeName types.NodeName) string {
	return string(nodeName)
}

// mapServerIDToNodeName maps a Brightbox Server ID to a nodename
// Again a simpl string cast
func mapServerIDToNodeName(name string) types.NodeName {
	return types.NodeName(name)
}

func mapProviderIDToNodeName(providerID string) types.NodeName {
	return mapServerIDToNodeName(mapProviderIDToServerID(providerID))
}

// getEnvVarWithDefault retrieves the value of the environment variable
// named by the key. If the variable is not present, return the default
//value instead.
func getenvWithDefault(key string, defaultValue string) string {
	if val, exists := os.LookupEnv(key); !exists {
		return defaultValue
	} else {
		return val
	}
}

//get a list of inserts and deletes that changes oldList into newList
func getSyncLists(oldList []string, newList []string) ([]string, []string) {
	sort.Strings(oldList)
	sort.Strings(newList)
	var x, y int
	var insList, delList []string
	for x < len(oldList) || y < len(newList) {
		switch {
		case y >= len(newList):
			delList = append(delList, oldList[x])
			x += 1
		case x >= len(oldList):
			insList = append(insList, newList[y])
			y += 1
		case oldList[x] < newList[y]:
			delList = append(delList, oldList[x])
			x += 1
		case oldList[x] > newList[y]:
			insList = append(insList, newList[y])
			y += 1
		default:
			y += 1
			x += 1
		}
	}
	return insList, delList
}

//Add the nasty hack to the load balancer name to trigger speed
func grokLoadBalancerName(name string) string {
	return name + " #type:container"
}
