package main

import (
	"errors"
	"testing"
)

func Test_update(t *testing.T) {
	sols := make(chan Solution, 2)
	s1 := Solution{
		Hash("312d31450cbdb857673e5417cedc3431cc4774483aa52b6e69eb5d8d4c67c99da5b124"),
		Transaction("1-1"),
		make(chan float64, 2),
	}
	s2 := Solution{
		Hash("def"),
		Transaction("1-2"),
		make(chan float64, 2),
	}

	type args struct {
		ledger    *Ledger
		solutions chan Solution
	}
	tests := []struct {
		name string
		args args
		want Hash
	}{
		{"appends the first received solution to the chain",
			args{
				NewLedger("xyz"),
				sols,
			},
			Hash("3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.solutions <- s1
			tt.args.solutions <- s2
			close(tt.args.solutions)

			update(tt.args.ledger, tt.args.solutions)

			if got, _ := tt.args.ledger.HashOf(); got != tt.want {
				t.Errorf("update() = %v, want %v", got, tt.want)
			}

			if got, err := tt.args.ledger.Get(s2.Hash); err == nil {
				t.Errorf("update() = %v, want %v", got, errors.New("This block should not be in the chain"))
			}
		})
	}
}
