package blockchain

func (bc *Blockchain) ValidateBlockchain() bool {
	if len(bc.Blocks) == 0 {
		return false
	}

	for i, currentBlock := range bc.Blocks {
		if currentBlock.Hash != CalculateHash(currentBlock) {
			return false
		}

		if i == 0 {
			if currentBlock.PreviousHash != "0000" {
				return false
			}
			continue
		}

		previousBlock := bc.Blocks[i-1]

		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}

		if !bc.isValidHash(currentBlock.Hash) {
			return false
		}
	}

	return true
}