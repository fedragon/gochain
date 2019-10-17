package main

import (
	"reflect"
	"testing"
)

func TestNewLedger(t *testing.T) {
	tests := []struct {
		name string
		seed Transaction
		want *Ledger
	}{
		{"generates a new ledger, containing a genesis block",
			Transaction("Hello world!"),
			&Ledger{&Block{
				Hash:     "c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a",
				Tx:       "Hello world!",
				Previous: nil,
				Next:     nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLedger(tt.seed); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLedger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedger_Add(t *testing.T) {
	genesis := &Block{
		Hash:     Hash("c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a"),
		Tx:       Transaction("Hello"),
		Previous: nil,
		Next:     nil,
	}

	type args struct {
		tx Transaction
	}
	tests := []struct {
		name       string
		ledger     Ledger
		args       args
		wantHash   Hash
		wantLedger *Ledger
		wantErr    bool
	}{
		{"appends a block to a non-empty ledger",
			Ledger{Genesis: genesis},
			args{Transaction("world")},
			Hash("776f726c64675de8ebf07b0ca1ed92f3cdce825df28d36d8fdc39904060d2c18b13c096edc"),
			&Ledger{
				Genesis: &Block{
					Hash:     genesis.Hash,
					Tx:       genesis.Tx,
					Previous: genesis.Previous,
					Next: &Block{
						Hash:     Hash("776f726c64675de8ebf07b0ca1ed92f3cdce825df28d36d8fdc39904060d2c18b13c096edc"),
						Tx:       Transaction("world"),
						Previous: genesis,
						Next:     nil}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Ledger{
				Genesis: genesis,
			}
			gotHash, err := l.Add(tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ledger.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHash != tt.wantHash {
				t.Errorf("Ledger.Add() = %v, want %v", gotHash, tt.wantHash)
			}
			if !reflect.DeepEqual(l, tt.wantLedger) {
				t.Errorf("Ledger.Add() =\ngot  %v,\nwant %v\n", l, tt.wantLedger)
			}
		})
	}
}
