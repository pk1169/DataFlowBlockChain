package network

import (
	"DataFlowBlockChain/core"
	"DataFlowBlockChain/accounts"
	"math/big"
	"net"
	"github.com/Unknwon/goconfig"
	"DataFlowBlockChain/log"
	"DataFlowBlockChain/ethdb"
	"crypto/rand"
	"sync"
)

type Node struct {
	NodeID     string
	NodeTable  map[string]string // key:nodeID  value: url
	View       *View
	TxPool     *core.TxPool
	Key        *accounts.Key
	BlockChain *core.BlockChain
	Peers      Peers
}

func NewNode(nodeID string) *Node {
	const viewID = 10000000000
	cfg, err := goconfig.LoadConfigFile("./config/config.ini")
	if err != nil {
		log.Fatal("load config file error", " error ", err)
	}
	keyList := cfg.GetKeyList("NodeTable")
	nodeTable := make(map[string]string)

	// read chainDir
	chainDir, err1 := cfg.GetValue("Database", "ChainDir")
	if err1 != nil {
		log.Fatal("dir key or section does not exist", "error", err)
	}

	//
	ldb, err2 := ethdb.NewLDBDatabase(chainDir, 1024, 1024)
	if err2 != nil {
		log.Fatal("new ldb err", "error", err)
	}
	bc, err3 := core.NewBlockChain(ldb)
	if err3 != nil {
		log.Fatal("new blockchain error", "error", err3)
	}

	for _, key := range keyList {
		value, err := cfg.GetValue("NodeTable", key)
		if err != nil {
			log.Fatal("key or section does not exist", "error", err)
		}
		nodeTable[key] = value
	}

	// read key, if not exist, generate key
	// read chainDir
	keyfile, err4 := cfg.GetValue("Database", "KeyFile")
	if err4 != nil {
		log.Fatal("dir key or section does not exist", "error", err)
	}

	key, err5 := accounts.GetKey(keyfile)
	if err5 != nil {
		log.Fatal("read key error", "error", err5)
	}
	if key == nil {
		var err error
		key, err = accounts.NewKey(rand.Reader)
		if err != nil {
			log.Fatal("new key error", "error", err)
		}
		_ = accounts.StoreKey(keyfile, key)
	}


	node := &Node{
		NodeID:    	nodeID,
		NodeTable: 	nodeTable,
		View: &View{
			ID: big.NewInt(viewID),
			Primary: keyList[0],
		},
		Key: key,
		TxPool: core.NewTxPool(),
		BlockChain: bc,
		Peers: make(Peers),
	}
	return node
}


type Peers map[string] *Peer

type Peer struct {
	NodeID string
	conn  net.Conn
}



func (peer *Peer) Close() {
	peer.conn.Close()
}

func (peer *Peer) SetConn(conn net.Conn) {
	peer.conn = conn
}

func (peers Peers) AddPeer(peer *Peer) {
	var lock sync.RWMutex
	defer lock.Unlock()
	peers[peer.NodeID] = peer
}

func (peers Peers) RemovePeer(nodeID string) {
	var lock sync.RWMutex
	defer lock.Unlock()
	peer := peers[nodeID]
	peer.Close()
	delete(peers, nodeID)
}

type View struct {
	ID  *big.Int
	Primary  string
}

