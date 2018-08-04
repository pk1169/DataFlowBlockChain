package core

import (
	"testing"
	"github.com/ethereum/go-ethereum/common"
	"DataFlowBlockChain/core/types"
	"math/big"
	"log"
	"DataFlowBlockChain/accounts"
	"fmt"
)

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


func TestTxPool_FindCommonTxs(t *testing.T) {
	var txs1 = make(map[common.Hash]*types.Transaction)
	var txs2 = make(map[common.Hash]*types.Transaction)
	txs := make(types.Transactions, 10000)
	key, err := accounts.GetKey("../accounts/keyfile")
	if err != nil {
		log.Fatal("new key error ", "key error", err)
	}

	for i := 0; i < len(txs); i++ {
		tx := types.NewTransaction(srcAddress, destAddress, srcPort, destPort, protocol, startTime, lastTime, size, key.PubKey[:])
		txs[i], err = tx.Sign(key.PrivateKey)

		txs1[txs[i].Hash()] = txs[i]

		if err != nil {
			log.Fatal("sign error", "error ", err)
		}
		size.Add(size, big.NewInt(1))
	}

	for i := 0; i < len(txs)-6000; i++ {

		txs2[txs[i].Hash()] = txs[i]
	}


	txpool := new(TxPool)
	bloom1 := types.TxsBloom(txs1)
	bloom2 := types.TxsBloom(txs2)
	var blooms = make([]types.Bloom, 0)
	blooms = append(blooms, bloom1, bloom2)

	commonBloom := types.CommonBloom(blooms)
	fmt.Println("common bl: ", commonBloom)
	txpool.QueuedTxs = txs1
	txpool.CommonTxs = make(map[common.Hash] *types.Transaction)
	txpool.FindCommonTxs(commonBloom)
	fmt.Println(len(txpool.QueuedTxs))
	fmt.Println(len(txpool.CommonTxs))

}



