package network

import (
	"DataFlowBlockChain/core"
	"DataFlowBlockChain/accounts"
	"math/big"
	"github.com/Unknwon/goconfig"
	"log"
	"DataFlowBlockChain/ethdb"
	"DataFlowBlockChain/core/rawdb"
	"DataFlowBlockChain/core/types"
	"github.com/ethereum/go-ethereum/common"
	"net"
	"DataFlowBlockChain/rlp"
	"time"
	"io/ioutil"
	"crypto/rand"
	"DataFlowBlockChain/crypto"
)

type PrimaryNode struct {
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

	// when node generate block, send the block to it
	newBlock  chan *types.Block
	newTx     chan *types.Transaction

	// signal the state of consensus
	blockNotVoting  chan struct{}
	txNotVoting     chan struct{}

}


// func InitChain is to init chain when new node
func InitChain(db ethdb.Database) {
	canonicalHash := rawdb.ReadCanonicalHash(db, 0)
	genesisBlock := rawdb.ReadBlock(db, canonicalHash, 0)
	if genesisBlock == nil {
		header := &types.Header{
			Version: big.NewInt(3),
			Number: big.NewInt(0),
			Time: big.NewInt(0),
			ParentHash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		}
		block := types.NewBlock(header, types.Transactions{&types.Transaction{}}, types.VoteCollection{&types.Vote{}})

		rawdb.WriteBlock(db, block)
		rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
		rawdb.WriteTxLookupEntries(db, block)
		rawdb.WriteVtLookupEntries(db, block)
		rawdb.WriteHeadHeaderHash(db, block.Hash())
		rawdb.WriteHeadBlockHash(db, block.Hash())
	}
}


