package types

import (
	"testing"
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

var json1 []byte

func TestHeader_MarshalJSON(t *testing.T) {
	var hash1 common.Hash
	header := Header{
		Version:3,
		ParentHash:hash1,
		TxHash:hash1,
		VotesRoot:hash1,
		Time:big.NewInt(3444444),
		Number:big.NewInt(344444),
	}
	json, _ := header.MarshalJSON()
	json1 = json
	fmt.Println(string(json))
	header1 := new(Header)
	header1.UnmarshalJSON(json1)
	fmt.Println(header1)
}

func TestHeader_UnmarshalJSON(t *testing.T) {
	header := new(Header)

	header.UnmarshalJSON(json1)

	fmt.Println(header)
}
