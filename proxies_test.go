package apigee

import (
  "testing"
	"io/ioutil"
  "path"
  "os"
  "time"
  "strings"
	"fmt"
	"math/rand"
)

const (
	proxyBundleDir = "testdata/proxybundles"
)

func getProxyZipFiles(t *testing.T) []os.FileInfo {
	entries, e := ioutil.ReadDir(proxyBundleDir)
	if e != nil {
		t.Errorf("while reading testdata directory, error:\n%#v\n", e)
    return nil
	}

	zipfiles := []os.FileInfo{}
	for i := range entries {
		if strings.HasSuffix(entries[i].Name(), ".zip") {
			zipfiles = append(zipfiles, entries[i])
		}
	}

	if len(zipfiles) <= 1 {
		t.Errorf("not enough zipfiles in the proxyBundles directory, error:\n%#v\n", e)
    return nil
	}
	return zipfiles
}


func getProxyBundleDirectories(t *testing.T) []os.FileInfo {
	entries, e := ioutil.ReadDir(proxyBundleDir)
	if e != nil {
		t.Errorf("while reading testdata directory, error:\n%#v\n", e)
    return nil
	}

	directories := []os.FileInfo{}
	for i := range entries {
		if (!strings.HasSuffix(entries[i].Name(), ".zip")) {
			directories = append(directories, entries[i])
		}
	}

	if len(directories) <= 1 {
		t.Errorf("not enough directories in the proxyBundles directory, error:\n%#v\n", e)
    return nil
	}

	return directories
}


func getEnvironments(t *testing.T, client *ApigeeClient) []string {
// get environments
  envlist, resp, e := client.Environments.List()
  if e != nil {
		t.Errorf("while listing environments, error:\n%#v\n", e)
    return nil
  }
	defer resp.Body.Close()
  if len(envlist) <= 0 {
		t.Errorf("no environments found")
    return nil
  }

	filteredEnvList := []string{}
	for i := range envlist {
		if !strings.HasPrefix(envlist[i], "portal") {
			filteredEnvList = append(filteredEnvList, envlist[i])
		}
	}
	t.Logf("environments: %#v", filteredEnvList)
	return filteredEnvList
}

func TestProxyImportFromZip(t *testing.T) {
	now := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02d-%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

	zipfiles := getProxyZipFiles(t)
  client := NewClientForTesting(t)

	for _, zipfile := range zipfiles {
		fullFileName := path.Join(proxyBundleDir, zipfile.Name())
		t.Logf(fullFileName)
		proxyName := fmt.Sprintf("%s-%s-%s", testPrefix, timestamp, zipfile.Name())

		proxyRev, resp, e := client.Proxies.Import(proxyName, fullFileName)
		if e != nil {
			t.Errorf("while importing, error:\n%#v\n", e)
			return
		}
		t.Logf("status: %s\n", resp.Status)
		if resp.Status != "201 Created" {
			t.Errorf("while importing, status: %#v\n", resp.Status)
			return
		}
		defer resp.Body.Close()
		t.Logf("proxyRev: %#v\n", proxyRev)
	}
}


func TestProxyImportFromDir(t *testing.T) {
  client := NewClientForTesting(t)

	now := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02d-%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

	directories := getProxyBundleDirectories(t)
	for _, directory := range directories {
		fullDirName := path.Join(proxyBundleDir, directory.Name())
		t.Logf(fullDirName)
		proxyName := fmt.Sprintf("%s-%s-%s", testPrefix, timestamp, directory.Name())

		proxyRev, resp, e := client.Proxies.Import(proxyName, fullDirName)
		if e != nil {
			t.Errorf("while importing, error:\n%#v\n", e)
			return
		}
		t.Logf("status: %s\n", resp.Status)
		if resp.Status != "201 Created" {
			t.Errorf("while importing, status: %#v\n", resp.Status)
			return
		}
		defer resp.Body.Close()
		t.Logf("proxyRev: %#v\n", proxyRev)
	}
}


func TestProxyList(t *testing.T) {
  client := NewClientForTesting(t)

	proxies, resp, e := client.Proxies.List()
	if e != nil {
		t.Errorf("while listing proxies, error:\n%#v\n", e)
    return
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		t.Errorf("while listing, status: %#v\n", resp.Status)
		return
	}
	if len(proxies) ==0 {
		t.Errorf("no proxies returned while listing")
		return
	}
}

func TestProxyGet(t *testing.T) {
  client := NewClientForTesting(t)
	proxies, resp, e := client.Proxies.List()
	if e != nil {
		t.Errorf("while listing proxies, error:\n%#v\n", e)
    return
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		t.Errorf("while listing, status: %#v\n", resp.Status)
		return
	}
	if len(proxies) ==0 {
		t.Errorf("no proxies returned while listing")
		return
	}

	readCount := 0
	for readCount < 8 {
		proxyName := proxies[rand.Intn(len(proxies))]

		proxy, resp, e := client.Proxies.Get(proxyName)
		if e != nil {
			t.Errorf("while getting, error:\n%#v\n", e)
			return
		}
		defer resp.Body.Close()
		//t.Logf("status: %s\n", resp.Status)
		if resp.Status != "200 OK" {
			t.Errorf("while getting, status: %#v\n", resp.Status)
			return
		}
		t.Logf("read: %#v\n", proxy)
		readCount++
	}

	if readCount < 1 {
		t.Errorf("did not GET any proxies")
	}
}

