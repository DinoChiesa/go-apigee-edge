package apigee

import (
  "encoding/json"
  "testing"
  "bytes"
)

const (
  cacheJson1 = `{
   "expiryDate": { "value": "09-22-2016" },
   "valuesNull" : false
 }`
  
  cacheJson2 = `{
   "timeoutInSec" : { "value" : "300" },  
   "valuesNull" : false
 }`
  
  cacheJson3 = `{
   "timeOfDay": { "value" : "14:30:00" },
   "valuesNull" : false
 }`
)


var (
  cacheExpiry1 = CacheExpiry{"expiryDate","09-22-2016",false}
  cacheExpiry2 = CacheExpiry{"timeoutInSec","300",false}
  cacheExpiry3 = CacheExpiry{"timeOfDay","14:30:00",false}
)


func TestCacheExpiry_Unmarshal(t *testing.T) {
  testCases := []struct {
    desc      string
    data      string
    expected  CacheExpiry
    wantErr   bool
    equal     bool
  }{
    {"cacheJson1", cacheJson1, cacheExpiry1, false, true},
    {"cacheJson2", cacheJson2, cacheExpiry2, false, true},
    {"cacheJson3", cacheJson3, cacheExpiry3, false, true},
  }
  for _, tc := range testCases {
    var got CacheExpiry
    err := json.Unmarshal([]byte(tc.data), &got)
    t.Logf("%s: got=%v", tc.desc, got)
    if gotErr := err != nil; gotErr != tc.wantErr {
      t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
      continue
    }
    equal := got.String() == tc.expected.String()
    if equal != tc.equal {
      t.Errorf("%14s: value[got=%#v, expected=%#v], equal[got=%v, expected=%v]",
        tc.desc, got, tc.expected, equal, tc.equal)
    }
  }
}


func TestCacheExpiry_Marshal(t *testing.T) {
  testCases := []struct {
    desc      string
    data      CacheExpiry
    expected  string
    wantErr   bool
    equal     bool
  }{
    {"cacheJson1", cacheExpiry1, cacheJson1, false, true},
    {"cacheJson2", cacheExpiry2, cacheJson2, false, true},
    {"cacheJson3", cacheExpiry3, cacheJson3, false, true},
  }
  
  for _, tc := range testCases {
    out, err := json.Marshal(tc.data)
    if gotErr := (err != nil); gotErr != tc.wantErr {
      t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
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

