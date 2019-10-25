package main

import (
	"fmt"
	"time"
)

// Hash represents a block's hash
type Hash string

// Data represents a block's contents
type Data string

// Block represents a block in the chain
type Block struct {
	Index     int
	Hash      Hash
	Data      Data
	Timestamp time.Time
	Previous  *Block `json:"-"`
	Next      *Block `json:",omitempty"`
}

// Create creates an unverified block for provided ledger
func Create(ledger *Ledger, data Data) (*Block, error) {
	last := ledger.Last()
	now := time.Now()

	index := 0
	if last != nil {
		index = last.Index
	}

	hash, err := CalculateHash(last, now, data)
	if err != nil {
		return nil, err
	}

	return &Block{
		Index:     index + 1,
		Hash:      hash,
		Data:      data,
		Timestamp: now,
		Previous:  last,
	}, nil
}

func (b Block) String() string {
	if b.Next == nil {
		return fmt.Sprintf("{ hash: %v, data: %v, next: nil }", b.Hash, b.Data)
	}

	return fmt.Sprintf("{ hash: %v, data: %v, next: %v }", b.Hash, b.Data, b.Next)
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
		"{%vHash: %v,%vData: %v%vNext: %v%v}\n", tt, shortHash, tt, b.Data, tt, next, tt)
}
