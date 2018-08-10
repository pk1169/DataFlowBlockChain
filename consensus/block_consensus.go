package consensus

import "DataFlowBlockChain/core"

type BlockConsensus struct {
	BlockValidator *core.BlockValidator
	State  		   *State
}



