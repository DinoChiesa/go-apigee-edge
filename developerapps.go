package apigee

import (
  "path"
  "errors"
)


// DeveloperappsService is an interface for interfacing with the Apigee Edge Admin API
// dealing with apps that belong to developers.
type DeveloperAppsService interface {
  Create(DeveloperApp) (*DeveloperApp, *Response, error)
  Delete(string) (*DeveloperApp, *Response, error)
  // List() ([]string, *Response, error)
  // Get( string) (*Developer, *Response, error)
  // Update(DeveloperApp) (*Developer, *Response, error)
  // Revoke(string) (*Response, error)
  // Approve(string) (*Response, error)
}

type DeveloperAppsServiceOp struct {
  client *EdgeClient
  developerId string
}

var _ DeveloperAppsService = &DeveloperAppsServiceOp{}

// DeveloperApp holds information about a registered DeveloperApp.
type DeveloperApp struct {
  Name             string      `json:"name,omitempty"`
  ApiProducts      []string    `json:"apiProducts,omitempty"`
  InitialKeyExpiry string      `json:"keyExpiresIn,omitempty"`
  Attributes       Attributes  `json:"attributes,omitempty"`
  Id               string      `json:"appId,omitempty"`
  DeveloperId      string      `json:"developerId,omitempty"`
  Scopes           []string    `json:"scopes,omitempty"`
  Status           string      `json:"status,omitempty"`
}

func (s *DeveloperAppsServiceOp) Create(app DeveloperApp) (*DeveloperApp, *Response, error) {
	if (app.Name == "") {
		return nil, nil, errors.New("cannot create a developerapp with no name")
	}
	appsPath := path.Join(developersPath, s.developerId, "apps") 
  req, e := s.client.NewRequest("POST", appsPath, app)
  if e != nil {
    return nil, nil, e
  }
  returnedApp := DeveloperApp{}
  resp, e := s.client.Do(req, &returnedApp)
  if e != nil {
    return nil, resp, e
  }
  return &returnedApp, resp, e
}

func (s *DeveloperAppsServiceOp) Delete(appName string) (*DeveloperApp, *Response, error) {
  path := path.Join(developersPath, s.developerId, "apps", appName)
  req, e := s.client.NewRequest("DELETE", path, nil)
  if e != nil {
    return nil, nil, e
  }
  deletedApp := DeveloperApp{}
  resp, e := s.client.Do(req, &deletedApp)
  if e != nil {
    return nil, resp, e
  }
  return &deletedApp, resp, e
}

