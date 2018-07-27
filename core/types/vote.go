package types

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
)

type Vote struct {
	// 在上块的时候
	TxHash 		common.Hash `json:"txHash"	gencodec:"required"`

	IsExist		*big.Int	`json:"isExist" gencodec:"required"`
	NodeID		string		`json:"nodeID"	gencodec:"required"`
	Func		common.Hash	`json:"func"	gencodec:"required"`

	V 			*big.Int	`json:"v" 		gencodec:"required"`
	R 			*big.Int	`json:"r"		gencodec:"required"`
	S    		*big.Int	`json:"s"		gencodec:"required"`
	PubKey 		[]byte   	`json:"pubKey"	gencodec:"required"`// 压缩公钥，用来签名或者对签名进行验证
}

func (v *Vote) Hash() common.Hash {
	return rlpHash(v)
}

// HashNoNonce returns the hash which is used as input for the proof-of-work search.
func (v *Vote) HashNoSig() common.Hash {
	return rlpHash([]interface{}{
		v.TxHash,
		v.IsExist,
		v.NodeID,
		v.Func,
	})
}

