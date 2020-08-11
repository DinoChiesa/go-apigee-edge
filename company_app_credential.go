package apigee

import (
	"path"
)

// CompanyAppCredentialService is an interface for interfacing with the Apigee Edge Admin API
// dealing with companyApp credentials.
type CompanyAppCredentialService interface {
	Create(string, string, CompanyAppCredential) (*Credential, *Response, error)
	Update(string, string, string, CompanyAppCredential) (*Credential, *Response, error)
	Get(string, string, string) (*Credential, *Response, error)
	Delete(string, string, string) (*Response, error)
	RemoveApiProduct(string, string, string, string) (*Response, error)
}

type CompanyAppCredentialServiceOp struct {
	client *EdgeClient
}

var _ CompanyAppCredentialService = &CompanyAppCredentialServiceOp{}

type CompanyAppCredential struct {
	ConsumerKey    string      `json:"consumerKey,omitempty"`
	ConsumerSecret string      `json:"consumerSecret,omitempty"`
	ApiProducts    []string    `json:"apiProducts,omitempty"`
	ExpiresAt      int         `json:"expiresAt,omitempty"`
	Attributes     []Attribute `json:"attributes,omitempty"`
	Scopes         []string    `json:"scopes,omitempty"`
}

func (s *CompanyAppCredentialServiceOp) Create(companyName string, appName string, companyAppCredential CompanyAppCredential) (*Credential, *Response, error) {

	uripath := path.Join("companies", companyName, "apps", appName, "keys", "create")

	req, e := s.client.NewRequest("POST", uripath, companyAppCredential, "")
	if e != nil {
		return nil, nil, e
	}

	returnedCompanyAppCredentials := Credential{}

	resp, e := s.client.Do(req, &returnedCompanyAppCredentials)
	if e != nil {
		return nil, resp, e
	}

	return &returnedCompanyAppCredentials, resp, e

}

func (s *CompanyAppCredentialServiceOp) Update(companyName string, appName string, consumerKey string, companyAppCredential CompanyAppCredential) (*Credential, *Response, error) {

	uripath := path.Join("companies", companyName, "apps", appName, "keys", consumerKey)

	req, e := s.client.NewRequest("POST", uripath, companyAppCredential, "")
	if e != nil {
		return nil, nil, e
	}

	returnedCompanyAppCredentials := Credential{}

	resp, e := s.client.Do(req, &returnedCompanyAppCredentials)
	if e != nil {
		return nil, resp, e
	}

	return &returnedCompanyAppCredentials, resp, e

}

func (s *CompanyAppCredentialServiceOp) Get(companyName string, appName string, consumerKey string) (*Credential, *Response, error) {

	uripath := path.Join("companies", companyName, "apps", appName, "keys", consumerKey)

	req, e := s.client.NewRequest("GET", uripath, nil, "")
	if e != nil {
		return nil, nil, e
	}
	returnedCompanyAppCredential := Credential{}
	resp, e := s.client.Do(req, &returnedCompanyAppCredential)
	if e != nil {
		return nil, resp, e
	}
	return &returnedCompanyAppCredential, resp, e

}

func (s *CompanyAppCredentialServiceOp) Delete(companyName string, appName string, consumerKey string) (*Response, error) {

	uripath := path.Join("companies", companyName, "apps", appName, "keys", consumerKey)

	req, e := s.client.NewRequest("DELETE", uripath, nil, "")
	if e != nil {
		return nil, e
	}

	resp, e := s.client.Do(req, nil)
	if e != nil {
		return resp, e
	}

	return resp, e

}

func (s *CompanyAppCredentialServiceOp) RemoveApiProduct(companyName string, appName string, consumerKey string, apiProductName string) (*Response, error) {

	uripath := path.Join("companies", companyName, "apps", appName, "keys", consumerKey, "apiproducts", apiProductName)

	req, e := s.client.NewRequest("DELETE", uripath, nil, "")
	if e != nil {
		return nil, e
	}

	resp, e := s.client.Do(req, nil)
	if e != nil {
		return resp, e
	}

	return resp, e

}
