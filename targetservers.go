package apigee

import (
	"errors"
	"fmt"
	"path"
)

const targetserversPath = "targetservers"

// TargetserversService is an interface for interfacing with the Apigee Edge Admin API
// dealing with Target servers
type TargetserversService interface {
	List(string) ([]string, *Response, error)
	Get(string, string) (*TargetServer, *Response, error)
	Create(TargetServer, string) (*TargetServer, *Response, error)
	Update(TargetServer, string) (*TargetServer, *Response, error)
	Delete(string, string) (*TargetServer, *Response, error)
}

// TargetserversServiceOp represents the target server service used to
// communicate with Apigee
type TargetserversServiceOp struct {
	client *ApigeeClient
}

var _ TargetserversService = &TargetserversServiceOp{}

// TargetServer contains information about a Target Server withing an environment in an Edge Organization
type TargetServer struct {
	Name      string `json:"name"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	IsEnabled bool   `json:"isEnabled"`
	SSLInfo   struct {
		Ciphers                []interface{} `json:"ciphers,omitempty"`
		ClientAuthEnabled      bool          `json:"clientAuthEnabled,omitempty"`
		Enabled                bool          `json:"enabled,omitempty"`
		IgnoreValidationErrors bool          `json:"ignoreValidationErrors,omitempty"`
		KeyAlias               string        `json:"keyAlias,omitempty"`
		KeyStore               string        `json:"keyStore,omitempty"`
		Protocols              []interface{} `json:"protocols,omitempty"`
		TrustStore             string        `json:"trustStore,omitempty"`
	} `json:"SSLInfo,omitempty"`
}

//List retrieves the list of Target servers from a specific environment env
func (s *TargetserversServiceOp) List(env string) ([]string, *Response, error) {
	var p1 string
	p1 = path.Join("e", env, targetserversPath)
	req, e := s.client.NewRequest("GET", p1, nil)
	if e != nil {
		return nil, nil, e
	}
	targetserverlist := make([]string, 0)
	resp, e := s.client.Do(req, &targetserverlist)
	if e != nil {
		return nil, resp, e
	}
	return targetserverlist, resp, e
}

//Get retrieves a specific Target server from a specific environment env
func (s *TargetserversServiceOp) Get(name, env string) (*TargetServer, *Response, error) {
	var p1 string
	p1 = path.Join("e", env, targetserversPath, name)
	req, e := s.client.NewRequest("GET", p1, nil)
	if e != nil {
		return nil, nil, e
	}
	returnedtargetserver := TargetServer{}
	resp, e := s.client.Do(req, &returnedtargetserver)
	if e != nil {
		return nil, resp, e
	}
	return &returnedtargetserver, resp, e
}

//Create creates a new target server in the given environment
func (s *TargetserversServiceOp) Create(targetserver TargetServer, env string) (*TargetServer, *Response, error) {
	var p1 string

	p1 = path.Join("e", env, targetserversPath)

	fmt.Printf("argumnets are %v, %v", targetserver, env)
	req, e := s.client.NewRequest("POST", p1, targetserver)
	if e != nil {
		return nil, nil, e
	}

	returnedtargetserver := TargetServer{}
	resp, e := s.client.Do(req, &returnedtargetserver)
	if e != nil {
		return nil, resp, e
	}
	return &returnedtargetserver, resp, e
}

//Update updates a target server in the given environment
func (s *TargetserversServiceOp) Update(targetserver TargetServer, env string) (*TargetServer, *Response, error) {
	var p1 string

	if targetserver.Name == "" || targetserver.Host == "" || targetserver.Port == 0 {
		return nil, nil, errors.New("Must specify the name, host and port of the target server to update")
	}

	p1 = path.Join("e", env, targetserversPath, targetserver.Name)
	req, e := s.client.NewRequest("PUT", p1, targetserver)
	if e != nil {
		return nil, nil, e
	}

	returnedtargetserver := TargetServer{}
	resp, e := s.client.Do(req, &returnedtargetserver)
	if e != nil {
		return nil, resp, e
	}
	return &returnedtargetserver, resp, e
}

//Delete deletes a target server in the given environment
func (s *TargetserversServiceOp) Delete(name, env string) (*TargetServer, *Response, error) {
	var p1 string
	p1 = path.Join("e", env, targetserversPath, name)
	req, e := s.client.NewRequest("DELETE", p1, nil)
	deletedtargetserver := TargetServer{}

	resp, e := s.client.Do(req, &deletedtargetserver)

	if e != nil {
		return nil, resp, e
	}
	return &deletedtargetserver, resp, e
}
