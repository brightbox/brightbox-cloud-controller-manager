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
	"errors"
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

	lbActive   = "active"
	lbCreating = "creating"
	cipMapped  = "mapped"
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
	//retrieves a detailed view of one cloud ip
	CloudIP(identifier string) (*brightbox.CloudIP, error)
	//Issues a request to map the cloud ip to the destination.
	MapCloudIP(identifier string, destination string) error
	//Creates a new Cloud IP
	CreateCloudIP(newCloudIP *brightbox.CloudIPOptions) (*brightbox.CloudIP, error)
	//adds servers to an existing server group
	AddServersToServerGroup(identifier string, serverIds []string) (*brightbox.ServerGroup, error)
	//removes servers from an existing server group
	RemoveServersFromServerGroup(identifier string, serverIds []string) (*brightbox.ServerGroup, error)
	// ServerGroups retrieves a list of all server groups
	ServerGroups() ([]brightbox.ServerGroup, error)
	//creates a new server group
	CreateServerGroup(newServerGroup *brightbox.ServerGroupOptions) (*brightbox.ServerGroup, error)
	//creates a new firewall policy
	CreateFirewallPolicy(policyOptions *brightbox.FirewallPolicyOptions) (*brightbox.FirewallPolicy, error)
	//creates a new firewall rule
	CreateFirewallRule(ruleOptions *brightbox.FirewallRuleOptions) (*brightbox.FirewallRule, error)
	//updates an existing firewall rule
	UpdateFirewallRule(ruleOptions *brightbox.FirewallRuleOptions) (*brightbox.FirewallRule, error)

	//retrieves a list of all firewall policies
	FirewallPolicies() ([]brightbox.FirewallPolicy, error)
	// DestroyServerGroup destroys an existing server group
	DestroyServerGroup(identifier string) error
	// DestroyFirewallPolicy issues a request to destroy the firewall policy
	DestroyFirewallPolicy(identifier string) error
	// DestroyLoadBalancer issues a request to destroy the load balancer
	DestroyLoadBalancer(identifier string) error
	// DestroyCloudIP issues a request to destroy the cloud ip
	DestroyCloudIP(identifier string) error
}

func (c *cloud) getServer(ctx context.Context, id string) (*brightbox.Server, error) {
	glog.V(4).Infof("getServer called for '%q'", id)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	srv, err := client.Server(id)
	if err != nil {
		if isNotFound(err) {
			return nil, cloudprovider.InstanceNotFound
		}
		return nil, err
	}
	return srv, nil
}

func isNotFound(e error) bool {
	switch v := e.(type) {
	case brightbox.ApiError:
		return v.StatusCode == http.StatusNotFound
	default:
		return false
	}
}

func isAlive(lb *brightbox.LoadBalancer) bool {
	return lb.Status == lbActive || lb.Status == lbCreating
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
		if isAlive(&lbList[i]) && lbName == lbList[i].Name {
			result = &lbList[i]
			break
		}
	}
	return result, nil
}

func (c *cloud) createLoadBalancer(newLB *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	glog.V(4).Infof("createLoadBalancer called for %q", *newLB.Name)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	return client.CreateLoadBalancer(newLB)
}

func (c *cloud) updateLoadBalancer(newLB *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	glog.V(4).Infof("updateLoadBalancer called for (%q, %q)", newLB.Id, *newLB.Name)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	return client.UpdateLoadBalancer(newLB)
}

// get a FirewallPolicy By Name
func (c *cloud) getFirewallPolicyByName(fpName string) (*brightbox.FirewallPolicy, error) {
	glog.V(4).Infof("getFirewallPolicyByName called for %q", fpName)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	fpList, err := client.FirewallPolicies()
	if err != nil {
		return nil, err
	}
	var result *brightbox.FirewallPolicy
	for i := range fpList {
		if fpName == fpList[i].Name {
			result = &fpList[i]
			break
		}
	}
	return result, nil
}

// get a serverGroup By Name
func (c *cloud) getServerGroupByName(sgName string) (*brightbox.ServerGroup, error) {
	glog.V(4).Infof("getServerGroupByName called for %q", sgName)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	sgList, err := client.ServerGroups()
	if err != nil {
		return nil, err
	}
	var result *brightbox.ServerGroup
	for i := range sgList {
		if sgName == sgList[i].Name {
			result = &sgList[i]
			break
		}
	}
	return result, nil
}

func (c *cloud) createServerGroup(name string) (*brightbox.ServerGroup, error) {
	glog.V(4).Infof("createServerGroup called for %q", name)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	return client.CreateServerGroup(&brightbox.ServerGroupOptions{Name: &name})
}

//Firewall Policy
func (c *cloud) createFirewallPolicy(group *brightbox.ServerGroup) (*brightbox.FirewallPolicy, error) {
	glog.V(4).Infof("createFirewallPolicy called for %q", group.Name)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	return client.CreateFirewallPolicy(&brightbox.FirewallPolicyOptions{Name: &group.Name, ServerGroup: &group.Id})
}

//Firewall Rules
func (c *cloud) createFirewallRule(newFR *brightbox.FirewallRuleOptions) (*brightbox.FirewallRule, error) {
	glog.V(4).Infof("createFirewallRule called for %q", *newFR.Description)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	return client.CreateFirewallRule(newFR)
}

