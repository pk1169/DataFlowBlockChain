package main

import (


	//"github.com/syndtr/goleveldb/leveldb"
	"DataFlowBlockChain/core/types"
	"math/big"
	//
	//"DataFlowBlockChain/rlp"
	//"DataFlowBlockChain/crypto"
	//"DataFlowBlockChain/common"
	//"DataFlowBlockChain/ethdb"

	//"DataFlowBlockChain/accounts"
	"log"
	"DataFlowBlockChain/crypto"
	"github.com/ethereum/go-ethereum/common"
	"fmt"
	"DataFlowBlockChain/core/rawdb"
	"DataFlowBlockChain/ethdb"
	//"time"
	"DataFlowBlockChain/core"
)

var ChainDir string
var blockChain *core.BlockChain
func init() {

}

var (
	srcAddress = "192.168.2.01"
	destAddress = "192.168.3.02"
	srcPort = big.NewInt(2000)
	destPort = big.NewInt(1002)
	protocol = big.NewInt(200)
	startTime = big.NewInt(20833)
	lastTime = big.NewInt(30900)
	size = big.NewInt(10)
)

var (
	nodeID = "Apple"
	funcName = "TimeAnalysis"
	isExist = big.NewInt(0)
	funcHash = crypto.Keccak256Hash([]byte(funcName))
)

func main() {
	//txs := make(types.Transactions, 100)
	//key, err := accounts.GetKey("accounts/keyfile")
	//if err != nil {
	//	log.Fatal("new key error ", "key error", err)
	//}
	//
	//for i := 0; i < len(txs); i++ {
	//	tx := types.NewTransaction(srcAddress, destAddress, srcPort, destPort, protocol, startTime, lastTime, size, key.PubKey[:])
	//	txs[i], err = tx.Sign(key.PrivateKey)
	//	if err != nil {
	//		log.Fatal("sign error", "error ", err)
	//	}
	//	size.Add(size, big.NewInt(1))
	//}
	//
	//votes := make(types.VoteCollection, 100)
	//
	//var txHash common.Hash
	//for i := 0; i < len(txs); i++ {
	//	txHash = txs[i].Hash()
	//	votes[i] = types.NewVote(txHash, isExist, nodeID, funcHash, key.PubKey[:])
	//}
	//
	//for i := 1; i < 10; i++ {
	//	header := &types.Header{
	//		Version: big.NewInt(3),
	//		Number: big.NewInt(int64(i)),
	//		Time: big.NewInt(time.Now().Unix()),
	//		ParentHash: blockChain.GetBlockByNumber(uint64(i-1)).Hash(),
	//	}
	//	txArray := txs[(i-1)*10: i*10]
	//	voteArray := votes[(i-1)*10: i*10]
	//	block := types.NewBlock(header, txArray, voteArray)
	//	blockChain.InsertBlock(block)
	//}
	//
	//block1 := blockChain.GetBlockByNumber(3)
	//block2 := blockChain.GetBlockByNumber(6)
	//
	//fmt.Println(block1.Header())
	//fmt.Println(block2.Header())
	ChainDir = "/Users/xiaozhang/workspace/go/src/DataFlowBlockChain/db/blockchain"
	ldb, err := ethdb.NewLDBDatabase(ChainDir, 1024, 1024)
	if err != nil {
		log.Fatal("new ldb error", "error ", err)
	}
	defer ldb.Close()
	//Create a test header to move around the database and make sure it's really new
	header := &types.Header{
		Version: big.NewInt(3),
		Number: big.NewInt(0),
		Time: big.NewInt(0),
		ParentHash: common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	}
	block := types.NewBlock(header, types.Transactions{&types.Transaction{}}, types.VoteCollection{&types.Vote{}})
	fmt.Println(block.Hash())
	rawdb.WriteBlock(ldb, block)
	fmt.Println("done")
	rawdb.WriteCanonicalHash(ldb, block.Hash(), block.NumberU64())
	rawdb.WriteTxLookupEntries(ldb, block)
	rawdb.WriteVtLookupEntries(ldb, block)
	rawdb.WriteHeadHeaderHash(ldb, block.Hash())
	rawdb.WriteHeadBlockHash(ldb, block.Hash())
}
