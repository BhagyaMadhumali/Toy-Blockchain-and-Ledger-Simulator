package blockchain

import "strings"

func (bc *Blockchain) isValidHash(hash string) bool {
	prefix := strings.Repeat("0", bc.Difficulty)
	return strings.HasPrefix(hash, prefix)
}

func MineBlock(block *Block, difficulty int) {
	prefix := strings.Repeat("0", difficulty)

	for {
		block.Hash = CalculateHash(*block)

		if strings.HasPrefix(block.Hash, prefix) {
			return
		}

		block.Nonce++
	}
}