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
	"flag"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"k8s.io/klog"
)

func init() {
	klog.InitFlags(nil)
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "4")
	flag.Parse()
}

func TestGetMetadataClient(t *testing.T) {
	client := MakeTestClient(nil, nil)
	mdc, err := client.MetadataClient()
	if err != nil {
		t.Errorf("Failed to get metadata client: %s", err.Error())
	}
	switch mdc.(type) {
	case (*ec2metadata.Client):
	default:
		t.Errorf("Returned incorrect metadata client")
	}
}

func TestInvalidCloudClient(t *testing.T) {
	ResetAuthEnvironment()
	defer ResetAuthEnvironment()
	client := MakeTestClient(nil, nil)
	_, err := client.CloudClient()
	if err == nil {
		t.Errorf("Expected account error")
	}
	SetAuthEnvClientID()
	SetAuthEnvUsername()
	_, err = client.CloudClient()
	if err == nil {
		t.Errorf("Expected User Credentials error")
	}
	SetAuthEnvPassword()
	_, err = client.CloudClient()
	if err == nil {
		t.Errorf("Expected User Credentials error")
	}
	ClearAuthEnvUsername()
	_, err = client.CloudClient()
	if err == nil {
		t.Errorf("Expected User Credentials error")
	}
	//	switch cc.(type) {
	//	case (*brightbox.Client):
	//	default:
	//		t.Errorf("Returned incorrect cloud client")
	//	}
}

func TestBadPasswordCloudClient(t *testing.T) {
	ts := GetAuthEnvTokenHandler(t)
	defer ResetAuthEnvironment()
	defer ts.Close()
	client := MakeTestClient(nil, nil)

	SetAuthEnvUsername()
	_, err := client.CloudClient()
	if err == nil {
		t.Errorf("Expected User Credentials error")
	}
}

func TestUsernameValidation(t *testing.T) {
	ts := GetAuthEnvTokenHandler(t)
	defer ResetAuthEnvironment()
	defer ts.Close()
	SetAuthEnvUsername()
	SetAuthEnvPassword()
	client := MakeTestClient(nil, nil)
	_, err := client.CloudClient()
	if err != nil {
		t.Errorf("Expected User Credentials validation, got %s", err.Error())
	}
}
