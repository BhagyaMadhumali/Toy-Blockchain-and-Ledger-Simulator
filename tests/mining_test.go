package tests

import (
	"testing"
	"toy-blockchain/blockchain"
)

func TestMiningCreatesValidProof(t *testing.T) {
	block := blockchain.Block{Index: 1, Timestamp: 1000, PreviousHash: "abc"}
	result, err := blockchain.MineBlock(&block, 2)
	if err != nil {
		t.Fatal(err)
	}
	if result.Attempts == 0 || !blockchain.HasValidProof(block.Hash, 2) {
		t.Fatal("invalid mining result")
	}
	if result.Attempts != block.Nonce+1 {
		t.Fatal("attempt count should equal nonce plus one")
	}
}
