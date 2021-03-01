package apigee

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// DeployableAsset contains information about an API Proxy or SharedFlow within an Apigee organization.
type DeployableAsset struct {
	Revisions []Revision         `json:"revision,omitempty"`
	Name      string             `json:"name,omitempty"`
	MetaData  DeployableMetadata `json:"metaData,omitempty"`
}

// DeployableRevision holds information about a revision of an API Proxy, or a SharedFlow.
type DeployableRevision struct {
	Name            string    `json:"name,omitempty"`
	DisplayName     string    `json:"displayName,omitempty"`
	Revision        Revision  `json:"revision,omitempty"`
	CreatedBy       string    `json:"createdBy,omitempty"`
	CreatedAt       Timestamp `json:"createdAt,omitempty"`
	LastModifiedBy  string    `json:"lastModifiedBy,omitempty"`
	LastModifiedAt  Timestamp `json:"lastModifiedAt,omitempty"`
	Description     string    `json:"description,omitempty"`
	ContextInfo     string    `json:"contextInfo,omitempty"`
	TargetEndpoints []string  `json:"targetEndpoints,omitempty"`
	TargetServers   []string  `json:"targetServers,omitempty"`
	Resources       []string  `json:"resources,omitempty"`
	ProxyEndpoints  []string  `json:"proxyEndpoints,omitempty"`
	SharedFlows     []string  `json:"sharedFlows,omitempty"`
	Policies        []string  `json:"policies,omitempty"`
	Type            string    `json:"type,omitempty"`
}

// ProxyMetadata contains information related to the creation and last modified
// time and actor for an API Proxy within an organization.
type DeployableMetadata struct {
	LastModifiedBy string    `json:"lastModifiedBy,omitempty"`
	CreatedBy      string    `json:"createdBy,omitempty"`
	LastModifiedAt Timestamp `json:"lastModifiedAt,omitempty"`
	CreatedAt      Timestamp `json:"createdAt,omitempty"`
}

// When Delete returns successfully, it returns a payload that contains very little useful
// information. This struct deserializes that information.
type DeletedItemInfo struct {
	Name string `json:"name,omitempty"`
}

// When inquiring the deployment status of an API PRoxy revision, even implicitly
// as when performing a Deploy or Undeploy, the response includes the deployment
// status for each particular Edge Server in the environment. This struct
// deserializes that information. It will normally not be useful at all. In rare
// cases, it may be useful in helping to diagnose problems.  For example, if there
// is a problem with a deployment change, as when a Message Processor is
// experiencing a problem and cannot undeploy, or more commonly, cannot deploy an
// API Proxy, this struct will hold relevant information.
type ApigeeServer struct {
	Status string   `json:"status,omitempty"`
	Uuid   string   `json:"uUID,omitempty"`
	Type   []string `json:"type,omitempty"`
}

// Deployment (nee ProxyDeployment) holds information about the deployment state of a
// all revisions of an API Proxy or SharedFlow.
type Deployment struct {
	Environments []EnvironmentDeployment `json:"environment,omitempty"`
	Name         string                  `json:"name,omitempty"`
	Organization string                  `json:"organization,omitempty"`
}

type EnvironmentDeployment struct {
	Name     string               `json:"name,omitempty"`
	Revision []RevisionDeployment `json:"revision,omitempty"`
}

type RevisionDeployment struct {
	Number  Revision       `json:"name,omitempty"`
	State   string         `json:"state,omitempty"`
	Servers []ApigeeServer `json:"server,omitempty"`
}

type Deployable struct{}

func (s *Deployable) List(client *ApigeeClient, uriPathElement string) ([]string, *Response, error) {
	req, e := client.NewRequest("GET", uriPathElement, nil)
	if e != nil {
		return nil, nil, e
	}
	namelist := make([]string, 0)
	resp, e := client.Do(req, &namelist)
	if e != nil {
		return nil, resp, e
	}
	return namelist, resp, e
}

func (s *Deployable) Get(client *ApigeeClient, uriPathElement, assetName string) (*DeployableAsset, *Response, error) {
	path := path.Join(uriPathElement, assetName)
	req, e := client.NewRequest("GET", path, nil)
	if e != nil {
		return nil, nil, e
	}
	returnedAsset := DeployableAsset{}
	resp, e := client.Do(req, &returnedAsset)
	if e != nil {
		return nil, resp, e
	}
	return &returnedAsset, resp, e
}

func smartFilter(path string) bool {
	if strings.HasSuffix(path, "~") {
		return false
	}
	if strings.HasSuffix(path, "#") && strings.HasPrefix(path, "#") {
		return false
	}
	return true
}

