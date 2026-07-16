package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"toy-blockchain/ledger"
)

type persistedBlockchain struct {
	Blocks              []Block              `json:"blocks"`
	PendingTransactions []ledger.Transaction `json:"pending_transactions"`
}

func (bc *Blockchain) SaveBlockchain(filename string) error {
	if err := bc.ValidateBlockchain(); err != nil {
		return fmt.Errorf("refusing to save invalid blockchain: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	state := persistedBlockchain{Blocks: bc.Blocks, PendingTransactions: bc.PendingTransactions}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadBlockchain(filename string) (*Blockchain, error) {
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		return NewBlockchain(), nil
	}
	if err != nil {
		return nil, err
	}

	var state persistedBlockchain
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	bc := &Blockchain{
		Blocks:              state.Blocks,
		PendingTransactions: state.PendingTransactions,
		Ledger:              ledger.NewLedger(),
	}
	if bc.PendingTransactions == nil {
		bc.PendingTransactions = []ledger.Transaction{}
	}
	if err := bc.ValidateBlockchain(); err != nil {
		return nil, fmt.Errorf("saved blockchain is invalid: %w", err)
	}
	if err := bc.RebuildLedger(); err != nil {
		return nil, err
	}

	// Also validate pending transactions sequentially so edited persisted data
	// cannot be mined later.
	temporary := bc.Ledger.Clone()
	for i, tx := range bc.PendingTransactions {
		if err := temporary.ApplyTransaction(tx); err != nil {
			return nil, fmt.Errorf("saved pending transaction %d is invalid: %w", i, err)
		}
	}
	return bc, nil
}

// LoadCandidateBlocks reads a competing chain from a normal blockchain JSON
// file. Pending transactions in that file are intentionally ignored because
// fork choice compares confirmed blocks only.
func LoadCandidateBlocks(filename string) ([]Block, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var state persistedBlockchain
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("decode candidate blockchain: %w", err)
	}
	if len(state.Blocks) == 0 {
		return nil, fmt.Errorf("candidate blockchain contains no blocks")
	}
	return cloneBlocks(state.Blocks), nil
}

// ExportBlockchain writes a complete snapshot that another node can use as a
// competing-chain candidate.
func (bc *Blockchain) ExportBlockchain(filename string) error {
	return bc.SaveBlockchain(filename)
}
