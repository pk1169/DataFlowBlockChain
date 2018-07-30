package core

import "DataFlowBlockChain/core/types"

type Validator interface {
	// ValidateBody validates the given block's content.
	ValidateBody(block *types.Block) error

	// ValidateState validates the given statedb and optionally the receipts and
	// gas used.
	ValidateHeader(block *types.Block) error
}
