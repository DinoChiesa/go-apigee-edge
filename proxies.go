package apigee

const uriPathElement = "apis"

// ProxiesService is an interface for interfacing with the Apigee Admin API
// dealing with apiproxies.
type ProxiesService interface {
	List() ([]string, *Response, error)
	Get(string) (*DeployableAsset, *Response, error)
	Import(string, string) (*DeployableRevision, *Response, error)
	Delete(string) (*DeletedItemInfo, *Response, error)
	DeleteRevision(string, Revision) (*DeployableRevision, *Response, error)
	Deploy(string, string, Revision) (*RevisionDeployment, *Response, error)
	DeployAtPath(string, string, string, Revision) (*RevisionDeployment, *Response, error)
	Undeploy(string, string, Revision) (*RevisionDeployment, *Response, error)
	Export(string, Revision) (string, *Response, error)
	GetDeployments(string) (*Deployment, *Response, error)
}

type ProxiesServiceOp struct {
	client     *ApigeeClient
	deployable Deployable
}

var _ ProxiesService = &ProxiesServiceOp{}

// // ProxyRevisionDeployment holds information about the deployment state of a
// // single revision of an API Proxy.
// type ProxyRevisionDeployment struct {
//   Name            string        `json:"aPIProxy,omitempty"`
//   Revision        Revision      `json:"revision,omitempty"`
//   Environment     string        `json:"environment,omitempty"`
//   Organization    string        `json:"organization,omitempty"`
//   State           string        `json:"state,omitempty"`
//   Servers         []ApigeeServer  `json:"server,omitempty"`
// }

// type proxiesRoot struct {
//   Proxies []Proxy `json:"proxies"`
// }

// retrieve the list of apiproxy names for the organization referred by the ApigeeClient.
func (s *ProxiesServiceOp) List() ([]string, *Response, error) {
	return s.deployable.List(s.client, uriPathElement)
}

// Get retrieves the information about an API Proxy in an organization, information including
// the list of available revisions, and the created and last modified dates and actors.
func (s *ProxiesServiceOp) Get(proxyName string) (*DeployableAsset, *Response, error) {
	return s.deployable.Get(s.client, uriPathElement, proxyName)
}

// Import an API proxy into an organization, creating a new API Proxy revision.
// The proxyName can be passed as "nil" in which case the name is derived from the source.
// The source can be either a filesystem directory containing an exploded apiproxy bundle, OR
// the path of a zip file containing an API Proxy bundle. Returns the API proxy revision information.
// This method does not deploy the imported proxy. See the Deploy method.
func (s *ProxiesServiceOp) Import(proxyName string, source string) (*DeployableRevision, *Response, error) {
	return s.deployable.Import(s.client, uriPathElement, proxyName, source)
}

// Export a revision of an API proxy within an organization, to a filesystem file.
func (s *ProxiesServiceOp) Export(proxyName string, rev Revision) (string, *Response, error) {
	return s.deployable.Export(s.client, uriPathElement, proxyName, rev)
}

// DeleteRevision deletes a specific revision of an API Proxy from an organization.
// The revision must exist, and must not be currently deployed.
func (s *ProxiesServiceOp) DeleteRevision(proxyName string, rev Revision) (*DeployableRevision, *Response, error) {
	return s.deployable.DeleteRevision(s.client, uriPathElement, proxyName, rev)
}

// Undeploy a specific revision of an API Proxy from a particular environment within an Edge organization.
func (s *ProxiesServiceOp) Undeploy(proxyName, env string, rev Revision) (*RevisionDeployment, *Response, error) {
	return s.deployable.Undeploy(s.client, uriPathElement, proxyName, env, rev)
}

// Deploy a revision of an API proxy to a specific environment within an organization.
func (s *ProxiesServiceOp) Deploy(proxyName, env string, rev Revision) (*RevisionDeployment, *Response, error) {
	return s.deployable.Deploy(s.client, uriPathElement, proxyName, "", env, rev)
}

// Deploy a revision of an API proxy to a specific environment within an organization.
func (s *ProxiesServiceOp) DeployAtPath(proxyName, basepath, env string, rev Revision) (*RevisionDeployment, *Response, error) {
	return s.deployable.Deploy(s.client, uriPathElement, proxyName, basepath, env, rev)
}

// Delete an API Proxy and all its revisions from an organization. This method
// will fail if any of the revisions of the named API Proxy are currently deployed
// in any environment.
func (s *ProxiesServiceOp) Delete(proxyName string) (*DeletedItemInfo, *Response, error) {
	return s.deployable.Delete(s.client, uriPathElement, proxyName)
}

// GetDeployments retrieves the information about deployments of an API Proxy in
// an organization, including the environment names and revision numbers.
func (s *ProxiesServiceOp) GetDeployments(proxyName string) (*Deployment, *Response, error) {
	return s.deployable.GetDeployments(s.client, uriPathElement, proxyName)
}
