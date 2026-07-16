package tests

import (
	"testing"
	"toy-blockchain/blockchain"
)

func TestConcurrentMiningCreatesValidProof(t *testing.T) {
	block := blockchain.Block{Index: 1, Timestamp: 1000, PreviousHash: "abc"}
	result, err := blockchain.MineBlockWithWorkers(&block, 2, 4)
	if err != nil {
		t.Fatal(err)
	}
	if result.Attempts == 0 {
		t.Fatal("mining should perform at least one hash attempt")
	}
	if result.Workers != 4 {
		t.Fatalf("expected 4 workers, got %d", result.Workers)
	}
	if !blockchain.HasValidProof(block.Hash, 2) {
		t.Fatal("concurrent mining produced an invalid proof")
	}
	if block.Hash != blockchain.CalculateHash(block) {
		t.Fatal("stored hash does not match the winning nonce")
	}
}

func TestSingleWorkerMiningIsSequential(t *testing.T) {
	block := blockchain.Block{Index: 1, Timestamp: 1000, PreviousHash: "abc"}
	result, err := blockchain.MineBlockWithWorkers(&block, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	if result.Attempts != block.Nonce+1 {
		t.Fatalf("single-worker attempts should equal nonce plus one: attempts=%d nonce=%d", result.Attempts, block.Nonce)
	}
}

func TestMiningRejectsInvalidWorkerCount(t *testing.T) {
	block := blockchain.Block{}
	if _, err := blockchain.MineBlockWithWorkers(&block, 2, 0); err == nil {
		t.Fatal("expected invalid worker count to be rejected")
	}
}
