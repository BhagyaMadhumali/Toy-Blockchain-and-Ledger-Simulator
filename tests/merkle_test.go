package tests

import (
	"testing"
	"toy-blockchain/blockchain"
	"toy-blockchain/ledger"
)

func TestMerkleRootIsDeterministic(t *testing.T) {
	transactions := []ledger.Transaction{
		signedTransaction(t, "Alice", "Bob", 10),
		signedTransaction(t, "Bob", "Charlie", 5),
	}

	first := blockchain.CalculateMerkleRoot(transactions)
	second := blockchain.CalculateMerkleRoot(transactions)
	if first != second {
		t.Fatal("same transactions must produce the same Merkle root")
	}
}

func TestMerkleRootChangesWhenTransactionChanges(t *testing.T) {
	transactions := []ledger.Transaction{
		signedTransaction(t, "Alice", "Bob", 10),
		signedTransaction(t, "Bob", "Charlie", 5),
	}
	original := blockchain.CalculateMerkleRoot(transactions)

	transactions[0].Amount = 11
	changed := blockchain.CalculateMerkleRoot(transactions)
	if original == changed {
		t.Fatal("changed transaction must change the Merkle root")
	}
}

func TestMerkleRootChangesWhenTransactionOrderChanges(t *testing.T) {
	first := signedTransaction(t, "Alice", "Bob", 10)
	second := signedTransaction(t, "Bob", "Charlie", 5)

	rootOne := blockchain.CalculateMerkleRoot([]ledger.Transaction{first, second})
	rootTwo := blockchain.CalculateMerkleRoot([]ledger.Transaction{second, first})
	if rootOne == rootTwo {
		t.Fatal("transaction order must affect the Merkle root")
	}
}

func TestOddMerkleTreeIsDeterministic(t *testing.T) {
	transactions := []ledger.Transaction{
		signedTransaction(t, "Alice", "Bob", 10),
		signedTransaction(t, "Bob", "Charlie", 5),
		signedTransaction(t, "Charlie", "Alice", 2),
	}

	root := blockchain.CalculateMerkleRoot(transactions)
	if root == "" {
		t.Fatal("odd transaction count must produce a Merkle root")
	}
	if root != blockchain.CalculateMerkleRoot(transactions) {
		t.Fatal("odd transaction tree must be deterministic")
	}
}

func TestBlockHashUsesStoredMerkleRoot(t *testing.T) {
	block := blockchain.Block{
		Index:        1,
		Timestamp:    1000,
		Transactions: []ledger.Transaction{{Sender: "Alice", Receiver: "Bob", Amount: 1}},
		MerkleRoot:   "fixed-root",
		PreviousHash: "previous",
		Difficulty:   3,
	}

	originalHash := blockchain.CalculateHash(block)
	block.Transactions[0].Amount = 999
	if originalHash != blockchain.CalculateHash(block) {
		t.Fatal("block hash must use MerkleRoot instead of hashing raw transactions")
	}
}
