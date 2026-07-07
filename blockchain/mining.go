package blockchain

import (
	"time"
	"toy-blockchain/ledger"
)

func (bc *Blockchain) MinePendingTransactions() {

	block := Block{
		Index:        len(bc.Blocks),
		Timestamp:    time.Now().Unix(),
		Transactions: bc.PendingTransactions,
		PreviousHash: bc.Blocks[len(bc.Blocks)-1].Hash,
		Nonce:        0,
	}

	for {
		hash := CalculateHash(block)
		if bc.isValidHash(hash) {
			block.Hash = hash
			break
		}
		block.Nonce++
	}

	// update ledger
	for _, tx := range bc.PendingTransactions {
		bc.Ledger.UpdateBalance(tx.Sender, tx.Receiver, tx.Amount)
	}

	bc.Blocks = append(bc.Blocks, block)
	bc.PendingTransactions = []ledger.Transaction{}
}