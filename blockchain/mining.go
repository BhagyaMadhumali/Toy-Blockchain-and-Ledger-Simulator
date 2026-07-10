package blockchain

// Mining logic is implemented as Blockchain.MinePendingTransactions
// in blockchain.go.
//
// This file is kept to match the Day 2 and Day 3 project structure.package blockchain

import (
	"fmt"
	"time"
	"toy-blockchain/ledger"
)

// MinePendingTransactions places all pending transactions
// into one new block and mines that block.
func (bc *Blockchain) MinePendingTransactions() error {
	// Mining cannot happen without transactions.
	if len(bc.PendingTransactions) == 0 {
		return fmt.Errorf(
			"no pending transactions to mine; add a transaction first",
		)
	}

	// Get the last block from the chain.
	previousBlock := bc.Blocks[len(bc.Blocks)-1]

	// Create a separate copy of pending transactions.
	transactions := append(
		[]ledger.Transaction(nil),
		bc.PendingTransactions...,
	)

	// Create the new block.
	newBlock := Block{
		Index:        len(bc.Blocks),
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PreviousHash: previousBlock.Hash,
		Nonce:        0,
	}

	// Perform Proof of Work.
	attempts := MineBlock(
		&newBlock,
		bc.Difficulty,
	)

	// Update balances only after successful mining.
	for _, tx := range newBlock.Transactions {
		bc.Ledger.ApplyTransaction(tx)
	}

	// Add the mined block to the blockchain.
	bc.Blocks = append(
		bc.Blocks,
		newBlock,
	)

	// Clear the pending transaction pool.
	bc.PendingTransactions = make(
		[]ledger.Transaction,
		0,
	)

	fmt.Printf(
		"Mining complete: nonce=%d, attempts=%d, hash=%s\n",
		newBlock.Nonce,
		attempts,
		newBlock.Hash,
	)

	return nil
}