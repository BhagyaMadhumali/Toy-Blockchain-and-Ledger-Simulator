package tests

import (
	"testing"
	"toy-blockchain/ledger"
)

func signedTransaction(t *testing.T, sender, receiver string, amount int) ledger.Transaction {
	t.Helper()
	tx := ledger.Transaction{Sender: sender, Receiver: receiver, Amount: amount}
	privateKey, err := ledger.DemoAccountPrivateKey(sender)
	if err != nil {
		t.Fatal(err)
	}
	if err := ledger.SignTransaction(&tx, privateKey); err != nil {
		t.Fatal(err)
	}
	return tx
}
