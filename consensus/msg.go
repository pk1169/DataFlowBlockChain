package consensus

import (
	"github.com/ethereum/go-ethereum/common"
	"DataFlowBlockChain/crypto/sha3"
	"DataFlowBlockChain/rlp"
)

// the value of MsgCode, flag the type of msg payload
const (
	TxMsg = 0x01
	BloomMsg = 0x02
	BlockMsg = 0x03
	VoteMsg = 0x04
)

// the value of Operation, flag the type of operation
const (
	SyncBlock = 0x10
	VoteBloom = 0x11
	VoteTx = 0x12
	VoteBlock = 0x13
)

// the value of MsgState, flag the state fo vote
const (
	PrepareMsg = 0x20
	CommitMsg = 0x21
)

type Msg struct {
	MsgCode 	uint64
	Operation  	uint64
	State 		uint64
	Payload     []byte
	Digest 		common.Hash
}

func (msg *Msg) Hash() common.Hash {
	return rlpHash([]interface{}{
		msg.MsgCode,
		msg.Operation,
		msg.State,
		msg.Payload,
	})
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}