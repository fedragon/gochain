package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Solution struct {
	Delay  time.Duration
	Hash   Hash
	NextTx Transaction
}

type Miner struct {
	Name      string
	NextTx    uint64
	Ledger    *Ledger
	Solutions chan<- Solution
}

// Mine simulates the resolution of a PoW-like puzzle; for simplicity, I'm returning
// as result the transaction that this miner would like to add to the blockchain, in
// case it's the first one to solve the puzzle
func (m *Miner) Mine() Solution {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	delay := time.Duration(10+rnd.Float64()*100) * time.Millisecond

	time.Sleep(delay)

	txID := fmt.Sprintf("%v-%v", m.Name, m.NextTx)
	m.NextTx++

	ledgerHash, _ := m.Ledger.HashOf()

	return Solution{delay, ledgerHash, Transaction(txID)}
}

func (m *Miner) Start() {
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for range ticker.C {
			m.Solutions <- m.Mine()
		}
	}()
}
