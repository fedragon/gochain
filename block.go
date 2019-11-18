package gochain

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

// Create creates an unverified block for provided chain
func Create(chain *Chain, data Data) (*Block, error) {
	last := chain.Last()
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
		return fmt.Sprintf("{ index: %v, hash: %v, data: %v, timestamp: %v, next: nil }", b.Index, b.Hash, b.Data, b.Timestamp)
	}

	return fmt.Sprintf("{ index: %v, hash: %v, data: %v, timestamp: %v, next: %v }", b.Index, b.Hash, b.Data, b.Timestamp, b.Next)
}
