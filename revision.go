package apigee

import (
  "strconv"
  "strings"
  "fmt"
)

// Revision represents a revision number. Edge returns rev numbers in string form. 
// This unmarshals them as int. 
type Revision int

func (r *Revision) MarshalJSON() ([]byte, error) {
  rev := fmt.Sprintf("%d", r)
  return []byte(rev), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Revision) UnmarshalJSON(b []byte) error {
  rev, e := strconv.ParseInt(strings.TrimSuffix(strings.TrimPrefix(string(b),"\""),"\""), 10, 32)
  if e != nil {
    return e
  }

  *r = Revision(rev)
  return nil
}

func (r Revision) String() string {
  return fmt.Sprintf("%d", r)
}

