package apigee

import (
  "encoding/json"
  "testing"
	"math/rand"
	"time"
)

const (
  developerJson1 = `{
  "attributes": [ {
    "name" : "tag1",
    "value" : "created by golang" }],
  "status": "active",
  "userName": "username",
  "lastName": "Martino",
  "firstName": "Dino",
  "email": "@apigee.com",
  "companies": [],
  "apps": []
}`
)


func randomDeveloperFromTemplate() (Developer, error) {
	got := Developer{}
	e := json.Unmarshal([]byte(developerJson1), &got)
	
	if e != nil {
		return got, e
	}
	// assign values
	tag := pretag + randomString(6)
	got.Email = tag + got.Email
	got.UserName = tag + "-" + got.UserName 
	got.FirstName = got.FirstName + "-" + tag
	return got, e
}


func TestDeveloperCreateDelete(t *testing.T) {
  opts := &EdgeClientOptions{Org: orgName, Auth: nil, Debug: false }
  client, e := NewEdgeClient(opts)
  if e != nil {
		t.Errorf("while initializing Edge client, error:\n%#v\n", e)
    return
  }
  //wait(1)

	dev, e := randomDeveloperFromTemplate()
  createdDeveloper, resp, e := client.Developers.Create(dev)
  if e != nil {
		t.Errorf("while creating Edge developer, error:\n%#v\n", e)
    return
  }
	t.Logf("Create: got=%+v", createdDeveloper)
	t.Logf("resp: got=%+v", resp)
	
  wait(1)

  deletedDeveloper, resp, e := client.Developers.Delete(createdDeveloper.Email)
  if e != nil {
		t.Errorf("while deleting Edge developer, error:\n%#v\n", e)
    return
  }
	t.Logf("Delete: got=%v", deletedDeveloper)
}


func TestDeveloperList(t *testing.T) {
  opts := &EdgeClientOptions{Org: orgName, Auth: nil, Debug: false }
  client, e := NewEdgeClient(opts)
  if e != nil {
		t.Errorf("while initializing Edge client, error:\n%#v\n", e)
    return
  }

  developerList, _, e := client.Developers.List()
  if e != nil {
		t.Errorf("while listing Edge developers, error:\n%#v\n", e)
    return
  }
	t.Logf("List: got=%+v", developerList)
}

func TestDeveloperGet(t *testing.T) {
  opts := &EdgeClientOptions{Org: orgName, Auth: nil, Debug: false }
  client, e := NewEdgeClient(opts)
  if e != nil {
		t.Errorf("while initializing Edge client, error:\n%#v\n", e)
    return
  }

  developerList, _, e := client.Developers.List()
  if e != nil {
		t.Errorf("while listing Edge developers, error:\n%#v\n", e)
    return
  }

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	selectedDevEmail := developerList[rand.Intn(len(developerList))]

	developerDetails, _, e := client.Developers.Get(selectedDevEmail)
  if e != nil {
		t.Errorf("while getting Edge developer, error:\n%#v\n", e)
    return
  }
	t.Logf("Get: selected=%+v", developerDetails)
	t.Logf("Get: email=%s", developerDetails.Email)
}


func TestDeveloperUpdate(t *testing.T) {
  opts := &EdgeClientOptions{Org: orgName, Auth: nil, Debug: false }
  client, e := NewEdgeClient(opts)
  if e != nil {
		t.Errorf("while initializing Edge client, error:\n%#v\n", e)
    return
  }

	dev, e := randomDeveloperFromTemplate()
  createdDeveloper, _, e := client.Developers.Create(dev)
  if e != nil {
		t.Errorf("while creating Edge developer, error:\n%#v\n", e)
    return
  }
	t.Logf("Create: got=%+v", createdDeveloper)
	
  wait(1)

	_, e = client.Developers.Revoke(createdDeveloper.Email)
  if e != nil {
		t.Errorf("while revoking Edge developer, error:\n%#v\n", e)
    return
  }
	
  wait(1)

	_, e = client.Developers.Approve(createdDeveloper.Email)
  if e != nil {
		t.Errorf("while approving Edge developer, error:\n%#v\n", e)
    return
  }
	
  wait(1)
  deletedDeveloper, _, e := client.Developers.Delete(createdDeveloper.Email)
  if e != nil {
		t.Errorf("while deleting Edge developer, error:\n%#v\n", e)
    return
  }
	t.Logf("Delete: got=%v", deletedDeveloper)
}
