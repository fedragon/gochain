package main

import (
	"testing"
)

func TestMiner_Mine(t *testing.T) {
	type fields struct {
		Address string
		Ledger  *Ledger
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Solution
		wantErr bool
	}{
		{"returns an error when the ledger is empty",
			fields{
				Address: "xyz",
				Ledger:  &Ledger{},
			},
			nil,
			true,
		},
		{"returns a solution otherwise",
			fields{
				Address: "xyz",
				Ledger:  NewLedger("abc"),
			},
			&Solution{
				Hash:   Hash("ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"),
				NextTx: Transaction("xyz-1"),
				Fees:   make(chan float64, 1),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMiner(
				tt.fields.Address,
				tt.fields.Ledger,
				nil,
			)

			got, err := m.Mine()
			if (err != nil) != tt.wantErr {
				t.Errorf("miner.Mine() = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (got.Hash != tt.want.Hash || got.NextTx != tt.want.NextTx) {
				t.Errorf("miner.Mine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiner_CollectFees(t *testing.T) {
	type fields struct {
		Address string
		Ledger  *Ledger
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{"sums collected fees",
			fields{
				Address: "xyz",
				Ledger:  NewLedger("abc"),
			},
			0.05,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMiner(
				tt.fields.Address,
				tt.fields.Ledger,
				nil,
			)

			m.Fees <- 0.05
			close(m.Fees)
			m.CollectFees()

			if got := m.CollectedFees; got != 0.05 {
				t.Errorf("miner.CollectFees() = %v, want %v", got, tt.want)
			}
		})
	}
}
