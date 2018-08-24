package network

import (
	"fmt"
	"os"
	"net"
	"DataFlowBlockChain/rlp"
	"log"
	"io/ioutil"
)

type Server struct {
	url string
	node
    msgBuffer MsgBuffer
}

type MsgBuffer chan *Msg


func NewServer(nodeID string) *Server {
	node := NewPrimaryNode(nodeID)
	url := node.NodeTable[nodeID]

	return &Server{
		node: node,
		url : url,
		msgBuffer: make(MsgBuffer, 1000),
	}
}

func (server *Server) EstablishConnect() {
}

func (server *Server) ListenAndServer() {

}


//
//func (server *Server) DialNode() {
//	nodeTable := (server.node).NodeTable
//	loop:
//	for {
//		if len(server.node.NodeTable) == len(server.node.Peers) {
//			break loop
//		}
//		for nodeID, url := range nodeTable {
//			if nodeID != (server.node).NodeID {
//				if _, ok := server.node.Peers[nodeID]; !ok {
//					tcpAddr, err := net.ResolveTCPAddr("tcp4", url) //获取一个TCP地址信息,TCPAddr
//					checkError(err)
//					conn, err := net.DialTCP("tcp", nil, tcpAddr) //创建一个TCP连接:TCPConn
//					checkError(err)
//					conn.Write([]byte(server.node.NodeID))
//					peer := new(Peer)
//					peer.NodeID = nodeID
//					peer.SetConn(conn)
//					server.node.Peers[nodeID] = peer
//				}
//			}
//		}
//	}
//}
//
//func (server *Server) Listen() {
//	tcpAddr, err := net.ResolveTCPAddr("tcp4", server.url) //获取一个tcpAddr
//	checkError(err)
//	listener, err := net.ListenTCP("tcp", tcpAddr) //监听一个端口
//	checkError(err)
//
//	for {
//		conn, err := listener.Accept()
//		if err != nil {
//			continue
//		}
//
//		peer := new(Peer)
//		peer.NodeID = nodeID
//		peer.SetConn(conn)
//		server.node.Peers[nodeID] = peer
//		if len(server.node.Peers) == len(server.node.NodeTable) {
//			break loop
//		}
//	}
//}

func (server *Server) BroadCast(msg *Msg) {
	for nodeID, url := range server.node.getNodeTable() {
		if nodeID != server.node.getNodeID() {
			tcpAddr, err := net.ResolveTCPAddr("tcp4", url) //获取一个TCP地址信息,TCPAddr
			checkError(err)
			conn, err := net.DialTCP("tcp", nil, tcpAddr) //创建一个TCP连接:TCPConn
			checkError(err)
			data, err := rlp.EncodeToBytes(msg)
			if err != nil {
				log.Fatal("mas rlp error", "error", err)
			}
			conn.Write(data)
			conn.Close()
		}
	}
}

func (server *Server) Msg(msg *Msg) {

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}