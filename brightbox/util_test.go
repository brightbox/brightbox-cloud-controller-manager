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
	"testing"
)

func TestMapProviderIDToServerID(t *testing.T) {
	testCases := map[string]struct {
		providerID string
		expected   string
	}{
		"no cloud prefix": {
			providerID: "srv-testy",
			expected:   "srv-testy",
		},
		"cloud prefix": {
			providerID: providerPrefix + "srv-testy",
			expected:   "srv-testy",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := mapProviderIDToServerID(tc.providerID)
			if result != tc.expected {
				t.Errorf("Expected server id %q, but got %q", tc.expected, result)
			}
		})
	}
}

func TestMapZoneHandleToRegion(t *testing.T) {
	testCases := map[string]struct {
		zoneHandle string
		expected   string
	}{
		"a zone": {
			zoneHandle: "gb1s-a",
			expected:   "gb1s",
		},
		"b zone": {
			zoneHandle: "gb1-b",
			expected:   "gb1",
		},
		"dodgy zone": {
			zoneHandle: "fred",
			expected:   "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := mapZoneHandleToRegion(tc.zoneHandle)
			if err != nil {
				if tc.expected != "" {
					t.Errorf(err.Error())
				}
			} else if result != tc.expected {
				t.Errorf("Expected server id %q, but got %q", tc.expected, result)
			}
		})
	}
}
