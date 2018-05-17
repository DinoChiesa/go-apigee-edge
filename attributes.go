package apigee

import (
 // "strconv"
 // "strings"
  "encoding/json"
  "sort"
)

// Attributes represents a revision number. Edge returns rev numbers in string form. 
// This marshals and unmarshals between that format and int.
type Attributes map[string]string


// This is just a wrapper struct to aid in serialization and de-serialization.
type PropertyWrapper struct {
  Property      Attributes  `json:"property,omitempty"`
}


// MarshalJSON implements the json.Marshaler interface. It marshals from
// an Attributes object (which is really a map[string]string) into a JSON that looks like
//    [ { "name" : "aaaaaa", "value" : "1234abcd"}, { "name" : "...", "value" : "..."} ]
func (attrs Attributes) MarshalJSON() ([]byte, error) {

  // According to the Golang spec, iterating over a map is
  // non-deterministic. This makes a problem for testing. To address that, we
  // sort the keys. The JSON returned by this fn will thus always be the same.
  var keys []string
  for k := range attrs {
    keys = append(keys, k)
  }
  sort.Strings(keys)

  var holder []map[string]string

  // iterate the map according a lexicographic sort of the keys
  for _, k := range keys {
    n := map[string]string{}
    n["name"]=k; n["value"]=attrs[k]
    holder = append(holder, n)
  }
  j,_ := json.Marshal(holder)
  
  return j, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface. It unmarshals from
// a string like "2" (including the quotes), into an integer 2.
func (attrs *Attributes) UnmarshalJSON(b []byte) error {
  var maybe []map[string]string
  e := json.Unmarshal(b, &maybe)
  if e == nil {
    //fmt.Printf("UnmarshalJSON: %v\n", maybe)
    a := Attributes{}
    for _, v := range maybe {
      if name, ok := v["name"]; ok {
        if val, ok := v["value"]; ok {
          a[name] = val
        }
      }
    }    
    *attrs = a
  }
  return nil
}

func (a Attributes) String() string {
  //return fmt.Sprintf("%v", string(a))
  v,_ := json.Marshal(a)
  return string(v)
}
