package apigee

import (
  "archive/zip"
  "strings"
  "path"
  "path/filepath"
  "os"
  "fmt"
  "net/url"
  "io"
  "io/ioutil"
  "errors"
)

const proxiesPath = "apis"


// ProxiesService is an interface for interfacing with the Apigee Edge Admin API
// dealing with apiproxies. 
type ProxiesService interface {
  List() ([]string, *Response, error)
  Get(string) (*Proxy, *Response, error)
  Import(string, string) (*ProxyRevision, *Response, error)
  Delete(string) (*DeletedProxyInfo, *Response, error)
  DeleteRevision(string, Revision) (*ProxyRevision, *Response, error)
  Deploy(string,string,Revision) (*ProxyRevisionDeployment, *Response, error)
  Undeploy(string,string,Revision) (*ProxyRevisionDeployment, *Response, error)
}

type ProxiesServiceOp struct {
  client *EdgeClient
}

var _ ProxiesService = &ProxiesServiceOp{}

type Proxy struct {
  Revisions   []Revision    `json:"revision,omitempty"`
  Name        string        `json:"name,omitempty"`
  MetaData    ProxyMetadata `json:"metaData,omitempty"`
}

type ProxyMetadata struct {
  LastModifiedBy  string     `json:"lastModifiedBy,omitempty"`
  CreatedBy       string     `json:"createdBy,omitempty"`
  LastModifiedAt  Timestamp  `json:"lastModifiedAt,omitempty"`
  CreatedAt       Timestamp  `json:"createdAt,omitempty"`
}

type ProxyRevision struct {
  CreatedBy       string     `json:"createdBy,omitempty"`
  CreatedAt       Timestamp  `json:"createdAt,omitempty"`
  Description     string     `json:"description,omitempty"`
  ContextInfo     string     `json:"contextInfo,omitempty"`
  DisplayName     string     `json:"displayName,omitempty"`
  Name            string     `json:"name,omitempty"`
  LastModifiedBy  string     `json:"lastModifiedBy,omitempty"`
  LastModifiedAt  Timestamp  `json:"lastModifiedAt,omitempty"`
  Revision        Revision   `json:"revision,omitempty"`
  TargetEndpoints []string   `json:"targetEndpoints,omitempty"`
  TargetServers   []string   `json:"targetServers,omitempty"`
  Resources       []string   `json:"resources,omitempty"`
  ProxyEndpoints  []string   `json:"proxyEndpoints,omitempty"`
  Policies        []string   `json:"policies,omitempty"`
  Type            string     `json:"type,omitempty"`
}

// {
//   "createdAt" : 1473206030269,
//   "createdBy" : "DChiesa@apigee.com",
//   "description" : "",
//   "displayName" : "ramnath-1",
//   "name" : "proxyname1",
//   "lastModifiedAt" : 1473206030269,
//   "lastModifiedBy" : "DChiesa@apigee.com",
//   "revision" : "1",
//   "targetEndpoints" : [ "default" ],
//   "targetServers" : [ ],
//   "resources" : [ "jsc://injectHeader.js", "jsc://maybeFormatFault.js", "jsc://setTargetUrl.js" ],
//   "proxyEndpoints" : [ "default" ],
//   "policies" : [ "AM-CleanResponseHeaders", "JS-InjectHeader", "JS-MaybeFormatFault", "JS-SetTargetUrl", "RF-UnknownRequest" ],
//   "contextInfo" : "Revision 1 of application proxyname1, in organization cap500",
//   "type" : "Application",
//
//   "resourceFiles" : {
//     "resourceFile" : [ {
//       "name" : "injectHeader.js",
//       "type" : "jsc"
//     }, {
//       "name" : "maybeFormatFault.js",
//       "type" : "jsc"
//     }, {
//       "name" : "setTargetUrl.js",
//       "type" : "jsc"
//     } ]
//   },
//
//   "configurationVersion" : {
//     "majorVersion" : 4,
//     "minorVersion" : 0
//   }
// }  


type ProxyRevisionDeployment struct {
  Name            string        `json:"aPIProxy,omitempty"`
  Revision        Revision      `json:"revision,omitempty"`
  Environment     string        `json:"environment,omitempty"`
  Organization    string        `json:"organization,omitempty"`
  State           string        `json:"state,omitempty"`
  Servers         []EdgeServer  `json:"server,omitempty"`
}

type EdgeServer struct {
  Status          string        `json:"status,omitempty"`
  Uuid            string        `json:"uUID,omitempty"`
  Type            string        `json:"type,omitempty"`
}

