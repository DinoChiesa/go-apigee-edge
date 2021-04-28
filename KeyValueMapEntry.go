package apigee

import "path"

// KeyValueMapEntryService is an interface for interfacing with the Apigee Edge Admin API
// dealing with KeyValueMapEntry.
type KeyValueMapEntryService interface {
	Get(string, string, string) (*KeyValueMapEntry, *Response, error)
	Create(string, KeyValueMapEntryKeys, string) (*KeyValueMapEntry, *Response, error)
	Delete(string, string, string) (*Response, error)
	Update(string, KeyValueMapEntryKeys, string) (*KeyValueMapEntry, *Response, error)
}

// KeyValueMapEntryServiceOp holds creds
type KeyValueMapEntryServiceOp struct {
	client *EdgeClient
}

var _ KeyValueMapEntryService = &KeyValueMapEntryServiceOp{}

// KeyValueMapEntryKeys to update
type KeyValueMapEntryKeys struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// KeyValueMapEntry Holds the Key value map
type KeyValueMapEntry struct {
	KVMName string                 `json:"kvmName,omitempty"`
	Entry   []KeyValueMapEntryKeys `json:"entry,omitempty"`
}

// Get the key value map entry
func (s *KeyValueMapEntryServiceOp) Get(keyValueMapName string, env string, keyValueMapEntry string) (*KeyValueMapEntry, *Response, error) {

	path := path.Join("environments", env, "keyvaluemaps", keyValueMapName, "entries", keyValueMapEntry)

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

// Create a new key value map entry
func (s *KeyValueMapEntryServiceOp) Create(keyValueMapName string, keyValueMapEntry KeyValueMapEntryKeys, env string) (*KeyValueMapEntry, *Response, error) {

	return postOrPutKeyValueMapEntry(keyValueMapName, keyValueMapEntry, env, "POST", s)
}

// Update an existing key value map entry
func (s *KeyValueMapEntryServiceOp) Update(keyValueMapName string, keyValueMapEntry KeyValueMapEntryKeys, env string) (*KeyValueMapEntry, *Response, error) {

	return postOrPutKeyValueMapEntry(keyValueMapName, keyValueMapEntry, env, "PUT", s)

}

// Delete an existing key value map entry
func (s *KeyValueMapEntryServiceOp) Delete(keyValueMapEntry string, keyValueMapName string, env string) (*Response, error) {

	path := path.Join("environments", env, "keyvaluemaps", keyValueMapName, "entries", keyValueMapEntry)

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

func postOrPutKeyValueMapEntry(keyValueMapName string, keyValueMapEntry KeyValueMapEntryKeys, env string, opType string, s *KeyValueMapEntryServiceOp) (*KeyValueMapEntry, *Response, error) {

	uripath := ""

	if opType == "PUT" {
		uripath = path.Join("environments", env, "keyvaluemaps", keyValueMapName, "entries", keyValueMapEntry.Name)
	} else {
		uripath = path.Join("environments", env, "keyvaluemaps", keyValueMapName, "entries")
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
