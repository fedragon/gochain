package gochain

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
		{"produces a hash by taking into account the block index, its data, and the block time",
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

func TestCalculateHash(t *testing.T) {
	now := time.Date(2019, 10, 20, 8, 55, 0, 0, time.UTC)

	type args struct {
		last      *Block
		timestamp time.Time
		data      Data
	}
	tests := []struct {
		name    string
		args    args
		want    Hash
		wantErr bool
	}{
		{"calculates the hash when the previous block is empty",
			args{
				last:      nil,
				timestamp: now,
				data:      "Hello",
			},
			"3a0db78eef0ec8dcf9e953d2694703fffdcae9ac129387651c37e5a5307d7238",
			false,
		},
		{"calculates the hash when the previous block is present",
			args{
				last:      &Block{Index: 1, Hash: "e6a025fb109f578a3f3517c036c68be724dbbef551377a498ea846d84ee080c1"},
				timestamp: now,
				data:      "World",
			},
			"8f441924ca8853a44e0016227d9eb99e15baf52cb8a70a971236c327d6b09dca",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateHash(tt.args.last, tt.args.timestamp, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
