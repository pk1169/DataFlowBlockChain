package network

import "DataFlowBlockChain/core"

type Node struct {
	NodeID  	string
	NodeTable	map[string]string
	View 		*View
	TxPool 		*core.TxPool

}

type View struct {
	ID  int64
	Primary  string
}

