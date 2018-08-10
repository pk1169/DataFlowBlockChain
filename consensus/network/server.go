package network

import (
	"DataFlowBlockChain/consensus"
	"fmt"
	"os"
	"net"
)

type Server struct {
	url string
	node *Node
    msgBuffer MsgBuffer

}

type MsgBuffer chan *consensus.Msg


func NewServer(nodeID string) *Server {
	node := NewNode(nodeID)
	url := node.NodeTable[nodeID]

	return &Server{
		node: node,
		url : url,
		msgBuffer: make(MsgBuffer, 1000),
	}
}

func (server *Server) EstablishConnect() {
	nodeTable := (server.node).NodeTable
	for nodeID, url := range nodeTable {
		if nodeID != (server.node).NodeID {
			if _, ok := server.node.Peers[nodeID]; !ok {
				tcpAddr, err := net.ResolveTCPAddr("tcp4", url) //获取一个TCP地址信息,TCPAddr
				checkError(err)
				conn, err := net.DialTCP("tcp", nil, tcpAddr) //创建一个TCP连接:TCPConn
				checkError(err)

				// init peer, set it into peers
				peer := new(Peer)
				peer.NodeID = nodeID
				peer.SetConn(conn)
				server.node.Peers[nodeID] = peer
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}