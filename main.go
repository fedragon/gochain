package main

import (
	"fmt"
	"strconv"
)

func main() {
	var ledger *Ledger = NewLedger("0-0")

	solutions := make(chan Solution, 10)

	for i := 0; i < 10; i++ {
		miner := Miner{
			strconv.FormatInt(int64(i+1), 10),
			1,
			ledger,
			solutions,
		}
		miner.Start()
	}

	maxConflicts := 10
	for sol := range solutions {
		fmt.Printf("[%v, %vms]\n", sol.NextTx, sol.Delay.Milliseconds())
		if lh, _ := ledger.HashOf(); lh == sol.Hash {
			ledger.Add(sol.NextTx)
			fmt.Printf("[%v] HAS been added to the ledger\n", sol.NextTx)
		} else {
			fmt.Printf("[%v] has NOT been added to the ledger\n", sol.NextTx)
			fmt.Printf("Ledger hash changed, someone else solved the puzzle before me!\n")
			fmt.Printf("Current ledger hash =  %x\n", lh)
			fmt.Printf("Returned ledger hash = %x\n", sol.Hash)

			maxConflicts--

			if maxConflicts == 0 {
				fmt.Printf("\nFinal version of the ledger\n%v\n", ledger)
				return
			}
		}
	}

}
