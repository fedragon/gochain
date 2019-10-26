package main

import (
	"errors"
	"fmt"
	"log"
)

// Chain represents a blockchain
type Chain struct {
	Genesis *Block
}

// IsEmpty returns true if the chain is empty, false otherwise
func (l Chain) IsEmpty() bool {
	return l.Genesis == nil
}

func (l Chain) String() string {
	if l.IsEmpty() {
		return "[ ]"
	}

	return fmt.Sprintf("[ genesis: %v ]\n", l.Genesis)
}

// Prettify returns a human-readable string representing the contents of the chain
func (l Chain) Prettify() string {
	if l.IsEmpty() {
		return "[ ]"
	}

	return fmt.Sprintf("[\n%v\n]", l.Genesis.Prettify(1))
}

// NewChain creates a chain containing a genesis block
func NewChain(data Data) *Chain {
	l := &Chain{}
	block, err := Create(l, data)

	if err != nil {
		log.Fatalf("Could not initialize chain: %v\n", err)
	}

	return &Chain{Genesis: block}
}

// Append appends a block to the end of the chain
func (l *Chain) Append(next *Block) {
	last := l.Last()

	if last == nil {
		l.Genesis = next
	}

	last.Next = next
}

// Get retrieves a block from the chain, if found; returns an err
func (l *Chain) Get(h Hash) (*Block, error) {
	if l.IsEmpty() {
		return nil, errors.New("This chain is empty")
	}

	block := l.Genesis
	hash := block.Hash

	if hash == h {
		return block, nil
	}

	for block.Next != nil {
		if hash == h {
			return block, nil
		}

		block = block.Next
		hash = block.Hash
	}

	if hash == h {
		return block, nil
	}

	return nil, nil
}

// Last returns the last block in the chain
func (l *Chain) Last() *Block {
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
func (l *Chain) HashOf() Hash {
	last := l.Last()

	if last == nil {
		return ""
	}

	return last.Hash
}