func zipDirectory(source string, target string, filter func(string) bool) error {
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
		if filter == nil || filter(path) {
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

func (s *Deployable) Import(client *ApigeeClient, uriPathElement, assetName, source string) (*DeployableRevision, *Response, error) {
	info, err := os.Stat(source)
	if err != nil {
		return nil, nil, err
	}
	zipfileName := source
	if info.IsDir() {
		// create a temporary zip file
		if assetName == "" {
			assetName = filepath.Base(source)
		}
		tempDir, e := ioutil.TempDir("", "go-apigee-")
		if e != nil {
			return nil, nil, errors.New(fmt.Sprintf("while creating temp dir, error: %#v", e))
		}
		zipfileName = path.Join(tempDir, "bundle.zip")
		var filePathElement string
		if uriPathElement == "apis" {
			filePathElement = "apiproxy"
		} else {
			filePathElement = "sharedflowbundle"
		}

		e = zipDirectory(path.Join(source, filePathElement), zipfileName, smartFilter)
		if e != nil {
			return nil, nil, errors.New(fmt.Sprintf("while creating temp dir, error: %#v", e))
		}
		fmt.Printf("zipped %s into %s\n\n", source, zipfileName)
		cleanup := func(filename string) {
			_ = os.Remove(filename)
			// if e != nil {
			// 	//..
			// }
		}
		defer cleanup(zipfileName)
	}

	if !strings.HasSuffix(zipfileName, ".zip") {
		return nil, nil, errors.New("source must be a zipfile")
	}

	info, err = os.Stat(zipfileName)
	if err != nil {
		return nil, nil, err
	}

	// append the query params
	origURL, err := url.Parse(uriPathElement)
	if err != nil {
		return nil, nil, err
	}
	q := origURL.Query()
	q.Add("action", "import")
	q.Add("name", assetName)
	origURL.RawQuery = q.Encode()
	path := origURL.String()

	ioreader, err := os.Open(zipfileName)
	if err != nil {
		return nil, nil, err
	}
	defer ioreader.Close()

	req, e := client.NewRequest("POST", path, ioreader)
	if e != nil {
		return nil, nil, e
	}
	returnedRevision := DeployableRevision{}
	resp, e := client.Do(req, &returnedRevision)
	if e != nil {
		return nil, resp, e
	}
	return &returnedRevision, resp, e
}

func (s *Deployable) Export(client *ApigeeClient, uriPathElement, assetName string, rev Revision) (string, *Response, error) {
	// curl -u USER:PASSWORD \
	//  http://MGMTSERVER/v1/o/ORGNAME/apis/APINAME/revisions/REVNUMBER?format=bundle > bundle.zip

	path := path.Join(uriPathElement, assetName, "revisions", fmt.Sprintf("%d", rev))
	// TODO: factor out method: appendQueryParams
	// append the required query param
	origURL, err := url.Parse(path)
	if err != nil {
		return "", nil, err
	}
	q := origURL.Query()
	q.Add("format", "bundle")
	origURL.RawQuery = q.Encode()
	path = origURL.String()

	req, e := client.NewRequest("GET", path, nil)
	if e != nil {
		return "", nil, e
	}
	req.Header.Del("Accept")

	var assetType string
	if uriPathElement == "apis" {
		assetType = "apiproxy"
	} else {
		assetType = "sharedflowbundle"
	}

	t := time.Now()
	filename := fmt.Sprintf("%s-%s-r%d-%d%02d%02d-%02d%02d%02d.zip",
		assetType, assetName,
		rev, t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	out, e := os.Create(filename)
	if e != nil {
		return "", nil, e
	}

	resp, e := client.Do(req, out)
	if e != nil {
		return "", resp, e
	}
	out.Close()
	return filename, resp, e
}

func (s *Deployable) DeleteRevision(client *ApigeeClient, uriPathElement, assetName string, rev Revision) (*DeployableRevision, *Response, error) {
	path := path.Join(uriPathElement, assetName, "revisions", fmt.Sprintf("%d", rev))
	req, e := client.NewRequest("DELETE", path, nil)
	if e != nil {
		return nil, nil, e
	}
	proxyRev := DeployableRevision{}
	resp, e := client.Do(req, &proxyRev)
	if e != nil {
		return nil, resp, e
	}
	return &proxyRev, resp, e
}

func (s *Deployable) Undeploy(client *ApigeeClient, uriPathElement, assetName, env string, rev Revision) (*RevisionDeployment, *Response, error) {
	path := path.Join(uriPathElement, assetName, "revisions", fmt.Sprintf("%d", rev), "deployments")
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

	req, e := client.NewRequest("POST", path, nil)
	if e != nil {
		return nil, nil, e
	}

	deployment := RevisionDeployment{}
	resp, e := client.Do(req, &deployment)
	if e != nil {
		return nil, resp, e
	}
	return &deployment, resp, e
}

func (s *Deployable) Deploy(client *ApigeeClient, uriPathElement, assetName, basepath, env string, rev Revision) (*RevisionDeployment, *Response, error) {
	path := path.Join(uriPathElement, assetName, "revisions", fmt.Sprintf("%d", rev), "deployments")
	// append the query params
	origURL, err := url.Parse(path)
	if err != nil {
		return nil, nil, err
	}
	q := origURL.Query()
	q.Add("action", "deploy")
	q.Add("override", "true")
	q.Add("delay", DeploymentDelay)
	q.Add("env", env)
	if basepath != "" {
		q.Add("basepath", basepath)
	}
	origURL.RawQuery = q.Encode()
	path = origURL.String()

	req, e := client.NewRequest("POST", path, nil)
	if e != nil {
		return nil, nil, e
	}

	deployment := RevisionDeployment{}
	resp, e := client.Do(req, &deployment)
	if e != nil {
		return nil, resp, e
	}
	return &deployment, resp, e
}

// Delete an API Proxy and all its revisions from an organization. This method
// will fail if any of the revisions of the named API Proxy are currently deployed
// in any environment.
func (s *Deployable) Delete(client *ApigeeClient, uriPathElement, assetName string) (*DeletedItemInfo, *Response, error) {
	path := path.Join(uriPathElement, assetName)
	req, e := client.NewRequest("DELETE", path, nil)
	if e != nil {
		return nil, nil, e
	}
	item := DeletedItemInfo{}
	resp, e := client.Do(req, &item)
	if e != nil {
		return nil, resp, e
	}
	return &item, resp, e
}

func (s *Deployable) GetDeployments(client *ApigeeClient, uriPathElement, assetName string) (*Deployment, *Response, error) {
	path := path.Join(uriPathElement, assetName, "deployments")
	req, e := client.NewRequest("GET", path, nil)
	if e != nil {
		return nil, nil, e
	}
	deployments := Deployment{}
	resp, e := client.Do(req, &deployments)
	if e != nil {
		return nil, resp, e
	}
	return &deployments, resp, e
}
