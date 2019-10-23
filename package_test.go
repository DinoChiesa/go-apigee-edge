package apigee

import (
  "fmt"
	"time"
  "math/rand"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	pretag = "go-test-"
	orgName = "gaccelerate3"
)

func init() {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
}

func wait(delay int) {  
  fmt.Printf("Waiting %ds...\n", delay)
  time.Sleep(time.Duration(delay)*time.Second)
}

func randomString(length int) string {
    b := make([]byte, length)
    for i := range b {
        b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
    }
    return string(b)
}

