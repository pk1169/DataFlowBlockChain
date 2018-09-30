package network

import (
	"DataFlowBlockChain/core"
	"DataFlowBlockChain/accounts"
	"DataFlowBlockChain/core/types"
	"github.com/Unknwon/goconfig"
	"log"
	"crypto/rand"
	"DataFlowBlockChain/ethdb"
	"io/ioutil"
	"net"
	"DataFlowBlockChain/rlp"
)

// CommonNode：对消息的处理机制，当收到的是Pre-prepare 区块消息的时候，那么
type CommonNode struct {
	NodeID     string
	NodeTable  map[string]string // key:nodeID  value: url
	View       *View
	TxPool     *core.TxPool
	Key        *accounts.Key
	BlockChain *core.BlockChain

	// when msg accepted, MsgEntrance store it
	MsgEntrance  chan *Msg

	// following three are consensus engine
	tc        *TxConsensus
	bkc       *BlockConsensus

	// when node receive a block, send it to the newBlock
	newBlock  chan *types.Block
	newTx     chan *types.Transaction

	// signal the state of consensus
	blockNotVoting  chan struct{}
	txNotVoting     chan struct{}

}

func NewCommonNode(nodeID string) *CommonNode {
	const viewID = 10000000000
	cfg, err := goconfig.LoadConfigFile("./config/config.ini")
	if err != nil {
		log.Fatal("load config file error", " error ", err)
	}
	keyList := cfg.GetKeyList("NodeTable")
	nodeTable := make(map[string]string)
	//
	for _, key := range keyList {
		value, err := cfg.GetValue("NodeTable", key)
		if err != nil {
			log.Fatal("key or section does not exist", "error", err)
		}
		nodeTable[key] = value
	}

	// read key, if not exist, generate key
	keyfile := "./db/key/keyfile"

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

	// read database
	chainDir, err := cfg.GetValue("Database", "ChainDir")
	if err != nil {
		log.Fatal("read chainDir error", "error ", err)
	}
	ldb, err := ethdb.NewLDBDatabase(chainDir, 1024, 1024)
	if err != nil {
		log.Fatal("instance ldb error", "error ", err)
	}
	InitChain(ldb)
	blockChain, err := core.NewBlockChain(ldb)
	if err != nil {
		log.Fatal("new blockchain error", "error ", err)
	}

	view := &View{
		ID: viewID,
		Primary: "Apple",
	}
	node := &CommonNode{
		NodeID:    	nodeID,
		NodeTable: 	nodeTable,
		View: view,
		Key: key,
		TxPool: core.NewTxPool(),
		BlockChain: blockChain,
		// Channels
		MsgEntrance: make(chan *Msg),
		bkc:     	 NewBlockCunsensus(blockChain, view),
		tc:          NewTxConsensus(view),
		newBlock:    make(chan *types.Block),
		newTx:       make(chan *types.Transaction),



		blockNotVoting: make(chan struct{}, 1),
		txNotVoting: make(chan struct{}, 1),
	}
	node.blockNotVoting <- struct{}{}
	node.txNotVoting <- struct{}{}


	go node.Listen()
	go node.Update()


	return node
}

func (node *CommonNode) Listen() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", node.NodeTable[node.NodeID]) //获取一个tcpAddr
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr) //监听一个端口
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		data, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Fatal("read conn data error", "error", err)
		}

		var msg *Msg
		rlp.DecodeBytes(data, msg)
		node.MsgEntrance <- msg
	}
}

func (node *CommonNode) Update() {
	for {
		select {
		case msg := <- node.MsgEntrance:
			node.handleMsg(msg)
		case block := <- node.newBlock:
			for {
				select {
				case <- node.blockNotVoting:
					node.bkc.Refresh()
					node.bkc.State.MsgLogs.VoteData = block
					go node.BlockConsensus(block)
				}
			}

		case tx := <- node.newTx:
			for {
				select {
				case <- node.txNotVoting:
					node.tc.Refresh()
					node.tc.State.MsgLogs.VoteData = tx
					go node.TxConsensus(tx)
				}
			}
		}
	}
}




