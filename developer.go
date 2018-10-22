package apigee

import (
	"path"
)

// DeveloperService is an interface for interfacing with the Apigee Edge Admin API
// dealing with developers.
type DeveloperService interface {
	Get(string) (*Developer, *Response, error)
	Create(Developer) (*Developer, *Response, error)
	Delete(string) (*Response, error)
	Update(Developer) (*Developer, *Response, error)
}

type DeveloperServiceOp struct {
	client *EdgeClient
}

var _ DeveloperService = &DeveloperServiceOp{}

type Developer struct {
	Email       string      `json:"email,omitempty"`
	FirstName   string      `json:"firstName,omitempty"`
	LastName    string      `json:"lastName,omitempty"`
	UserName    string      `json:"userName,omitempty"`
	Attributes  []Attribute `json:"attributes,omitempty"`
	DeveloperId string      `json:"developerId,omitempty"`

	Apps   []string `json:"apps,omitempty"`
	Status string   `json:"status,omitempty"`
}

func (s *DeveloperServiceOp) Get(email string) (*Developer, *Response, error) {

	path := path.Join("developers", email)

	req, e := s.client.NewRequest("GET", path, nil, "")
	if e != nil {
		return nil, nil, e
	}
	returnedDeveloper := Developer{}
	resp, e := s.client.Do(req, &returnedDeveloper)
	if e != nil {
		return nil, resp, e
	}
	return &returnedDeveloper, resp, e

}

func (s *DeveloperServiceOp) Create(developer Developer) (*Developer, *Response, error) {

	return postOrPutDeveloper(developer, "POST", s)

}

func (s *DeveloperServiceOp) Update(developer Developer) (*Developer, *Response, error) {

	return postOrPutDeveloper(developer, "PUT", s)

}

func (s *DeveloperServiceOp) Delete(email string) (*Response, error) {

	path := path.Join("developers", email)

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

func postOrPutDeveloper(developer Developer, opType string, s *DeveloperServiceOp) (*Developer, *Response, error) {

	uripath := ""

	if opType == "PUT" {
		uripath = path.Join("developers", developer.Email)
	} else {
		uripath = path.Join("developers")
	}

	req, e := s.client.NewRequest(opType, uripath, developer, "")
	if e != nil {
		return nil, nil, e
	}

	returnedDeveloper := Developer{}

	resp, e := s.client.Do(req, &returnedDeveloper)
	if e != nil {
		return nil, resp, e
	}

	return &returnedDeveloper, resp, e

}
