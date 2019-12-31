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
	"github.com/brightbox/k8ssdk"
	"k8s.io/apimachinery/pkg/types"
)

// mapNodeNameToServerID maps a k8s NodeName to a Brightbox Server ID
// This is a simple string cast.
func MapNodeNameToServerID(nodeName types.NodeName) string {
	return string(nodeName)
}

// mapServerIDToNodeName maps a Brightbox Server ID to a nodename
// Again a simple string cast
func MapServerIDToNodeName(name string) types.NodeName {
	return types.NodeName(name)
}

func MapProviderIDToNodeName(providerID string) types.NodeName {
	return MapServerIDToNodeName(k8ssdk.MapProviderIDToServerID(providerID))
}

func MapNodeNameToProviderID(nodeName types.NodeName) string {
	return k8ssdk.MapServerIDToProviderID(MapNodeNameToServerID(nodeName))
}
