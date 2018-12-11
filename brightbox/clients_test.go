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
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestInvalidCloudClient(t *testing.T) {
	resetAuthEnvironment()
	defer resetAuthEnvironment()
	client := &cloud{}
	_, err := client.cloudClient()
	if err == nil {
		t.Errorf("Expected account error")
	}
	setAuthEnvClientID()
	setAuthEnvUsername()
	_, err = client.cloudClient()
	if err == nil {
		t.Errorf("Expected User Credentials error")
	}
	setAuthEnvPassword()
	_, err = client.cloudClient()
	if err == nil {
		t.Errorf("Expected User Credentials error")
	}
	clearAuthEnvUsername()
	_, err = client.cloudClient()
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
	ts := getAuthEnvTokenHandler(t)
	defer resetAuthEnvironment()
	defer ts.Close()
	client := &cloud{}
	setAuthEnvUsername()
	_, err := client.cloudClient()
	if err == nil {
		t.Errorf("Expected User Credentials error")
	}
}

func TestUsernameValidation(t *testing.T) {
	ts := getAuthEnvTokenHandler(t)
	defer resetAuthEnvironment()
	defer ts.Close()
	setAuthEnvUsername()
	setAuthEnvPassword()
	client := &cloud{}
	_, err := client.cloudClient()
	if err != nil {
		t.Errorf("Expected User Credentials validation, got %s", err.Error())
	}
}

func resetAuthEnvironment() {
	vars := []string{
		clientEnvVar,
		clientSecretEnvVar,
		usernameEnvVar,
		passwordEnvVar,
		accountEnvVar,
		apiUrlEnvVar,
	}
	for _, envvar := range vars {
		os.Unsetenv(envvar)
	}
}

func setAuthEnvClientID() {
	os.Setenv(clientSecretEnvVar, "not default")
}

func setAuthEnvUsername() {
	os.Setenv(usernameEnvVar, "itsy@bitzy.com")
}

func setAuthEnvPassword() {
	os.Setenv(passwordEnvVar, "madeuppassword")
}

func setAuthEnvAPIURL(value string) {
	os.Setenv(apiUrlEnvVar, value)
}

func setAuthEnvAccount() {
	os.Setenv(accountEnvVar, "acc-testy")
}

func clearAuthEnvUsername() {
	os.Unsetenv(usernameEnvVar)
}

func getAuthEnvTokenHandler(t *testing.T) *httptest.Server {
	resetAuthEnvironment()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		expected := "/token"
		if r.URL.String() != expected {
			t.Errorf("URL = %q; want %q", r.URL, expected)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed reading request body: %s.", err)
		}
		headerContentType := r.Header.Get("Content-Type")
		expected = "application/x-www-form-urlencoded"
		if headerContentType != expected {
			t.Errorf("Content-Type header = %q; want %q", headerContentType, expected)
		}
		headerAuth := r.Header.Get("Authorization")
		expected = "Basic YXBwLWRrbWNoOnVvZ29lbHpndDBud2F3Yg=="
		if headerAuth != expected {
			t.Errorf("Authorization header = %q; want %q", headerAuth, expected)
		}
		switch string(body) {
		case "grant_type=password&password=madeuppassword&scope=infrastructure&username=itsy%40bitzy.com":
			w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
			w.Write([]byte("access_token=90d64460d14870c08c81352a05dedd3465940a7c&scope=user&token_type=bearer"))
		case "grant_type=password&password=&scope=infrastructure&username=itsy%40bitzy.com":
			w.WriteHeader(http.StatusUnauthorized)
		default:
			t.Errorf("Unexpected res.Body = %q", string(body))
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	setAuthEnvAPIURL(ts.URL)
	setAuthEnvAccount()
	return ts
}
