package gochain

import "fmt"

func append(chain *Chain, verified chan Block) {
	for b := range verified {
		block := b
		chain.Append(&block)
		fmt.Println("Appended block", block.Hash, "to the chain")
	}
}
