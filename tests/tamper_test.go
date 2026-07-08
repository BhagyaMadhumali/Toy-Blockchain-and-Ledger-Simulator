package tests

import (
	"testing"
	"toy-blockchain/blockchain"
	"toy-blockchain/ledger"
)

func TestTamperingDetectionAmount(t *testing.T) {
	bc := blockchain.NewBlockchain()

	_ = bc.AddTransaction(ledger.Transaction{
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   20,
	})

	_ = bc.MinePendingTransactions()

	bc.Blocks[1].Transactions[0].Amount = 999

	if bc.ValidateBlockchain() {
		t.Errorf("expected tampering to be detected after amount change")
	}
}

func TestTamperingDetectionPreviousHash(t *testing.T) {
	bc := blockchain.NewBlockchain()

	_ = bc.AddTransaction(ledger.Transaction{
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   20,
	})

	_ = bc.MinePendingTransactions()

	bc.Blocks[1].PreviousHash = "fakehash"

	if bc.ValidateBlockchain() {
		t.Errorf("expected tampering to be detected after previous hash change")
	}
}