package apigee

import (
  "encoding/json"
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
  "proxies" : [ "dino-test" ],
  "scopes" : [ "" ]
}`
)

func randomProductFromTemplate() (ApiProduct, error) {
	got := ApiProduct{}
	e := json.Unmarshal([]byte(productJson1), &got)
	
	if e != nil {
		return got, e
	}
	// assign values
	tag := randomString(14)
	got.Name = tag + "-" + got.Name
	got.DisplayName = tag + "-" + got.DisplayName 
	got.Description = tag + " " + randomString(8) + " " + randomString(18)
  got.Scopes = []string { randomString(1), randomString(2), }
	return got, e
}


func TestProductCreateDelete(t *testing.T) {
  orgName := "cap500"
  opts := &EdgeClientOptions{Org: orgName, Auth: nil, Debug: false }
  client, e := NewEdgeClient(opts)
  if e != nil {
		t.Errorf("while initializing Edge client, error:\n%#v\n", e)
    return
  }

	product, e := randomProductFromTemplate()
  createdProduct, resp, e := client.Products.Create(product)
  if e != nil {
		t.Errorf("while creating Edge product, error:\n%#v\n", e)
    return
  }
	t.Logf("Create: got=%v", createdProduct)
	t.Logf("resp: got=%v", resp)
	
  wait(1)

  deletedProduct, resp, e := client.Products.Delete(createdProduct.Name)
  if e != nil {
		t.Errorf("while deleting Edge product, error:\n%#v\n", e)
    return
  }
	t.Logf("Delete: got=%v", deletedProduct)
}