func (c *cloud) updateFirewallRule(newFR *brightbox.FirewallRuleOptions) (*brightbox.FirewallRule, error) {
	glog.V(4).Infof("updateFirewallRule called for (%q, %q)", newFR.Id, *newFR.Description)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	return client.UpdateFirewallRule(newFR)
}

func (c *cloud) retryMapCip(cipId, lbId string) error {
	glog.V(4).Infof("retryMapCip (%q, %q)", cipId, lbId)
	client, err := c.cloudClient()
	if err != nil {
		return err
	}
	retryFunc := backoff.ExecuteFunc(func(_ context.Context) error {
		glog.V(4).Infof("attempting to map CloudIP %q", cipId)
		err := client.MapCloudIP(cipId, lbId)
		//Make sure we return a normal error
		if err != nil {
			return errors.New(err.Error())
		}
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeoutSeconds*time.Second)
	defer cancel()
	p := backoff.NewExponential()
	return backoff.Retry(ctx, p, retryFunc)
}

func (c *cloud) retryWaitForLbMap(cipId, lbId string) error {
	glog.V(4).Infof("retryWaitForLbMap (%q, %q)", cipId, lbId)
	client, err := c.cloudClient()
	if err != nil {
		return err
	}
	waitForMapped := backoff.ExecuteFunc(func(_ context.Context) error {
		glog.V(4).Infof("Waiting for mapping to occur (%q, %q)", cipId, lbId)
		cip, err := client.CloudIP(cipId)
		//Make sure we return a normal error
		if err != nil {
			return errors.New(err.Error())
		}
		if alreadyMapped(cip, lbId) {
			return nil
		}
		return errors.New("Not mapped")
	})
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeoutSeconds*time.Second)
	defer cancel()
	p := backoff.NewExponential()
	return backoff.Retry(ctx, p, waitForMapped)
}

// backoff retry mapping the cloudip to a load balancer
func (c *cloud) ensureMappedCip(lb *brightbox.LoadBalancer, cip *brightbox.CloudIP) error {
	glog.V(4).Infof("ensureMappedCip called for (%q, %q)", lb.Id, cip.Id)
	if alreadyMapped(cip, lb.Id) {
		return nil
	} else if cip.Status == cipMapped {
		return fmt.Errorf("Unexplained mapping of %q (%q)", cip.Id, cip.PublicIP)
	}
	err := c.retryMapCip(cip.Id, lb.Id)
	if err != nil {
		return err
	}
	return c.retryWaitForLbMap(cip.Id, lb.Id)
}

func alreadyMapped(cip *brightbox.CloudIP, lbId string) bool {
	return cip.LoadBalancer != nil && cip.LoadBalancer.Id == lbId
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

//Get a cloudIp by id
func (c *cloud) getCloudIP(identifier string) (*brightbox.CloudIP, error) {
	glog.V(4).Infof("getCloudIP (%q)", identifier)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	return client.CloudIP(identifier)
}

//Destroy things

func (c *cloud) destroyLoadBalancer(id string) error {
	glog.V(4).Infof("destroyLoadBalancer %q", id)
	client, err := c.cloudClient()
	if err != nil {
		return err
	}
	return client.DestroyLoadBalancer(id)
}

func (c *cloud) destroyServerGroup(id string) error {
	glog.V(4).Infof("destroyServerGroup %q", id)
	client, err := c.cloudClient()
	if err != nil {
		return err
	}
	return client.DestroyServerGroup(id)
}

func (c *cloud) destroyFirewallPolicy(id string) error {
	glog.V(4).Infof("destroyFirewallPolicy %q", id)
	client, err := c.cloudClient()
	if err != nil {
		return err
	}
	return client.DestroyFirewallPolicy(id)
}

// backoff retry removing the cloudip
func (c *cloud) retryDestroyCloudIP(id string) error {
	glog.V(4).Infof("retryDestroyCloudIP called")
	client, err := c.cloudClient()
	if err != nil {
		return err
	}
	retryFunc := backoff.ExecuteFunc(func(_ context.Context) error {
		glog.V(4).Infof("attempting to remove CloudIP %q", id)
		err := client.DestroyCloudIP(id)
		//Ensure we return a normal error to stop backoff complaining
		if err != nil {
			return errors.New(err.Error())
		}
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeoutSeconds*time.Second)
	defer cancel()
	p := backoff.NewExponential()
	return backoff.Retry(ctx, p, retryFunc)
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

func mapServersToServerIds(servers []brightbox.Server) []string {
	result := make([]string, len(servers))
	for i := range servers {
		result[i] = servers[i].Id
	}
	return result
}

func (c *cloud) syncServerGroup(group *brightbox.ServerGroup, newIds []string) (*brightbox.ServerGroup, error) {
	oldIds := mapServersToServerIds(group.Servers)
	glog.V(4).Infof("syncServerGroup called for (%v, %v, %v)", group.Id, oldIds, newIds)
	client, err := c.cloudClient()
	if err != nil {
		return nil, err
	}
	insIds, delIds := getSyncLists(oldIds, newIds)
	result := group
	if len(insIds) > 0 {
		glog.V(4).Infof("Adding Servers %v", insIds)
		result, err = client.AddServersToServerGroup(group.Id, insIds)
	}
	if err == nil && len(delIds) > 0 {
		glog.V(4).Infof("Removing Servers %v", delIds)
		result, err = client.RemoveServersFromServerGroup(group.Id, delIds)
	}
	return result, err
}
