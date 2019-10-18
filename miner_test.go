package main

import (
	"reflect"
	"testing"
)

func TestMiner_Mine(t *testing.T) {
	fees := make(chan float64)

	type fields struct {
		Address       string
		NextTx        uint64
		CollectedFees float64
		Ledger        *Ledger
		Solutions     chan<- Solution
		Fees          chan float64
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Solution
		wantErr bool
	}{
		{"returns an error when the ledger is empty",
			fields{
				Address:       "xyz",
				NextTx:        1,
				CollectedFees: 0,
				Ledger:        &Ledger{},
				Solutions:     make(chan<- Solution),
				Fees:          fees,
			},
			nil,
			true,
		},
		{"returns a solution otherwise",
			fields{
				Address:       "xyz",
				NextTx:        1,
				CollectedFees: 0,
				Ledger:        NewLedger("abc"),
				Solutions:     make(chan<- Solution),
				Fees:          fees,
			},
			&Solution{
				Hash:   Hash("ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"),
				NextTx: Transaction("xyz-1"),
				Fees:   fees,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Miner{
				Address:       tt.fields.Address,
				NextTx:        tt.fields.NextTx,
				CollectedFees: tt.fields.CollectedFees,
				Ledger:        tt.fields.Ledger,
				Solutions:     tt.fields.Solutions,
				Fees:          tt.fields.Fees,
			}

			got, err := m.Mine()
			if (err != nil) != tt.wantErr {
				t.Errorf("Miner.Mine() = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miner.Mine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiner_CollectFees(t *testing.T) {
	fees := make(chan float64, 1)

	type fields struct {
		Address       string
		NextTx        uint64
		CollectedFees float64
		Ledger        *Ledger
		Solutions     chan<- Solution
		Fees          chan float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{"sums collected fees",
			fields{
				Address:       "xyz",
				NextTx:        1,
				CollectedFees: 0,
				Ledger:        NewLedger("abc"),
				Solutions:     make(chan<- Solution),
				Fees:          fees,
			},
			0.05,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Miner{
				Address:       tt.fields.Address,
				NextTx:        tt.fields.NextTx,
				CollectedFees: tt.fields.CollectedFees,
				Ledger:        tt.fields.Ledger,
				Solutions:     tt.fields.Solutions,
				Fees:          tt.fields.Fees,
			}

			fees <- 0.05
			close(fees)
			m.CollectFees()

			if got := m.CollectedFees; got != 0.05 {
				t.Errorf("Miner.CollectFees() = %v, want %v", got, tt.want)
			}
		})
	}
}