// func NewPrimaryNode is create a new Primary node
func NewPrimaryNode(nodeID string) *PrimaryNode {
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
	node := &PrimaryNode{
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
	go node.Check()


	return node
}

func (pmnode *PrimaryNode) getNodeTable() map[string]string{
	return pmnode.NodeTable
}

func (pmnode *PrimaryNode) getNodeID() string {
	return pmnode.NodeID
}

type View struct {
	ID  uint64
	Primary  string
}


func (node *PrimaryNode) handleMsg(msg *Msg) {
	if msg.Digest != msg.Hash() || msg.ViewID != node.View.ID{
		return
	}
	switch msg.Operation{
	case VoteTx:
		switch msg.State{
		case PrePrepareMsg:
			var tx = new(types.Transaction)
			rlp.DecodeBytes(msg.Payload, tx)

			if tx.VerifySig() {
				node.tc.State.MsgLogs.VoteData = tx
			}
		case PrepareMsg:
			var vote = new(types.Vote)
			rlp.DecodeBytes(msg.Payload, vote)
			checkFunc := node.tc.tv.funcMap[vote.Func]
			tx := node.tc.State.MsgLogs.VoteData.(*types.Transaction)

			if checkFunc(tx) == vote.IsInnocent {
				if vote.VerifySig() && tx.Hash() == vote.DataHash {
					node.tc.State.MsgLogs.PrepareMsgs[vote.NodeID] = vote
				}
			}

		case CommitMsg:
			var vote = new(types.Vote)
			rlp.DecodeBytes(msg.Payload, vote)
			checkFunc := node.tc.tv.funcMap[vote.Func]
			tx := node.tc.State.MsgLogs.VoteData.(*types.Transaction)

			if checkFunc(tx) == vote.IsInnocent {
				if vote.VerifySig() && tx.Hash() == vote.DataHash {
					node.tc.State.MsgLogs.CommitMsgs[vote.NodeID] = vote
				}
			}
		}

	case VoteBlock:
		switch msg.State{
		case PrePrepareMsg:
			var block = new(types.Block)
			rlp.DecodeBytes(msg.Payload, block)
			if node.bkc.BlockValidator.ValidateHeader(block) == nil {
				node.bkc.State.MsgLogs.VoteData = block
			}
		case PrepareMsg:
			var vote = new(types.Vote)
			rlp.DecodeBytes(msg.Payload, vote)
			block := node.tc.State.MsgLogs.VoteData.(*types.Block)
			if vote.VerifySig() && block.Hash() == vote.DataHash {
				node.tc.State.MsgLogs.PrepareMsgs[vote.NodeID] = vote
			}
		case CommitMsg:
			var vote = new(types.Vote)
			rlp.DecodeBytes(msg.Payload, vote)
			block := node.tc.State.MsgLogs.VoteData.(*types.Block)
			if vote.VerifySig() && block.Hash() == vote.DataHash {
				node.tc.State.MsgLogs.CommitMsgs[vote.NodeID] = vote
			}
		}
	}
}

func (node *PrimaryNode) Update() {
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

func (node *PrimaryNode) Check() {
	for {
		// if tx number > 1000
		if len(node.TxPool.VotedTxs) >= 1000 {
			go node.GenerateBlock()
		}

		select {
		case tx := <- node.TxPool.PendingTxs :
			node.newTx <- tx
		}
	}
}

func (node *PrimaryNode) BlockConsensus(block *types.Block) {

	go node.bkc.State.CheckStage()

	// new vote
	vote := &types.Vote{
		DataHash: block.Hash(),
		NodeID:   node.NodeID,
		PubKey:   node.Key.PubKey[:],
	}
	vote1, err := vote.Sign(node.Key.PrivateKey)

	if err != nil {
		log.Fatal("vote sig error", " error ", err)
	}

	for node.bkc.State.CurrentStage != Committed {
		if node.bkc.State.CurrentStage == Idle {
			msg := NewMsg(node.View.ID, VoteBlock, PrePrepareMsg, block)
			node.BroadCast(msg)
		}
        // if PrePrepared, send prepare msg
		if node.bkc.State.CurrentStage == PrePrepared {

			// log self's vote
			node.bkc.State.MsgLogs.CommitMsgs[vote1.NodeID] = vote1
			msg := NewMsg(node.View.ID, VoteBlock, PrepareMsg, vote1)
			node.BroadCast(msg)
		}

		// if Prepared, send commit msg
		if node.bkc.State.CurrentStage == Prepared {
			// log self's msg
			node.bkc.State.MsgLogs.CommitMsgs[vote1.NodeID] = vote1

			msg := NewMsg(node.View.ID, VoteBlock, CommitMsg, vote1)
			node.BroadCast(msg)
		}
	}
	node.BlockChain.InsertBlock(block)
	node.bkc.Refresh()
	node.blockNotVoting <- struct{}{}
}

func (node *PrimaryNode) Listen() {
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

func (node *PrimaryNode) TxConsensus(tx *types.Transaction) {

	go node.tc.State.CheckStage()
	c := time.Now().Unix()
	name := node.tc.tv.funcs[c%3]
	checkFunc := node.tc.tv.funcMap[crypto.Keccak256Hash([]byte(name))]
	r := checkFunc(tx)

	vote := &types.Vote{
		DataHash: tx.Hash(),
		IsExist: 1,
		IsInnocent:r,
		NodeID:node.NodeID,
		PubKey: node.Key.PubKey[:],
	}

	vote1, err := vote.Sign(node.Key.PrivateKey)
	if err != nil {
		log.Fatal("sig vote err", " error ", err)
	}

	for node.tc.State.CurrentStage != Committed {
		if node.tc.State.CurrentStage == Idle {
			msg := NewMsg(node.View.ID, VoteTx, PrePrepareMsg, tx)
			node.BroadCast(msg)
		}

		if node.tc.State.CurrentStage == PrePrepared {
			// log self's msg
			node.tc.State.MsgLogs.PrepareMsgs[node.NodeID] = vote1

			msg := NewMsg(node.View.ID, VoteTx, PrepareMsg, vote1)
			node.BroadCast(msg)
		}


		if node.tc.State.CurrentStage == Prepared {
			// log self's msg
			node.tc.State.MsgLogs.PrepareMsgs[node.NodeID] = vote1

			msg := NewMsg(node.View.ID, VoteTx, CommitMsg, tx)
			node.BroadCast(msg)
		}
	}

	// collect votes of tx
	countVote1 := make([]*types.Vote, 0)
	for _, vote := range node.tc.State.MsgLogs.CommitMsgs {
		if vote.IsExist == 1 {
			countVote1 = append(countVote1, vote)
		}
	}

	if len(countVote1) > 2*f {
		count := make(map[uint64]int)
		for _, vote := range countVote1 {
			count[vote.IsInnocent]++
		}

		if count[0] > count[1] {
			tx.SetResult(0)
		} else {
			tx.SetResult(1)
		}
	}
	node.TxPool.AddVotedTx(tx)
	node.tc.Refresh()
	node.blockNotVoting <- struct{}{}
}

func (node *PrimaryNode) BroadCast(msg *Msg) {
	for nodeID, url := range node.NodeTable {
		if nodeID != node.NodeID {
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

func (node *PrimaryNode) GenerateBlock()  {
	version := big.NewInt(3)
	timeStamp := time.Now().Unix()
	parentHash := node.BlockChain.CurrentBlock().Hash()
	pnumber := node.BlockChain.CurrentBlock().Number()
	number := pnumber.Add(pnumber, big.NewInt(1))



	header := &types.Header{
		Version: version,
		Time: big.NewInt(timeStamp),
		ParentHash: parentHash,
		Number: number,
		PubKey: node.Key.PubKey[:],
	}

	txs := node.TxPool.PopVotedTxs()
	votes := make([]*types.Vote, 0)

	for _, tx := range txs {
		votelist := node.TxPool.PopVote(tx.Hash())
		votes = append(votes, votelist...)
	}

	block := types.NewBlock(header, txs, votes)

	node.newBlock <- block
}

func (node *PrimaryNode) AddNewTx() {

}