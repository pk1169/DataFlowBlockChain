package accounts

import (
	"encoding/json"
	"os"
)


func GetKey(filename string) (*Key, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	key := new(Key)
	if err := json.NewDecoder(fd).Decode(key); err != nil {
		return nil, err
	}

	return key, nil
}

func StoreKey(filename string, key *Key) error {
	content, err := json.Marshal(key)
	if err != nil {
		return err
	}
	return writeKeyFile(filename, content)
}

