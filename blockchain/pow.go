package blockchain

import (
	"fmt"
	"strings"
	"time"
)

type MiningResult struct {
	Attempts uint64
	Elapsed  time.Duration
}

func HasValidProof(hash string, difficulty int) bool {
	if difficulty < 1 || difficulty > 64 {
		return false
	}
	return strings.HasPrefix(hash, strings.Repeat("0", difficulty))
}

// MineBlock increments the nonce until the block satisfies proof-of-work.
func MineBlock(block *Block, difficulty int) (MiningResult, error) {
	if difficulty < 1 || difficulty > 6 {
		return MiningResult{}, fmt.Errorf("difficulty must be between 1 and 6")
	}

	block.Difficulty = difficulty
	block.Nonce = 0
	start := time.Now()
	var attempts uint64

	for {
		block.Hash = CalculateHash(*block)
		attempts++
		if HasValidProof(block.Hash, difficulty) {
			return MiningResult{Attempts: attempts, Elapsed: time.Since(start)}, nil
		}
		block.Nonce++
	}
}
