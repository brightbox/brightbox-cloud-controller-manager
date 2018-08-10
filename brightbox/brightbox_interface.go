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
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/brightbox/gobrightbox"
	"github.com/golang/glog"
	"github.com/lestrrat-go/backoff"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

const (
	defaultClientID     = "app-dkmch"
	defaultClientSecret = "uogoelzgt0nwawb"
	clientEnvVar        = "BRIGHTBOX_CLIENT"
	clientSecretEnvVar  = "BRIGHTBOX_CLIENT_SECRET"
	usernameEnvVar      = "BRIGHTBOX_USER_NAME"
	passwordEnvVar      = "BRIGHTBOX_PASSWORD"
	accountEnvVar       = "BRIGHTBOX_ACCOUNT"
	apiUrlEnvVar        = "BRIGHTBOX_API_URL"

	defaultTimeoutSeconds = 10
)

var infrastructureScope = []string{"infrastructure"}

type authdetails struct {
	APIClient string
	APISecret string
	UserName  string
	password  string
	Account   string
	APIURL    string
}

// CloudAccess is an abstraction over the Brightbox API to allow testing
type CloudAccess interface {
	//Fetch a server
	Server(identifier string) (*brightbox.Server, error)
	//Fetch a list of LoadBalancers
	LoadBalancers() ([]brightbox.LoadBalancer, error)
	//Creates a new load balancer
	CreateLoadBalancer(newLB *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error)
	//Updates an existing load balancer
	UpdateLoadBalancer(newLB *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error)
	//Retrieves a list of all cloud IPs
	CloudIPs() ([]brightbox.CloudIP, error)
	//Issues a request to map the cloud ip to the destination.
	MapCloudIP(identifier string, destination string) error
	//Creates a new Cloud IP
	CreateCloudIP(newCloudIP *brightbox.CloudIPOptions) (*brightbox.CloudIP, error)
}

func (c *cloud) getServer(ctx context.Context, id string) (*brightbox.Server, error) {
	glog.V(4).Infof("getServer called for '%q'", id)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	srv, err := client.Server(id)
	if err != nil {
		if isNotFound(err.(brightbox.ApiError)) {
			return nil, cloudprovider.InstanceNotFound
		}
		return nil, err
	}
	return srv, nil
}

func isNotFound(e brightbox.ApiError) bool {
	return e.StatusCode == http.StatusNotFound
}

func isActive(lb *brightbox.LoadBalancer) bool {
	return lb.Status == "Active"
}

// get a loadbalancer by name
func (c *cloud) getLoadBalancerByName(lbName string) (*brightbox.LoadBalancer, error) {
	glog.V(4).Infof("getLoadBalancerByName called for %q", lbName)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	lbList, err := client.LoadBalancers()
	if err != nil {
		return nil, err
	}
	var result *brightbox.LoadBalancer
	for i := range lbList {
		if isActive(&lbList[i]) && lbName == lbList[i].Name {
			result = &lbList[i]
			break
		}
	}
	return result, nil
}

func (c *cloud) createLoadBalancer(newLB *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	glog.V(4).Infof("createLoadBalancer called for %q", newLB.Name)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	return client.CreateLoadBalancer(newLB)
}

func (c *cloud) updateLoadBalancer(newLB *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	glog.V(4).Infof("updateLoadBalancer called for (%q, %q)", newLB.Id, newLB.Name)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	return client.UpdateLoadBalancer(newLB)
}

// backoff retry mapping the cloudip to a load balancer
func (c *cloud) ensureMappedCip(lb *brightbox.LoadBalancer, cip *brightbox.CloudIP) error {
	if alreadyMapped(lb, cip) {
		return nil
	}
	glog.V(4).Infof("ensureMappedCip called for (%q, %q)", lb.Id, cip.Id)
	client, err := c.cloudClient()
	if err != nil {
		return err
	}
	retryFunc := backoff.ExecuteFunc(func(_ context.Context) error {
		glog.V(4).Infof("attempting to map CloudIP %q", cip.Id)
		return client.MapCloudIP(cip.Id, lb.Id)
	})
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeoutSeconds*time.Second)
	defer cancel()
	p := backoff.NewExponential()
	return backoff.Retry(ctx, p, retryFunc)
}

