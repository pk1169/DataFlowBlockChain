package network

import (
	"testing"
	"DataFlowBlockChain/core/types"
	"fmt"
)

func TestNewMsg(t *testing.T) {
	var bloom types.Bloom
	msg := NewMsg(0,VoteMsg,0,bloom)
	fmt.Println(msg.MsgCode)
}
