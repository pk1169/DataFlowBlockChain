package core

import "github.com/ethereum/go-ethereum/common"

type TxPool struct {
	PendingTx	map[common.Hash]Transaction
	VotedTx		map[common.Hash]Transaction

}
