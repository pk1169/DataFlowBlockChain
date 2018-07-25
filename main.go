package main

import (

	"fmt"

	//"github.com/syndtr/goleveldb/leveldb"
	"DataFlowBlockChain/core/types"
	"math/big"
	//
	//"DataFlowBlockChain/rlp"
	//"DataFlowBlockChain/crypto"
	//"DataFlowBlockChain/common"
	//"DataFlowBlockChain/ethdb"

	"DataFlowBlockChain/ethdb"
	"DataFlowBlockChain/core/rawdb"
)



func main() {
	//privKey, err1 := crypto.GenerateKey()
	//
	//if err1 != nil {
	//	fmt.Println(err1)
	//}
	//
	//data := [][]byte{
	//	[]byte{1, 2, 3, 4},
	//	[]byte{5, 6, 7, 8},
	//}
	//
	//rlpdata, _ := rlp.EncodeToBytes(data)
	//fmt.Println(rlpdata)
	//hash := crypto.Keccak256(rlpdata)
	//
	//sig, err2 := crypto.Sign(hash, privKey)
	//
	//if err2 != nil {
	//	fmt.Println(err2)
	//}
	//
	//fmt.Println(len(sig))
	//
	//db, err := leveldb.OpenFile("./db", nil)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer db.Close()
	//
	//data1, err3 := db.Get(hash, nil)
	//if err3 != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(data1)
	//}
	//
	//var source [][]byte
	//rlp.DecodeBytes(data1, &source)
	//fmt.Println(source)
	//
	//var hash1 common.Hash
	//c := append(hash1[:0], hash...)
	//fmt.Println(c)
	//header := types.Header{
	//	Version:3,
	//	ParentHash:hash1,
	//	TxHash:hash1,
	//	VotesRoot:hash1,
	//	Time:big.NewInt(3444444),
	//	Number:big.NewInt(344444),
	//}
	//rawHash := header.HashNoSig()
	//sig1, err6 := crypto.Sign(rawHash[:], privKey)
	//if err6 != nil {
	//	fmt.Println(err6)
	//}
	////fmt.Println(sig1[64])
	//header.V = new(big.Int).SetBytes(sig1[64:])
	//header.R = new(big.Int).SetBytes(sig1[:32])
	//header.S = new(big.Int).SetBytes(sig1[32:64])


	//
	//pubKey := crypto.CompressPubkey(&privKey.PublicKey)
	//
	//isCorrect := crypto.VerifySignature(pubKey, rawHash[:], sig1[:64])
	//fmt.Println(isCorrect)
	//
	//json, _ := header.MarshalJSON()
	//fmt.Println(string(json))

	//db, err := ethdb.NewLDBDatabase("./ldb", 0, 0)
	//if err != nil {
	//	fmt.Println(err)
	//}

	var txs types.Transactions

	txs = types.Transactions{
		&types.Transaction{
			V:big.NewInt(0),
			R:big.NewInt(0),
			S:big.NewInt(0),

			SrcAddress:"11fdfdf",
			DestAddress:"22fdf",
		},
		&types.Transaction{
			V:big.NewInt(0),
			R:big.NewInt(3),
			S:big.NewInt(0),

			SrcAddress:"11fdfdf",
			DestAddress:"22fdf",
		},
	}

	votes := types.VoteCollection{
		&types.Vote{
			V:big.NewInt(0),
			R:big.NewInt(1),
			S:big.NewInt(2),

			NodeID:"Apple",
			TxIndex:big.NewInt(0),
		},
		&types.Vote{
			V:big.NewInt(0),
			R:big.NewInt(3),
			S:big.NewInt(4),

			NodeID:"IBM",
			TxIndex:big.NewInt(1),
		},
	}
	//txshash := types.DeriveSha(txs)
	//voteshash := types.DeriveSha(votes)
	//fmt.Println(hash)

	//
	//db := ethdb.NewMemDatabase()
	ldb, _ := ethdb.NewLDBDatabase("./db", 128, 1024)
	defer ldb.Close()

	 //Create a test header to move around the database and make sure it's really new
	header := &types.Header{Number: big.NewInt(42)}
	block := types.NewBlock(header, txs, votes)
	fmt.Println(block.Hash())
	rawdb.WriteBlock(ldb, block)
	fmt.Println("done")
	rawdb.WriteTxLookupEntries(ldb, block)
	rawdb.WriteVtLookupEntries(ldb, block)
	blockHash, blockNumber, txIndex:= rawdb.ReadTxLookupEntry(ldb, txs[1].Hash())
//	tx0, _, _, _ := rawdb.ReadTransaction(ldb, txs[0].Hash())
//	fmt.Println(tx0)
	vote0, _, _, _ := rawdb.ReadVote(ldb, txs[0].Hash(), "Apple")
	fmt.Println(vote0)

	fmt.Println(blockHash, "\n blocknumber: ", blockNumber, "index: ", txIndex)
	readheader := rawdb.ReadHeader(ldb, block.Hash(), 42)
	fmt.Println(readheader.VotesRoot)
}
