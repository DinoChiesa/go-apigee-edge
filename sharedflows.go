package apigee

//const sfUriPathElement = "sharedflows"

// SharedFlowsService is an interface for interfacing with the Apigee Admin API
// dealing with apiproxies.
type SharedFlowsService interface {
  List() ([]string, *Response, error)
  Get(string) (*DeployableAsset, *Response, error)
  Import(string, string) (*DeployableRevision, *Response, error)
  Delete(string) (*DeletedItemInfo, *Response, error)
  DeleteRevision(string, Revision) (*DeployableRevision, *Response, error)
  Deploy(string,string,Revision) (*RevisionDeployment, *Response, error)
  Undeploy(string,string,Revision) (*RevisionDeployment, *Response, error)
  Export(string, Revision) (string, *Response, error)
  GetDeployments(string) (*Deployment, *Response, error)
}

type SharedFlowsServiceOp struct {
  client *ApigeeClient
  deployable Deployable
}

var _ SharedFlowsService = &SharedFlowsServiceOp{}

func (s *SharedFlowsServiceOp) List() ([]string, *Response, error) {
	return s.deployable.List(s.client, sfUriPathElement)
}

func (s *SharedFlowsServiceOp) Get(proxyName string) (*DeployableAsset, *Response, error) {
	return s.deployable.Get(s.client, sfUriPathElement, proxyName)
}

func (s *SharedFlowsServiceOp) Import(proxyName string, source string) (*DeployableRevision, *Response, error) {
	return s.deployable.Import(s.client, sfUriPathElement, proxyName, source)
}

func (s *SharedFlowsServiceOp) Export(proxyName string, rev Revision) (string, *Response, error) {
	return s.deployable.Export(s.client, sfUriPathElement, proxyName, rev)
}

func (s *SharedFlowsServiceOp) DeleteRevision(proxyName string, rev Revision) (*DeployableRevision, *Response, error) {
	return s.deployable.DeleteRevision(s.client, sfUriPathElement, proxyName, rev)
}

func (s *SharedFlowsServiceOp) Undeploy(proxyName, env string, rev Revision) (*RevisionDeployment, *Response, error) {
	return s.deployable.Undeploy(s.client, sfUriPathElement, proxyName, env, rev)
}

func (s *SharedFlowsServiceOp) Deploy(proxyName, env string, rev Revision) (*RevisionDeployment, *Response, error) {
	return s.deployable.Deploy(s.client, sfUriPathElement, proxyName, "", env, rev)
}

func (s *SharedFlowsServiceOp) Delete(proxyName string) (*DeletedItemInfo, *Response, error) {
	return s.deployable.Delete(s.client, sfUriPathElement, proxyName)
}

func (s *SharedFlowsServiceOp) GetDeployments(proxyName string) (*Deployment, *Response, error) {
	return s.deployable.GetDeployments(s.client, sfUriPathElement, proxyName)
}
