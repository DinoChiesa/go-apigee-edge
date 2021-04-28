package main

import (
	"flag"
	"fmt"

	"<fullpath>/go-apigee-edge"
)

func usage() {
	fmt.Printf("test-keyValuemap -org <org name> -env test\n")
}

func main() {
	orgPtr := flag.String("org", "", "an Edge Organization")
	envPtr := flag.String("env", "", "an Edge environment")
	flag.Parse()

	if *orgPtr == "" {
		usage()
		return
	}

	if *envPtr == "" {
		usage()
		return
	}

	// Specify creds like so
	auth := &apigee.EdgeAuth{AccessToken: "<a bear token that you have previously generated>"}

	opts := &apigee.EdgeClientOptions{Org: *orgPtr, Auth: auth, Debug: true}
	client, e := apigee.NewEdgeClient(opts)
	if e != nil {
		fmt.Printf("while initializing Edge client, error:\n%#v\n", e)
		return
	}

	fmt.Printf(client.BaseURL.EscapedPath())

	//Get an existing keyvalumap
	fmt.Printf("Getting KeyValueMap list that should always exist...\n")

	KeyValueMapValue, resp, e := client.KeyValueMap.Get(*mapPtr, *envPtr)
	if e != nil {
		fmt.Printf("while getting KeyValueMap value, error:\n%#v\n", e)
		return
	}

	fmt.Printf("status: %d\n", resp.StatusCode)
	fmt.Printf("status: %s\n", resp.Status)
	defer resp.Body.Close()
	fmt.Printf("KeyValueMapValue: %#v\n", KeyValueMapValue)

	//Create a keyvalumap
	fmt.Printf("Creating a KeyValueMap...\n")

	//A KVM that will get created as part of the tests
	createKVM := apigee.KeyValueMap{
		Name:      "CreatedKeymap",
		Encrypted: false,
		Entry: []apigee.EntryStruct{
			apigee.EntryStruct{
				Name:  "Key1",
				Value: "value1",
			},
			apigee.EntryStruct{
				Name:  "Key2",
				Value: "value2",
			},
			apigee.EntryStruct{
				Name:  "Key3",
				Value: "value3",
			},
		},
	}

	KeyValueMapValue3, resp3, e3 := client.KeyValueMap.Create(createKVM, *envPtr)
	if e3 != nil {
		fmt.Printf("while getting KeyValueMap value, error:\n%#v\n", e)
	} else {
		fmt.Printf("status: %d\n", resp3.StatusCode)
		fmt.Printf("status: %s\n", resp3.Status)
		defer resp3.Body.Close()
		fmt.Printf("KeyValueMapValue: %#v\n", KeyValueMapValue3)
	}

	//delete a key value map
	resp5, e5 := client.KeyValueMap.Delete("CreatedKeymap", *envPtr)
	if e5 != nil {
		fmt.Printf("while getting KeyValueMap value, error:\n%#v\n", e)
	} else {
		fmt.Printf("status: %d\n", resp5.StatusCode)
		fmt.Printf("status: %s\n", resp5.Status)
		defer resp5.Body.Close()
	}
}
