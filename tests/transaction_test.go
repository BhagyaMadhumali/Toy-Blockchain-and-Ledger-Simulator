package tests

import (
	"testing"
	"toy-blockchain/ledger"
)

func TestApplyValidTransaction(t *testing.T) {
	l := ledger.NewLedger()
	if err := l.Credit("Alice", 100); err != nil {
		t.Fatal(err)
	}
	if err := l.ApplyTransaction(ledger.Transaction{Sender: "Alice", Receiver: "Bob", Amount: 20}); err != nil {
		t.Fatal(err)
	}
	if l.GetBalance("Alice") != 80 || l.GetBalance("Bob") != 20 {
		t.Fatal("unexpected balances")
	}
}

func TestInvalidTransactions(t *testing.T) {
	cases := []ledger.Transaction{
		{Sender: "Alice", Receiver: "Bob", Amount: 0},
		{Sender: "Alice", Receiver: "Bob", Amount: -1},
		{Sender: "Alice", Receiver: "Alice", Amount: 1},
		{Sender: "", Receiver: "Bob", Amount: 1},
		{Sender: "Alice", Receiver: "", Amount: 1},
		{Sender: ledger.SystemAccount, Receiver: "Alice", Amount: 1},
	}
	for _, tx := range cases {
		l := ledger.NewLedger()
		_ = l.Credit("Alice", 100)
		if err := l.ApplyTransaction(tx); err == nil {
			t.Fatalf("expected rejection for %+v", tx)
		}
	}
}
