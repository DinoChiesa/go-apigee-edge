
package apigee

import (
	"path"
)

// CompanyService is an interface for interfacing with the Apigee Edge Admin API
// dealing with companys.
type CompanyService interface {
	Get(string) (*Company, *Response, error)
	Create(Company) (*Company, *Response, error)
	Delete(string) (*Response, error)
	Update(Company) (*Company, *Response, error)
}

type CompanyServiceOp struct {
	client *EdgeClient
}

var _ CompanyService = &CompanyServiceOp{}

type Company struct {
	Name						string			`json:"name,omitempty"`
	DisplayName					string			`json:"displayName,omitempty"`
	Attributes					[]Attribute		`json:"attributes,omitempty"`

	Status						string			`json:"status,omitempty"`
	Apps 						[]string		`json:"apps,omitempty"`
}

func (s *CompanyServiceOp) Get(name string) (*Company, *Response, error) {

	path := path.Join("companies", name)

	req, e := s.client.NewRequest("GET", path, nil, "")
	if e != nil {
		return nil, nil, e
	}
	returnedCompany := Company{}
	resp, e := s.client.Do(req, &returnedCompany)
	if e != nil {
		return nil, resp, e
	}
	return &returnedCompany, resp, e

}

func (s *CompanyServiceOp) Create(company Company) (*Company, *Response, error) {

	return postOrPutCompany(company, "POST", s)

}


func (s *CompanyServiceOp) Update(company Company) (*Company, *Response, error) {

	return postOrPutCompany(company, "PUT", s)

}

func (s *CompanyServiceOp) Delete(name string) (*Response, error) {

	path := path.Join("companies", name)

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

func postOrPutCompany(company Company, opType string, s *CompanyServiceOp) (*Company, *Response, error) {

	uripath := ""

	if opType == "PUT" {
		uripath = path.Join("companies", company.Name)
	} else {
		uripath = path.Join("companies")
	}

	req, e := s.client.NewRequest(opType, uripath, company, "")
	if e != nil {
		return nil, nil, e
	}

	returnedCompany := Company{}

	resp, e := s.client.Do(req, &returnedCompany)
	if e != nil {
		return nil, resp, e
	}

	return &returnedCompany, resp, e

}