package main

import (
	"errors"
	"fmt"
	"log"
)

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
func NewLedger(data Data) *Ledger {
	l := &Ledger{}
	block, err := Create(l, data)

	if err != nil {
		log.Fatalf("Could not initialize ledger: %v\n", err)
	}

	return &Ledger{Genesis: block}
}

// Append appends a block to the end of the chain
func (l *Ledger) Append(next *Block) {
	last := l.Last()

	if last == nil {
		l.Genesis = next
	}

	last.Next = next
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

// Last returns the last block in the ledger
func (l *Ledger) Last() *Block {
	if l.IsEmpty() {
		return nil
	}

	block := l.Genesis
	for block.Next != nil {
		block = block.Next
	}

	return block
}

// HashOf returns the hash of the chain (= the hash of its last block)
func (l *Ledger) HashOf() Hash {
	last := l.Last()

	if last == nil {
		return ""
	}

	return last.Hash
}
