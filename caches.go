package apigee

import (
	"path"
)

const cachesPath = "caches"

// CachesService is an interface for interfacing with the Apigee Edge Admin API
// dealing with caches.
type CachesService interface {
	List(string) ([]string, *Response, error)
	Get(string, string) (*Cache, *Response, error)
}

type CachesServiceOp struct {
	client *ApigeeClient
}

var _ CachesService = &CachesServiceOp{}

// Cache contains information about a cache within an Edge organization.
type Cache struct {
	Name                string      `json:"name,omitempty"`
	Description         string      `json:"description,omitempty"`
	OverflowToDisk      bool        `json:"overflowToDisk,omitempty"`
	Persistent          bool        `json:"persistent,omitempty"`
	Distributed         bool        `json:"distributed,omitempty"`
	DiskSizeInMB        int         `json:"diskSizeInMB,omitempty"`
	InMemorySizeInKB    int         `json:"inMemorySizeInKB,omitempty"`
	MaxElementsInMemory int         `json:"maxElementsInMemory,omitempty"`
	MaxElementsOnDisk   int         `json:"maxElementsOnDisk,omitempty"`
	Expiry              CacheExpiry `json:"expirySettings,omitempty"`
}

// List retrieves the list of cache names for the organization referred by the ApigeeClient,
// or a set of cache names for a specific environment within an organization.
func (s *CachesServiceOp) List(env string) ([]string, *Response, error) {
	var p1 string
	if env == "" {
		p1 = cachesPath
	} else {
		p1 = path.Join("e", env, cachesPath)
	}
	req, e := s.client.NewRequest("GET", p1, nil)
	if e != nil {
		return nil, nil, e
	}
	namelist := make([]string, 0)
	resp, e := s.client.Do(req, &namelist)
	if e != nil {
		return nil, resp, e
	}
	return namelist, resp, e
}

// Get retrieves the information about a Cache in an organization, or about a
// cache in an environment within an organization. This information includes the
// properties, and the created and last modified details.
func (s *CachesServiceOp) Get(name, env string) (*Cache, *Response, error) {
	var p1 string
	if env == "" {
		p1 = path.Join(cachesPath, env)
	} else {
		p1 = path.Join("e", env, cachesPath)
	}
	req, e := s.client.NewRequest("GET", p1, nil)
	if e != nil {
		return nil, nil, e
	}
	returnedCache := Cache{}
	resp, e := s.client.Do(req, &returnedCache)
	if e != nil {
		return nil, resp, e
	}
	return &returnedCache, resp, e
}
