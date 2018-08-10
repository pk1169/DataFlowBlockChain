package consensus

import (
	"math/big"
)

type State struct {
	ViewID  		*big.Int
	MsgLogs 		*MsgLogs
	LastSequenceID 	*big.Int
	CurrentStage 	Stage
}

type Stage uint64
const (
	Idle  Stage = iota  // Node is created
	PrePrepared       	// Node is pre-prepared
	Prepared			// Node is Prepared
	Committed 			// Node is Committed
)

type MsgLogs struct {
	ReqMsg 		  *Msg
	PrepareMsgs	  map[string]*Msg
	CommitMsgs	  map[string]*Msg
}

//
const f = 1

//
func NewState(viewID, lastSequenceID *big.Int) *State {
	return &State{
		ViewID: viewID,
		MsgLogs: &MsgLogs{
			ReqMsg: nil,
			PrepareMsgs: make(map[string] *Msg),
			CommitMsgs:	 make(map[string] *Msg),
		},
		LastSequenceID: lastSequenceID,
		CurrentStage: Idle,
	}
}