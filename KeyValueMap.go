package apigee

import "path"

// KeyValueMapService is an interface for interfacing with the Apigee Edge Admin API
// dealing with KeyValueMap.
type KeyValueMapService interface {
	Get(string, string) (*KeyValueMap, *Response, error)
	Create(KeyValueMap, string) (*KeyValueMap, *Response, error)
	Delete(string, string) (*Response, error)
	//update is not implemented as the API is being deprectated. See KeyValueMapEntry.
	//	Update(KeyValueMap, string) (*KeyValueMap, *Response, error)
}

// KeyValueMapServiceOp holds creds
type KeyValueMapServiceOp struct {
	client *EdgeClient
}

var _ KeyValueMapService = &KeyValueMapServiceOp{}

// EntryStruct Holds the Key value map entry
type EntryStruct struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// KeyValueMap Holds the Key value map
type KeyValueMap struct {
	Name      string        `json:"name,omitempty"`
	Encrypted bool          `json:"encrypted,omitempty"`
	Entry     []EntryStruct `json:"entry,omitempty"`
}

// Get the Keyvaluemap
func (s *KeyValueMapServiceOp) Get(name string, env string) (*KeyValueMap, *Response, error) {

	path := path.Join("environments", env, "keyvaluemaps", name)

	req, e := s.client.NewRequest("GET", path, nil, "")
	if e != nil {
		return nil, nil, e
	}
	returnedKeyValueMap := KeyValueMap{}
	resp, e := s.client.Do(req, &returnedKeyValueMap)
	if e != nil {
		return nil, resp, e
	}
	return &returnedKeyValueMap, resp, e

}

// Create a new key value map
func (s *KeyValueMapServiceOp) Create(keyValueMap KeyValueMap, env string) (*KeyValueMap, *Response, error) {

	return postOrPutKeyValueMap(keyValueMap, env, "POST", s)
}

// Update an existing key value map
//func (s *KeyValueMapServiceOp) Update(keyValueMap KeyValueMap, env string) (*KeyValueMap, *Response, error) {
//	return postOrPutKeyValueMap(keyValueMap, env, "PUT", s)
//}

// Delete an existing key value map
func (s *KeyValueMapServiceOp) Delete(name string, env string) (*Response, error) {

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

func postOrPutKeyValueMap(keyValueMap KeyValueMap, env string, opType string, s *KeyValueMapServiceOp) (*KeyValueMap, *Response, error) {

	uripath := ""

	if opType == "PUT" {
		uripath = path.Join("environments", env, "keyvaluemaps", keyValueMap.Name)
	} else {
		uripath = path.Join("environments", env, "keyvaluemaps")
	}

	req, e := s.client.NewRequest(opType, uripath, keyValueMap, "")
	if e != nil {
		return nil, nil, e
	}

	returnedKeyValueMap := KeyValueMap{}

	resp, e := s.client.Do(req, &returnedKeyValueMap)
	if e != nil {
		return nil, resp, e
	}

	return &returnedKeyValueMap, resp, e

}
