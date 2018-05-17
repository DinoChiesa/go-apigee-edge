package apigee

import (
  "encoding/json"
  "testing"
  "bytes"
)

const (
  attrsJson1 = `[ {
    "name" : "access",
    "value" : "private"
  } ]`
  attrsJson2 = `[ {
    "name" : "access",
    "value" : "private"
  } , {
    "name" : "creator",
    "value" : "Brahma"
  } ]`
  attrsJson3 = `[ {
    "name" : "access",
    "value" : "private"
  }, {
    "name" : "creator",
    "value" : "Brahma"
  }, {
    "name" : "lastModified",
    "value" : "Wednesday,  7 September 2016, 14:45"
  } ]`
)

var (
  attrsMap1 = Attributes{"access":"private"}
  attrsMap2 = Attributes{"access":"private", "creator":"Brahma"}
  attrsMap3 = Attributes{"access":"private", "creator":"Brahma","lastModified":"Wednesday,  7 September 2016, 14:45"}
)

func TestAttributes_Unmarshal(t *testing.T) {
  testCases := []struct {
    desc         string
    data         string
    expected     Attributes
    expectError  bool
    equal        bool
  }{
    {"one member   ", attrsJson1, attrsMap1, false, true},
    {"two members  ", attrsJson2, attrsMap2, false, true},
    {"three members", attrsJson3, attrsMap3, false, true},
  }
  for _, tc := range testCases {
    var got Attributes
    err := json.Unmarshal([]byte(tc.data), &got)
    t.Logf("%s: got=%v", tc.desc, got)
    if gotErr := err != nil; gotErr != tc.expectError {
      t.Errorf("%s: gotErr=%v, expectError=%v, err=%v", tc.desc, gotErr, tc.expectError, err)
      continue
    }
    equal := got.String() == tc.expected.String()
    if equal != tc.equal {
      t.Errorf("%14s: value[got=%#v, expected=%#v], equal[got=%v, expected=%v]",
        tc.desc, got, tc.expected, equal, tc.equal)
    }
  }
}


func TestAttributes_Marshal(t *testing.T) {
  testCases := []struct {
    desc         string
    data         Attributes
    expected     string
    expectError  bool
    equal        bool
  }{
    {"case 1 ", attrsMap1, attrsJson1, false, true},
    {"case 2 ", attrsMap2, attrsJson2, false, true},
    {"case 3 ", attrsMap3, attrsJson3, false, true},
  }
  
  for _, tc := range testCases {
    out, err := json.Marshal(tc.data)
    if gotErr := (err != nil); gotErr != tc.expectError {
      t.Errorf("%s: gotErr=%v, expectError=%v, err=%v", tc.desc, gotErr, tc.expectError, err)
    }
    got := string(out)
    buffer := new(bytes.Buffer)
    if err := json.Compact(buffer, []byte(tc.expected)); err != nil {
      t.Errorf("%s: %v", tc.desc, err)
    }
    equal := got == buffer.String() 
    if equal != tc.equal {
      t.Errorf("%s: value[actual=%s, expected=%s], equal[actual=%v, expected=%v]",
        tc.desc, got, tc.expected, equal, tc.equal)
    }
  }
}


// func TestTimstamp_MarshalReflexivity(t *testing.T) {
//   testCases := []struct {
//     desc string
//     data Timestamp
//   }{
//     {"Reference", Timestamp{referenceTime}},
//     {"WorkStart", Timestamp{workStartDate}},
//     {"UnixOrigin", Timestamp{unixOriginTime}},
//     {"Empty", Timestamp{}}, // degenerate case.  I don't really care about this; it will never happen.
//   }
//   for _, tc := range testCases {
//     data, err := json.Marshal(tc.data)
//     if err != nil {
//       t.Errorf("%s: Marshal err=%v", tc.desc, err)
//     }
//     var got Timestamp
//     err = json.Unmarshal(data, &got)
//     t.Logf("%s: %+v ?= %s", tc.desc, got, string(data))
//     if got.String() != tc.data.String() {
//       t.Errorf("%s: %+v != %+v", tc.desc, got, data)
//     }
//   }
// }

