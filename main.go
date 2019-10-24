package main

import "fmt"

func verify(ledger *Ledger, unverified *Block) bool {
	last := ledger.Last()

	hash, err := hashOf(last.Index, last.Hash, unverified.Timestamp, unverified.Data)
	if err != nil {
		return false
	}

	return hash == unverified.Hash
}

func receive(ledger *Ledger, updates chan<- Ledger, submissions <-chan Block) {
	for b := range submissions {
		block := b
		fmt.Println("Received block", block.Hash)

		if verify(ledger, &block) {
			ledger.Append(&block)
			fmt.Println("New ledger hash", ledger.HashOf())
		}

		updates <- *ledger
	}
}

func main() {
	ledger := NewLedger("We ❤️ blockchains")

	updates := make(chan Ledger)
	submissions := make(chan Block)
	node := &Node{
		Updates:     updates,
		Submissions: submissions,
	}

	go node.Run()

	updates <- *ledger

	receive(ledger, updates, submissions)
}
