package accounts

import (
	"crypto/ecdsa"
	"github.com/pborman/uuid"
	"encoding/hex"
	"DataFlowBlockChain/crypto"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"io/ioutil"
)

type PubKey []byte

type Key struct {
	Id 	uuid.UUID
	PrivateKey  	*ecdsa.PrivateKey
	PubKey			PubKey
}


type plainKeyJSON struct {
	PubKey string `json:"pubKey"`
	PrivateKey 	string `json:"privateKey"`
	Id 			string `json:"id"`
}

func (k *Key) MarshalJSON() (j []byte, err error) {
	jStruct := plainKeyJSON{
		hex.EncodeToString(k.PubKey[:]),
		hex.EncodeToString(crypto.FromECDSA(k.PrivateKey)),
		k.Id.String(),
	}
	j, err = json.Marshal(jStruct)
	return j, err
}

func (k *Key) UnMarshalJSON(j []byte) (err error) {
	keyJSON := new(plainKeyJSON)
	err = json.Unmarshal(j, &keyJSON)
	if err != nil {
		return err
	}

	u := new(uuid.UUID)
	*u = uuid.Parse(keyJSON.Id)
	k.Id = *u
	pubKey, err := hex.DecodeString(keyJSON.PubKey)
	if err != nil {
		return err
	}

	privKey, err := crypto.HexToECDSA(keyJSON.PrivateKey)
	if err != nil {
		return err
	}

	k.PubKey = pubKey
	k.PrivateKey = privKey

	return nil
}

func newKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *Key {
	id := uuid.NewRandom()
	key := &Key{
		Id:         id,
		PubKey:    crypto.CompressPubkey(&privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}

func NewKey(rand io.Reader) (*Key, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
	if err != nil {
		return nil, err
	}
	return newKeyFromECDSA(privateKeyECDSA), nil
}

func writeKeyFile(file string, content []byte) error {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}
	f.Close()
	return os.Rename(f.Name(), file)
}