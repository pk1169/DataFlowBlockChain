package core

import (
	"github.com/ethereum/go-ethereum/common"
	"DataFlowBlockChain/core/types"
)

type TxPool struct {
	PendingTx	map[common.Hash]*types.Transaction
	CommonTx 	map[common.Hash]*types.Transaction
	VotedTx		map[common.Hash]*types.Transaction

	CommonBloomCh  chan<- *
}

func (txp *TxPool) AddPendingTx(transaction *types.Transaction) {
	txp.PendingTx[transaction.Hash()] = transaction
}

func (txp *TxPool) RemovePendingTx(transaction *types.Transaction) {
	delete(txp.PendingTx, transaction.Hash())
}

func (txp *TxPool) AddCo
