package network

import (
	"DataFlowBlockChain/core/types"
	"github.com/ethereum/go-ethereum/common"
	"DataFlowBlockChain/crypto"
)

//
//
//// State notes the state of consenesus
//type State struct {
//	ViewID 		int64
//	MsgLogs  	*MsgLogs
//}
//
//// MsgLogs logs the msg for consensus
//types MsgLogs
//
//
//// Stage marks the present stage of consensus
//type Stage int
//
//const (
//	Idle  Stage = iota
//	PrePrepared
//	Prepared
//	Committed
//)
//
//// the number of nodes can be tolerated to be attacked
//const f = 1
//
//// la
//// func CreateState(viewID int64, lastSequenceID  int64)

type TxConsensus struct {
	State  *State
	tv *TxVoteValidator
}
func NewTxConsensus(view *View) *TxConsensus {
	return &TxConsensus{
		State: NewState(view.ID),
		tv:    NewTxVoteValidator(),
	}
}

func (txc *TxConsensus) Refresh() {
	viewID := txc.State.ViewID
	txc.State = NewState(viewID)
}

type CheckFunc func(txvote *types.Vote) uint64

type TxVoteValidator struct {
	funcMap map[common.Hash] CheckFunc
}

func NewTxVoteValidator() *TxVoteValidator{
	funcTable := make(map[common.Hash] CheckFunc)
	funcs := []string{"Check1", "Check2", "Check3"}
	tv := new(TxVoteValidator)
	for _, f := range funcs {
		hash := crypto.Keccak256Hash([]byte(f))
		switch f {
		case "Check1":
			funcTable[hash] = Check1
		case "Check2":
			funcTable[hash] = Check2
		case "Check3":
			funcTable[hash] = Check3
		}
	}
	tv.funcMap = funcTable
	return tv
}

func Check1(txvote *types.Vote) uint64 {

	return 1
}

func Check2(txvote *types.Vote) uint64 {

	return 1
}

func Check3(txvote *types.Vote) uint64 {

	return 1
}

