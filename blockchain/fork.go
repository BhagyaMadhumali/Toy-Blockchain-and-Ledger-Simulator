package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"toy-blockchain/ledger"
)

// ForkResolutionResult describes whether a competing chain replaced the local
// chain and how pending/orphaned transactions were reconciled.
type ForkResolutionResult struct {
	Replaced            bool
	OldLength           int
	CandidateLength     int
	CommonAncestorIndex int
	Requeued            int
	Discarded           int
}

// ResolveFork applies the longest-valid-chain rule.
//
// A candidate replaces the local chain only when:
//  1. the candidate is fully valid under all consensus rules; and
//  2. it contains more blocks than the local chain.
//
// Transactions from orphaned local blocks and the old pending pool are then
// reconsidered against the winning chain. Transactions already present in the
// winning chain, duplicates, or transactions that are no longer valid are not
// requeued.
func (bc *Blockchain) ResolveFork(candidate []Block) (ForkResolutionResult, error) {
	result := ForkResolutionResult{
		OldLength:           len(bc.Blocks),
		CandidateLength:     len(candidate),
		CommonAncestorIndex: -1,
	}

	if err := bc.ValidateBlockchain(); err != nil {
		return result, fmt.Errorf("local blockchain is invalid: %w", err)
	}

	candidateCopy := cloneBlocks(candidate)
	candidateChain := &Blockchain{
		Blocks:              candidateCopy,
		PendingTransactions: []ledger.Transaction{},
		Ledger:              ledger.NewLedger(),
	}
	if err := candidateChain.ValidateBlockchain(); err != nil {
		return result, fmt.Errorf("candidate blockchain is invalid: %w", err)
	}

	if len(candidateCopy) <= len(bc.Blocks) {
		return result, nil
	}

	commonAncestor := findCommonAncestor(bc.Blocks, candidateCopy)
	result.CommonAncestorIndex = commonAncestor
	if commonAncestor < 0 {
		return result, fmt.Errorf("candidate blockchain does not share the trusted genesis block")
	}

	winningLedger, err := ReplayBlocks(candidateCopy)
	if err != nil {
		return result, fmt.Errorf("replay candidate blockchain: %w", err)
	}

	included := transactionSet(candidateCopy)
	possiblePending := make([]ledger.Transaction, 0)

	// Transactions from blocks after the common ancestor became orphaned.
	for i := commonAncestor + 1; i < len(bc.Blocks); i++ {
		for _, tx := range bc.Blocks[i].Transactions {
			if tx.Sender != ledger.SystemAccount {
				possiblePending = append(possiblePending, tx)
			}
		}
	}
	possiblePending = append(possiblePending, bc.PendingTransactions...)

	reconciled := make([]ledger.Transaction, 0, len(possiblePending))
	seenPending := make(map[string]struct{})
	temporary := winningLedger.Clone()

	for _, tx := range possiblePending {
		id := transactionID(tx)
		if _, exists := included[id]; exists {
			result.Discarded++
			continue
		}
		if _, duplicate := seenPending[id]; duplicate {
			result.Discarded++
			continue
		}
		if err := temporary.ApplyTransaction(tx); err != nil {
			result.Discarded++
			continue
		}
		seenPending[id] = struct{}{}
		reconciled = append(reconciled, tx)
		result.Requeued++
	}

	bc.Blocks = candidateCopy
	bc.PendingTransactions = reconciled
	bc.Ledger = winningLedger
	result.Replaced = true
	return result, nil
}

func findCommonAncestor(local, candidate []Block) int {
	limit := len(local)
	if len(candidate) < limit {
		limit = len(candidate)
	}

	ancestor := -1
	for i := 0; i < limit; i++ {
		if local[i].Hash != candidate[i].Hash {
			break
		}
		ancestor = i
	}
	return ancestor
}

func cloneBlocks(blocks []Block) []Block {
	cloned := make([]Block, len(blocks))
	for i, block := range blocks {
		cloned[i] = block
		cloned[i].Transactions = append([]ledger.Transaction(nil), block.Transactions...)
	}
	return cloned
}

func transactionSet(blocks []Block) map[string]struct{} {
	set := make(map[string]struct{})
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			if tx.Sender != ledger.SystemAccount {
				set[transactionID(tx)] = struct{}{}
			}
		}
	}
	return set
}

func transactionID(tx ledger.Transaction) string {
	encoded, err := json.Marshal(tx)
	if err != nil {
		panic("transaction serialization failed: " + err.Error())
	}
	hash := sha256.Sum256(encoded)
	return hex.EncodeToString(hash[:])
}
