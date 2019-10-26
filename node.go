package main

import (
	"fmt"
)

// Node represents a node submitting new blocks to the chain
type Node struct {
	Updates    <-chan Chain
	Unverified <-chan Block
	Verified   chan<- Block
}

// verifies a block before it's accepted in the chain
func verify(chain *Chain, unverified Block) bool {
	last := chain.Last()

	hash, err := hashOf(last.Index, last.Hash, unverified.Timestamp, unverified.Data)
	if err != nil {
		return false
	}

	return hash == unverified.Hash
}

// Run executes the main loop of a node, periodically submitting new blocks and
// receiving chain updates
func (n *Node) Run() {
	var chain *Chain

	for {
		select {
		case block := <-n.Unverified:
			if verify(chain, block) {
				fmt.Println("Verified block", block.Hash)
				n.Verified <- block
			}
		case updatedChain := <-n.Updates:
			chain = &updatedChain
		}
	}
}
