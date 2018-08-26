package core

import (
	"testing"
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




func TestTxPool_AddPendingTx(t *testing.T) {
	txp := NewTxPool()
	key, err := accounts.GetKey("../accounts/keyfile")
	if err != nil {
		log.Fatal("new key error ", "key error", err)
	}

	for i := 0; i < 10; i++ {
		tx := types.NewTransaction(srcAddress, destAddress, srcPort, destPort, protocol, startTime, lastTime, size, key.PubKey[:])
		tx, err = tx.Sign(key.PrivateKey)
		if err != nil {
			log.Fatal("sign error", "error ", err)
		}
		txp.AddPendingTx(tx)
		size.Add(size, big.NewInt(1))
	}

	close(txp.PendingTxs)
	for pendingTx := range txp.PendingTxs {
		fmt.Println(pendingTx)
		if _, ok := txp.TxLog[pendingTx.Hash()]; ok {
			fmt.Println("tx is exist")
		}
	}

}

func TestTxPool_AddVotedTx(t *testing.T) {
	txp := NewTxPool()
	key, err := accounts.GetKey("../accounts/keyfile")
	if err != nil {
		log.Fatal("new key error ", "key error", err)
	}

	for i := 0; i < 10; i++ {
		tx := types.NewTransaction(srcAddress, destAddress, srcPort, destPort, protocol, startTime, lastTime, size, key.PubKey[:])
		tx, err = tx.Sign(key.PrivateKey)
		if err != nil {
			log.Fatal("sign error", "error ", err)
		}
		txp.AddPendingTx(tx)
		size.Add(size, big.NewInt(1))
	}

	close(txp.PendingTxs)
	for pendingTx := range txp.PendingTxs {
		fmt.Println(pendingTx)
	}

}

