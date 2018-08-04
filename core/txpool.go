package core

import (
	"github.com/ethereum/go-ethereum/common"
	"DataFlowBlockChain/core/types"
	"sync"
)

type TxPool struct {
	PendingTxs	map[common.Hash]*types.Transaction // txs received from network
	QueuedTxs	map[common.Hash]*types.Transaction // txs used to generate bloom
	CommonTxs 	map[common.Hash]*types.Transaction // txs filtered by bloom
	VotedTxs	map[common.Hash]*types.Transaction // txs voted by nodes
	Bloom  		types.Bloom


	CommonBloomCh  chan types.Bloom // inform txpool to find common txs
	QueuedEvent	 chan struct{}  // inform the txpool to queue txs from pendingTxs
	StartVoteEvent  chan struct{} // to sign txs in txpool prepared for voting
}

// func
func (txp *TxPool) AddPendingTx(transaction *types.Transaction) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	txp.PendingTxs[transaction.Hash()] = transaction
}

func (txp *TxPool) RemovePendingTx(txHash common.Hash) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	delete(txp.PendingTxs, txHash)
}

//
func (txp *TxPool) GenerateQueuedTxs() {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	for k, v := range txp.PendingTxs {
		txp.QueuedTxs[k] = v
	}
	txp.PendingTxs = make(map[common.Hash] *types.Transaction)
}

func (txp *TxPool) AddingQueuedTx(hash common.Hash, transaction *types.Transaction) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	txp.QueuedTxs[hash] = transaction
}

func (txp *TxPool) RemoveQueuedTx(hash common.Hash) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	delete(txp.QueuedTxs, hash)
}

func (txp *TxPool) FindCommonTxs(commonBloom types.Bloom) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()
	for hash, tx := range txp.QueuedTxs {
		isCommon := types.BloomLookup(commonBloom, hash)
		//// @ mode
		//fmt.Println(isCommon)
		if isCommon {
			txp.CommonTxs[hash] = tx
		}
	}
}



func (txp *TxPool) RemoveCommonTx(hash common.Hash) {
	var lock sync.RWMutex
	lock.Lock()
	defer lock.Unlock()

}

func (txp *TxPool) GenerateTxBloom() types.Bloom{
	txp.Bloom  = types.TxsBloom(txp.QueuedTxs)
	return txp.Bloom
}


func (txp *TxPool) Update() {
	select {
	case <- txp.QueuedEvent:
		txp.GenerateQueuedTxs()
		case commonBloom := <-txp.CommonBloomCh:
			txp.FindCommonTxs(commonBloom)
	}
}