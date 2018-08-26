package network

import (
	"DataFlowBlockChain/core/types"
	"sync"
)

type State struct {
	ViewID  		uint64
	MsgLogs 		*MsgLogs
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
	VoteData 	  interface{}
	PrepareMsgs	  map[string]*types.Vote
	CommitMsgs	  map[string]*types.Vote
}

//
const f = 1

//
func NewState(viewID uint64) *State {
	return &State{
		ViewID: viewID,
		MsgLogs: &MsgLogs{
			VoteData: nil,
			PrepareMsgs: make(map[string] *types.Vote),
			CommitMsgs:	 make(map[string] *types.Vote),
		},
		CurrentStage: Idle,
	}
}

func (state *State) SetVoteData(v interface{}) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	state.MsgLogs.VoteData = v
}

func (state *State) AddPrepareMsg(vote *types.Vote) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	state.MsgLogs.PrepareMsgs[vote.NodeID] = vote
}

func (state *State) AddCommitMsg(vote *types.Vote) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	state.MsgLogs.CommitMsgs[vote.NodeID] = vote
}

func (state *State) SetStage(stage Stage) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	state.CurrentStage = stage
}

func (state *State) CheckStage() {
	for {
		if state.MsgLogs.VoteData != nil {
			state.CurrentStage = PrePrepared
		}

		if len(state.MsgLogs.PrepareMsgs) > 2*f {
			state.CurrentStage = PrePrepared
		}

		if len(state.MsgLogs.CommitMsgs) > 2*f {
			state.CurrentStage = Committed
			break
		}
	}
}

