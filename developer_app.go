
package apigee

import (
	"path"
)

// DeveloperAppService is an interface for interfacing with the Apigee Edge Admin API
// dealing with developerApps.
type DeveloperAppService interface {
	Get(string, string) (*DeveloperApp, *Response, error)
	Create(string, DeveloperApp) (*DeveloperApp, *Response, error)
	Delete(string, string) (*Response, error)
	Update(string, DeveloperApp) (*DeveloperApp, *Response, error)
}

type DeveloperAppServiceOp struct {
	client *EdgeClient
}

var _ DeveloperAppService = &DeveloperAppServiceOp{}

type DeveloperApp struct {
	Name						string			`json:"name,omitempty"`
	ApiProducts 				[]string		`json:"apiProducts,omitempty"`
	KeyExpiresIn				int				`json:"keyExpiresIn,omitempty"`
	Attributes					[]Attribute		`json:"attributes,omitempty"`
	Scopes						[]string		`json:"scopes,omitempty"`
	CallbackUrl					string			`json:"callbackUrl,omitempty"`
	Credentials					[]Credential	`json:"credentials,omitempty"`
	AppId						string			`json:"appId,omitempty"`
	DeveloperId					string			`json:"developerId,omitempty"`
	AppFamily					string			`json:"appFamily,omitempty"`
	Status						string			`json:"status,omitempty"`

}

func (s *DeveloperAppServiceOp) Get(email string, name string) (*DeveloperApp, *Response, error) {

	path := path.Join("developers", email, "apps", name)

	req, e := s.client.NewRequest("GET", path, nil, "")
	if e != nil {
		return nil, nil, e
	}
	returnedDeveloperApp := DeveloperApp{}
	resp, e := s.client.Do(req, &returnedDeveloperApp)
	if e != nil {
		return nil, resp, e
	}
	return &returnedDeveloperApp, resp, e

}

func (s *DeveloperAppServiceOp) Create(email string, developerApp DeveloperApp) (*DeveloperApp, *Response, error) {

	return postOrPutDeveloperApp(email, developerApp, "POST", s)

}


func (s *DeveloperAppServiceOp) Update(email string, developerApp DeveloperApp) (*DeveloperApp, *Response, error) {

	return postOrPutDeveloperApp(email, developerApp, "PUT", s)

}

func (s *DeveloperAppServiceOp) Delete(email string, name string) (*Response, error) {

	path := path.Join("developers", email, "apps", name)

	req, e := s.client.NewRequest("DELETE", path, nil, "")
	if e != nil {
		return nil, e
	}

	resp, e := s.client.Do(req, nil)
	if e != nil {
		return resp, e
	}

	return resp, e

}

func postOrPutDeveloperApp(email string, developerApp DeveloperApp, opType string, s *DeveloperAppServiceOp) (*DeveloperApp, *Response, error) {

	uripath := ""

	if opType == "PUT" {
		uripath = path.Join("developers", email, "apps", developerApp.Name)
	} else {
		uripath = path.Join("developers", email, "apps")
	}

	req, e := s.client.NewRequest(opType, uripath, developerApp, "")
	if e != nil {
		return nil, nil, e
	}

	returnedDeveloperApp := DeveloperApp{}

	resp, e := s.client.Do(req, &returnedDeveloperApp)
	if e != nil {
		return nil, resp, e
	}

	return &returnedDeveloperApp, resp, e

}