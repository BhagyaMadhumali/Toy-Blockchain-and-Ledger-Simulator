package blockchain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"toy-blockchain/ledger"
)

// SaveBlockchain saves the blockchain as formatted JSON.
func (bc *Blockchain) SaveBlockchain(filename string) error {
	directory := filepath.Dir(filename)

	if err := os.MkdirAll(
		directory,
		os.ModePerm,
	); err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		bc,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		filename,
		data,
		0644,
	)
}

// LoadBlockchain loads blockchain data from a JSON file.
//
// If the file does not exist, a new blockchain is created.
func LoadBlockchain(filename string) (*Blockchain, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return NewBlockchain(), nil
	}

	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var bc Blockchain

	if err := json.Unmarshal(
		data,
		&bc,
	); err != nil {
		return nil, err
	}

	// Rebuild the ledger if it was missing from the file.
	if bc.Ledger == nil || bc.Ledger.Balances == nil {
		bc.RebuildLedger()
	}

	// Prevent a nil pending transaction slice.
	if bc.PendingTransactions == nil {
		bc.PendingTransactions = []ledger.Transaction{}
	}

	// Use the default difficulty if no difficulty was saved.
	if bc.Difficulty == 0 {
		bc.Difficulty = DefaultDifficulty
	}

	return &bc, nil
}