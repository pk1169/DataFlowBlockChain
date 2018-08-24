package network

import "DataFlowBlockChain/core/types"

type BloomConsensus struct {
	Blooms []*types.Bloom
}

func NewBloomConsensus() *BloomConsensus {
	blooms := make([]*types.Bloom, 0)
	return &BloomConsensus{
		Blooms: blooms,
	}
}

func (blc *BloomConsensus) Refresh() {
	blooms := make([]*types.Bloom, 0)
	blc.Blooms = blooms
}