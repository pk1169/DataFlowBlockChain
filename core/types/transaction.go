package types

import (
	"math/big"

	"DataFlowBlockChain/rlp"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common"
	"sync/atomic"
)


//go:generate gencodec -type txdata -field-override txdataMarshaling -out gen_tx_json.go
// 交易是流量数据信息的载体
type txdata struct {
	// Signature values
	V *big.Int	`json:"v"	gencodec:"required"`
	R *big.Int	`json:"r"	gencodec:"required"`
	S *big.Int	`json:"s"	gencodec:"required"`

	// 以下是流量数据的信息
	SrcAddress string	`json:"srcAddress"	gencodec:"required"`
	DestAddress string	`json:"destAddress"	gencodec:"required"`
	SrcPort *big.Int			`json:"srcPort" 	gencodec:"required"`
	DestPort *big.Int		`json:"destPort"	gencodec:"required"`
	Protocol *big.Int		`json:"protocol"	gencodec:"required"`
	StartTime *big.Int	`json:"startTime"	gencodec:"required"`
	LastTime *big.Int	`json:"lastTime"	gencodec:"required"`
	Size *big.Int		`json:"size"		gencodec:"required"` // 流量大小

	PubKey 		[]byte	`json:"pubKey"		gencodec:"required"`
}

type txdataMarshaling struct {
	StartTime       *hexutil.Big
	LastTime		*hexutil.Big
	Size   			*hexutil.Big
	V            *hexutil.Big
	R            *hexutil.Big
	S            *hexutil.Big
	// This is only used when marshaling to JSON.
	Hash common.Hash `json:"hash"`
}

type Transaction struct {
	data txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

func (t *Transaction) Hash() common.Hash {
	return rlpHash(t)
}


// HashNoNonce returns the hash which is used as input for the proof-of-work search.
func (t *Transaction) HashNoSig() common.Hash {
	return rlpHash([]interface{}{
		t.data.SrcAddress,
		t.data.DestAddress,
		t.data.SrcPort,
		t.data.DestPort,
		t.data.Protocol,
		t.data.StartTime,
		t.data.LastTime,
		t.Size,
	})
}

// Size returns the true RLP encoded storage size of the transaction, either by
// encoding and returning it, or returning a previsouly cached value.
func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &tx.data)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// DataInformation is the information of dataflow
type DataInformation struct {
	SrcAddress string
	DestAddress string
	SrcPort int
	DestPort int
	Protocol int
	StartTime big.Int
	LastTime big.Int

	Size big.Int
}

// Transactions is a Transaction slice type for basic sorting.
type Transactions []*Transaction

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}


