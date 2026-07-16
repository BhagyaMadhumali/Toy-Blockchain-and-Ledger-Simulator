package tests

import (
	"strings"
	"testing"
	"toy-blockchain/blockchain"
)

func minedChain(t *testing.T) *blockchain.Blockchain {
	t.Helper()
	bc := blockchain.NewBlockchain()
	if err := bc.AddTransaction(signedTransaction(t, "Alice", "Bob", 20)); err != nil {
		t.Fatal(err)
	}
	if _, err := bc.MinePendingTransactions(); err != nil {
		t.Fatal(err)
	}
	return bc
}

func TestHashIntegrityTampering(t *testing.T) {
	bc := minedChain(t)
	bc.Blocks[1].Transactions[0].Amount = 999
	if err := bc.ValidateBlockchain(); err == nil || !strings.Contains(err.Error(), "Merkle root") {
		t.Fatalf("expected Merkle root error, got %v", err)
	}
}

func TestPreviousHashBranch(t *testing.T) {
	bc := minedChain(t)
	bc.Blocks[1].PreviousHash = "fake"
	if _, err := blockchain.MineBlock(&bc.Blocks[1], blockchain.DefaultDifficulty); err != nil {
		t.Fatal(err)
	}
	if err := bc.ValidateBlockchain(); err == nil || !strings.Contains(err.Error(), "previous-hash") {
		t.Fatalf("expected linkage error, got %v", err)
	}
}

func TestProofOfWorkBranch(t *testing.T) {
	bc := minedChain(t)
	block := &bc.Blocks[1]
	for nonce := uint64(0); ; nonce++ {
		block.Nonce = nonce
		block.Hash = blockchain.CalculateHash(*block)
		if !blockchain.HasValidProof(block.Hash, block.Difficulty) {
			break
		}
	}
	if err := bc.ValidateBlockchain(); err == nil || !strings.Contains(err.Error(), "proof-of-work") {
		t.Fatalf("expected proof error, got %v", err)
	}
}

func TestIndexAndTimestampValidation(t *testing.T) {
	bc := minedChain(t)
	bc.Blocks[1].Index = 9
	if _, err := blockchain.MineBlock(&bc.Blocks[1], blockchain.DefaultDifficulty); err != nil {
		t.Fatal(err)
	}
	if err := bc.ValidateBlockchain(); err == nil || !strings.Contains(err.Error(), "incorrect index") {
		t.Fatalf("expected index error, got %v", err)
	}

	bc = minedChain(t)
	bc.Blocks[1].Timestamp = bc.Blocks[0].Timestamp - 1
	if _, err := blockchain.MineBlock(&bc.Blocks[1], blockchain.DefaultDifficulty); err != nil {
		t.Fatal(err)
	}
	if err := bc.ValidateBlockchain(); err == nil || !strings.Contains(err.Error(), "timestamp") {
		t.Fatalf("expected timestamp error, got %v", err)
	}
}

func TestTrustedDifficultyValidation(t *testing.T) {
	bc := minedChain(t)
	if _, err := blockchain.MineBlock(&bc.Blocks[1], 1); err != nil {
		t.Fatal(err)
	}
	if err := bc.ValidateBlockchain(); err == nil || !strings.Contains(err.Error(), "retargeted difficulty") {
		t.Fatalf("expected difficulty error, got %v", err)
	}
}

func TestNegativeBalanceReplayDetection(t *testing.T) {
	bc := minedChain(t)
	bc.Blocks[1].Transactions[0] = signedTransaction(t, "Alice", "Bob", 5000)
	if _, err := blockchain.MineBlock(&bc.Blocks[1], blockchain.DefaultDifficulty); err != nil {
		t.Fatal(err)
	}
	if err := bc.ValidateBlockchain(); err == nil || !strings.Contains(err.Error(), "insufficient") {
		t.Fatalf("expected replay balance error, got %v", err)
	}
}
