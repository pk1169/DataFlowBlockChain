package core

import (
	"github.com/ethereum/go-ethereum/common"
	"DataFlowBlockChain/core/types"
	"sync"
)

type TxPool struct {
	PendingTxs  chan *types.Transaction//
	VotedTxs	map[common.Hash]*types.Transaction // txs voted by nodes
	Votes       map[common.Hash](map[string]*types.Vote) // the votes of votedTxs
}

func NewTxPool() *TxPool{
	return &TxPool{
		PendingTxs: make(chan *types.Transaction, 10000),
		VotedTxs:   make(map[common.Hash]*types.Transaction),
		Votes:      make(map[common.Hash](map[string]*types.Vote)),
	}
}


// func
func (txp *TxPool) AddPendingTx(tx *types.Transaction) {
	txp.PendingTxs <- tx
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
func (txp *TxPool) PopVotedTxs() (txlist []*types.Transaction){
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	txlist = make([]*types.Transaction, 0)
	for _, tx := range txp.VotedTxs{
		txlist = append(txlist, tx)
	}
	txp.VotedTxs = make(map[common.Hash]*types.Transaction)
	return
}

func (txp *TxPool) AddTxVote(vote *types.Vote) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	txp.Votes[vote.DataHash][vote.NodeID] = vote
}

func (txp *TxPool) PopVote(hash common.Hash)  []*types.Vote {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	voteList := make([]*types.Vote, 0)
	for _, vote := range txp.Votes[hash] {
		voteList = append(voteList, vote)
	}
	txp.Votes = make(map[common.Hash](map[string]*types.Vote))
	return voteList
}







