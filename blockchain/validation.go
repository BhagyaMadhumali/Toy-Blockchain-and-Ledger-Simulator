package blockchain

import (
	"fmt"
	"toy-blockchain/ledger"
)

type ValidationError struct {
	BlockIndex       int
	TransactionIndex int
	Reason           string
}

func (e *ValidationError) Error() string {
	if e.TransactionIndex >= 0 {
		return fmt.Sprintf("block %d, transaction %d: %s", e.BlockIndex, e.TransactionIndex, e.Reason)
	}
	return fmt.Sprintf("block %d: %s", e.BlockIndex, e.Reason)
}

func validationError(blockIndex int, reason string) error {
	return &ValidationError{BlockIndex: blockIndex, TransactionIndex: -1, Reason: reason}
}

// ValidateBlockchain identifies the first offending block instead of returning
// only true/false.
func (bc *Blockchain) ValidateBlockchain() error {
	if len(bc.Blocks) == 0 {
		return validationError(0, "blockchain is empty")
	}

	expectedGenesis := NewGenesisBlock()
	genesis := bc.Blocks[0]
	if genesis.Hash != expectedGenesis.Hash || CalculateHash(genesis) != expectedGenesis.Hash {
		return validationError(0, "genesis block does not match the fixed trusted genesis block")
	}

	if _, err := ReplayBlocks(bc.Blocks); err != nil {
		return err
	}
	return nil
}

// ReplayBlocks validates structure, hashes, proof-of-work and transactions while
// deriving all balances from the chain.
func ReplayBlocks(blocks []Block) (*ledger.Ledger, error) {
	if len(blocks) == 0 {
		return nil, validationError(0, "blockchain is empty")
	}

	replayed := ledger.NewLedger()
	for i, block := range blocks {
		if block.Index != i {
			return nil, validationError(i, fmt.Sprintf("incorrect index %d; expected %d", block.Index, i))
		}
		if block.Difficulty != DefaultDifficulty {
			return nil, validationError(i, fmt.Sprintf("untrusted difficulty %d; expected %d", block.Difficulty, DefaultDifficulty))
		}
		if block.Hash != CalculateHash(block) {
			return nil, validationError(i, "stored hash does not match block contents")
		}
		if !HasValidProof(block.Hash, DefaultDifficulty) {
			return nil, validationError(i, "proof-of-work does not satisfy trusted difficulty")
		}

		if i == 0 {
			if block.Timestamp != GenesisTimestamp || block.PreviousHash != GenesisPreviousHash {
				return nil, validationError(0, "genesis constants are invalid")
			}
			if len(block.Transactions) != len(genesisTransactions) {
				return nil, validationError(0, "genesis allocations are invalid")
			}
			for txIndex, tx := range block.Transactions {
				expected := genesisTransactions[txIndex]
				if tx != expected {
					return nil, &ValidationError{BlockIndex: 0, TransactionIndex: txIndex, Reason: "genesis allocation does not match trusted allocation"}
				}
				if err := replayed.Credit(tx.Receiver, tx.Amount); err != nil {
					return nil, &ValidationError{BlockIndex: 0, TransactionIndex: txIndex, Reason: err.Error()}
				}
			}
			continue
		}

		previous := blocks[i-1]
		if block.PreviousHash != previous.Hash {
			return nil, validationError(i, "previous-hash link is invalid")
		}
		if block.Timestamp < previous.Timestamp {
			return nil, validationError(i, "timestamp is earlier than the previous block")
		}
		for txIndex, tx := range block.Transactions {
			if err := replayed.ApplyTransaction(tx); err != nil {
				return nil, &ValidationError{BlockIndex: i, TransactionIndex: txIndex, Reason: err.Error()}
			}
		}
	}
	return replayed, nil
}
