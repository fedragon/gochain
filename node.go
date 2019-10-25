package main

import (
	"fmt"
)

// Node represents a node submitting new blocks to the ledger
type Node struct {
	Updates    <-chan Ledger
	Unverified <-chan Block
	Verified   chan<- Block
}

// verifies a block before it's accepted in the ledger
func verify(ledger *Ledger, unverified Block) bool {
	last := ledger.Last()

	hash, err := hashOf(last.Index, last.Hash, unverified.Timestamp, unverified.Data)
	if err != nil {
		return false
	}

	return hash == unverified.Hash
}

// Run executes the main loop of a node, periodically submitting new blocks and
// receiving ledger updates
func (n *Node) Run() {
	var ledger *Ledger

	for {
		select {
		case block := <-n.Unverified:
			if verify(ledger, block) {
				fmt.Println("Verified block", block.Hash)
				n.Verified <- block
			}
		case updatedLedger := <-n.Updates:
			ledger = &updatedLedger
		}
	}
}
