# golang client library for Apigee administrative API

Use this from Go-lang programs to invoke administrative operations on Apigee.

The goal is to allow golang programs to easiy do these things:

| entity type   | actions             |
| :------------ | :------------------ |
| apis          | list, query, inquire revisions, inquire deployment status, import, export, delete, delete revision, deploy, undeploy
| sharedflows   | list, query, inquire revisions, inquire deployment status, import, export, delete, delete revision, deploy, undeploy
| apiproducts   | list, query, create, delete, change quota, modify public/private, modify description, modify approvalType, modify scopes, add or remove proxy, modify custom attrs
| developers    | list, query, create, delete, make active or inactive, modify custom attrs
| developerapps | list, query, create, delete, revoke, approve, add new credential, remove credential, modify custom attrs
| credential    | list, revoke, approve, add apiproduct, remove apiproduct
| kvm           | list, query, create, delete, get all entries, get entry, add entry, modify entry, remove entry
| cache         | list, query, create, delete, clear
| environment   | list, query

The Apigee administrative API is just a REST-ful API, so of course any go program could invoke it directly. This library will provide a wrapper, which will make it easier.

Not yet in scope:

- OAuth2.0 tokens - Listing, Querying, Approving, Revoking, Deleting, or Updating
- TargetServers: list, create, edit, etc
- keystores, truststores: adding certs, listing certs
- data masks
- specs
- analytics or custom reports
- DebugSessions (trace)
- OPDK-specific things.  Like starting or stopping services, manipulating pods, adding servers into environments, etc.

These items may be added later as need and demand warrants.

## This is not an official Google product

This library and any example tools included here are not an official Google product, nor are they part of an official Google product.
Support is available on a best-effort basis via github or [community.apigee.com](https://community.apigee.com) .

## Copyright and License

This code is [Copyright (c) 2016 Apigee Corp, 2017-2020 Google LLC](NOTICE). it is licensed under the [Apache 2.0 Source Licese](LICENSE).


## Status

This project is a work-in-progress. Here's the status:

| entity type   | implemented              | not implemented yet
| :------------ | :----------------------- | :--------------------
| apis          | list, query, inquire revisions, import, export, delete, delete revision, deploy, undeploy, inquire deployment status |
| sharedflows   | list, query, inquire revisions, import, export, delete, delete revision, deploy, undeploy, inquire deployment status |
| apiproducts   | list, query, create, delete modify description, modify approvalType, modify scopes, add or remove proxy, add or remove custom attrs, modify public/private, change quota | |
| developers    | list, query, create, update, delete, modify custom attrs, make active or inactive, modify custom attrs |
| developerapps | list, query, create, delete, revoke, approve, modify custom attrs | add new credential, remove credential
| credential    | | list, revoke, approve, add apiproduct, remove apiproduct |
| kvm           | | list, query, create, delete, get all entries, get entry, add entry, modify entry, remove entry
| cache         | list, query | create, delete, clear |
| environment   | list, query | |

Pull requests are welcomed.


## Usage Examples

## The Import

This shows the reference to import.

```go
package main

import (
  "github.com/DinoChiesa/go-apigee-edge"
)
```

### List environments

This example connects to Apigee SaaS using credentials obtained from .netrc (because Auth:nil). It will use an OAuth token (WantToken: true).
That token gets stashed in a local file, and the library will re-use the token on subsequent runs, until the token expires.

```go
  opts := &apigee.ApigeeClientOptions{Org: "myorg", Auth: nil, Debug: true, WantToken: true }
  client, e := apigee.NewApigeeClient(opts)
  if e != nil {
    fmt.Printf("while initializing Edge client, error:\n%#v\n", e)
    return
  }

  fmt.Printf("\nListing...\n")
  list, resp, e := client.Environments.List()
  if e != nil {
    fmt.Printf("while listing, error:\n%#v\n", e)
    return
  }
  // full response is available in resp if necessary
  showStatus(resp)
  fmt.Printf("environments: %#v\n", list)
```

### Importing a Proxy

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
  orgPtr := flag.String("org", "", "an Apigee Organization")
  flag.Parse()

  if *namePtr != "" {
    proxyName = *namePtr
  }

  if *srcPtr == "" || *orgPtr == "" {
    usage()
    return
  }

  var auth *apigee.AdminAuth = nil

  // Specifying nil for Auth implies "read from .netrc"
  // Specify a password explicitly like so:
  // auth := apigee.AdminAuth{Username: "user@example.org", Password: "Secret*123"}

  opts := &apigee.ApigeeClientOptions{Org: *orgPtr, Auth: auth, Debug: false }
  client, e := apigee.NewApigeeClient(opts)
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
  fmt.Printf("status: %s\n", resp.Status)
  defer resp.Body.Close()
  fmt.Printf("proxyRev: %#v\n", proxyRev)
}

```

### Deleting a specific API Proxy Revision

```go
func main() {
  opts := &apigee.ApigeeClientOptions{Org: *orgPtr, Auth: nil, Debug: false }
  client, e := apigee.NewApigeeClient(opts)
  if e != nil {
    fmt.Printf("while initializing Apigee client, error:\n%#v\n", e)
    return
  }
  fmt.Printf("Deleting...\n")
  deletedRev, resp, e := client.Proxies.DeleteRevision(proxyName, Revision{2})
  if e != nil {
    fmt.Printf("while deleting, error:\n%#v\n", e)
    return
  }
  fmt.Printf("status: %s\n", resp.Status)
  defer resp.Body.Close()
  fmt.Printf("proxyRev: %#v\n", deletedRev)
}
```

### Deleting all revisions of an API Proxy
```go
func main() {
  opts := &apigee.ApigeeClientOptions{Org: *orgPtr, Auth: nil, Debug: false }
  client, e := apigee.NewApigeeClient(opts)
  if e != nil {
    fmt.Printf("while initializing Apigee client, error:\n%#v\n", e)
    return
  }
  fmt.Printf("Deleting...\n")
  // works only if no revisions are deployed
  deletedItem, resp, e := client.Proxies.Delete(proxyName)
  if e != nil {
    fmt.Printf("while deleting, error:\n%#v\n", e)
    return
  }
  fmt.Printf("status: %s\n", resp.Status)
  defer resp.Body.Close()
  fmt.Printf("deleted: %#v\n", deletedItem)
}
```

### Listing API Products

```go
func main() {
  orgPtr := flag.String("org", "", "an Edge Organization")
  flag.Parse()
  if *orgPtr == "" {
    usage()
    return
  }

  opts := &apigee.ApigeeClientOptions{Org: *orgPtr, Auth: nil, Debug: false }
  client, e := apigee.NewApigeeClient(opts)
  if e != nil {
    fmt.Printf("while initializing Edge client, error:\n%#v\n", e)
    return
  }

  fmt.Printf("\nListing...\n")
  list, resp, e := client.Products.List()
  if e != nil {
    fmt.Printf("while listing, error:\n%#v\n", e)
    return
  }
  showStatus(resp)
  fmt.Printf("products: %#v\n", list)
  resp.Body.Close()

  for _, element := range list {
    product, resp, e := client.Products.Get(element)
    if e != nil {
      fmt.Printf("while getting, error:\n%#v\n", e)
      return
    }
    showStatus(resp)
    fmt.Printf("product: %#v\n", product)
    resp.Body.Close()
  }

  fmt.Printf("\nall done.\n")
}
```

## Bugs

* The function is incomplete.

* There tests are incomplete.

* The examples are thin.

* There is no package versioning strategy (eg, no use of GoPkg.in)

* When deploying a proxy, there's no way to specify the override and delay parameters.
