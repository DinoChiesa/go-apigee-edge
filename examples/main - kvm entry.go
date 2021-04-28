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

	//Get an existing keyvalumapentry - YOU MUST CREATE A KVM FIRST, this assumes its called helloThere
	fmt.Printf("Getting KeyValueMapEntry list that should always exist...\n")

	KeyValueMapValue, resp, e := client.KeyValueMapEntry.Get("helloThere", *envPtr, "1t")
	if e != nil {
		fmt.Printf("while getting KeyValueMapEntry value, error:\n%#v\n", e)
		return
	}

	fmt.Printf("status: %d\n", resp.StatusCode)
	fmt.Printf("status: %s\n", resp.Status)
	defer resp.Body.Close()
	fmt.Printf("KeyValueMapValue: %#v\n", KeyValueMapValue)

	//Create a keyvaluemapentry
	fmt.Printf("Creating a KeyValueMap entry...\n")

	//A KVM entry that will get created as part of the tests
	createKVM := apigee.KeyValueMapEntryKeys{
		Name:  "General",
		Value: "Kenobi",
	}

	KeyValueMapValue2, resp2, e2 := client.KeyValueMapEntry.Create("helloThere", createKVM, *envPtr)
	if e2 != nil {
		fmt.Printf("while getting KeyValueMap value, error:\n%#v\n", e)
	} else {
		fmt.Printf("status: %d\n", resp2.StatusCode)
		fmt.Printf("status: %s\n", resp2.Status)
		defer resp2.Body.Close()
		fmt.Printf("KeyValueMapValue: %#v\n", KeyValueMapValue2)
	}

	//update an existing kvm entry
	updateKVM := apigee.KeyValueMapEntryKeys{
		Name:  "General",
		Value: "BoldOne",
	}

	KeyValueMapValue3, resp3, e3 := client.KeyValueMapEntry.Update("helloThere", updateKVM, *envPtr)
	if e3 != nil {
		fmt.Printf("while getting KeyValueMap value, error:\n%#v\n", e)
	} else {
		fmt.Printf("status: %d\n", resp3.StatusCode)
		fmt.Printf("status: %s\n", resp3.Status)
		defer resp3.Body.Close()
		fmt.Printf("KeyValueMapValue: %#v\n", KeyValueMapValue3)
	}

	//delete a key value map
	resp5, e5 := client.KeyValueMapEntry.Delete("Key1", "helloThere", *envPtr)
	if e5 != nil {
		fmt.Printf("while getting KeyValueMap value, error:\n%#v\n", e)
	} else {
		fmt.Printf("status: %d\n", resp5.StatusCode)
		fmt.Printf("status: %s\n", resp5.Status)
		defer resp5.Body.Close()
	}
}
