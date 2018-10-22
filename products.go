package apigee

import (
	"path"
)

// ProductsService is an interface for interfacing with the Apigee Edge Admin API
// dealing with apiproducts.
type ProductsService interface {
	Get(string) (*Product, *Response, error)
	Create(Product) (*Product, *Response, error)
	Delete(string) (*Response, error)
	Update(Product) (*Product, *Response, error)
}

type ProductsServiceOp struct {
	client *EdgeClient
}

var _ ProductsService = &ProductsServiceOp{}

type Product struct {
	Name          string      `json:"name,omitempty"`
	DisplayName   string      `json:"displayName,omitempty"`
	ApprovalType  string      `json:"approvalType,omitempty"` //manual or auto
	Attributes    []Attribute `json:"attributes,omitempty"`
	Description   string      `json:"description,omitempty"`
	ApiResources  []string    `json:"apiResources,omitempty"`
	Proxies       []string    `json:"proxies,omitempty"`
	Quota         string      `json:"quota,omitempty"`
	QuotaInterval string      `json:"quotaInterval,omitempty"`
	QuotaTimeUnit string      `json:"quotaTimeUnit,omitempty"`
	Scopes        []string    `json:"scopes,omitempty"`
	Environments  []string    `json:"environments,omitempty"`
}

func (s *ProductsServiceOp) Get(name string) (*Product, *Response, error) {

	path := path.Join("apiproducts", name)

	req, e := s.client.NewRequest("GET", path, nil, "")
	if e != nil {
		return nil, nil, e
	}
	returnedProduct := Product{}
	resp, e := s.client.Do(req, &returnedProduct)
	if e != nil {
		return nil, resp, e
	}
	return &returnedProduct, resp, e

}

func (s *ProductsServiceOp) Create(product Product) (*Product, *Response, error) {

	return postOrPutProduct(product, "POST", s)

}

func (s *ProductsServiceOp) Update(product Product) (*Product, *Response, error) {

	return postOrPutProduct(product, "PUT", s)

}

func (s *ProductsServiceOp) Delete(name string) (*Response, error) {

	path := path.Join("apiproducts", name)

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

func postOrPutProduct(product Product, opType string, s *ProductsServiceOp) (*Product, *Response, error) {

	uripath := ""

	if opType == "PUT" {
		uripath = path.Join("apiproducts", product.Name)
	} else {
		uripath = path.Join("apiproducts")
	}

	req, e := s.client.NewRequest(opType, uripath, product, "")
	if e != nil {
		return nil, nil, e
	}

	returnedProduct := Product{}

	resp, e := s.client.Do(req, &returnedProduct)
	if e != nil {
		return nil, resp, e
	}

	return &returnedProduct, resp, e

}
