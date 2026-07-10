package blockchain

import (
	"fmt"
	"toy-blockchain/ledger"
)

// MinePendingTransactions re-validates every pending transaction against a
// freshly replayed ledger before any block is mined or balance is changed.
func (bc *Blockchain) MinePendingTransactions() (MiningResult, error) {
	if len(bc.PendingTransactions) == 0 {
		return MiningResult{}, fmt.Errorf("no pending transactions to mine")
	}
	if err := bc.RebuildLedger(); err != nil {
		return MiningResult{}, fmt.Errorf("cannot mine invalid blockchain: %w", err)
	}

	temporary := bc.Ledger.Clone()
	transactions := append([]ledger.Transaction(nil), bc.PendingTransactions...)
	for i, tx := range transactions {
		if err := temporary.ApplyTransaction(tx); err != nil {
			return MiningResult{}, fmt.Errorf("pending transaction %d is invalid: %w", i, err)
		}
	}

	newBlock := newCandidateBlock(bc.Blocks[len(bc.Blocks)-1], transactions)
	result, err := MineBlock(&newBlock, DefaultDifficulty)
	if err != nil {
		return MiningResult{}, err
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	bc.PendingTransactions = []ledger.Transaction{}
	bc.Ledger = temporary
	return result, nil
}
