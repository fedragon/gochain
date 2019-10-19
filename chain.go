package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

// Hash represents a block's hash
type Hash string

// Transaction represents a block's tx
type Transaction string

// Block represents a block in the chain
type Block struct {
	Hash     Hash
	Tx       Transaction
	Previous *Block
	Next     *Block
}

func (b Block) String() string {
	if b.Next == nil {
		return fmt.Sprintf("{ hash: %v, tx: %v, next: nil }", b.Hash, b.Tx)
	}

	return fmt.Sprintf("{ hash: %v, tx: %v, next: %v }", b.Hash, b.Tx, b.Next)
}

// Prettify returns a human-readable string representing the contents of the block
func (b Block) Prettify(tabs int) string {
	next := ""

	if b.Next != nil {
		next = b.Next.Prettify(tabs + 1)
	}

	shortHash := fmt.Sprintf("%x", b.Hash[:4])

	tt := "\n"
	for i := 0; i < tabs; i++ {
		tt += "\t"
	}

	return fmt.Sprintf(
		"{%vHash: %v,%vTx: %v%vNext: %v%v}\n", tt, shortHash, tt, b.Tx, tt, next, tt)
}

// Ledger represents a blockchain
type Ledger struct {
	Genesis *Block
}

// IsEmpty returns true if the chain is empty, false otherwise
func (l Ledger) IsEmpty() bool {
	return l.Genesis == nil
}

func (l Ledger) String() string {
	if l.IsEmpty() {
		return "[ ]"
	}

	return fmt.Sprintf("[ genesis: %v ]\n", l.Genesis)
}

// Prettify returns a human-readable string representing the contents of the chain
func (l Ledger) Prettify() string {
	if l.IsEmpty() {
		return "[ ]"
	}

	return fmt.Sprintf("[\n%v\n]", l.Genesis.Prettify(1))
}

// NewLedger creates a ledger containing a genesis block
func NewLedger(seed Transaction) *Ledger {
	hasher := sha256.New()
	hasher.Write([]byte(seed))
	hash := Hash(fmt.Sprintf("%x", hasher.Sum(nil)))

	return &Ledger{Genesis: &Block{Hash: hash, Tx: seed, Previous: nil, Next: nil}}
}

// Add appends a block to the end of the chain
func (l *Ledger) Add(tx Transaction) (Hash, error) {
	hasher := sha256.New()

	if l.Genesis == nil {
		return "", errors.New("This ledger is empty")
	}

	block := l.Genesis
	hasher.Write([]byte(block.Hash))

	for block.Next != nil {
		hasher.Sum([]byte(block.Hash))
		block = block.Next
	}

	hash := Hash(fmt.Sprintf("%x", hasher.Sum([]byte(tx))))
	block.Next = &Block{Hash: hash, Tx: tx, Previous: block, Next: nil}

	return hash, nil
}

// Get retrieves a block from the chain, if found; returns an err
func (l *Ledger) Get(h Hash) (*Block, error) {
	if l.IsEmpty() {
		return nil, errors.New("This ledger is empty")
	}

	block := l.Genesis
	hash := block.Hash

	for block.Next != nil {
		if hash == h {
			return block, nil
		}

		block = block.Next
		hash = block.Hash
	}

	return nil, errors.New("Block not found")
}

// HashOf returns the hash of the chain (= the hash of its last block)
func (l *Ledger) HashOf() (Hash, error) {
	if l.IsEmpty() {
		return "", errors.New("This ledger is empty")
	}

	block := l.Genesis
	hash := block.Hash

	for block.Next != nil {
		block = block.Next
		hash = block.Hash
	}

	return hash, nil
}
