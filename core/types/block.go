package types

import (
	"math/big"

	"DataFlowBlockChain/crypto/sha3"
	"DataFlowBlockChain/rlp"
	"sync/atomic"
	"time"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
)

var (
	EmptyRootHash  = DeriveSha(Transactions{})
	EmptyVotesHash = DeriveSha(VoteCollection{})
)

//go:generate gencodec -type Header -field-override headerMarshaling -out gen_header_json.go
type Header struct {
	Version *big.Int 			`json:"version" gencodec:"required"`
	ParentHash common.Hash 	`json:"parentHash" gencodec:"required"`
	TxHash common.Hash 		`json:"txHash" gencodec:"required"`
	VotesRoot common.Hash 	`json:"votesRoot" gencodec:"required"`
	Time *big.Int 			`json:"time" gencodec:"required"`
	Number *big.Int 		`json:"number" gencodec:"required"`

	// Signature values of the node who generates the block
	V *big.Int 				`json:"v" 		gencodec:"required"`
	R *big.Int 				`json:"r" 		gencodec:"required"`
	S *big.Int 				`json:"s" 		gencodec:"required"`
	PubKey []byte			`json:"pubKey"	gencodec:"required"`
}

// field type overrides for gencodec
type headerMarshaling struct {
	Number     *hexutil.Big
	Time       *hexutil.Big
	V            *hexutil.Big
	R            *hexutil.Big
	S            *hexutil.Big
	Hash       common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
}

func (h *Header) Hash() common.Hash {
	return rlpHash(h)
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// HashNoNonce returns the hash which is used as input for the proof-of-work search.
func (h *Header) HashNoSig() common.Hash {
	return rlpHash([]interface{}{
		h.Version,
		h.ParentHash,
		h.TxHash,
		h.VotesRoot,
		h.Time,
		h.Number,
	})
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}


// Block represents an entire block in the Ethereum blockchain.
type Block struct {
	header       *Header
	transactions Transactions
	votes		VoteCollection

	// caches
	hash atomic.Value
	size atomic.Value


	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

func (b *Block) Transaction(hash common.Hash) *Transaction {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}

func (b *Block) Transactions() Transactions { return b.transactions }
func (b *Block) VoteCollection() VoteCollection {return b.votes}

// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the
// block.
//
// The values of TxHash, UncleHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs, uncles
// and receipts.
func NewBlock(header *Header, txs Transactions, votes VoteCollection) *Block {
	b := &Block{header: CopyHeader(header), }

	if len(txs) == 0 {
		b.header.TxHash = EmptyRootHash
	} else {
		b.header.TxHash = DeriveSha(txs)
		b.transactions = make(Transactions, len(txs))
		copy(b.transactions, txs)
	}

	if len(votes) == 0 {
		b.header.VotesRoot = EmptyVotesHash
	} else {
		// need to
		b.header.VotesRoot = crypto.Keccak256Hash()
		b.votes = make(VoteCollection, len(votes))
		copy(b.votes, votes)
	}



	return b
}


// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	cpy := *h
	if cpy.Time = new(big.Int); h.Time != nil {
		cpy.Time.Set(h.Time)
	}
	if cpy.Number = new(big.Int); h.Number != nil {
		cpy.Number.Set(h.Number)
	}

	cpy.Version = h.Version

	return &cpy
}

// VoteCollection collects the vote
type VoteCollection []*Vote

// Len returns the length of s.
func (s VoteCollection) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s VoteCollection) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (s VoteCollection) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

type Body struct {
	Transactions 	Transactions
	Votes			VoteCollection
}


// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: CopyHeader(header)}
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

// ParentHash returns the parentHash of this block
func (b *Block) ParentHash() common.Hash {
	return b.header.ParentHash
}

// WithBody returns a new block with the given transaction and uncle contents.
func (b *Block) WithBody(transactions []*Transaction, votes VoteCollection) *Block {
	block := &Block{
		header:       CopyHeader(b.header),
		transactions: make([]*Transaction, len(transactions)),
		votes:       make(VoteCollection, len(VoteCollection{})),
	}
	copy(block.transactions, transactions)
	for k := range votes {
		block.votes[k] = votes[k]
	}

	return block
}


// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Header *Header
	Txs    []*Transaction
	Votes VoteCollection
}


// EncodeRLP serializes b into the Ethereum RLP block format.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extblock{
		Header: b.header,
		Txs:    b.transactions,
		Votes: b.votes,
	})
}

// Body returns the non-header content of the block.
func (b *Block) Body() *Body { return &Body{b.transactions, b.votes} }

func (b *Block) Header() *Header { return CopyHeader(b.header) }
func (b *Block) Number() *big.Int     { return new(big.Int).Set(b.header.Number) }
func (b *Block) NumberU64() uint64        { return b.header.Number.Uint64() }
//func (b *Vote) NumberU64()	uint64 {return b.TxIndex.Uint64()}