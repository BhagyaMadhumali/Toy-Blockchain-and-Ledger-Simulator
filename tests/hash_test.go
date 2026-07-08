package tests

import (
	"testing"
	"toy-blockchain/blockchain"
)

func TestHashConsistency(t *testing.T) {
	block := blockchain.Block{
		Index:        1,
		Timestamp:    1000,
		PreviousHash: "abc",
		Nonce:        0,
	}

	hash1 := blockchain.CalculateHash(block)
	hash2 := blockchain.CalculateHash(block)

	if hash1 != hash2 {
		t.Errorf("expected same block to generate same hash")
	}
}

func TestDifferentBlocksHaveDifferentHashes(t *testing.T) {
	block1 := blockchain.Block{
		Index:        1,
		Timestamp:    1000,
		PreviousHash: "abc",
		Nonce:        0,
	}

	block2 := blockchain.Block{
		Index:        2,
		Timestamp:    1000,
		PreviousHash: "abc",
		Nonce:        0,
	}

	if blockchain.CalculateHash(block1) == blockchain.CalculateHash(block2) {
		t.Errorf("expected different blocks to generate different hashes")
	}
}