type ImportedProxyStruct struct {
	proxyName string
	env string
	rev Revision
}

func TestProxyDeployUndeploy(t *testing.T) {
	now := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02d-%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

  client := NewClientForTesting(t)
	zipfiles := getProxyZipFiles(t)
	environmentList := getEnvironments(t, client)
	importedProxies := []ImportedProxyStruct{}

	for _, zipfile := range zipfiles {
		fullFileName := path.Join(proxyBundleDir, zipfile.Name())
		t.Logf(fullFileName)
		proxyName := fmt.Sprintf("%s-deploy-%s-%s", testPrefix, timestamp, zipfile.Name())
		proxyRev, resp, e := client.Proxies.Import(proxyName, fullFileName)
		if e != nil {
			t.Errorf("while importing, error:\n%#v\n", e)
			return
		}
		t.Logf("status: %s\n", resp.Status)
		if resp.Status != "201 Created" {
			t.Errorf("while importing, status: %#v\n", resp.Status)
			return
		}
		defer resp.Body.Close()

		randomPath := fmt.Sprintf("/%s-deploy", timestamp)
		env := environmentList[rand.Intn(len(environmentList))]
		importedProxies = append(importedProxies, ImportedProxyStruct{proxyName:proxyName, env:env, rev:proxyRev.Revision})
		deployment, resp, e := client.Proxies.DeployAtPath(proxyName, randomPath, env, proxyRev.Revision)
		if e != nil {
			t.Errorf("while deploying, error:\n%#v\n", e)
			return
		}
		defer resp.Body.Close()
		if resp.Status != "200 OK" {
			t.Errorf("while deploying, status: %#v\n", resp.Status)
			return
		}
		t.Logf("deployed: %#v\n", deployment)
	}

	for _, item := range importedProxies {
		deployment, resp, e := client.Proxies.Undeploy(item.proxyName, item.env, item.rev)
		if e != nil {
			t.Errorf("while undeploying, error:\n%#v\n", e)
			return
		}
		defer resp.Body.Close()
		if resp.Status != "200 OK" {
			t.Errorf("while undeploying, status: %#v\n", resp.Status)
			return
		}
		t.Logf("undeployed: %#v\n", deployment)
	}
}


func TestProxyInquireDeployment(t *testing.T) {
	now := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02d-%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
  client := NewClientForTesting(t)
	zipfiles := getProxyZipFiles(t)
	selectedZip := zipfiles[rand.Intn(len(zipfiles))]
	environmentList := getEnvironments(t, client)
	env := environmentList[rand.Intn(len(environmentList))]

	// import
	proxyName := fmt.Sprintf("%s-deployment-%s-%s", testPrefix, timestamp, selectedZip.Name())
	fullFileName := path.Join(proxyBundleDir, selectedZip.Name())
	proxyRev, resp, e := client.Proxies.Import(proxyName, fullFileName)
	if e != nil {
		t.Errorf("while importing, error:\n%#v\n", e)
		return
	}
	t.Logf("status: %s\n", resp.Status)
	if resp.Status != "201 Created" {
		t.Errorf("while importing, status: %#v\n", resp.Status)
		return
	}
	defer resp.Body.Close()

	// inquire deployment
	deployment, resp, e := client.Proxies.GetDeployments(proxyName)
	if e != nil {
		t.Errorf("while inquiring deployments, error:\n%#v\n", e)
		return
	}
	t.Logf("status: %s\n", resp.Status)
	if resp.Status != "200 OK" {
		t.Errorf("while inquiring deployments, status: %#v\n", resp.Status)
		return
	}
	defer resp.Body.Close()
	if len(deployment.Environments) != 0 {
		t.Errorf("found unexpected deployments: %#v\n", deployment.Environments)
		return
	}

	// deploy
	randomPath := fmt.Sprintf("/%s-inquiredeployment", timestamp)
	revisionDeployment, resp, e := client.Proxies.DeployAtPath(proxyName, randomPath, env, proxyRev.Revision)
	if e != nil {
		t.Errorf("while deploying, error:\n%#v\n", e)
		return
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		t.Errorf("while deploying, status: %#v\n", resp.Status)
		return
	}
	t.Logf("deployed: %#v\n", deployment)

	// inquire deployment
	deployment, resp, e = client.Proxies.GetDeployments(proxyName)
	if e != nil {
		t.Errorf("while inquiring deployments, error:\n%#v\n", e)
		return
	}
	t.Logf("status: %s\n", resp.Status)
	if resp.Status != "200 OK" {
		t.Errorf("while inquiring deployments, status: %#v\n", resp.Status)
		return
	}
	defer resp.Body.Close()
	if len(deployment.Environments) != 1 {
		t.Errorf("found unexpected number of deployments: %#v\n", deployment.Environments)
		return
	}

	// undeploy
	revisionDeployment, resp, e = client.Proxies.Undeploy(proxyName, env, proxyRev.Revision)
	if e != nil {
		t.Errorf("while undeploying, error:\n%#v\n", e)
		return
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		t.Errorf("while undeploying, status: %#v\n", resp.Status)
		return
	}
	t.Logf("undeployed: %#v\n", revisionDeployment)
}


