package blockchain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"toy-blockchain/ledger"
)

func (bc *Blockchain) SaveBlockchain(filename string) error {
	dir := filepath.Dir(filename)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	data, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func LoadBlockchain(filename string) (*Blockchain, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return NewBlockchain(), nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var bc Blockchain

	if err := json.Unmarshal(data, &bc); err != nil {
		return nil, err
	}

	if bc.Ledger == nil || bc.Ledger.Balances == nil {
		bc.RebuildLedger()
	}

	if bc.PendingTransactions == nil {
		bc.PendingTransactions = []ledger.Transaction{}
	}

	if bc.Difficulty == 0 {
		bc.Difficulty = 4
	}

	return &bc, nil
}