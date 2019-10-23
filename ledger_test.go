package main

import (
	"testing"
)

func DeepEqual(l, m *Ledger) bool {
	return DeepEqualNoHash(l.Genesis, m.Genesis)
}

func DeepEqualNoHash(b, c *Block) bool {
	if b == nil && c == nil {
		return true
	}

	if (b == nil) != (c == nil) {
		return false
	}

	return b.Index == c.Index &&
		b.Data == c.Data &&
		DeepEqualNoHash(b.Next, c.Next)
}

func TestNewLedger(t *testing.T) {
	tests := []struct {
		name string
		data Data
		want *Ledger
	}{
		{"generates a new ledger, containing a genesis block",
			"Hello world!",
			&Ledger{&Block{
				Index:    1,
				Hash:     "c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a",
				Data:     "Hello world!",
				Previous: nil,
				Next:     nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLedger(tt.data); !DeepEqual(got, tt.want) {
				t.Errorf("NewLedger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedger_Append(t *testing.T) {
	genesis := &Block{
		Hash:     Hash("c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a"),
		Data:     Data("Hello"),
		Previous: nil,
		Next:     nil,
	}

	type args struct {
		b *Block
	}
	tests := []struct {
		name       string
		ledger     Ledger
		args       args
		wantLedger *Ledger
	}{
		{"appends a block to a non-empty ledger",
			Ledger{Genesis: genesis},
			args{
				&Block{
					Hash:     Hash("foo"),
					Data:     Data("world"),
					Previous: genesis,
					Next:     nil},
			},
			&Ledger{
				Genesis: &Block{
					Hash:     genesis.Hash,
					Data:     genesis.Data,
					Previous: genesis.Previous,
					Next: &Block{
						Hash:     Hash("776f726c64675de8ebf07b0ca1ed92f3cdce825df28d36d8fdc39904060d2c18b13c096edc"),
						Data:     Data("world"),
						Previous: genesis,
						Next:     nil}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Ledger{
				Genesis: genesis,
			}

			l.Append(tt.args.b)
			if !DeepEqual(l, tt.wantLedger) {
				t.Errorf("Ledger.Append() =\ngot  %v,\nwant %v\n", l, tt.wantLedger)
			}
		})
	}
}

func TestLedger_HashOf(t *testing.T) {
	genesis := &Block{
		Hash:     Hash("c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a"),
		Data:     Data("Hello"),
		Previous: nil,
		Next:     nil,
	}

	type fields struct {
		Genesis *Block
	}
	tests := []struct {
		name    string
		fields  fields
		want    Hash
		wantErr bool
	}{
		{"returns an error with an empty ledger",
			fields{nil},
			"",
			true,
		},
		{"returns the hash of the genesis block, when there is exactly one block in the ledger",
			fields{genesis},
			genesis.Hash,
			false,
		},
		{"returns the hash of last block, in all other cases",
			fields{&Block{
				Hash:     genesis.Hash,
				Data:     genesis.Data,
				Previous: nil,
				Next: &Block{
					Hash:     Hash("776f726c64675de8ebf07b0ca1ed92f3cdce825df28d36d8fdc39904060d2c18b13c096edc"),
					Data:     Data("world"),
					Previous: genesis,
					Next:     nil,
				},
			}},
			Hash("776f726c64675de8ebf07b0ca1ed92f3cdce825df28d36d8fdc39904060d2c18b13c096edc"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Ledger{
				Genesis: tt.fields.Genesis,
			}
			got, err := l.HashOf()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ledger.HashOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ledger.HashOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedger_IsEmpty(t *testing.T) {
	genesis := &Block{
		Hash:     Hash("c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a"),
		Data:     Data("Hello"),
		Previous: nil,
		Next:     nil,
	}

	type fields struct {
		Genesis *Block
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"returns true if the ledger has no blocks",
			fields{nil},
			true,
		},
		{"returns false otherwise",
			fields{genesis},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Ledger{
				Genesis: tt.fields.Genesis,
			}
			if got := l.IsEmpty(); got != tt.want {
				t.Errorf("Ledger.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeepEqualNoHash(t *testing.T) {
	type args struct {
		b *Block
		c *Block
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"returns true when two genesis blocks hold the same information, except for their hashes",
			args{
				&Block{
					Index: 1,
					Hash:  "foo",
					Data:  "Hello",
				},
				&Block{
					Index: 1,
					Hash:  "bar",
					Data:  "Hello",
				},
			},
			true,
		},
		{"returns true when two blocks hold the same information, except for their hashes",
			args{
				&Block{
					Index: 2,
					Hash:  "foo",
					Data:  "world",
					Previous: &Block{
						Index: 1,
						Hash:  "xyz",
						Data:  "Hello",
					},
					Next: &Block{
						Index: 3,
						Hash:  "abc",
						Data:  "!",
					},
				},
				&Block{
					Index: 2,
					Hash:  "erm",
					Data:  "world",
					Previous: &Block{
						Index: 1,
						Hash:  "xyz",
						Data:  "Hello",
					},
					Next: &Block{
						Index: 3,
						Hash:  "abc",
						Data:  "!",
					},
				},
			},
			true,
		},
		{"returns false when the first block is nil",
			args{
				nil,
				&Block{
					Index: 1,
					Hash:  "bar",
					Data:  "Hello",
					Next: &Block{
						Index: 2,
						Hash:  "def",
						Data:  "world",
					},
				},
			},
			false,
		},
		{"returns false when the second block is nil",
			args{
				&Block{
					Index: 1,
					Hash:  "bar",
					Data:  "Hello",
					Next: &Block{
						Index: 2,
						Hash:  "def",
						Data:  "world",
					},
				},
				nil,
			},
			false,
		},
		{"returns false when the blocks hold different information, besides their hashes",
			args{
				&Block{
					Index: 1,
					Hash:  "bar",
					Data:  "Hello",
					Next: &Block{
						Index: 2,
						Hash:  "def",
						Data:  "world",
					},
				},
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeepEqualNoHash(tt.args.b, tt.args.c); got != tt.want {
				t.Errorf("DeepEqualNoHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
