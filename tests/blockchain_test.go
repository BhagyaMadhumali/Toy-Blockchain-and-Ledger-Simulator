package tests

import (
	"testing"
	"toy-blockchain/blockchain"
)

func TestGenesisIsDeterministic(t *testing.T) {
	if blockchain.NewGenesisBlock().Hash != blockchain.NewGenesisBlock().Hash {
		t.Fatal("genesis block must be deterministic")
	}
}

func TestBalancesAreReplayedFromChain(t *testing.T) {
	bc := blockchain.NewBlockchain()
	if bc.Ledger.GetBalance("Alice") != 100 || bc.Ledger.GetBalance("Bob") != 50 || bc.Ledger.GetBalance("Charlie") != 75 {
		t.Fatal("genesis allocations were not replayed")
	}
	if err := bc.AddTransaction(signedTransaction(t, "Alice", "Bob", 20)); err != nil {
		t.Fatal(err)
	}
	if _, err := bc.MinePendingTransactions(); err != nil {
		t.Fatal(err)
	}
	bc.Ledger.Balances["Alice"] = 999999
	if err := bc.RebuildLedger(); err != nil {
		t.Fatal(err)
	}
	if bc.Ledger.GetBalance("Alice") != 80 {
		t.Fatalf("expected replayed balance 80, got %d", bc.Ledger.GetBalance("Alice"))
	}
}

func TestPendingPoolDoubleSpend(t *testing.T) {
	bc := blockchain.NewBlockchain()
	if err := bc.AddTransaction(signedTransaction(t, "Alice", "Bob", 80)); err != nil {
		t.Fatal(err)
	}
	if err := bc.AddTransaction(signedTransaction(t, "Alice", "Charlie", 30)); err == nil {
		t.Fatal("expected second pending transaction to be rejected")
	}
}

func TestMineRevalidatesInjectedPendingTransaction(t *testing.T) {
	bc := blockchain.NewBlockchain()
	bc.PendingTransactions = append(bc.PendingTransactions, signedTransaction(t, "Bob", "Alice", 5000))
	if _, err := bc.MinePendingTransactions(); err == nil {
		t.Fatal("expected mining to reject injected overspending transaction")
	}
}
