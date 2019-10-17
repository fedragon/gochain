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

// Ledger represents a blockchain
type Ledger struct {
	Genesis *Block
}

func (l Ledger) String() string {
	if l.Genesis == nil {
		return "[ ]"
	}

	return fmt.Sprintf("[ genesis: %v ]\n", l.Genesis)
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

// Get retrieves a block from the chain, if found; returns an error otherwise
func (l *Ledger) Get(h Hash) (*Block, error) {
	if l.Genesis == nil {
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

// HashOf return the hash of the chain (= the hash of its last block)
func (l *Ledger) HashOf() (Hash, error) {
	if l.Genesis == nil {
		return "", errors.New("This ledger is empty")
	}

	var hash Hash
	block := l.Genesis
	for block.Next != nil {
		hash = block.Hash
		block = block.Next
	}

	return hash, nil
}