func TestProxyDelete(t *testing.T) {
  client := NewClientForTesting(t)
	proxies, resp, e := client.Proxies.List()
	if e != nil {
		t.Errorf("while listing proxies, error:\n%#v\n", e)
    return
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		t.Errorf("while listing, status: %#v\n", resp.Status)
		return
	}

	deletedCount := 0
	for _, proxyName := range proxies {
		if strings.HasPrefix(proxyName, testPrefix) {

			// inquire deployment
			deployment, resp, e := client.Proxies.GetDeployments(proxyName)
			if e != nil {
				t.Errorf("while inquiring deployments, error:\n%#v\n", e)
				return
			}
			if resp.Status != "200 OK" {
				t.Errorf("while inquiring deployments, status: %#v\n", resp.Status)
				return
			}
			defer resp.Body.Close()

			// conditionally undeploy
			if len(deployment.Environments) != 0 {
				count := 0
				for count < len(deployment.Environments) {
					revisionDeployment, resp, e := client.Proxies.Undeploy(proxyName,
						deployment.Environments[count].Name,
						deployment.Environments[count].Revision[0].Number)
					if e != nil {
						t.Errorf("while undeploying, error:\n%#v\n", e)
						return
					}
					defer resp.Body.Close()
					if resp.Status != "200 OK" {
						t.Errorf("while undeploying, status: %#v\n", resp.Status)
						return
					}
					t.Logf("undeployed: %#v\n", revisionDeployment)
					count++
				}
			}

			// delete
			t.Logf("deleting %#v", proxyName)
			deletedItem, resp, e := client.Proxies.Delete(proxyName)
			if e != nil {
				t.Errorf("while deleting, error:\n%#v\n", e)
				return
			}
			t.Logf("status: %s\n", resp.Status)
			if resp.Status != "200 OK" {
				t.Errorf("while deleting, status: %#v\n", resp.Status)
				return
			}
			t.Logf("deleted: %#v\n", deletedItem)
			defer resp.Body.Close()
			deletedCount++
		}
	}

	if deletedCount == 0 {
		t.Errorf("no proxies deleted")
	}
}


func TestProxyDeleteFail(t *testing.T) {
	now := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02d-%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

  client := NewClientForTesting(t)
	zipfiles := getProxyZipFiles(t)
	selectedZip := zipfiles[rand.Intn(len(zipfiles))]
	environmentList := getEnvironments(t, client)
	env := environmentList[rand.Intn(len(environmentList))]

	fullFileName := path.Join(proxyBundleDir, selectedZip.Name())
	t.Logf(fullFileName)
	proxyName := fmt.Sprintf("%s-deletefail-%s-%s", testPrefix, timestamp, selectedZip.Name())
	proxyRev, resp, e := client.Proxies.Import(proxyName, fullFileName)
	if e != nil {
		t.Errorf("while importing, error:\n%#v\n", e)
		return
	}
	t.Logf("status: %s\n", resp.Status)
	if resp.Status != "201 Created" {
		t.Errorf("while importing, status: %#v\n", resp.Status)
		return
	}
	defer resp.Body.Close()

	randomPath := fmt.Sprintf("/%s-deletefail", timestamp)
	deployment, resp, e := client.Proxies.DeployAtPath(proxyName, randomPath, env, proxyRev.Revision)
	if e != nil {
		t.Errorf("while deploying, error:\n%#v\n", e)
		return
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		t.Errorf("while deploying, status: %#v\n", resp.Status)
		return
	}
	t.Logf("deployed: %#v\n", deployment)
	wait(1)
	t.Logf("attempting delete %#v", proxyName)
	deletedItem, resp, e := client.Proxies.Delete(proxyName)
	if e == nil {
		t.Errorf("while attempting delete, expected an error\n")
		return
	}
	// if e.Error() != "something" {
	// 	t.Errorf("unexpected error during delete: %#v\n", e.Error())
	// 	return
	// }
	defer resp.Body.Close()

	// undeploy
	deployment, resp, e = client.Proxies.Undeploy(proxyName, env, proxyRev.Revision)
	if e != nil {
		t.Errorf("while undeploying, error:\n%#v\n", e)
		return
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		t.Errorf("while undeploying, status: %#v\n", resp.Status)
		return
	}
	wait(1)

	// then delete succeeds
	t.Logf("attempting delete %#v", proxyName)
	deletedItem, resp, e = client.Proxies.Delete(proxyName)
	if e != nil {
		t.Errorf("while deleting, error:\n%#v\n", e)
		return
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		t.Errorf("while deleting, status: %#v\n", resp.Status)
		return
	}
	t.Logf("deleted %#v", deletedItem)

}
