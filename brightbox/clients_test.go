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

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/brightbox/gobrightbox"
)

func TestGetMetadataClient(t *testing.T) {
	client := &cloud{}
	mdc, err := client.metadataClient()
	if err != nil {
		t.Errorf("Failed to get metadata client: %s", err.Error())
	}
	switch mdc.(type) {
	case (*ec2metadata.EC2Metadata):
	default:
		t.Errorf("Returned incorrect metadata client")
	}
}

func TestGetCloudClient(t *testing.T) {
	client := &cloud{}
	cc, err := client.cloudClient()
	if err != nil {
		t.Errorf("Failed to get cloud client: %s", err.Error())
	}
	switch cc.(type) {
	case (*brightbox.Client):
	default:
		t.Errorf("Returned incorrect cloud client")
	}
}
