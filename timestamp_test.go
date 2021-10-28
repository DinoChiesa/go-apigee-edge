package apigee

import (
	"encoding/json"
	//"fmt"
	"testing"
	"time"
)

const (
	lastModifiedMs   = `1444426707423`
	originMs         = `0`
	workStartMs      = `1343752707000`
	referenceTimeStr = `1473275339334`
)

var (
	lastModifiedTime = time.Date(2015, 10, 9, 21, 38, 27, 423*1000000, time.UTC)
	unixOriginTime   = time.Unix(0, 0).In(time.UTC)
	workStartDate    = time.Date(2012, 7, 31, 16, 38, 27, 0, time.UTC)
	referenceTime    = time.Date(2016, 9, 7, 19, 8, 59, 334*1000000, time.UTC)
)

func TestTimestamp_Marshal(t *testing.T) {
	testCases := []struct {
		desc     string
		data     Timestamp
		expected string
		wantErr  bool
		equal    bool
	}{
		{"lastModified ", Timestamp{lastModifiedTime}, lastModifiedMs, false, true},
		{"origin       ", Timestamp{unixOriginTime}, originMs, false, true},
		{"workStartDate", Timestamp{workStartDate}, originMs, false, false},
		{"workStartDate", Timestamp{workStartDate}, workStartMs, false, true},
	}

	for _, tc := range testCases {
		out, err := json.Marshal(tc.data)
		if gotErr := (err != nil); gotErr != tc.wantErr {
			t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
		}
		got := string(out)
		equal := got == tc.expected
		if (got == tc.expected) != tc.equal {
			t.Errorf("%s: value[actual=%s, expected=%s], equal[actual=%v, expected=%v]", tc.desc, got, tc.expected, equal, tc.equal)
		}
	}
}

func TestTimestamp_Unmarshal(t *testing.T) {
	testCases := []struct {
		desc     string
		data     string
		expected Timestamp
		wantErr  bool
		equal    bool
	}{
		{"Reference    ", referenceTimeStr, Timestamp{referenceTime}, false, true},
		{"Mismatch     ", referenceTimeStr, Timestamp{}, false, false},
	}
	for _, tc := range testCases {
		var got Timestamp
		err := json.Unmarshal([]byte(tc.data), &got)
		t.Logf("%s: got=%v", tc.desc, got)
		t.Logf("%s: got=%v", tc.desc, got.Time.String())
		if gotErr := err != nil; gotErr != tc.wantErr {
			t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
			continue
		}
		equal := got.Equal(tc.expected)
		if equal != tc.equal {
			t.Errorf("%s: values[got=%#v, expected=%#v], equal[got=%v, expected=%v]", tc.desc, got, tc.expected, equal, tc.equal)
		}
	}
}

func TestTimestamp_MarshalReflexivity(t *testing.T) {
	testCases := []struct {
		desc string
		data Timestamp
	}{
		{"Reference", Timestamp{referenceTime}},
		{"WorkStart", Timestamp{workStartDate}},
		{"UnixOrigin", Timestamp{unixOriginTime}},
		{"Empty", Timestamp{}}, // degenerate case.  I don't really care about this; it will never happen.
	}
	for _, tc := range testCases {
		data, err := json.Marshal(tc.data)
		if err != nil {
			t.Errorf("%s: Marshal err=%v", tc.desc, err)
		}
		var got Timestamp
		err = json.Unmarshal(data, &got)
		t.Logf("%s: %+v ?= %s", tc.desc, got, string(data))
		if got.String() != tc.data.String() {
			t.Errorf("%s: %+v != %+v", tc.desc, got, data)
		}
	}
}
