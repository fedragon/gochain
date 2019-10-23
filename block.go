package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

// Hash represents a block's hash
type Hash string

// Data represents a block's contents
type Data string

// Block represents a block in the chain
type Block struct {
	Index    int
	Hash     Hash
	Data     Data
	Previous *Block
	Next     *Block
}

// Create creates an unverified block for provided ledger
func Create(ledger *Ledger, data Data) (*Block, error) {
	var index int
	var lastHash Hash
	var last *Block

	if !ledger.IsEmpty() {
		last = ledger.Genesis

		for last.Next != nil {
			last = last.Next
		}

		index = last.Index
		lastHash = last.Hash
	}

	hash, err := hashOf(index, lastHash, time.Now(), data)

	if err != nil {
		return nil, err
	}

	return &Block{
		Index:    index + 1,
		Hash:     hash,
		Data:     data,
		Previous: last,
	}, nil
}

func hashOf(index int, previous Hash, timestamp time.Time, data Data) (Hash, error) {
	hasher := sha256.New()

	hasher.Write([]byte(strconv.Itoa(index)))
	hasher.Write([]byte(previous))
	hasher.Write([]byte(data))
	ts, err := timestamp.In(time.UTC).MarshalBinary()

	if err != nil {
		return "", err
	}

	hasher.Write(ts)
	hash := Hash(fmt.Sprintf("%x", hasher.Sum(nil)))

	return hash, nil
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
