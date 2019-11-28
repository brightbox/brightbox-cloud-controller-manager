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

package k8ssdk

import (
	"testing"

	"github.com/go-test/deep"
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
			providerID: ProviderPrefix + "srv-testy",
			expected:   "srv-testy",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := MapProviderIDToServerID(tc.providerID)
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
			result, err := MapZoneHandleToRegion(tc.zoneHandle)
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

func TestGetSyncLists(t *testing.T) {
	testCases := map[string]struct {
		old []string
		new []string
		ins []string
		del []string
	}{
		"all nil": {
			old: nil,
			new: nil,
			ins: nil,
			del: nil,
		},
		"no change": {
			old: []string{"srv-12345", "srv-67890"},
			new: []string{"srv-67890", "srv-12345"},
			ins: nil,
			del: nil,
		},
		"one insert": {
			old: []string{"srv-12345", "srv-67890"},
			new: []string{"srv-67890", "srv-12345", "srv-testy"},
			ins: []string{"srv-testy"},
			del: nil,
		},
		"one delete": {
			old: []string{"srv-12345", "srv-testy", "srv-67890"},
			new: []string{"srv-67890", "srv-12345"},
			ins: nil,
			del: []string{"srv-testy"},
		},
		"nil to something": {
			old: nil,
			new: []string{"srv-67890", "srv-12345"},
			ins: []string{"srv-12345", "srv-67890"},
			del: nil,
		},
		"something to nil": {
			old: []string{"srv-67890", "srv-12345"},
			new: nil,
			ins: nil,
			del: []string{"srv-12345", "srv-67890"},
		},
		"change": {
			old: []string{"srv-12345", "srv-testy", "srv-67890"},
			new: []string{"srv-67890", "srv-fasty", "srv-newly"},
			ins: []string{"srv-fasty", "srv-newly"},
			del: []string{"srv-12345", "srv-testy"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ins, del := getSyncLists(tc.old, tc.new)
			if diff := deep.Equal(ins, tc.ins); diff != nil {
				t.Error(diff)
			}
			if diff := deep.Equal(del, tc.del); diff != nil {
				t.Error(diff)
			}
		})
	}
}
