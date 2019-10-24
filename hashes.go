package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

// CalculateHash calculates the hash using the hash of the last block in the ledger,
// the timestamp, and data provided
func CalculateHash(last *Block, timestamp time.Time, data Data) (Hash, error) {
	var index int = 1
	var hash Hash

	if last != nil {
		index = last.Index
		hash = last.Hash
	}
	return hashOf(index, hash, timestamp, data)
}

func hashOf(index int, base Hash, timestamp time.Time, data Data) (Hash, error) {
	hasher := sha256.New()

	hasher.Write([]byte(strconv.Itoa(index)))
	hasher.Write([]byte(base))
	hasher.Write([]byte(data))
	ts, err := timestamp.In(time.UTC).MarshalBinary()

	if err != nil {
		return "", err
	}

	hasher.Write(ts)
	hash := Hash(fmt.Sprintf("%x", hasher.Sum(nil)))

	return hash, nil
}
