package apigee

import (
	"path"
)

// TargetServersService is an interface for interfacing with the Apigee Edge Admin API
// dealing with target servers.
type TargetServersService interface {
	Get(string, string) (*TargetServer, *Response, error)
	Create(TargetServer, string) (*TargetServer, *Response, error)
	Delete(string, string) (*Response, error)
	Update(TargetServer, string) (*TargetServer, *Response, error)
}

type TargetServersServiceOp struct {
	client *EdgeClient
}

var _ TargetServersService = &TargetServersServiceOp{}

type TargetServer struct {
	Name    string   `json:"name,omitempty"`
	Host    string   `json:"host,omitempty"`
	Enabled bool     `json:"isEnabled"`
	Port    int      `json:"port,omitempty"`
	SSLInfo *SSLInfo `json:"sSLInfo,omitempty"`
}

// For some reason Apigee returns SOME bools as strings and others a bools.
type SSLInfo struct {
	SSLEnabled             string   `json:"enabled,omitempty"`
	ClientAuthEnabled      string   `json:"clientAuthEnabled,omitempty"`
	KeyStore               string   `json:"keyStore,omitempty"`
	TrustStore             string   `json:"trustStore,omitempty"`
	KeyAlias               string   `json:"keyAlias,omitempty"`
	Ciphers                []string `json:"ciphers,omitempty"`
	IgnoreValidationErrors bool     `json:"ignoreValidationErrors"`
	Protocols              []string `json:"protocols,omitempty"`
}

func (s *TargetServersServiceOp) Get(name string, env string) (*TargetServer, *Response, error) {

	path := path.Join("environments", env, "targetservers", name)

	req, e := s.client.NewRequest("GET", path, nil, "")
	if e != nil {
		return nil, nil, e
	}
	returnedTargetServer := TargetServer{}
	resp, e := s.client.Do(req, &returnedTargetServer)
	if e != nil {
		return nil, resp, e
	}
	return &returnedTargetServer, resp, e

}

func (s *TargetServersServiceOp) Create(targetServer TargetServer, env string) (*TargetServer, *Response, error) {

	return postOrPutTargetServer(targetServer, env, "POST", s)

}

func (s *TargetServersServiceOp) Update(targetServer TargetServer, env string) (*TargetServer, *Response, error) {

	return postOrPutTargetServer(targetServer, env, "PUT", s)

}

func (s *TargetServersServiceOp) Delete(name string, env string) (*Response, error) {

	path := path.Join("environments", env, "targetservers", name)

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

func postOrPutTargetServer(targetServer TargetServer, env string, opType string, s *TargetServersServiceOp) (*TargetServer, *Response, error) {

	uripath := ""

	if opType == "PUT" {
		uripath = path.Join("environments", env, "targetservers", targetServer.Name)
	} else {
		uripath = path.Join("environments", env, "targetservers")
	}

	req, e := s.client.NewRequest(opType, uripath, targetServer, "")
	if e != nil {
		return nil, nil, e
	}

	returnedTargetServer := TargetServer{}

	resp, e := s.client.Do(req, &returnedTargetServer)
	if e != nil {
		return nil, resp, e
	}

	return &returnedTargetServer, resp, e

}
