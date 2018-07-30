package accounts

import (
	"crypto/ecdsa"
	"github.com/pborman/uuid"
)

type Key struct {
	Id 	uuid.UUID
	PrivateKey  	ecdsa.PrivateKey
	PubKey			string
}

type keyStore interface {
	//
}


type plainKeyJSON struct {
	PubKey string `json:"pubKey"`
	PrivateKey 	string `json:"privateKey"`

}