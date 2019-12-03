package k8ssdk

import brightbox "github.com/brightbox/gobrightbox"

// CloudAccess is an abstraction over the Brightbox API to allow testing
type CloudAccess interface {
	//Fetch a server
	Server(identifier string) (*brightbox.Server, error)

	//Fetch a list of LoadBalancers
	LoadBalancers() ([]brightbox.LoadBalancer, error)

	//Retrieves a detailed view of one load balancer
	LoadBalancer(identifier string) (*brightbox.LoadBalancer, error)

	//Creates a new load balancer
	CreateLoadBalancer(newDetails *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error)

	//Updates an existing load balancer
	UpdateLoadBalancer(newDetails *brightbox.LoadBalancerOptions) (*brightbox.LoadBalancer, error)

	//Retrieves a list of all cloud IPs
	CloudIPs() ([]brightbox.CloudIP, error)

	//retrieves a detailed view of one cloud ip
	CloudIP(identifier string) (*brightbox.CloudIP, error)

	//Issues a request to map the cloud ip to the destination
	MapCloudIP(identifier string, destination string) error

	//UnMapCloudIP issues a request to unmap the cloud ip
	UnMapCloudIP(identifier string) error

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
