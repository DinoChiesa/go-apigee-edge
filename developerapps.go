package apigee

import (
  "path"
  "net/url"
  "errors"
)


// DeveloperAppsService is an interface for interfacing with the Apigee Edge Admin API
// dealing with apps that belong to a particular developer.
type DeveloperAppsService interface {
  Create(DeveloperApp) (*DeveloperApp, *Response, error)
  Delete(string) (*DeveloperApp, *Response, error)
  Revoke(string) (*Response, error)
  Approve(string) (*Response, error)
  List() ([]string, *Response, error)
  Get( string) (*DeveloperApp, *Response, error)
  Update(DeveloperApp) (*DeveloperApp, *Response, error)
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


func updateAppStatus (s DeveloperAppsServiceOp, appName string, desiredStatus string) (*Response, error) {

  appPath := path.Join(developersPath, s.developerId, "apps", appName)

  // append the necessary query param
  origURL, e := url.Parse(appPath)
  if e != nil {
     return nil, e
  }
  q := origURL.Query()
  q.Add("action", desiredStatus)
  origURL.RawQuery = q.Encode()
  appPath = origURL.String()

	req, e := s.client.NewRequest("POST", appPath, nil)
  if e != nil {
    return nil, e
  }
  resp, e := s.client.Do(req, nil)
  if e != nil {
    return resp, e
  }
  return resp, e
}

func (s *DeveloperAppsServiceOp) Revoke(appName string) (*Response, error) {
	return updateAppStatus(*s, appName, "revoke")
}

func (s *DeveloperAppsServiceOp) Approve(appName string) (*Response, error) {
	return updateAppStatus(*s, appName, "approve")
}

func (s *DeveloperAppsServiceOp) List() ([]string, *Response, error) {
  appsPath := path.Join(developersPath, s.developerId, "apps")
  req, e := s.client.NewRequest("GET", appsPath, nil)
  if e != nil {
    return nil, nil, e
  }
  namelist := make([]string,0)
  resp, e := s.client.Do(req, &namelist)
  if e != nil {
    return nil, resp, e
  }
  return namelist, resp, e
}

func (s *DeveloperAppsServiceOp) Get(appName string) (*DeveloperApp, *Response, error) {
  appPath := path.Join(developersPath, s.developerId, "apps", appName)
  req, e := s.client.NewRequest("GET", appPath, nil)
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

func (s *DeveloperAppsServiceOp) Update(app DeveloperApp) (*DeveloperApp, *Response, error) {
	if app.Name == "" {
    return nil, nil, errors.New("missing the Name of the App to update")
	}
	appPath := path.Join(developersPath, s.developerId, "apps", app.Name)
	
  req, e := s.client.NewRequest("POST", appPath, app)
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

