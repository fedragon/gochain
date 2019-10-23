package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_hashOf(t *testing.T) {
	now := time.Date(2019, 10, 20, 8, 55, 0, 0, time.UTC)
	type args struct {
		index     int
		previous  Hash
		timestamp time.Time
		data      Data
	}
	tests := []struct {
		name    string
		args    args
		want    Hash
		wantErr bool
	}{
		{"produces a hash by taking into account the block index, its data, and the current time",
			args{
				previous:  "",
				timestamp: now,
				data:      "Hello",
			},
			"e6a025fb109f578a3f3517c036c68be724dbbef551377a498ea846d84ee080c1",
			false,
		},
		{"and the previous block's hash, if present",
			args{
				previous:  "000005fb109f578a3f3517c036c68be724dbbef551377a498ea846d84ee080c1",
				timestamp: now,
				data:      "Hello",
			},
			"f8bcbca4d6ce5f12a2f47b193f8a6783838732f8d4095b1b45e7d61bc95aef5d",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hashOf(tt.args.index, tt.args.previous, tt.args.timestamp, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hashOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		ledger *Ledger
		data   Data
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
		{"creates an unverified block for provided ledger",
			args{
				&Ledger{},
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
			got, err := Create(tt.args.ledger, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Index != tt.want.index || got.Data != tt.want.data {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
