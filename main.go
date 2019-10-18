package main

import (
	"fmt"
	"strconv"
	"time"
)

func update(ledger *Ledger, solutions chan Solution) {
	maxConflicts := 10

	for sol := range solutions {
		fmt.Printf("[%v] %v\n", time.Now(), sol.NextTx)
		if lh, _ := ledger.HashOf(); lh == sol.Hash {
			ledger.Add(sol.NextTx)
			sol.Fees <- 0.001
			fmt.Printf("[%v] %v HAS been added to the ledger\n", time.Now(), sol.NextTx)
		} else {
			fmt.Printf("[%v] %v has NOT been added to the ledger\n", time.Now(), sol.NextTx)

			maxConflicts--

			if maxConflicts == 0 {
				fmt.Println(ledger.Prettify())
				return
			}
		}
	}
}

func main() {
	var ledger *Ledger = NewLedger("0-0")
	solutions := make(chan Solution, 10)

	for i := 0; i < 10; i++ {
		miner := Miner{
			strconv.FormatInt(int64(i+1), 10),
			1,
			0.0,
			ledger,
			solutions,
			make(chan float64),
		}
		miner.Start()
	}

	update(ledger, solutions)
}
