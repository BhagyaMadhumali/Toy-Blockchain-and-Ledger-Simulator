package blockchain

import (
	"fmt"
	"time"
	"toy-blockchain/ledger"
)

// DefaultDifficulty controls how many zeroes a mined hash must start with.
const DefaultDifficulty = 3

// Blockchain stores the blocks, pending transactions,
// mining difficulty, and account ledger.
type Blockchain struct {
	Blocks              []Block              `json:"blocks"`
	PendingTransactions []ledger.Transaction `json:"pending_transactions"`
	Difficulty          int                  `json:"difficulty"`
	Ledger              *ledger.Ledger       `json:"ledger"`
}

// NewBlockchain creates a new blockchain.
func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Blocks:              make([]Block, 0),
		PendingTransactions: make([]ledger.Transaction, 0),
		Difficulty:          DefaultDifficulty,
		Ledger:              ledger.NewLedger(),
	}

	bc.InitializeAccounts()

	genesisBlock := NewGenesisBlock()
	bc.Blocks = append(bc.Blocks, genesisBlock)

	return bc
}

// NewGenesisBlock creates the first block in the blockchain.
func NewGenesisBlock() Block {
	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: make([]ledger.Transaction, 0),
		PreviousHash: "0000",
		Nonce:        0,
	}

	genesis.Hash = CalculateHash(genesis)

	return genesis
}

// InitializeAccounts creates sample blockchain accounts.
func (bc *Blockchain) InitializeAccounts() {
	bc.Ledger.AddAccount("Alice", 100)
	bc.Ledger.AddAccount("Bob", 50)
	bc.Ledger.AddAccount("Charlie", 75)
}

// AddTransaction validates and adds a transaction
// to the pending transaction pool.
func (bc *Blockchain) AddTransaction(tx ledger.Transaction) error {
	availableBalance := bc.availableBalance(tx.Sender)

	if err := ledger.ValidateTransactionWithBalance(
		tx,
		availableBalance,
	); err != nil {
		return err
	}

	bc.PendingTransactions = append(
		bc.PendingTransactions,
		tx,
	)

	return nil
}

// availableBalance calculates the amount a sender can still spend.
//
// It subtracts transactions that are already waiting to be mined.
// This prevents a sender from adding several pending transactions
// that together exceed the account balance.
func (bc *Blockchain) availableBalance(sender string) int {
	balance := bc.Ledger.GetBalance(sender)

	for _, tx := range bc.PendingTransactions {
		if tx.Sender == sender {
			balance -= tx.Amount
		}
	}

	return balance
}

// PrintChain prints all blocks and transactions.
func (bc *Blockchain) PrintChain() {
	for _, block := range bc.Blocks {
		fmt.Println("--------------------------------")
		fmt.Println("Index:", block.Index)
		fmt.Println("Timestamp:", block.Timestamp)
		fmt.Println("Previous Hash:", block.PreviousHash)
		fmt.Println("Nonce:", block.Nonce)
		fmt.Println("Hash:", block.Hash)
		fmt.Println("Transactions:")

		if len(block.Transactions) == 0 {
			fmt.Println("  No transactions")
		}

		for _, tx := range block.Transactions {
			fmt.Printf(
				"  %s -> %s : %d\n",
				tx.Sender,
				tx.Receiver,
				tx.Amount,
			)
		}
	}
}

// PrintPendingTransactions displays transactions
// that have not yet been mined.
func (bc *Blockchain) PrintPendingTransactions() {
	if len(bc.PendingTransactions) == 0 {
		fmt.Println("No pending transactions.")
		return
	}

	fmt.Println("Pending Transactions:")

	for index, tx := range bc.PendingTransactions {
		fmt.Printf(
			"%d. %s -> %s : %d\n",
			index+1,
			tx.Sender,
			tx.Receiver,
			tx.Amount,
		)
	}
}

// PrintBalances displays every account balance.
func (bc *Blockchain) PrintBalances() {
	fmt.Println("Account Balances:")

	for user, balance := range bc.Ledger.Balances {
		fmt.Printf("%s: %d\n", user, balance)
	}
}

// RebuildLedger reconstructs account balances
// using all mined transactions.
func (bc *Blockchain) RebuildLedger() {
	newLedger := ledger.NewLedger()

	newLedger.AddAccount("Alice", 100)
	newLedger.AddAccount("Bob", 50)
	newLedger.AddAccount("Charlie", 75)

	for blockIndex, block := range bc.Blocks {
		// Skip the genesis block.
		if blockIndex == 0 {
			continue
		}

		for _, tx := range block.Transactions {
			newLedger.ApplyTransaction(tx)
		}
	}

	bc.Ledger = newLedger
}