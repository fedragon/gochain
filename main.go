package main

import "fmt"

func receive(ledger *Ledger, updates chan<- Ledger, submissions <-chan Block) {
	hash, _ := ledger.HashOf()

	for b := range submissions {
		block := b
		fmt.Println("Received block", block)

		if hash == block.Previous.Hash {
			ledger.Append(&block)
			hash, _ = ledger.HashOf()
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
