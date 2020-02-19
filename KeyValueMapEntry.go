package apigee

import "path"

// KeyValueMapEntryService is an interface for interfacing with the Apigee Edge Admin API
// dealing with KeyValueMapEntry.
type KeyValueMapEntryService interface {
	Get(string, string) (*KeyValueMapEntry, *Response, error)
	Create(KeyValueMapEntry, string) (*KeyValueMapEntry, *Response, error)
	Delete(string, string) (*Response, error)
	Update(KeyValueMapEntry, string) (*KeyValueMapEntry, *Response, error)
}

// KeyValueMapEntryServiceOp holds creds
type KeyValueMapEntryServiceOp struct {
	client *EdgeClient
}

var _ KeyValueMapEntryService = &KeyValueMapEntryServiceOp{}

// KeyValueMapEntry Holds the Key value map
type KeyValueMapEntry struct {
	Name  string `json:"name,omitempty"`
	Value bool   `json:"value,omitempty"`
}

// Get the Keyvaluemap
func (s *KeyValueMapEntryServiceOp) Get(name string, env string) (*KeyValueMapEntry, *Response, error) {

	path := path.Join("environments", env, "keyvaluemaps", name)

	req, e := s.client.NewRequest("GET", path, nil, "")
	if e != nil {
		return nil, nil, e
	}
	returnedKeyValueMapEntry := KeyValueMapEntry{}
	resp, e := s.client.Do(req, &returnedKeyValueMapEntry)
	if e != nil {
		return nil, resp, e
	}
	return &returnedKeyValueMapEntry, resp, e

}

// Create a new key value map
func (s *KeyValueMapEntryServiceOp) Create(keyValueMapEntry KeyValueMapEntry, env string) (*KeyValueMapEntry, *Response, error) {

	return postOrPutKeyValueMapEntry(keyValueMapEntry, env, "POST", s)
}

// Update an existing key value map
func (s *KeyValueMapEntryServiceOp) Update(keyValueMapEntry KeyValueMapEntry, env string) (*KeyValueMapEntry, *Response, error) {

	return postOrPutKeyValueMapEntry(keyValueMapEntry, env, "PUT", s)

}

// Delete an existing key value map
func (s *KeyValueMapEntryServiceOp) Delete(name string, env string) (*Response, error) {

	path := path.Join("environments", env, "keyvaluemaps", name)

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

func postOrPutKeyValueMapEntry(keyValueMapEntry KeyValueMapEntry, env string, opType string, s *KeyValueMapEntryServiceOp) (*KeyValueMapEntry, *Response, error) {

	uripath := ""

	if opType == "PUT" {
		uripath = path.Join("environments", env, "keyvaluemaps", keyValueMapEntry.Name)
	} else {
		uripath = path.Join("environments", env, "keyvaluemaps")
	}

	req, e := s.client.NewRequest(opType, uripath, keyValueMapEntry, "")
	if e != nil {
		return nil, nil, e
	}

	returnedKeyValueMapEntry := KeyValueMapEntry{}

	resp, e := s.client.Do(req, &returnedKeyValueMapEntry)
	if e != nil {
		return nil, resp, e
	}

	return &returnedKeyValueMapEntry, resp, e

}
