package network

import (
	"github.com/ethereum/go-ethereum/common"
	"DataFlowBlockChain/crypto/sha3"
	"DataFlowBlockChain/rlp"
	"DataFlowBlockChain/core/types"
	"log"
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
	PrePrepareMsg = 0x20
	PrepareMsg = 0x21
	CommitMsg = 0x22
)

type Msg struct {
	ViewID      uint64
	MsgCode 	uint64
	Operation  	uint64
	State 		uint64
	Payload     []byte
	Digest 		common.Hash
}

func (msg *Msg) Hash() common.Hash {
	return rlpHash([]interface{}{
		msg.ViewID,
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

func NewMsg(viewID, operation, state uint64, v interface{}) *Msg {
	msg := new(Msg)
	msg.ViewID = viewID
	msg.Operation = operation
	msg.State = state

	var err error
	msg.Payload, err = rlp.EncodeToBytes(v)
	if err != nil {
		log.Fatal("encode payload error", "error", err)
	}
	switch v.(type) {
	case *types.Block:
		msg.MsgCode = BlockMsg
	case *types.Transaction:
		msg.MsgCode = TxMsg
	case *types.Vote:
		msg.MsgCode = VoteMsg
	case types.Bloom:
		msg.MsgCode = BloomMsg
	}

	msg.Digest = msg.Hash()
	return msg
}