package apigee

import (
  "path"
  "errors"
)

const productsPath = "apiproducts"

// ProductsService is an interface for interfacing with the Apigee Edge Admin API
// dealing with apiproducts.
type ProductsService interface {
  List() ([]string, *Response, error)
  Get(string) (*ApiProduct, *Response, error)
  Create(ApiProduct) (*ApiProduct, *Response, error)
  Update(ApiProduct) (*ApiProduct, *Response, error)
  Delete(string) (*ApiProduct, *Response, error)
}

type ProductsServiceOp struct {
  client *EdgeClient
}

var _ ProductsService = &ProductsServiceOp{}

// ApiProduct contains information about an API Product within an Edge organization.
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

func reallyUpdateProduct(s ProductsServiceOp, product ApiProduct) (*ApiProduct, *Response, error) {
  path := path.Join(productsPath, product.Name)
  req, e := s.client.NewRequest("POST", path, product)
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


func (s *ProductsServiceOp) Update(product ApiProduct) (*ApiProduct, *Response, error) {
	if product.Name == "" {
    return nil, nil, errors.New("must specify Name of ApiProduct to update")
	}

	if product.ApprovalType == "" || product.DisplayName == "" || product.Environments == nil {
		// The request is lacking some required information.
		// Must get the apiproduct first, to fill in these "required" parameters.
    retrievedProduct, resp, e := s.Get(product.Name)
		if e != nil {
			return nil, resp, e
		}
		if product.ApprovalType == "" {
			product.ApprovalType = retrievedProduct.ApprovalType
		}
		if product.DisplayName == "" {
			product.DisplayName = retrievedProduct.DisplayName
		}
		if product.Environments == nil {
			product.Environments = retrievedProduct.Environments
		}
	}

	// We have all required information...
	// If the caller has omitted the list of api proxies from the product,
	// this call will update the product to have no proxies!  Likewise
	// attributes.
	return reallyUpdateProduct(*s, product);
}


func (s *ProductsServiceOp) Create(product ApiProduct) (*ApiProduct, *Response, error) {
  req, e := s.client.NewRequest("POST", productsPath, product)
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


func (s *ProductsServiceOp) Delete(productName string) (*ApiProduct, *Response, error) {
  path := path.Join(productsPath, productName)
  req, e := s.client.NewRequest("DELETE", path, nil)
  if e != nil {
    return nil, nil, e
  }
  deletedProduct := ApiProduct{}
  resp, e := s.client.Do(req, &deletedProduct)
  if e != nil {
    return nil, resp, e
  }
  return &deletedProduct, resp, e
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
func (s *ProductsServiceOp) Get(productName string) (*ApiProduct, *Response, error) {
  path := path.Join(productsPath, productName)
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
