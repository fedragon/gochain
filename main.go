package main

import (
	"fmt"
	"strconv"
)

func main() {
	var ledger *Ledger = NewLedger("0-0")
	solutions := make(chan Solution, 10)

	for i := 0; i < 10; i++ {
		fee := make(chan float64)

		miner := Miner{
			strconv.FormatInt(int64(i+1), 10),
			1,
			0.0,
			ledger,
			solutions,
			fee,
		}
		miner.Start()
	}

	maxConflicts := 10

	for sol := range solutions {
		fmt.Printf("[%v, %vms]\n", sol.NextTx, sol.Delay.Milliseconds())
		if lh, _ := ledger.HashOf(); lh == sol.Hash {
			ledger.Add(sol.NextTx)
			sol.Fees <- 0.001
			fmt.Printf("[%v] HAS been added to the ledger\n", sol.NextTx)
		} else {
			fmt.Printf("[%v] has NOT been added to the ledger: ledger hash does not match\n", sol.NextTx)

			maxConflicts--

			if maxConflicts == 0 {
				fmt.Println(ledger.Prettify())
				return
			}
		}
	}
}
