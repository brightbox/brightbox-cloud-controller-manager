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
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/brightbox/gobrightbox"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"k8s.io/klog"
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

	LbActive              = "active"
	LbCreating            = "creating"
	cipMapped             = "mapped"
	ValidAcmeDomainStatus = "valid"
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

func (c *Cloud) GetServer(ctx context.Context, id string, notFoundError error) (*brightbox.Server, error) {
	klog.V(4).Infof("getServer (%q)", id)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	srv, err := client.Server(id)
	if err != nil {
		if isNotFound(err) {
			return nil, notFoundError
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

func (c *Cloud) CreateServer(newDetails *brightbox.ServerOptions) (*brightbox.Server, error) {
	klog.V(4).Infof("CreateServer (%q)", *newDetails.Name)
	klog.V(6).Infof("%+v", newDetails)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.CreateServer(newDetails)
}

func isAlive(lb *brightbox.LoadBalancer) bool {
	return lb.Status == LbActive || lb.Status == LbCreating
}

func (c *Cloud) GetLoadBalancerByName(name string) (*brightbox.LoadBalancer, error) {
	klog.V(4).Infof("GetLoadBalancerByName (%q)", name)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	lbList, err := client.LoadBalancers()
	if err != nil {
		return nil, err
	}
	for i := range lbList {
		if isAlive(&lbList[i]) && name == lbList[i].Name {
			return &lbList[i], nil
		}
	}
	return nil, nil
}

func (c *Cloud) GetLoadBalancerById(id string) (*brightbox.LoadBalancer, error) {
	klog.V(4).Infof("GetLoadBalancerById (%q)", id)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.LoadBalancer(id)
}

func (c *Cloud) CreateLoadBalancer(newDetails *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	klog.V(4).Infof("CreateLoadBalancer (%q)", *newDetails.Name)
	klog.V(6).Infof("%+v", newDetails)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.CreateLoadBalancer(newDetails)
}

func (c *Cloud) UpdateLoadBalancer(newDetails *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error) {
	klog.V(4).Infof("UpdateLoadBalancer (%q, %q)", newDetails.Id, *newDetails.Name)
	klog.V(6).Infof("%+v", newDetails)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.UpdateLoadBalancer(newDetails)
}

// get a FirewallPolicy By Name
func (c *Cloud) GetFirewallPolicyByName(name string) (*brightbox.FirewallPolicy, error) {
	klog.V(4).Infof("getFirewallPolicyByName (%q)", name)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	firewallPolicyList, err := client.FirewallPolicies()
	if err != nil {
		return nil, err
	}
	var result *brightbox.FirewallPolicy
	for i := range firewallPolicyList {
		if name == firewallPolicyList[i].Name {
			result = &firewallPolicyList[i]
			break
		}
	}
	return result, nil
}

func (c *Cloud) GetServerGroups() ([]brightbox.ServerGroup, error) {
	klog.V(4).Info("GetServerGroups")
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.ServerGroups()
}

func (c *Cloud) GetServerGroup(identifier string) (*brightbox.ServerGroup, error) {
	klog.V(4).Infof("GetServerGroup %q", identifier)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.ServerGroup(identifier)
}

// get a serverGroup By Name
func (c *Cloud) GetServerGroupByName(name string) (*brightbox.ServerGroup, error) {
	klog.V(4).Infof("GetServerGroupByName (%q)", name)
	serverGroupList, err := c.GetServerGroups()
	if err != nil {
		return nil, err
	}
	var result *brightbox.ServerGroup
	for i := range serverGroupList {
		if name == serverGroupList[i].Name {
			result = &serverGroupList[i]
			break
		}
	}
	return result, nil
}

func (c *Cloud) CreateServerGroup(name string) (*brightbox.ServerGroup, error) {
	klog.V(4).Infof("CreateServerGroup (%q)", name)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.CreateServerGroup(&brightbox.ServerGroupOptions{Name: &name})
}

//Firewall Policy
func (c *Cloud) CreateFirewallPolicy(group *brightbox.ServerGroup) (*brightbox.FirewallPolicy, error) {
	klog.V(4).Infof("createFirewallPolicy (%q)", group.Name)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.CreateFirewallPolicy(&brightbox.FirewallPolicyOptions{Name: &group.Name, ServerGroup: &group.Id})
}

//Firewall Rules
func (c *Cloud) CreateFirewallRule(newDetails *brightbox.FirewallRuleOptions) (*brightbox.FirewallRule, error) {
	klog.V(4).Infof("createFirewallRule (%q)", *newDetails.Description)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.CreateFirewallRule(newDetails)
}

func (c *Cloud) UpdateFirewallRule(newDetails *brightbox.FirewallRuleOptions) (*brightbox.FirewallRule, error) {
	klog.V(4).Infof("updateFirewallRule (%q, %q)", newDetails.Id, *newDetails.Description)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.UpdateFirewallRule(newDetails)
}

func (c *Cloud) EnsureMappedCloudIP(lb *brightbox.LoadBalancer, cip *brightbox.CloudIP) error {
	klog.V(4).Infof("EnsureMappedCloudIP (%q, %q)", lb.Id, cip.Id)
	if alreadyMapped(cip, lb.Id) {
		return nil
	} else if cip.Status == cipMapped {
		return fmt.Errorf("Unexplained mapping of %q (%q)", cip.Id, cip.PublicIP)
	}
	client, err := c.CloudClient()
	if err != nil {
		return err
	}
	return client.MapCloudIP(cip.Id, lb.Id)
}

func alreadyMapped(cip *brightbox.CloudIP, loadBalancerId string) bool {
	return cip.LoadBalancer != nil && cip.LoadBalancer.Id == loadBalancerId
}

func (c *Cloud) AllocateCloudIP(name string) (*brightbox.CloudIP, error) {
	klog.V(4).Infof("AllocateCloudIP %q", name)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	opts := &brightbox.CloudIPOptions{
		Name: &name,
	}
	return client.CreateCloudIP(opts)
}

func (c *Cloud) GetCloudIPs() ([]brightbox.CloudIP, error) {
	klog.V(4).Infof("GetCloudIPs")
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.CloudIPs()
}

//Get a cloudIp by id
func (c *Cloud) getCloudIP(id string) (*brightbox.CloudIP, error) {
	klog.V(4).Infof("getCloudIP (%q)", id)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	return client.CloudIP(id)
}

//Destroy things

func (c *Cloud) DestroyLoadBalancer(id string) error {
	klog.V(4).Infof("DestroyLoadBalancer %q", id)
	client, err := c.CloudClient()
	if err != nil {
		return err
	}
	return client.DestroyLoadBalancer(id)
}

func (c *Cloud) DestroyServer(id string) error {
	klog.V(4).Infof("DestroyServer %q", id)
	client, err := c.CloudClient()
	if err != nil {
		return err
	}
	return client.DestroyServer(id)
}

func (c *Cloud) DestroyServerGroup(id string) error {
	klog.V(4).Infof("DestroyServerGroup %q", id)
	client, err := c.CloudClient()
	if err != nil {
		return err
	}
	return client.DestroyServerGroup(id)
}

func (c *Cloud) DestroyFirewallPolicy(id string) error {
	klog.V(4).Infof("DestroyFirewallPolicy %q", id)
	client, err := c.CloudClient()
	if err != nil {
		return err
	}
	return client.DestroyFirewallPolicy(id)
}

func (c *Cloud) DestroyCloudIP(id string) error {
	klog.V(4).Infof("DestroyCloudIP (%q)", id)
	client, err := c.CloudClient()
	if err != nil {
		return err
	}
	return client.DestroyCloudIP(id)
}

func (c *Cloud) unmapCloudIP(id string) error {
	klog.V(4).Infof("unmapCloudIP (%q)", id)
	client, err := c.CloudClient()
	if err != nil {
		return err
	}
	return client.UnMapCloudIP(id)
}

//Destroy CloudIPs matching 'name' from a supplied list of cloudIPs
func (c *Cloud) DestroyCloudIPs(cloudIpList []brightbox.CloudIP, name string) error {
	klog.V(4).Infof("DestroyCloudIPs (%q)", name)
	for i := range cloudIpList {
		if cloudIpList[i].Name == name {
			err := c.DestroyCloudIP(cloudIpList[i].Id)
			if err != nil {
				klog.V(4).Infof("Error destroying CloudIP %q", cloudIpList[i].Id)
				return err
			}
		}
	}
	return nil
}

// EnsureOldCloudIPsDeposed unmaps any CloudIPs mapped to the loadbalancer
// that are not the allocated cloud ip.
func (c *Cloud) EnsureOldCloudIPsDeposed(lb *brightbox.LoadBalancer, cip *brightbox.CloudIP, name string) error {
	klog.V(4).Infof("EnsureOldCloudIPsDeposed (%q, %q, %q)", lb.Id, cip.Id, name)
	deposedCloudIPList := getDeposedCloudIpList(lb.CloudIPs, cip.Id)
	for i := range deposedCloudIPList {
		if err := c.unmapCloudIP(deposedCloudIPList[i].Id); err != nil {
			return err
		}
	}
	return nil
}

func getDeposedCloudIpList(cloudIPList []brightbox.CloudIP, id string) []brightbox.CloudIP {
	deposedCloudIpList := make([]brightbox.CloudIP, 0, len(cloudIPList))
	for i := range cloudIPList {
		if cloudIPList[i].Id != id {
			deposedCloudIpList = append(deposedCloudIpList, cloudIPList[i])
		}
	}
	return deposedCloudIpList
}

// obtainCloudClient creates a new Brightbox client using details from
// the environment
func obtainCloudClient() (*brightbox.Client, error) {
	klog.V(4).Infof("obtainCloudClient")
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
	klog.V(4).Infof("validateConfig")
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
			TokenURL:  authd.tokenURL(),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}
	klog.V(4).Infof("Obtaining authentication for user %s", authd.UserName)
	klog.V(4).Infof("Speaking to %s", authd.tokenURL())
	token, err := conf.PasswordCredentialsToken(ctx, authd.UserName, authd.password)
	if err != nil {
		return nil, err
	}
	klog.V(4).Infof("Refreshing current token as required")
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
	klog.V(4).Infof("Obtaining API client authorisation for client %s", authd.APIClient)
	klog.V(4).Infof("Speaking to %s", authd.tokenURL())
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

// SyncServerGroup ensures a Brightbox Server Group contains the supplied
// list of Servers and nothing else
func (c *Cloud) SyncServerGroup(group *brightbox.ServerGroup, newServerIds []string) (*brightbox.ServerGroup, error) {
	oldServerIds := mapServersToServerIds(group.Servers)
	klog.V(4).Infof("SyncServerGroup (%v, %v, %v)", group.Id, oldServerIds, newServerIds)
	client, err := c.CloudClient()
	if err != nil {
		return nil, err
	}
	serverIdsToInsert, serverIdsToDelete := getSyncLists(oldServerIds, newServerIds)
	result := group
	if len(serverIdsToInsert) > 0 {
		klog.V(4).Infof("Adding Servers %v", serverIdsToInsert)
		result, err = client.AddServersToServerGroup(group.Id, serverIdsToInsert)
	}
	if err == nil && len(serverIdsToDelete) > 0 {
		klog.V(4).Infof("Removing Servers %v", serverIdsToDelete)
		result, err = client.RemoveServersFromServerGroup(group.Id, serverIdsToDelete)
	}
	return result, err
}

// IsUpdateLoadBalancerRequired checks whether a set of LoadBalancerOptions
// warrants an API update call.
func IsUpdateLoadBalancerRequired(lb *brightbox.LoadBalancer, newDetails brightbox.LoadBalancerOptions) bool {
	klog.V(6).Infof("Update LoadBalancer Required (%v, %v)", *newDetails.Name, lb.Name)
	return (newDetails.Name != nil && *newDetails.Name != lb.Name) ||
		(newDetails.Healthcheck != nil && isUpdateLoadBalancerHealthcheckRequired(newDetails.Healthcheck, &lb.Healthcheck)) ||
		isUpdateLoadBalancerNodeRequired(newDetails.Nodes, lb.Nodes) ||
		isUpdateLoadBalancerListenerRequired(newDetails.Listeners, lb.Listeners) ||
		isUpdateLoadBalancerDomainsRequired(newDetails.Domains, lb.Acme)
}

func isUpdateLoadBalancerHealthcheckRequired(newHealthCheck *brightbox.LoadBalancerHealthcheck, oldHealthCheck *brightbox.LoadBalancerHealthcheck) bool {
	klog.V(6).Infof("Update LoadBalancer Healthcheck Required (%#v, %#v)", *newHealthCheck, *oldHealthCheck)
	return (newHealthCheck.Type != oldHealthCheck.Type) ||
		(newHealthCheck.Port != oldHealthCheck.Port) ||
		(newHealthCheck.Request != oldHealthCheck.Request)
}

func isUpdateLoadBalancerNodeRequired(a []brightbox.LoadBalancerNode, b []brightbox.Server) bool {
	klog.V(6).Infof("Update LoadBalancer Node Required (%v, %v)", a, b)
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return true
	}
	if len(a) != len(b) {
		return true
	}
	for i := range a {
		if a[i].Node != b[i].Id {
			return true
		}
	}
	return false
}

func isUpdateLoadBalancerListenerRequired(a []brightbox.LoadBalancerListener, b []brightbox.LoadBalancerListener) bool {
	klog.V(6).Infof("Update LoadBalancer Listener Required (%v, %v)", a, b)
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return true
	}
	if len(a) != len(b) {
		return true
	}
	for i := range a {
		if (a[i].Protocol != b[i].Protocol) ||
			(a[i].In != b[i].In) ||
			(a[i].Out != b[i].Out) ||
			(a[i].Timeout != 0 && b[i].Timeout != 0 && a[i].Timeout != b[i].Timeout) ||
			(a[i].ProxyProtocol != b[i].ProxyProtocol) {
			return true
		}
	}
	return false
}

func isUpdateLoadBalancerDomainsRequired(a []string, acme *brightbox.LoadBalancerAcme) bool {
	klog.V(6).Infof("Update LoadBalancer Domains Required (%v)", a)
	if acme == nil {
		return a != nil
	}
	b := make([]string, len(acme.Domains))
	for i, domain := range acme.Domains {
		b[i] = domain.Identifier
	}
	return !sameStringSlice(a, b)
}

func ErrorIfNotErased(lb *brightbox.LoadBalancer) error {
	switch {
	case lb == nil:
		return nil
	case lb.CloudIPs != nil && len(lb.CloudIPs) > 0:
		return fmt.Errorf("CloudIPs still mapped to load balancer %q", lb.Id)
	case !isAlive(lb):
		return nil
	}
	return fmt.Errorf("Unknown reason why %q has not deleted", lb.Id)
}

func ErrorIfNotComplete(lb *brightbox.LoadBalancer, cipId, name string) error {
	switch {
	case lb == nil:
		return fmt.Errorf("Load Balancer for %q is missing", name)
	case !isAlive(lb):
		return fmt.Errorf("Load Balancer %q still building", lb.Id)
	case len(lb.CloudIPs) > 1:
		return fmt.Errorf("Unmapping of deposed CloudIPs to %q not complete", lb.Id)
	case len(lb.CloudIPs) <= 0 || lb.CloudIPs[0].Id != cipId:
		return fmt.Errorf("Mapping of CloudIP %q to %q not complete", cipId, lb.Id)
	}
	return ErrorIfAcmeNotComplete(lb.Acme)
}

func ErrorIfAcmeNotComplete(acme *brightbox.LoadBalancerAcme) error {
	if acme != nil {
		for _, domain := range acme.Domains {
			if domain.Status != ValidAcmeDomainStatus {
				return fmt.Errorf("Domain %q has not yet been validated for SSL use (%q:%q)", domain.Identifier, domain.Status, domain.LastMessage)
			}
		}
	}
	return nil
}
