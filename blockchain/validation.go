package blockchain

// ValidateBlockchain checks every block in the chain.
func (bc *Blockchain) ValidateBlockchain() bool {
	if len(bc.Blocks) == 0 {
		return false
	}

	for index, currentBlock := range bc.Blocks {
		// Recalculate and verify the current block hash.
		calculatedHash := CalculateHash(currentBlock)

		if currentBlock.Hash != calculatedHash {
			return false
		}

		// Special checks for the genesis block.
		if index == 0 {
			if currentBlock.PreviousHash != "0000" {
				return false
			}

			continue
		}

		previousBlock := bc.Blocks[index-1]

		// Verify the connection with the previous block.
		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}

		// Verify Proof of Work.
		if !bc.isValidHash(currentBlock.Hash) {
			return false
		}
	}

	return true
}