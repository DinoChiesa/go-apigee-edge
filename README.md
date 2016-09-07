# golang client library for Apigee Edge administrative API

Use this from Go-lang programs to invoke administrative operations on Apigee Edge.

## Copyright and License

This code is [Copyright (c) 2016 Apigee Corp](NOTICE). it is licensed under the [Apache 2.0 Source Licese](LICENSE).

## Status

This is a work-in-progress.  Currently the only entity implemented is apiproxies.


## Usage Examples

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
  userPtr := flag.String("user", "", "an administrator in that Edge Organization")
  flag.Parse()

  if *namePtr != "" {
    proxyName = *namePtr
  } 
  
  if *srcPtr == "" || *userPtr == "" || *orgPtr == "" {
    usage()
    return
  }
  
  auth := apigee.EdgeAuth{Username: *userPtr}
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

