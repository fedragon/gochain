package gochain

import (
	"reflect"
	"testing"
)

func DeepEqual(l, m *Chain) bool {
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

func TestNewChain(t *testing.T) {
	tests := []struct {
		name string
		data Data
		want *Chain
	}{
		{"generates a new chain, containing a genesis block",
			"Hello world!",
			&Chain{&Block{
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
			if got := NewChain(tt.data); !DeepEqual(got, tt.want) {
				t.Errorf("NewChain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChain_Append(t *testing.T) {
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
		name      string
		chain     Chain
		args      args
		wantChain *Chain
	}{
		{"appends a block to a non-empty chain",
			Chain{Genesis: genesis},
			args{
				&Block{
					Hash:     Hash("foo"),
					Data:     Data("world"),
					Previous: genesis,
					Next:     nil},
			},
			&Chain{
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
			l := &Chain{
				Genesis: genesis,
			}

			l.Append(tt.args.b)
			if !DeepEqual(l, tt.wantChain) {
				t.Errorf("Chain.Append() =\ngot  %v,\nwant %v\n", l, tt.wantChain)
			}
		})
	}
}

func TestChain_HashOf(t *testing.T) {
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
		want   Hash
	}{
		{"returns an empty string with an empty chain",
			fields{nil},
			"",
		},
		{"returns the hash of the genesis block, when there is exactly one block in the chain",
			fields{genesis},
			genesis.Hash,
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Chain{
				Genesis: tt.fields.Genesis,
			}
			got := l.HashOf()
			if got != tt.want {
				t.Errorf("Chain.HashOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChain_IsEmpty(t *testing.T) {
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
		{"returns true if the chain has no blocks",
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
			l := Chain{
				Genesis: tt.fields.Genesis,
			}
			if got := l.IsEmpty(); got != tt.want {
				t.Errorf("Chain.IsEmpty() = %v, want %v", got, tt.want)
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

func TestChain_Last(t *testing.T) {
	genesis := &Block{
		Hash:     Hash("c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a"),
		Data:     Data("Hello"),
		Previous: nil,
		Next:     nil,
	}
	other := &Block{
		Hash:     Hash("776f726c64675de8ebf07b0ca1ed92f3cdce825df28d36d8fdc39904060d2c18b13c096edc"),
		Data:     Data("world"),
		Previous: genesis,
		Next:     nil,
	}

	type fields struct {
		Genesis *Block
	}
	tests := []struct {
		name   string
		fields fields
		want   *Block
	}{
		{"returns nil, if the chain is empty",
			fields{nil},
			nil,
		},
		{"returns the genesis block, if there is only one block in the chain",
			fields{genesis},
			genesis,
		},
		{"returns the last block in the chain",
			fields{&Block{
				Hash: genesis.Hash,
				Data: genesis.Data,
				Next: other,
			}},
			other,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Chain{
				Genesis: tt.fields.Genesis,
			}
			if got := l.Last(); !DeepEqualNoHash(got, tt.want) {
				t.Errorf("Chain.Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChain_Get(t *testing.T) {
	genesis := &Block{
		Hash:     Hash("c0535e4be2b79ffd93291305436bf889314e4a3faec05ecffcbb7df31ad9e51a"),
		Data:     Data("Hello"),
		Previous: nil,
		Next:     nil,
	}
	other := &Block{
		Hash:     Hash("776f726c64675de8ebf07b0ca1ed92f3cdce825df28d36d8fdc39904060d2c18b13c096edc"),
		Data:     Data("world"),
		Previous: genesis,
		Next:     nil,
	}

	type fields struct {
		Genesis *Block
	}
	type args struct {
		h Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Block
		wantErr bool
	}{
		{"returns an error if the chain has not been initialized yet",
			fields{nil},
			args{"irrelevant"},
			nil,
			true,
		},
		{"returns nil if a block with provided hash is not found",
			fields{genesis},
			args{"notthegenesishash"},
			nil,
			false,
		},
		{"returns the genesis block, if provided its hash",
			fields{genesis},
			args{genesis.Hash},
			genesis,
			false,
		},
		{"returns the block matching provided hash, if found",
			fields{&Block{
				Hash: genesis.Hash,
				Data: genesis.Data,
				Next: other,
			}},
			args{other.Hash},
			other,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Chain{
				Genesis: tt.fields.Genesis,
			}
			got, err := l.Get(tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chain.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chain.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
