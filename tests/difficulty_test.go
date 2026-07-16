package tests

import (
	"testing"
	"toy-blockchain/blockchain"
)

func TestDifficultyRetargeting(t *testing.T) {
	genesis := blockchain.Block{Timestamp: 1000, Difficulty: blockchain.DefaultDifficulty}
	first := blockchain.Block{Timestamp: 2000, Difficulty: blockchain.DefaultDifficulty}

	if got := blockchain.CalculateNextDifficulty([]blockchain.Block{genesis}); got != blockchain.DefaultDifficulty {
		t.Fatalf("first block should keep default difficulty: got %d", got)
	}
	if got := blockchain.CalculateNextDifficulty([]blockchain.Block{genesis, first}); got != blockchain.DefaultDifficulty {
		t.Fatalf("second block should ignore the genesis interval: got %d", got)
	}

	fast := blockchain.Block{Timestamp: 2005, Difficulty: blockchain.DefaultDifficulty}
	if got := blockchain.CalculateNextDifficulty([]blockchain.Block{genesis, first, fast}); got != blockchain.DefaultDifficulty+1 {
		t.Fatalf("fast block should increase difficulty: got %d", got)
	}

	onTarget := blockchain.Block{Timestamp: 2030, Difficulty: blockchain.DefaultDifficulty}
	if got := blockchain.CalculateNextDifficulty([]blockchain.Block{genesis, first, onTarget}); got != blockchain.DefaultDifficulty {
		t.Fatalf("on-target block should keep difficulty: got %d", got)
	}

	slow := blockchain.Block{Timestamp: 2100, Difficulty: blockchain.DefaultDifficulty}
	if got := blockchain.CalculateNextDifficulty([]blockchain.Block{genesis, first, slow}); got != blockchain.DefaultDifficulty-1 {
		t.Fatalf("slow block should decrease difficulty: got %d", got)
	}
}

func TestDifficultyRetargetingBounds(t *testing.T) {
	fastChain := []blockchain.Block{
		{Timestamp: 0, Difficulty: blockchain.DefaultDifficulty},
		{Timestamp: 1000, Difficulty: blockchain.MaxDifficulty},
		{Timestamp: 1001, Difficulty: blockchain.MaxDifficulty},
	}
	if got := blockchain.CalculateNextDifficulty(fastChain); got != blockchain.MaxDifficulty {
		t.Fatalf("difficulty exceeded maximum: %d", got)
	}

	slowChain := []blockchain.Block{
		{Timestamp: 0, Difficulty: blockchain.DefaultDifficulty},
		{Timestamp: 1000, Difficulty: blockchain.MinDifficulty},
		{Timestamp: 2000, Difficulty: blockchain.MinDifficulty},
	}
	if got := blockchain.CalculateNextDifficulty(slowChain); got != blockchain.MinDifficulty {
		t.Fatalf("difficulty fell below minimum: %d", got)
	}
}
