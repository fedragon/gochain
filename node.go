package main

import (
	"log"
	"time"
)

// Node represents a node submitting new blocks to the ledger
type Node struct {
	Address     string
	Updates     <-chan Ledger
	Submissions chan<- Block
}

// Run executes the main loop of a node, periodically submitting new blocks and
// receiving ledger updates
func (n *Node) Run() {
	t := time.NewTicker(time.Millisecond * 500)
	var ledger *Ledger

	for {
		select {
		case <-t.C:
			if ledger != nil {
				block, err := Create(ledger, "foobar")

				if err != nil {
					log.Println("Unable to create block")
				}

				n.Submissions <- *block
			}
		case updatedLedger := <-n.Updates:
			ledger = &updatedLedger
		}
	}
}
