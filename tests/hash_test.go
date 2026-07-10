package tests

import (
	"testing"
	"toy-blockchain/blockchain"
)

func TestHashConsistency(t *testing.T) {
	block := blockchain.Block{Index: 1, Timestamp: 1000, PreviousHash: "abc", Difficulty: 3}
	if blockchain.CalculateHash(block) != blockchain.CalculateHash(block) {
		t.Fatal("same block must produce the same hash")
	}
}

func TestHashSerializationIsUnambiguous(t *testing.T) {
	first := blockchain.Block{Index: 1, Timestamp: 23, PreviousHash: "abc", Difficulty: 3}
	second := blockchain.Block{Index: 12, Timestamp: 3, PreviousHash: "abc", Difficulty: 3}
	if blockchain.CalculateHash(first) == blockchain.CalculateHash(second) {
		t.Fatal("different blocks produced the same hash")
	}
}
