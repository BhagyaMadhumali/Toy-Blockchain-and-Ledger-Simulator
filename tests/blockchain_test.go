package tests

import (
	"testing"
	"toy-blockchain/blockchain"
	"toy-blockchain/ledger"
)

func TestGenesisBlockCreated(t *testing.T) {
	bc := blockchain.NewBlockchain()

	if len(bc.Blocks) != 1 {
		t.Errorf(
			"expected blockchain to start with one genesis block",
		)
	}

	if bc.Blocks[0].Index != 0 {
		t.Errorf(
			"expected genesis block index to be 0",
		)
	}
}

func TestAddTransaction(t *testing.T) {
	bc := blockchain.NewBlockchain()

	tx := ledger.Transaction{
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   20,
	}

	if err := bc.AddTransaction(tx); err != nil {
		t.Errorf(
			"expected valid transaction, got error: %v",
			err,
		)
	}

	if len(bc.PendingTransactions) != 1 {
		t.Errorf(
			"expected one pending transaction",
		)
	}
}

func TestBlockchainValidationAfterMining(t *testing.T) {
	bc := blockchain.NewBlockchain()

	tx := ledger.Transaction{
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   20,
	}

	if err := bc.AddTransaction(tx); err != nil {
		t.Fatalf(
			"failed to add transaction: %v",
			err,
		)
	}

	if err := bc.MinePendingTransactions(); err != nil {
		t.Fatalf(
			"failed to mine transactions: %v",
			err,
		)
	}

	if !bc.ValidateBlockchain() {
		t.Errorf(
			"expected blockchain to be valid",
		)
	}
}