package apigee

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const sharedFlows = "sharedflows"

// SharedFlowService is an interface for interfacing with the Apigee Edge Admin API
// dealing with shardFlows.
type SharedFlowService interface {
	List() ([]string, *Response, error)
	Get(string) (*SharedFlow, *Response, error)
	Deploy(string, string, Revision, int, bool) (*SharedFlowRevisionDeployment, *Response, error)
	Import(string, string) (*SharedFlowRevision, *Response, error)
	Delete(string) (*DeletedSharedFlowInfo, *Response, error)
	GetDeployments(string) (*SharedFlowDeployment, *Response, error)
	ReDeploy(string, string, Revision, int, bool) (*SharedFlowRevisionDeployments, *Response, error)
	Undeploy(string, string, Revision) (*SharedFlowRevisionDeployment, *Response, error)
}

// SharedFlowRevision holds information about a revision of an API Proxy.
type SharedFlowRevision struct {
	CreatedBy       string    `json:"createdBy,omitempty"`
	CreatedAt       Timestamp `json:"createdAt,omitempty"`
	Description     string    `json:"description,omitempty"`
	ContextInfo     string    `json:"contextInfo,omitempty"`
	DisplayName     string    `json:"displayName,omitempty"`
	Name            string    `json:"name,omitempty"`
	LastModifiedBy  string    `json:"lastModifiedBy,omitempty"`
	LastModifiedAt  Timestamp `json:"lastModifiedAt,omitempty"`
	Revision        Revision  `json:"revision,omitempty"`
	TargetEndpoints []string  `json:"targetEndpoints,omitempty"`
	TargetServers   []string  `json:"targetServers,omitempty"`
	Resources       []string  `json:"resources,omitempty"`
	Policies        []string  `json:"policies,omitempty"`
	Type            string    `json:"type,omitempty"`
}

type SharedFlowServiceOp struct {
	client *EdgeClient
}

type SharedFlow struct {
	Revisions []Revision         `json:"revision,omitempty"`
	Name      string             `json:"name,omitempty"`
	MetaData  SharedFlowMetadata `json:"metaData,omitempty"`
}

// SharedFlowMetadata contains information related to the creation and last modified
// time and actor for an API Proxy within an organization.
type SharedFlowMetadata struct {
	LastModifiedBy string    `json:"lastModifiedBy,omitempty"`
	CreatedBy      string    `json:"createdBy,omitempty"`
	LastModifiedAt Timestamp `json:"lastModifiedAt,omitempty"`
	CreatedAt      Timestamp `json:"createdAt,omitempty"`
}

// SharedFlowRevisionDeployment holds information about the deployment state of a
// single revision of a shared flow.
type SharedFlowRevisionDeployment struct {
	Name         string       `json:",omitempty"`
	Revision     Revision     `json:"revision,omitempty"`
	Environment  string       `json:"environment,omitempty"`
	Organization string       `json:"organization,omitempty"`
	State        string       `json:"state,omitempty"`
	Servers      []EdgeServer `json:"server,omitempty"`
}

// SharedFlowRevisionDeployments holds information about the deployment state of a
// single revision of a shared flow across environments
type SharedFlowRevisionDeployments struct {
	Name         string                         `json:"aPIProxy,omitempty"`
	Environments []SharedFlowRevisionDeployment `json:"environment,omitempty"`
	Organization string                         `json:"organization,omitempty"`
}

// SharedFlowDeployment holds information about the deployment state of
// all revisions of a shared flow
type SharedFlowDeployment struct {
	Environments []EnvironmentDeployment `json:"environment,omitempty"`
	Name         string                  `json:"name,omitempty"`
	Organization string                  `json:"organization,omitempty"`
}

// DeletedSharedFlowInfo is a  payload that contains very little useful
// information. This struct deserializes that information.
type DeletedSharedFlowInfo struct {
	Name string `json:"name,omitempty"`
}

// Get retrieves the information about a SharedFlow in an organization, information including
// the list of available revisions, and the created and last modified dates and actors.
func (s *SharedFlowServiceOp) Get(name string) (*SharedFlow, *Response, error) {
	path := path.Join(sharedFlows, name)
	req, err := s.client.NewRequest("GET", path, nil, "")
	if err != nil {
		return nil, nil, err
	}
	sharedFlow := &SharedFlow{}
	resp, err := s.client.Do(req, sharedFlow)
	if err != nil {
		return nil, nil, err
	}

	return sharedFlow, resp, err
}

// List retrieves the list of sharedFlow names for the organization referred by the EdgeClient.
func (s *SharedFlowServiceOp) List() ([]string, *Response, error) {
	req, err := s.client.NewRequest("GET", sharedFlows, nil, "")
	if err != nil {
		return nil, nil, err
	}
	namelist := make([]string, 0)
	resp, err := s.client.Do(req, &namelist)
	if err != nil {
		return nil, resp, err
	}
	return namelist, resp, err
}

// Deploy a revision of a ShareFlow to a specific environment within an organization.
func (s *SharedFlowServiceOp) Deploy(name, env string, rev Revision, delay int, override bool) (*SharedFlowRevisionDeployment, *Response, error) {
	// TODO test this after creating a new one
	deployURL, err := url.Parse(path.Join("environments", env, sharedFlows, name, "revisions", fmt.Sprintf("%d", rev), "deployments"))
	if err != nil {
		return nil, nil, nil
	}
	q := deployURL.Query()
	q.Add("override", strconv.FormatBool(override))
	q.Add("delay", fmt.Sprintf("%d", delay))
	deployURL.RawQuery = q.Encode()
	path := deployURL.String()
	req, err := s.client.NewRequest("POST", path, nil, "application/x-www-form-urlencoded")
	if err != nil {
		return nil, nil, err
	}

	deployment := SharedFlowRevisionDeployment{}
	resp, e := s.client.Do(req, &deployment)

	return &deployment, resp, e
}

