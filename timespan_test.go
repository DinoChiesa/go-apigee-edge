package apigee

import (
  "encoding/json"
  //"fmt"
  "testing"
)


func TestTimespan_New(t *testing.T) {
  testCases := []struct {
    input     string
    expected  string
  }{
    {"1s", "1000"},
    {"1m", "60000"},
    {"3m", "180000"},
    {"1h", "3600000"},
    {"1d", "86400000"},
    {"10d", "864000000"},
  }

  for _, tc := range testCases {
    ts := NewTimespan(tc.input)
    // if (e != nil) {
    //   t.Errorf("%s: error=%v", tc.input, e)
    // }
    //actualOutput := ts.String()
    actualOutput, e := json.Marshal(ts)
    if (e != nil) {
      t.Errorf("%s: error=%v", tc.input, e)
    }
    got := string(actualOutput)
    if got != tc.expected {
      t.Errorf("%s: value[actual=%s, expected=%s]", tc.input, got, tc.expected)
    }
  }
}
