package core

import (
	"fmt"
	"DataFlowBlockChain/core/types"
	"errors"
	"DataFlowBlockChain/crypto"
)

type BlockValidator struct {
	bc 		*BlockChain
}

func NewBlockValidator(blockchain *BlockChain) *BlockValidator {
	bv := &BlockValidator{
		bc:	blockchain,
	}
	return bv
}

func (bv *BlockValidator) ValidateBody(block *types.Block) error {
	// Check whether the block's known, and if not, that it's linkable
	if bv.bc.HasBlock(block.Hash(), block.NumberU64()) {
		return errors.New("Error Known block")
	}
	if !bv.bc.HasBlock(block.ParentHash(), block.NumberU64()-1) {
		if !bv.bc.HasBlock(block.ParentHash(), block.NumberU64()-1) {
			return errors.New("Error Known parent block")
		}
	}
	// Header validity is known at this point, check the uncles and transactions
	header := block.Header()

	if hash := types.DeriveSha(block.Transactions()); hash != header.TxHash {
		return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash, header.TxHash)
	}

	if hash := types.DeriveSha(block.VoteCollection()); hash != header.VotesRoot {
		return fmt.Errorf("votes root hash mismatch: have %x, want %x", hash, header.VotesRoot)
	}

	return nil
}

func (bv *BlockValidator) ValidateHeader(block *types.Block) error {
	// Check wether the header's signature is correct, and if not, it's bad block
	hash := block.Header().HashNoSig()
	pubKey := block.Header().PubKey

	sig := make([]byte, 64)
	r, s := block.Header().R.Bytes(), block.Header().S.Bytes()
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)

	isCorrect := crypto.VerifySignature(pubKey, hash[:], sig[:])

	if !isCorrect {
		return errors.New("header's sig is not correct")
	}

	return nil
}