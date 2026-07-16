package blockchain

import (
	"fmt"
	"sort"
	"time"
	"toy-blockchain/ledger"
)

const (
	DefaultDifficulty   = 3
	GenesisTimestamp    = int64(1704067200) // 2024-01-01 00:00:00 UTC
	GenesisPreviousHash = "0"
)

var genesisTransactions = []ledger.Transaction{
	{Sender: ledger.SystemAccount, Receiver: "Alice", Amount: 100},
	{Sender: ledger.SystemAccount, Receiver: "Bob", Amount: 50},
	{Sender: ledger.SystemAccount, Receiver: "Charlie", Amount: 75},
}

// Blockchain persists only blocks and pending transactions.
// Ledger is always rebuilt from the chain and is never trusted from JSON.
type Blockchain struct {
	Blocks              []Block              `json:"blocks"`
	PendingTransactions []ledger.Transaction `json:"pending_transactions"`
	Ledger              *ledger.Ledger       `json:"-"`
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Blocks:              []Block{NewGenesisBlock()},
		PendingTransactions: []ledger.Transaction{},
		Ledger:              ledger.NewLedger(),
	}
	if err := bc.RebuildLedger(); err != nil {
		panic(err)
	}
	return bc
}

func NewGenesisBlock() Block {
	block := Block{
		Index:        0,
		Timestamp:    GenesisTimestamp,
		Transactions: append([]ledger.Transaction(nil), genesisTransactions...),
		PreviousHash: GenesisPreviousHash,
		Difficulty:   DefaultDifficulty,
	}
	if _, err := MineBlockWithWorkers(&block, DefaultDifficulty, 1); err != nil {
		panic(err)
	}
	return block
}

func (bc *Blockchain) AddTransaction(tx ledger.Transaction) error {
	if err := bc.RebuildLedger(); err != nil {
		return fmt.Errorf("cannot add transaction to invalid chain: %w", err)
	}

	temporary := bc.Ledger.Clone()
	for i, pending := range bc.PendingTransactions {
		if err := temporary.ApplyTransaction(pending); err != nil {
			return fmt.Errorf("existing pending transaction %d is invalid: %w", i, err)
		}
	}
	if err := temporary.ApplyTransaction(tx); err != nil {
		return err
	}

	bc.PendingTransactions = append(bc.PendingTransactions, tx)
	return nil
}

func (bc *Blockchain) RebuildLedger() error {
	rebuilt, err := ReplayBlocks(bc.Blocks)
	if err != nil {
		return err
	}
	bc.Ledger = rebuilt
	return nil
}

func (bc *Blockchain) PrintChain() {
	for _, block := range bc.Blocks {
		fmt.Println("--------------------------------")
		fmt.Println("Index:", block.Index)
		fmt.Println("Timestamp:", block.Timestamp)
		fmt.Println("Previous Hash:", block.PreviousHash)
		fmt.Println("Merkle Root:", block.MerkleRoot)
		fmt.Println("Difficulty:", block.Difficulty)
		fmt.Println("Nonce:", block.Nonce)
		fmt.Println("Hash:", block.Hash)
		fmt.Println("Transactions:")
		for _, tx := range block.Transactions {
			if tx.Sender == ledger.SystemAccount {
				fmt.Printf("  %s -> %s : %d (trusted genesis allocation)\n", tx.Sender, tx.Receiver, tx.Amount)
			} else {
				fmt.Printf("  %s -> %s : %d (signed by %s)\n", tx.Sender, tx.Receiver, tx.Amount, tx.PublicKeyFingerprint())
			}
		}
	}
}

func (bc *Blockchain) PrintPendingTransactions() {
	if len(bc.PendingTransactions) == 0 {
		fmt.Println("No pending transactions.")
		return
	}
	for i, tx := range bc.PendingTransactions {
		fmt.Printf("%d. %s -> %s : %d (signed by %s)\n", i+1, tx.Sender, tx.Receiver, tx.Amount, tx.PublicKeyFingerprint())
	}
}

func (bc *Blockchain) PrintBalances() error {
	if err := bc.RebuildLedger(); err != nil {
		return err
	}
	accounts := make([]string, 0, len(bc.Ledger.Balances))
	for account := range bc.Ledger.Balances {
		accounts = append(accounts, account)
	}
	sort.Strings(accounts)
	fmt.Println("Account Balances:")
	for _, account := range accounts {
		fmt.Printf("%s: %d\n", account, bc.Ledger.GetBalance(account))
	}
	return nil
}

func newCandidateBlock(previous Block, transactions []ledger.Transaction, difficulty int) Block {
	timestamp := time.Now().Unix()
	if timestamp < previous.Timestamp {
		timestamp = previous.Timestamp
	}
	return Block{
		Index:        previous.Index + 1,
		Timestamp:    timestamp,
		Transactions: append([]ledger.Transaction(nil), transactions...),
		PreviousHash: previous.Hash,
		Difficulty:   difficulty,
	}
}
