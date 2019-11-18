package gochain

import (
	"fmt"
	"testing"
)

func Test_append(t *testing.T) {
	v := make(chan Block, 1)
	type args struct {
		chain    *Chain
		verified chan Block
	}
	tests := []struct {
		name string
		args args
	}{
		{"appends a block to the chain",
			args{
				NewChain("Hello world!"),
				v,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := Create(tt.args.chain, "Heyo!")
			if err != nil {
				t.Errorf("Append() error = %v", err)
				return
			}

			go append(tt.args.chain, tt.args.verified)

			tt.args.verified <- *block
			fmt.Println("sent block")

			if last := tt.args.chain.Last(); !DeepEqualNoHash(last, block) {
				t.Errorf("Append() = %v, want %v", last, block)
			}
		})
	}
}
