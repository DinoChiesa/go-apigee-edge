# golang client library for Apigee Edge administrative API

Use this from Go-lang programs to invoke administrative operations on Apigee Edge.

The goal is to allow golang programs to easiy do these things:

| entity type   | actions             |
| :------------ | :------------------ |
| apis          | list, query, inquire revisions, inquire deployment status, import, export, delete, delete revision, deploy, undeploy
| apiproducts   | list, query, create, delete, change quota, modify public/private, modify description, modify approvalType, modify scopes, add or remove proxy, modify custom attrs
| developers    | list, query, create, delete, make active or inactive, modify custom attrs
| developer app | list, query, create, delete, revoke, approve, add new credential, remove credential, modify custom attrs
| credential    | list, revoke, approve, add apiproduct, remove apiproduct
| kvm           | list, query, create, delete, get all entries, get entry, add entry, modify entry, remove entry
| cache         | list, query, create, delete, clear
| environment   | list, query

The Apigee Edge administrative API is just a REST-ful API, so of course any go program could invoke it directly. This library will provide a wrapper, which will make it easier.


Not in scope:

- OAuth2.0 tokens - Listing, Querying, Approving, Revoking, Deleting, or Updating
- TargetServers: list, create, edit, etc
- keystores, truststores: adding certs, listing certs
- data masks
- apimodels
- shared flows or flow hooks (for now; we will deliver this when shared flows are final)
- analytics or custom reports
- DebugSessions (trace)
- anything in BaaS
- OPDK-specific things.  Like starting or stopping services, manipulating pods, adding servers into environments, etc.

These items may be added later as need and demand warrants.

## Copyright and License

This code is [Copyright (c) 2016 Apigee Corp](NOTICE). it is licensed under the [Apache 2.0 Source Licese](LICENSE).


## Status

This project is a work-in-progress. Here's the status:

| entity type   | implemented              | not implemented yet
| :------------ | :----------------------- | :--------------------
| apis          | list, query, inquire revisions, import, export, delete, delete revision, deploy, undeploy, inquire deployment status |
| apiproducts   | | list, query, create, delete, modify description, modify approvalType, modify scopes, add or remove proxy, add or remove custom attrs, modify public/private, change quota |
| developers    | | list, query, make active or inactive, create, delete, modify custom attrs |
| developer app | | list, query, create, delete, revoke, approve, add new credential, remove credential | modify custom attrs
| credential    | | list, revoke, approve, add apiproduct, remove apiproduct |
| kvm           | query, create, delete,  get entry, add entry, modify entry, remove entry | list, get all entries
| cache         | | list, query, create, delete, clear | 
| environment   | | list, query |

Pull requests are welcomed.


## Usage Example

```go
package main

import (
  "fmt"
  "flag"
  "time"
  "github.com/DinoChiesa/go-apigee-edge"
)

func usage() {
  fmt.Printf("import-proxy -user dino@example.org -org cap500 -name foobar -src /path/to/apiproxy\n\n")
}


func main() {
  proxyName := ""
  namePtr := flag.String("name", "", "name for the API Proxy")
  srcPtr := flag.String("src", "", "a directory containing an exploded apiproxy bundle, or a zipped bundle")
  orgPtr := flag.String("org", "", "an Edge Organization")
  flag.Parse()

  if *namePtr != "" {
    proxyName = *namePtr
  }

  if *srcPtr == "" || *orgPtr == "" {
    usage()
    return
  }

  var auth *apigee.EdgeAuth = nil

  // Specifying nil for Auth implies "read from .netrc"
  // Specify a password explicitly like so:
  // auth := apigee.EdgeAuth{Username: "user@example.org", Password: "Secret*123"}

  opts := &apigee.EdgeClientOptions{Org: *orgPtr, Auth: auth, Debug: false }
  client, e := apigee.NewEdgeClient(opts)
  if e != nil {
    fmt.Printf("while initializing Edge client, error:\n%#v\n", e)
    return
  }

  fmt.Printf("\nImporting...\n")
  proxyRev, resp, e := client.Proxies.Import(proxyName, *srcPtr)
  if e != nil {
    fmt.Printf("while importing, error:\n%#v\n", e)
    return
  }
  fmt.Printf("status: %d\n", resp.StatusCode)
  fmt.Printf("status: %s\n", resp.Status)
  defer resp.Body.Close()
  fmt.Printf("proxyRev: %#v\n", proxyRev)

  // TODO: Deploy the proxy revision with override = 10

  // TODO: Undeploy the proxy revision

  fmt.Printf("\nWaiting...\n")
  time.Sleep(3 * time.Second)

  fmt.Printf("\nDeleting...\n")
  deletedRev, resp, e := client.Proxies.DeleteRevision(proxyRev.Name, proxyRev.Revision)
  if e != nil {
    fmt.Printf("while deleting, error:\n%#v\n", e)
    return
  }
  fmt.Printf("status: %d\n", resp.StatusCode)
  fmt.Printf("status: %s\n", resp.Status)
  defer resp.Body.Close()
  fmt.Printf("proxyRev: %#v\n", deletedRev)
}

```

## Bugs

* There are embarrassingly few tests.

* When importing from a source directory, the library creates a temporary zip file, but doesn't delete the file.

* There is no working code for example clients, included in the distribution here.

* There is no package versioning strategy (eg, no use of GoPkg.in)

* When deploying a proxy, there's no way to specify the override and delay parameters.
