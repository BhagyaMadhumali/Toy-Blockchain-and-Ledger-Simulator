package tests

import (
	"testing"
	"toy-blockchain/blockchain"
	"toy-blockchain/ledger"
)

func TestTamperingDetectionAmount(t *testing.T) {
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
			"failed to mine transaction: %v",
			err,
		)
	}

	// Deliberately modify a mined transaction.
	bc.Blocks[1].Transactions[0].Amount = 999

	if bc.ValidateBlockchain() {
		t.Errorf(
			"expected tampering to be detected after amount change",
		)
	}
}

func TestTamperingDetectionPreviousHash(t *testing.T) {
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
			"failed to mine transaction: %v",
			err,
		)
	}

	// Deliberately break the previous-hash connection.
	bc.Blocks[1].PreviousHash = "fakehash"

	if bc.ValidateBlockchain() {
		t.Errorf(
			"expected tampering to be detected after previous hash change",
		)
	}
}