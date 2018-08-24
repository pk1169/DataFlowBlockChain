package network

import (
	"DataFlowBlockChain/core"

)

type BlockConsensus struct {
	BlockValidator *core.BlockValidator // used to validate blockMsg
	State  		   *State
}

func NewBlockCunsensus(blockchain *core.BlockChain, view *View) *BlockConsensus {
	bv := core.NewBlockValidator(blockchain)
	state := NewState(view.ID)

	return &BlockConsensus{
		BlockValidator: bv,
		State: state,
	}
}

func (bkc *BlockConsensus) Refresh() {
	viewID := bkc.State.ViewID
	bkc.State = NewState(viewID)
}


//func (bc *BlockConsensus) StartConsensus(msg *Msg) {
//    // save ReqMsgs to its logs
//	bc.State.MsgLogs.ReqMsg = msg
//
//	//
//	digest := msg.Hash()
//
//}



