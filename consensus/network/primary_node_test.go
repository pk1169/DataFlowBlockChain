package network

import (
	"testing"
	"fmt"
)

func TestNewNode(t *testing.T) {
	node := NewNode("Apple")
	fmt.Println(node.NodeTable)
}
