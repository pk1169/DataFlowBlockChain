package types

import (
	"math/big"

	"DataFlowBlockChain/rlp"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common"
	"sync/atomic"
	"io"
	"fmt"
	"crypto/ecdsa"
	"DataFlowBlockChain/crypto"
	"log"
)


//go:generate gencodec -type txdata -field-override txdataMarshaling -out gen_tx_json.go
// 交易是流量数据信息的载体
type txdata struct {
	// Signature values
	V *big.Int	`json:"v"	gencodec:"required"`
	R *big.Int	`json:"r"	gencodec:"required"`
	S *big.Int	`json:"s"	gencodec:"required"`
	Abnormal  uint64 `json:"abnormal gencodec:"required"`


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

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type txdataMarshaling struct {
	StartTime       *hexutil.Big
	LastTime		*hexutil.Big
	Size   			*hexutil.Big
	V            *hexutil.Big
	R            *hexutil.Big
	S            *hexutil.Big
}

type Transaction struct {
	data txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

func NewTransaction(srcAddress, destAddress string, srcPort, destPort, protocol, startTime, lastTime, size *big.Int, pubKey []byte) *Transaction {
	return newTransaction(srcAddress ,destAddress, srcPort, destPort, protocol, startTime, lastTime, size, pubKey)
}

func newTransaction(srcAddress, destAddress string, srcPort, destPort, protocol, startTime, lastTime, size *big.Int, pubKey []byte) *Transaction {
	d := txdata{
		SrcAddress:  srcAddress,
		DestAddress: destAddress,
		SrcPort:     new(big.Int),
		DestPort:    new(big.Int),
		Protocol:    new(big.Int),
		StartTime:   new(big.Int),
		LastTime:    new(big.Int),
		Size:		 new(big.Int),
		PubKey: 	 pubKey,
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
	if srcPort != nil {
		d.SrcPort.Set(srcPort)
	}
	if destPort != nil {
		d.DestPort.Set(destPort)
	}
	if protocol != nil {
		d.Protocol.Set(protocol)
	}
	if startTime != nil {
		d.StartTime.Set(startTime)
	}
	if lastTime != nil {
		d.LastTime.Set(lastTime)
	}
	if size != nil {
		d.Size.Set(size)
	}
	if pubKey != nil && len(pubKey) == 33 {
		copy(d.PubKey, pubKey)
	}

	return &Transaction{data: d}
}

func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
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

// EncodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &tx.data)
}

// DecodeRLP implements rlp.Decoder
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&tx.data)
	if err == nil {
		tx.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

// MarshalJSON encodes the web3 RPC transaction format.
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	hash := tx.Hash()
	data := tx.data
	data.Hash = &hash
	return data.MarshalJSON()
}

// UnmarshalJSON decodes the web3 RPC transaction format.
func (tx *Transaction) UnmarshalJSON(input []byte) error {
	var dec txdata
	if err := dec.UnmarshalJSON(input); err != nil {
		return err
	}

	*tx = Transaction{data: dec}
	return nil
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func (tx *Transaction) WithSignature(sig []byte) (*Transaction, error) {
	if len(sig) != 65 {
		panic(fmt.Sprint("wrong size for signature: got %d, want 65", len(sig)))
	}
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	v := new(big.Int).SetBytes([]byte{sig[64]})
	cpy := &Transaction{data: tx.data}
	cpy.data.R, cpy.data.S, cpy.data.V = r, s, v
	return cpy, nil
}

// Sign returns signed transaction from rawtransaction
func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey) (*Transaction, error){
	hash := tx.HashNoSig()
	sig, err := crypto.Sign(hash[:], privKey)

	if err != nil {
		log.Fatal("sig error", "error", err)
	}

	return tx.WithSignature(sig)
}

func (tx *Transaction) RawSignatureValues() (*big.Int, *big.Int, *big.Int) {
	return tx.data.V, tx.data.R, tx.data.S
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

func (tx *Transaction) VerifySig() bool{
	sig := make([]byte, 64)
	r, s := tx.data.R.Bytes(), tx.data.S.Bytes()
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)

	hash := tx.HashNoSig()
	isCorrect := crypto.VerifySignature(tx.data.PubKey, hash[:], sig[:])

	return isCorrect
}

func (tx *Transaction) SetResult(r uint64) {
	tx.data.Abnormal = r
}

