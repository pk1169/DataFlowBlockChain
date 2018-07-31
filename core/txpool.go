package core

import (
	"github.com/ethereum/go-ethereum/common"
	"DataFlowBlockChain/core/types"
)

type TxPool struct {
	PendingTx	map[common.Hash]*types.Transaction
	CommonTx 	map[common.Hash]*types.Transaction
	VotedTx		map[common.Hash]*types.Transaction
}