// {
//   "aPIProxy" : "ramnath-1",
//   "name" : "6",
//   "environment" : "test",
//   "organization" : "cap500",
//   "revision" : "6",
//   "state" : "undeployed"
//   "server" : [ {
//     "status" : "undeployed",
//     "type" : [ "message-processor" ],
//     "uUID" : "a4850e3b-6ce9-482a-9521-d9869be8482e"
//   }, {
//     "status" : "undeployed",
//     "type" : [ "router" ],
//     "uUID" : "f0a80e0e-572c-46ad-b064-acac9ed1d870"
//   } ],
//
//   "configuration" : {
//     "basePath" : "/",
//     "steps" : [ ]
//   },
// }



type DeletedProxyInfo struct {
  Name   string   `json:"name,omitempty"`
}
// {
//   "configurationVersion" : {
//     "majorVersion" : 4,
//     "minorVersion" : 0
//   },
//   "contextInfo" : "Revision null of application -NA-, in organization -NA-",
//   "name" : "ramnath-1",
//   "policies" : [ ],
//   "proxyEndpoints" : [ ],
//   "resourceFiles" : {
//     "resourceFile" : [ ]
//   },
//   "resources" : [ ],
//   "targetEndpoints" : [ ],
//   "targetServers" : [ ],
//   "type" : "Application"
// }


type proxiesRoot struct {
  Proxies []Proxy `json:"proxies"`
}


func (s *ProxiesServiceOp) List() ([]string, *Response, error) {
  req, e := s.client.NewRequest("GET", proxiesPath, nil)
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



func (s *ProxiesServiceOp) Get(proxy string) (*Proxy, *Response, error) {
  path := path.Join(proxiesPath, proxy)
  req, e := s.client.NewRequest("GET", path, nil)
  if e != nil {
    return nil, nil, e
  }
  returnedProxy := Proxy{}
  resp, e := s.client.Do(req, &returnedProxy)
  if e != nil {
    return nil, resp, e
  }
  return &returnedProxy, resp, e
}

// func (s *ProxiesServiceOp) ListExpanded() ([]Proxy, *Response, error) {
//   root := new(proxiesRoot)
//   origURL, err := url.Parse(proxiesPath)
//   if err != nil {
//     return root.Proxies, nil, err
//   }
//   q := origURL.Query()
//   q.Add("expand", "true")
//   origURL.RawQuery = q.Encode()
//   path := origURL.String()
// 
//   req, e := s.client.NewRequest("GET", path, nil)
//   if e != nil {
//     return nil, nil, e
//   }
// 
//   resp, e := s.client.Do(req, &root)
//   if e != nil {
//     return nil, resp, e
//   }
//   return root.Proxies, resp, e
// }


func smartFilter(path string) bool {
  if strings.HasSuffix(path, "~") {
    return false
  }
  if strings.HasSuffix(path, "#") && strings.HasPrefix(path, "#") {
    return false
  }
  return true
}


func zipDirectory (source string, target string, filter func(string) bool) error {
  zipfile, err := os.Create(target)
  if err != nil {
    return err
  }
  defer zipfile.Close()

  archive := zip.NewWriter(zipfile)
  defer archive.Close()

  info, err := os.Stat(source)
  if err != nil {
    return nil
  }

  var baseDir string
  if info.IsDir() {
    baseDir = filepath.Base(source)
  }

  filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
    if filter != nil && filter(path) {
      if err != nil {
        return err
      }

      header, err := zip.FileInfoHeader(info)
      if err != nil {
        return err
      }

      if baseDir != "" {
        header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
      }

      // This archive will be unzipped by a Java process.  When ZIP64 extensions
      // are used, Java insists on having Deflate as the compression method (0x08)
      // even for directories.
      header.Method = zip.Deflate
      
      if info.IsDir() {
        header.Name += "/"
      } 

      writer, err := archive.CreateHeader(header)
      if err != nil {
        return err
      }

      if info.IsDir() {
        return nil
      }

      file, err := os.Open(path)
      if err != nil {
        return err
      }
      defer file.Close()
      _, err = io.Copy(writer, file)
    }
    return err
  })

  return err
}



