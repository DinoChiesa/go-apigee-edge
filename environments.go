package apigee

import (
	"path"
)

const environmentsPath = "environments"

// EnvironmentsService is an interface for interfacing with the Apigee Edge Admin API
// querying Edge environments.
type EnvironmentsService interface {
	List() ([]string, *Response, error)
	Get(string) (*Environment, *Response, error)
}

type EnvironmentsServiceOp struct {
	client *ApigeeClient
}

var _ EnvironmentsService = &EnvironmentsServiceOp{}

// Environment contains information about an environment within an Edge organization.
type Environment struct {
	Name           string          `json:"name,omitempty"`
	CreatedBy      string          `json:"createdBy,omitempty"`
	CreatedAt      Timestamp       `json:"createdAt,omitempty"`
	LastModifiedBy string          `json:"lastModifiedBy,omitempty"`
	LastModifiedAt Timestamp       `json:"lastModifiedAt,omitempty"`
	Properties     PropertyWrapper `json:"properties,omitempty"`
}

// List retrieves the list of environment names for the organization referred by the ApigeeClient.
func (s *EnvironmentsServiceOp) List() ([]string, *Response, error) {
	req, e := s.client.NewRequest("GET", environmentsPath, nil)
	if e != nil {
		return nil, nil, e
	}
	namelist := make([]string, 0)
	resp, e := s.client.Do(req, &namelist)
	if e != nil {
		return nil, resp, e
	}
	return namelist, resp, e
}

// Get retrieves the information about an Environment in an organization, information including
// the properties, and the created and last modified details.
func (s *EnvironmentsServiceOp) Get(env string) (*Environment, *Response, error) {
	path := path.Join(environmentsPath, env)
	req, e := s.client.NewRequest("GET", path, nil)
	if e != nil {
		return nil, nil, e
	}
	returnedEnv := Environment{}
	resp, e := s.client.Do(req, &returnedEnv)
	if e != nil {
		return nil, resp, e
	}
	return &returnedEnv, resp, e
}
