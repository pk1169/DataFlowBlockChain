package core

import (
	"github.com/ethereum/go-ethereum/common"
	"DataFlowBlockChain/core/types"
	"sync"
	"fmt"
)

type TxPool struct {
	PendingTxs  chan *types.Transaction//
	// Txlog is log if tx exist in txpool
	TxLog       map[common.Hash]struct{}
	// votedTxs store txs voted by consensus engine
	VotedTxs	map[common.Hash]*types.Transaction // txs voted by nodes
	Votes       map[common.Hash](map[string]*types.Vote) // the votes of votedTxs
}

func NewTxPool() *TxPool{
	return &TxPool{
		PendingTxs: make(chan *types.Transaction, 10000),
		VotedTxs:   make(map[common.Hash]*types.Transaction),
		Votes:      make(map[common.Hash](map[string]*types.Vote)),

		TxLog:      make(map[common.Hash]struct{}),
	}
}

func (txp *TxPool) LogTx(hash common.Hash) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	txp.TxLog[hash] = struct{}{}
}

func (txp *TxPool) UnLogTx(hash common.Hash) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	delete(txp.TxLog, hash)
}

// func
func (txp *TxPool) AddPendingTx(tx *types.Transaction) {
	txp.PendingTxs <- tx
	txp.LogTx(tx.Hash())
}


//
func (txp *TxPool) PopPendingTx() *types.Transaction {
	tx := <- txp.PendingTxs
	return tx
}

// func AddVotedTx is to add voted tx to txp
func (txp *TxPool) AddVotedTx(tx *types.Transaction) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()

	txp.VotedTxs[tx.Hash()] = tx
}

func (txp *TxPool) ReadVotedTx(hash common.Hash) *types.Transaction {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()

	return txp.VotedTxs[hash]
}

//
func (txp *TxPool) PopVotedTxs() ([]*types.Transaction){
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	txlist := make([]*types.Transaction, 0)
	for _, tx := range txp.VotedTxs{
		txlist = append(txlist, tx)
		txp.UnLogTx(tx.Hash())
	}

	// refresh the store space of votedTxs
	txp.VotedTxs = make(map[common.Hash]*types.Transaction)
	return txlist
}

func (txp *TxPool) AddTxVote(vote *types.Vote) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()

	if txp.Votes[vote.DataHash] == nil {
		txp.Votes[vote.DataHash] = make(map[string]*types.Vote)
		txp.Votes[vote.DataHash][vote.NodeID] = vote

	} else {
		txp.Votes[vote.DataHash][vote.NodeID] = vote
	}

	// @mode
	fmt.Println(txp.Votes[vote.DataHash][vote.NodeID])
}

func (txp *TxPool) PopVote(hash common.Hash)  []*types.Vote {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	voteList := make([]*types.Vote, 0)
	voteMap := txp.Votes[hash]
	for _, vote := range voteMap {
		voteList = append(voteList, vote)
	}
	return voteList
}

func (txp *TxPool) RefreshVotes() {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	txp.Votes = make(map[common.Hash](map[string]*types.Vote))
}






