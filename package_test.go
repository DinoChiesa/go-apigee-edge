package apigee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	testPrefix     = "go-test"
	testConfigFile = "testdata/test_config.json"
)

var apigeeClient ApigeeClient

var testSettings struct {
	Orgname string `json:"orgname"`
	Notes   string `json:"notes"`
}

func init() {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	file, e := ioutil.ReadFile(testConfigFile)
	if e != nil {
		fmt.Printf("reading configuration file: %#v\n", e)
	}

	e = json.Unmarshal(file, &testSettings)
	if e != nil {
		fmt.Printf("unmarshaling configuration: %#v\n", e)
	}
}

func NewClientForTesting(t *testing.T) *ApigeeClient {
	opts := &ApigeeClientOptions{Org: testSettings.Orgname, Auth: nil, Debug: false}
	client, e := NewApigeeClient(opts)
	if e != nil {
		t.Errorf("while initializing Edge client, error:\n%#v\n", e)
		return nil
	}
	return client
}

func wait(delay int) {
	fmt.Printf("Waiting %ds...\n", delay)
	time.Sleep(time.Duration(delay) * time.Second)
}

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
