package accounts

import (
	"testing"
	"crypto/rand"
	"log"
	"fmt"
)

func TestStoreKey(t *testing.T) {
	key, err := NewKey(rand.Reader)
	if err != nil {
		log.Fatal("new key error ", "error", err)
	}

	StoreKey("./keyfile", key)
}

func TestGetKey(t *testing.T) {
	key, err := GetKey("./keyfile")
	if err != nil {
		log.Fatal("new key error ", "error", err)
	}
	fmt.Println(key)
}
