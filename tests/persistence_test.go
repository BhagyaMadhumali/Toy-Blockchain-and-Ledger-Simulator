package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"toy-blockchain/blockchain"
)

func TestPersistenceRebuildsLedgerAndDoesNotSaveBalances(t *testing.T) {
	path := filepath.Join(t.TempDir(), "blockchain.json")
	bc := blockchain.NewBlockchain()
	if err := bc.AddTransaction(signedTransaction(t, "Alice", "Bob", 15)); err != nil {
		t.Fatal(err)
	}
	if _, err := bc.MinePendingTransactions(); err != nil {
		t.Fatal(err)
	}
	if err := bc.SaveBlockchain(path); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "balances") || strings.Contains(string(data), "ledger") {
		t.Fatal("derived ledger must not be persisted")
	}

	loaded, err := blockchain.LoadBlockchain(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Ledger.GetBalance("Alice") != 85 || loaded.Ledger.GetBalance("Bob") != 65 {
		t.Fatal("loaded balances were not rebuilt from chain")
	}
}

func TestLoadRejectsInjectedPendingOverspend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "blockchain.json")
	bc := blockchain.NewBlockchain()
	if err := bc.SaveBlockchain(path); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	modified := strings.Replace(string(data), `"pending_transactions": []`, `"pending_transactions": [{"sender":"Bob","receiver":"Alice","amount":5000}]`, 1)
	if err := os.WriteFile(path, []byte(modified), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := blockchain.LoadBlockchain(path); err == nil {
		t.Fatal("expected invalid pending pool to be rejected")
	}
}
