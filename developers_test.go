package apigee

import (
	"encoding/json"
	"math/rand"
	"testing"
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
	tag := testPrefix + randomString(6)
	got.Email = tag + got.Email
	got.UserName = tag + "-" + got.UserName
	got.FirstName = got.FirstName + "-" + tag
	return got, e
}

func TestDeveloperCreateDelete(t *testing.T) {
	client := NewClientForTesting(t)
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
	client := NewClientForTesting(t)
	developerList, _, e := client.Developers.List()
	if e != nil {
		t.Errorf("while listing Edge developers, error:\n%#v\n", e)
		return
	}
	t.Logf("List: got=%+v", developerList)
}

func TestDeveloperGet(t *testing.T) {
	client := NewClientForTesting(t)
	developerList, _, e := client.Developers.List()
	if e != nil {
		t.Errorf("while listing Edge developers, error:\n%#v\n", e)
		return
	}

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
	client := NewClientForTesting(t)
	dev, e := randomDeveloperFromTemplate()
	createdDeveloper, _, e := client.Developers.Create(dev)
	if e != nil {
		t.Errorf("while creating Edge developer, error:\n%#v\n", e)
		return
	}
	t.Logf("Create: got=%+v", createdDeveloper)
	teardown := func(t *testing.T) {
		deletedDeveloper, _, e := client.Developers.Delete(createdDeveloper.Email)
		if e != nil {
			t.Errorf("while deleting Edge developer, error:\n%#v\n", e)
			return
		}
		t.Logf("Delete: got=%v", deletedDeveloper)
	}
	defer teardown(t)
	wait(1)

	_, e = client.Developers.Revoke(createdDeveloper.Email)
	if e != nil {
		t.Errorf("while revoking Edge developer, error:\n%#v\n", e)
		return
	}
	t.Logf("Revoke")
	wait(1)

	_, e = client.Developers.Approve(createdDeveloper.Email)
	if e != nil {
		t.Errorf("while approving Edge developer, error:\n%#v\n", e)
		return
	}
	t.Logf("Approve")
	wait(1)
}
