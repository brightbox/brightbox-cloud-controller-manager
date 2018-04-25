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