func alreadyMapped(lb *brightbox.LoadBalancer, cip *brightbox.CloudIP) bool {
	return cip.LoadBalancer != nil && cip.LoadBalancer.Id == lb.Id
}

func (c *cloud) allocateCip(name string) (*brightbox.CloudIP, error) {
	glog.V(4).Infof("allocateCip %q", name)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	opts := &brightbox.CloudIPOptions{
		Name: &name,
	}
	return client.CreateCloudIP(opts)
}

func (c *cloud) getCloudIPs() ([]brightbox.CloudIP, error) {
	glog.V(4).Infof("getCloudIPs")
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	return client.CloudIPs()
}

// Obtain a Brightbox cloud client anew
func obtainCloudClient() (*brightbox.Client, error) {
	glog.V(4).Infof("obtainCloudClient called ")
	config := &authdetails{
		APIClient: getenvWithDefault(clientEnvVar,
			defaultClientID),
		APISecret: getenvWithDefault(clientSecretEnvVar,
			defaultClientSecret),
		UserName: os.Getenv(usernameEnvVar),
		password: os.Getenv(passwordEnvVar),
		Account:  os.Getenv(accountEnvVar),
		APIURL:   os.Getenv(apiUrlEnvVar),
	}
	err := config.validateConfig()
	if err != nil {
		return nil, err
	}
	return config.authenticatedClient()
}

// Validate account config entries
func (authd *authdetails) validateConfig() error {
	glog.V(4).Infof("validateConfig called ")
	if authd.APIClient == defaultClientID &&
		authd.APISecret == defaultClientSecret {
		if authd.Account == "" {
			return fmt.Errorf("Must specify Account with User Credentials")
		}
	} else {
		if authd.UserName != "" || authd.password != "" {
			return fmt.Errorf("User Credentials not used with API Client.")
		}
	}
	return nil
}

// Authenticate the details and return a client
func (authd *authdetails) authenticatedClient() (*brightbox.Client, error) {
	ctx := context.Background()
	switch {
	case authd.UserName != "" || authd.password != "":
		return authd.tokenisedAuth(ctx)
	default:
		return authd.apiClientAuth(ctx)
	}
}

func (authd *authdetails) tokenURL() string {
	return authd.APIURL + "/token"
}

func (authd *authdetails) tokenisedAuth(ctx context.Context) (*brightbox.Client, error) {
	conf := oauth2.Config{
		ClientID:     authd.APIClient,
		ClientSecret: authd.APISecret,
		Scopes:       infrastructureScope,
		Endpoint: oauth2.Endpoint{
			TokenURL: authd.tokenURL(),
		},
	}
	glog.V(4).Infof("Obtaining authentication for user %s", authd.UserName)
	glog.V(4).Infof("Speaking to %s", authd.tokenURL())
	token, err := conf.PasswordCredentialsToken(ctx, authd.UserName, authd.password)
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Refreshing current token as required")
	oauthConnection := conf.Client(ctx, token)
	return brightbox.NewClient(authd.APIURL, authd.Account, oauthConnection)
}

func (authd *authdetails) apiClientAuth(ctx context.Context) (*brightbox.Client, error) {
	conf := clientcredentials.Config{
		ClientID:     authd.APIClient,
		ClientSecret: authd.APISecret,
		Scopes:       infrastructureScope,
		TokenURL:     authd.tokenURL(),
	}
	glog.V(4).Infof("Obtaining API client authorisation for client %s", authd.APIClient)
	glog.V(4).Infof("Speaking to %s", authd.tokenURL())
	oauthConnection := conf.Client(ctx)
	return brightbox.NewClient(authd.APIURL, authd.Account, oauthConnection)
}