// Import an SharedFlow into an organization, creating a new API Proxy revision.
// The proxyName can be passed as "nil" in which case the name is derived from the source.
// The source can be either a filesystem directory containing an exploded apiproxy bundle, OR
// the path of a zip file containing an SharedFlow bundle. Returns the API proxy revision information.
// This method does not deploy the imported proxy. See the Deploy method.
func (s *SharedFlowServiceOp) Import(name string, source string) (*SharedFlowRevision, *Response, error) {
	info, err := os.Stat(source)
	if err != nil {
		return nil, nil, err
	}
	zipfileName := source

	log.Printf("[INFO] *** Import *** isDir: %#v\n", info.IsDir())

	if info.IsDir() {
		// create a temporary zip file
		if name == "" {
			name = filepath.Base(source)
		}
		log.Printf("[INFO] *** Import *** proxyName: %#v\n", name)
		tempDir, err := ioutil.TempDir("", "go-apigee-edge-")
		if err != nil {
			log.Printf("[ERROR] *** Import *** error: %#v\n", err)
			return nil, nil, fmt.Errorf("while creating temp dir, error: %#v", err)
		}
		log.Printf("[INFO] *** Import *** tempDir: %#v\n", tempDir)
		log.Printf("[INFO] *** Import *** sourceDir: %#v\n", source)
		zipfileName = path.Join(tempDir, "sharedflow.zip")
		err = zipDirectory(path.Join(source, "sharedflowbundle"), zipfileName, smartFilter)
		if err != nil {
			return nil, nil, fmt.Errorf("while creating temp dir, error: %#v", err)
		}
		log.Printf("[INFO] *** zipped %s into %s\n\n", source, zipfileName)
	}

	if !strings.HasSuffix(zipfileName, ".zip") {
		return nil, nil, errors.New("source must be a zipfile")
	}

	info, err = os.Stat(zipfileName)
	if err != nil {
		return nil, nil, err
	}

	origURL, err := url.Parse(sharedFlows)
	if err != nil {
		return nil, nil, err
	}
	q := origURL.Query()
	q.Add("action", "import")
	q.Add("name", name)
	origURL.RawQuery = q.Encode()
	path := origURL.String()

	ioreader, err := os.Open(zipfileName)
	if err != nil {
		return nil, nil, err
	}
	defer ioreader.Close()

	req, err := s.client.NewRequest("POST", path, ioreader, "")
	if err != nil {
		return nil, nil, err
	}
	sharedFlowRevision := SharedFlowRevision{}
	resp, err := s.client.Do(req, &sharedFlowRevision)
	if err != nil {
		return nil, resp, err
	}
	return &sharedFlowRevision, resp, err
}

// Delete an SharedFlow and all its revisions from an organization. This method
// will fail if any of the revisions of the named API Proxy are currently deployed
// in any environment.
func (s *SharedFlowServiceOp) Delete(name string) (*DeletedSharedFlowInfo, *Response, error) {
	path := path.Join(sharedFlows, name)
	req, e := s.client.NewRequest("DELETE", path, nil, "")
	if e != nil {
		return nil, nil, e
	}
	sharedFlow := DeletedSharedFlowInfo{}
	resp, err := s.client.Do(req, &sharedFlow)
	if err != nil {
		return nil, resp, err
	}
	return &sharedFlow, resp, err
}

// GetDeployments retrieves the information about deployments of a shared flow in
// an organization, including the environment names and revision numbers.
func (s *SharedFlowServiceOp) GetDeployments(name string) (*SharedFlowDeployment, *Response, error) {
	path := path.Join(sharedFlows, name, "deployments")
	req, e := s.client.NewRequest("GET", path, nil, "")
	if e != nil {
		return nil, nil, e
	}
	deployments := SharedFlowDeployment{}
	resp, e := s.client.Do(req, &deployments)
	if e != nil {
		return nil, resp, e
	}
	return &deployments, resp, e
}

func (s *SharedFlowServiceOp) ReDeploy(sharedFlowName, env string, rev Revision, delay int, override bool) (*SharedFlowRevisionDeployments, *Response, error) {

	req, e := prepareDeployRequest(sharedFlowName, env, sharedFlows, rev, delay, override, s.client)

	deployment := SharedFlowRevisionDeployments{}
	resp, e := s.client.Do(req, &deployment)

	return &deployment, resp, e

}

// Undeploy a specific revision of a shared flow from a particular environment within an Edge organization.
func (s *SharedFlowServiceOp) Undeploy(sharedFlowName, env string, rev Revision) (*SharedFlowRevisionDeployment, *Response, error) {
	path := path.Join(sharedFlows, sharedFlowName, "revisions", fmt.Sprintf("%d", rev), "deployments")
	// append the query params
	origURL, err := url.Parse(path)
	if err != nil {
		return nil, nil, err
	}
	q := origURL.Query()
	q.Add("action", "undeploy")
	q.Add("env", env)
	origURL.RawQuery = q.Encode()
	path = origURL.String()

	req, e := s.client.NewRequest("POST", path, nil, "")
	if e != nil {
		return nil, nil, e
	}

	deployment := SharedFlowRevisionDeployment{}
	resp, e := s.client.Do(req, &deployment)
	if e != nil {
		return nil, resp, e
	}
	return &deployment, resp, e
}
