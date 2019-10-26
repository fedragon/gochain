package main

import (
	"testing"
)

func TestCreate(t *testing.T) {
	type args struct {
		chain *Chain
		data  Data
	}
	type result struct {
		index int
		data  Data
	}
	tests := []struct {
		name    string
		args    args
		want    result
		wantErr bool
	}{
		{"creates an unverified block for provided chain",
			args{
				&Chain{},
				"Hello",
			},
			result{
				index: 1,
				data:  "Hello",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.chain, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Index != tt.want.index || got.Data != tt.want.data {
				t.Errorf("Create() = {%v %v}, want %v", got.Index, got.Data, tt.want)
			}
		})
	}
}
