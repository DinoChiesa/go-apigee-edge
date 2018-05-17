package apigee

import (
  "fmt"
	"time"
  "math/rand"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

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

