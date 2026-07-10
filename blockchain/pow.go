package blockchain

import "strings"

// isValidHash checks whether a hash satisfies
// the blockchain difficulty.
func (bc *Blockchain) isValidHash(hash string) bool {
	requiredPrefix := strings.Repeat(
		"0",
		bc.Difficulty,
	)

	return strings.HasPrefix(
		hash,
		requiredPrefix,
	)
}

// MineBlock repeatedly changes the nonce until
// the block hash begins with the required zeroes.
//
// It returns the total number of hash attempts.
func MineBlock(block *Block, difficulty int) int {
	requiredPrefix := strings.Repeat(
		"0",
		difficulty,
	)

	attempts := 0

	for {
		block.Hash = CalculateHash(*block)
		attempts++

		if strings.HasPrefix(
			block.Hash,
			requiredPrefix,
		) {
			return attempts
		}

		block.Nonce++
	}
}