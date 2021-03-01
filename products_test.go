package apigee

import (
	"encoding/json"
	"math/rand"
	"testing"
)

const (
	productJson1 = `{
  "apiResources" : [ ],
  "approvalType" : "auto",
  "attributes" : [ {
    "name" : "category",
    "value" : "Flight Operations, Reservations"
  }, {
    "name" : "image_url",
    "value" : "http://i.cdn.turner.com/cnn/2010/TRAVEL/11/18/flying.driving.travel.apps/t1larg.flight.status.jpg"
  } ],
  "description" : "",
  "displayName" : "FlightStatus",
  "environments" : [ "test" ],
  "name" : "FlightStatus",
  "proxies" : [ "will-be-replaced" ],
  "scopes" : [ "" ]
}`
)

func randomProductFromTemplate(proxyname string) (ApiProduct, error) {

	// just a way to quickly set values to a few defaults
	got := ApiProduct{}
	e := json.Unmarshal([]byte(productJson1), &got)

	if e != nil {
		return got, e
	}
	// assign values
	tag := testPrefix + randomString(7)
	got.Name = tag + "-" + got.Name
	got.Proxies = []string{proxyname}
	got.DisplayName = tag + "-" + got.DisplayName
	got.Description = tag + " " + randomString(8) + " " + randomString(18)
	got.Scopes = []string{randomString(1), randomString(2)}
	return got, e
}

func TestProductCreateDelete(t *testing.T) {
	client := NewClientForTesting(t)
	namelist, resp, e := client.Proxies.List()
	if e != nil {
		t.Errorf("while listing proxies, error:\n%#v\n", e)
		return
	}
	if len(namelist) <= 0 {
		t.Errorf("no proxies found")
		return
	}

	selectedProxy := namelist[rand.Intn(len(namelist))]

	product, e := randomProductFromTemplate(selectedProxy)
	createdProduct, resp, e := client.Products.Create(product)
	if e != nil {
		t.Errorf("while creating Apigee product, error:\n%#v\n", e)
		return
	}
	t.Logf("Create: got=%v", createdProduct)
	t.Logf("resp: got=%v", resp)

	wait(1)

	deletedProduct, resp, e := client.Products.Delete(createdProduct.Name)
	if e != nil {
		t.Errorf("while deleting Apigee product, error:\n%#v\n", e)
		return
	}
	t.Logf("Delete: got=%v", deletedProduct)
}
