package accounts

import (
	"testing"
	"crypto/rand"
	"log"
)

func TestStoreKey(t *testing.T) {
	key, err := NewKey(rand.Reader)
	if err != nil {
		log.Fatal("new key error ", "error", err)
	}

	StoreKey("./keyfile", key)
}