func (s *ProxiesServiceOp) Import(proxyName string, source string) (*ProxyRevision, *Response, error) {
  info, err := os.Stat(source)
  if err != nil {
    return nil, nil, err
  }
  zipfileName := source
  if info.IsDir() {
    // create a temporary zip file
    if proxyName == "" {
      proxyName = filepath.Base(source) 
    }
    tempDir, e := ioutil.TempDir("", "golang-") 
    if e != nil {
      return nil, nil, errors.New(fmt.Sprintf("while creating temp dir, error: %#v", e))
    }
    zipfileName = path.Join(tempDir, "apiproxy.zip")
    e = zipDirectory (path.Join(source, "apiproxy"), zipfileName, smartFilter)
    if e != nil {
      return nil, nil, errors.New(fmt.Sprintf("while creating temp dir, error: %#v", e))
    }
    fmt.Printf("zipped %s into %s\n\n", source, zipfileName)
  }

  
  if !strings.HasSuffix(zipfileName,".zip") {
    return nil, nil, errors.New("source must be a zipfile")
  }
  
  info, err = os.Stat(zipfileName)
  if err != nil {
    return nil, nil, err
  }

  // append the query params
  origURL, err := url.Parse(proxiesPath)
  if err != nil {
     return nil, nil, err
  }
  q := origURL.Query()
  q.Add("action", "import")
  q.Add("name", proxyName)
  origURL.RawQuery = q.Encode()
  path := origURL.String()
  
  ioreader, err := os.Open(zipfileName)
  if err != nil {
     return nil, nil, err
  }
  defer ioreader.Close()
  
  req, e := s.client.NewRequest("POST", path, ioreader)
  if e != nil {
    return nil, nil, e
  }
  returnedProxyRevision := ProxyRevision{}
  resp, e := s.client.Do(req, &returnedProxyRevision)
  if e != nil {
    return nil, resp, e
  }
  return &returnedProxyRevision, resp, e
}


func (s *ProxiesServiceOp) DeleteRevision(proxyName string, rev Revision) (*ProxyRevision, *Response, error) {
  path := path.Join(proxiesPath, proxyName, "revisions", fmt.Sprintf("%d",rev))
  req, e := s.client.NewRequest("DELETE", path, nil)
  if e != nil {
    return nil, nil, e
  }
  proxyRev := ProxyRevision{}
  resp, e := s.client.Do(req, &proxyRev)
  if e != nil {
    return nil, resp, e
  }
  return &proxyRev, resp, e
}

func (s *ProxiesServiceOp) Undeploy(proxyName, env string, rev Revision) (*ProxyRevisionDeployment, *Response, error) {
  path := path.Join(proxiesPath, proxyName, "revisions", fmt.Sprintf("%d",rev), "deployments")
  // append the query params
  origURL, err := url.Parse(path)
  if err != nil {
     return nil, nil, err
  }
  q := origURL.Query()
  q.Add("action", "undeploy")
  q.Add("env", env)
  origURL.RawQuery = q.Encode()
  path = origURL.String()
  
  req, e := s.client.NewRequest("POST", path, nil)
  if e != nil {
    return nil, nil, e
  }
  
  deployment := ProxyRevisionDeployment{}
  resp, e := s.client.Do(req, &deployment)
  if e != nil {
    return nil, resp, e
  }
  return &deployment, resp, e
}


func (s *ProxiesServiceOp) Deploy(proxyName, env string, rev Revision) (*ProxyRevisionDeployment, *Response, error) {
  path := path.Join(proxiesPath, proxyName, "revisions", fmt.Sprintf("%d",rev), "deployments")
  // append the query params
  origURL, err := url.Parse(path)
  if err != nil {
     return nil, nil, err
  }
  q := origURL.Query()
  q.Add("action", "deploy")
  q.Add("override", "true")
  q.Add("delay", "60")
  q.Add("env", env)
  origURL.RawQuery = q.Encode()
  path = origURL.String()

  req, e := s.client.NewRequest("POST", path, nil)
  if e != nil {
    return nil, nil, e
  }
  
  deployment := ProxyRevisionDeployment{}
  resp, e := s.client.Do(req, &deployment)
  if e != nil {
    return nil, resp, e
  }
  return &deployment, resp, e
}


func (s *ProxiesServiceOp) Delete(proxyName string) (*DeletedProxyInfo, *Response, error) {
  path := path.Join(proxiesPath, proxyName)
  req, e := s.client.NewRequest("DELETE", path, nil)
  if e != nil {
    return nil, nil, e
  }
  proxy := DeletedProxyInfo{}
  resp, e := s.client.Do(req, &proxy)
  if e != nil {
    return nil, resp, e
  }
  return &proxy, resp, e
}
