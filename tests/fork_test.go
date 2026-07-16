package tests

import (
	"testing"

	"toy-blockchain/blockchain"
)

func mineSigned(t *testing.T, bc *blockchain.Blockchain, sender, receiver string, amount int) {
	t.Helper()
	if err := bc.AddTransaction(signedTransaction(t, sender, receiver, amount)); err != nil {
		t.Fatal(err)
	}
	if _, err := bc.MinePendingTransactions(); err != nil {
		t.Fatal(err)
	}
}

func TestResolveForkAcceptsLongerValidChain(t *testing.T) {
	local := blockchain.NewBlockchain()
	mineSigned(t, local, "Alice", "Bob", 10)

	candidate := blockchain.NewBlockchain()
	mineSigned(t, candidate, "Alice", "Charlie", 5)
	mineSigned(t, candidate, "Charlie", "Bob", 2)

	result, err := local.ResolveFork(candidate.Blocks)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Replaced {
		t.Fatal("expected longer valid candidate to replace local chain")
	}
	if len(local.Blocks) != len(candidate.Blocks) {
		t.Fatalf("got chain length %d, want %d", len(local.Blocks), len(candidate.Blocks))
	}
	if result.CommonAncestorIndex != 0 {
		t.Fatalf("got common ancestor %d, want 0", result.CommonAncestorIndex)
	}
	if result.Requeued != 1 {
		t.Fatalf("got %d requeued transactions, want 1", result.Requeued)
	}
	if len(local.PendingTransactions) != 1 {
		t.Fatalf("got %d pending transactions, want 1", len(local.PendingTransactions))
	}
	if err := local.ValidateBlockchain(); err != nil {
		t.Fatalf("resolved chain is invalid: %v", err)
	}
}

func TestResolveForkKeepsEqualOrShorterChain(t *testing.T) {
	local := blockchain.NewBlockchain()
	mineSigned(t, local, "Alice", "Bob", 10)

	equal := blockchain.NewBlockchain()
	mineSigned(t, equal, "Alice", "Charlie", 5)

	originalTip := local.Blocks[len(local.Blocks)-1].Hash
	result, err := local.ResolveFork(equal.Blocks)
	if err != nil {
		t.Fatal(err)
	}
	if result.Replaced {
		t.Fatal("equal-length candidate must not replace local chain")
	}
	if local.Blocks[len(local.Blocks)-1].Hash != originalTip {
		t.Fatal("local chain changed after equal-length fork")
	}
}

func TestResolveForkRejectsInvalidLongerChain(t *testing.T) {
	local := blockchain.NewBlockchain()

	candidate := blockchain.NewBlockchain()
	mineSigned(t, candidate, "Alice", "Bob", 5)
	mineSigned(t, candidate, "Bob", "Charlie", 2)

	candidate.Blocks[1].Transactions[0].Amount = 500

	result, err := local.ResolveFork(candidate.Blocks)
	if err == nil {
		t.Fatal("expected invalid candidate to be rejected")
	}
	if result.Replaced {
		t.Fatal("invalid candidate must not replace local chain")
	}
	if len(local.Blocks) != 1 {
		t.Fatal("local chain changed after invalid candidate")
	}
}

func TestResolveForkDoesNotRequeueTransactionAlreadyInWinner(t *testing.T) {
	sharedTx := signedTransaction(t, "Alice", "Bob", 10)

	// The local chain confirms sharedTx in its first normal block.
	local := blockchain.NewBlockchain()
	if err := local.AddTransaction(sharedTx); err != nil {
		t.Fatal(err)
	}
	if _, err := local.MinePendingTransactions(); err != nil {
		t.Fatal(err)
	}

	// Force the candidate to diverge immediately after genesis. The candidate
	// confirms sharedTx later, so the local copy becomes orphaned but is already
	// present in the winning chain.
	candidate := blockchain.NewBlockchain()
	mineSigned(t, candidate, "Alice", "Charlie", 5)
	if err := candidate.AddTransaction(sharedTx); err != nil {
		t.Fatal(err)
	}
	if _, err := candidate.MinePendingTransactions(); err != nil {
		t.Fatal(err)
	}

	result, err := local.ResolveFork(candidate.Blocks)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Replaced {
		t.Fatal("expected candidate to replace local chain")
	}
	if result.CommonAncestorIndex != 0 {
		t.Fatalf("got common ancestor %d, want 0", result.CommonAncestorIndex)
	}
	if len(local.PendingTransactions) != 0 {
		t.Fatal("transaction already confirmed in winner was requeued")
	}
	if result.Discarded != 1 {
		t.Fatalf("got discarded %d, want 1", result.Discarded)
	}
}
