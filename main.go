package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	ledger := NewLedger("We ❤️ blockchains")

	updates := make(chan Ledger)
	unverified := make(chan Block)
	verified := make(chan Block)
	node := &Node{
		Updates:    updates,
		Unverified: unverified,
		Verified:   verified,
	}

	go node.Run()

	updates <- *ledger

	t := time.NewTicker(time.Millisecond * 500)

	for {
		select {
		case <-t.C:
			block, err := Create(ledger, "foobar")
			if err != nil {
				log.Println("Unable to create block")
			}

			unverified <- *block
		case block := <-verified:
			ledger.Append(&block)
			fmt.Println("Appended block", block.Hash, "to the ledger")
		}
	}
}
