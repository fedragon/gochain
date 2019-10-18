package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Solution represents the PoW solution computed by a miner
type Solution struct {
	Hash   Hash
	NextTx Transaction
	Fees   chan float64
}

// Miner represents a node mining the chain
type Miner struct {
	Address       string
	NextTx        uint64
	CollectedFees float64
	Ledger        *Ledger
	Solutions     chan<- Solution
	Fees          chan float64
}

// Mine simulates the resolution of a PoW-like puzzle; for simplicity, I'm returning
// as result the transaction that this miner would like to add to the blockchain, in
// case it's the first one to solve the puzzle
func (m *Miner) Mine() (*Solution, error) {
	if m.Ledger.IsEmpty() {
		return nil, errors.New("There is no solution to an empty ledger")
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	delay := time.Duration(10+rnd.Float64()*100) * time.Millisecond

	time.Sleep(delay)

	txID := fmt.Sprintf("%v-%v", m.Address, m.NextTx)
	m.NextTx++

	ledgerHash, _ := m.Ledger.HashOf()

	return &Solution{ledgerHash, Transaction(txID), m.Fees}, nil
}

// CollectFees collects fees for this miner and adds them to the total
func (m *Miner) CollectFees() {
	for fee := range m.Fees {
		fmt.Printf("[%v] Collected fee %v\n", m.Address, fee)
		m.CollectedFees += fee
		fmt.Printf("[%v] Total collected fees %v\n", m.Address, m.CollectedFees)
	}
}

// Start starts a loop that periodically triggers a miner to send a PoW solution
func (m *Miner) Start() {
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for range ticker.C {
			sol, _ := m.Mine()

			m.Solutions <- *sol
		}
	}()

	go m.CollectFees()
}
