package apigee

import (
  "path"
)

const productsPath = "apiproducts"

// ProductsService is an interface for interfacing with the Apigee Edge Admin API
// dealing with apiproducts.
type ProductsService interface {
  List() ([]string, *Response, error)
  Get(string) (*ApiProduct, *Response, error)
}

type ProductsServiceOp struct {
  client *EdgeClient
}

var _ ProductsService = &ProductsServiceOp{}

// Proxy contains information about an API Proxy within an Edge organization.
type ApiProduct struct {
  Name            string      `json:"name,omitempty"`
  ApiResources    []string    `json:"apiResources,omitempty"`
  ApprovalType    string      `json:"approvalType,omitempty"`
  Attributes      Attributes  `json:"attributes,omitempty"`
  CreatedBy       string      `json:"createdBy,omitempty"`
  CreatedAt       Timestamp   `json:"createdAt,omitempty"`
  Description     string      `json:"description,omitempty"`
  DisplayName     string      `json:"displayName,omitempty"`
  LastModifiedBy  string      `json:"lastModifiedBy,omitempty"`
  LastModifiedAt  Timestamp   `json:"lastModifiedAt,omitempty"`
  Environments    []string    `json:"environments,omitempty"`
  Proxies         []string    `json:"proxies,omitempty"`
  Scopes          []string    `json:"scopes,omitempty"`
}

// List retrieves the list of apiproduct names for the organization referred by the EdgeClient.
func (s *ProductsServiceOp) List() ([]string, *Response, error) {
  req, e := s.client.NewRequest("GET", productsPath, nil)
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

// Get retrieves the information about an API Product in an organization, information including
// the list of API Proxies, the scopes, the quota, and other attributes. 
func (s *ProductsServiceOp) Get(product string) (*ApiProduct, *Response, error) {
  path := path.Join(productsPath, product)
  req, e := s.client.NewRequest("GET", path, nil)
  if e != nil {
    return nil, nil, e
  }
  returnedProduct := ApiProduct{}
  resp, e := s.client.Do(req, &returnedProduct)
  if e != nil {
    return nil, resp, e
  }
  return &returnedProduct, resp, e
}

