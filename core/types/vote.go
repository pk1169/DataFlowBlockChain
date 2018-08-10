package types

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"fmt"
	"crypto/ecdsa"
	"DataFlowBlockChain/crypto"
	"log"
	"go-ethereum1/common/hexutil"
)

//go:generate gencodec -type Vote -field-override VoteMarshaling -out gen_vote_json.go

type Vote struct {
	// 在上块的时候
	DataHash 		common.Hash `json:"txHash"	gencodec:"required"`

	// following three attributes are for TxVote
	IsExist		*big.Int	`json:"isExist" gencodec:"required"`
	NodeID		string		`json:"nodeID"	gencodec:"required"`
	Func		common.Hash	`json:"func"	gencodec:"required"`

	V 			*big.Int	`json:"v" 		gencodec:"required"`
	R 			*big.Int	`json:"r"		gencodec:"required"`
	S    		*big.Int	`json:"s"		gencodec:"required"`
	PubKey 		[]byte   	`json:"pubKey"	gencodec:"required"`// 压缩公钥，用来签名或者对签名进行验证
}

type VoteMarshaling struct {
	IsExist 	*hexutil.Big
	V            *hexutil.Big
	R            *hexutil.Big
	S            *hexutil.Big
	Hash       common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
}


func NewVote(txHash common.Hash, isExist *big.Int, nodeID string, funcHash common.Hash, pubKey []byte) *Vote {
	return &Vote{
		DataHash: txHash,
		IsExist: isExist,
		NodeID: nodeID,
		Func: funcHash,
		PubKey: pubKey,
	}
}

// WithSignature returns a new Header with the given signature.
// This signature needs to be formatted).
func (vote *Vote) WithSignature(sig []byte) (*Vote, error) {
	if len(sig) != 65 {
		panic(fmt.Sprint("wrong size for signature: got %d, want 65", len(sig)))
	}
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	v := new(big.Int).SetBytes([]byte{sig[64]})
	cpy := CopyVote(vote)
	cpy.R, cpy.S, cpy.V = r, s, v
	return cpy, nil
}

// func Sign reutruns a new signed Vote from nosig Vote
func (vote *Vote) Sign(privKey *ecdsa.PrivateKey) (*Vote, error) {
	hash := vote.HashNoSig()
	sig, err := crypto.Sign(hash[:], privKey)

	if err != nil {
		log.Fatal("sig error", "error", err)
	}

	return vote.WithSignature(sig)
}

func CopyVote(vote *Vote) (*Vote) {
	cpy := &Vote{
		DataHash: vote.DataHash,
		IsExist: vote.IsExist,
		NodeID: vote.NodeID,
		Func: vote.Func,
		PubKey: vote.PubKey,
	}
	return cpy
}

func (v *Vote) Hash() common.Hash {
	return rlpHash(v)
}

// HashNoNonce returns the hash which is used as input for the proof-of-work search.
func (v *Vote) HashNoSig() common.Hash {
	return rlpHash([]interface{}{
		v.DataHash,
		v.IsExist,
		v.NodeID,
		v.Func,
	})
}

