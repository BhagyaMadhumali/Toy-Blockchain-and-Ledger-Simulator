package tests

import (
	"testing"
	"toy-blockchain/ledger"
)

func TestValidTransaction(t *testing.T) {
	testLedger := ledger.NewLedger()
	testLedger.AddAccount("Alice", 100)

	tx := ledger.Transaction{
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   20,
	}

	if err := ledger.ValidateTransaction(
		testLedger,
		tx,
	); err != nil {
		t.Errorf(
			"expected valid transaction, got error: %v",
			err,
		)
	}
}

func TestNegativeAmountTransaction(t *testing.T) {
	testLedger := ledger.NewLedger()
	testLedger.AddAccount("Alice", 100)

	tx := ledger.Transaction{
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   -10,
	}

	if err := ledger.ValidateTransaction(
		testLedger,
		tx,
	); err == nil {
		t.Errorf(
			"expected negative amount to be invalid",
		)
	}
}

func TestZeroAmountTransaction(t *testing.T) {
	testLedger := ledger.NewLedger()
	testLedger.AddAccount("Alice", 100)

	tx := ledger.Transaction{
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   0,
	}

	if err := ledger.ValidateTransaction(
		testLedger,
		tx,
	); err == nil {
		t.Errorf(
			"expected zero amount to be invalid",
		)
	}
}

func TestInsufficientBalanceTransaction(t *testing.T) {
	testLedger := ledger.NewLedger()
	testLedger.AddAccount("Alice", 10)

	tx := ledger.Transaction{
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   50,
	}

	if err := ledger.ValidateTransaction(
		testLedger,
		tx,
	); err == nil {
		t.Errorf(
			"expected insufficient balance to be invalid",
		)
	}
}