package blockchain

import (
	"fmt"
	"time"
	"toy-blockchain/ledger"
)

type Blockchain struct {
	Blocks              []Block              `json:"blocks"`
	PendingTransactions []ledger.Transaction `json:"pending_transactions"`
	Difficulty          int                  `json:"difficulty"`
	Ledger              *ledger.Ledger       `json:"ledger"`
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Blocks:              []Block{},
		Difficulty:          4,
		Ledger:              ledger.NewLedger(),
		PendingTransactions: []ledger.Transaction{},
	}

	bc.InitializeAccounts()

	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []ledger.Transaction{},
		PreviousHash: "0000",
		Nonce:        0,
	}

	genesis.Hash = CalculateHash(genesis)
	bc.Blocks = append(bc.Blocks, genesis)

	return bc
}

func (bc *Blockchain) InitializeAccounts() {
	bc.Ledger.AddAccount("Alice", 100)
	bc.Ledger.AddAccount("Bob", 50)
	bc.Ledger.AddAccount("Charlie", 75)
}

func (bc *Blockchain) AddTransaction(tx ledger.Transaction) error {
	if err := ledger.ValidateTransaction(bc.Ledger, tx); err != nil {
		return err
	}

	bc.PendingTransactions = append(bc.PendingTransactions, tx)
	return nil
}

func (bc *Blockchain) MinePendingTransactions() error {
	if len(bc.PendingTransactions) == 0 {
		return fmt.Errorf("no pending transactions to mine")
	}

	block := Block{
		Index:        len(bc.Blocks),
		Timestamp:    time.Now().Unix(),
		Transactions: bc.PendingTransactions,
		PreviousHash: bc.Blocks[len(bc.Blocks)-1].Hash,
		Nonce:        0,
	}

	MineBlock(&block, bc.Difficulty)

	for _, tx := range bc.PendingTransactions {
		bc.Ledger.ApplyTransaction(tx)
	}

	bc.Blocks = append(bc.Blocks, block)
	bc.PendingTransactions = []ledger.Transaction{}

	return nil
}

func (bc *Blockchain) PrintChain() {
	for _, b := range bc.Blocks {
		fmt.Println("---------------")
		fmt.Println("Index:", b.Index)
		fmt.Println("Time:", b.Timestamp)
		fmt.Println("Prev:", b.PreviousHash)
		fmt.Println("Nonce:", b.Nonce)
		fmt.Println("Hash:", b.Hash)
		fmt.Println("Transactions:")

		if len(b.Transactions) == 0 {
			fmt.Println("   No transactions")
		}

		for _, tx := range b.Transactions {
			fmt.Println("  ", tx.Sender, "->", tx.Receiver, tx.Amount)
		}
	}
}

func (bc *Blockchain) PrintPendingTransactions() {
	if len(bc.PendingTransactions) == 0 {
		fmt.Println("No pending transactions.")
		return
	}

	fmt.Println("Pending Transactions:")
	for _, tx := range bc.PendingTransactions {
		fmt.Println(tx.Sender, "->", tx.Receiver, tx.Amount)
	}
}

func (bc *Blockchain) PrintBalances() {
	fmt.Println("Account Balances:")

	for user, balance := range bc.Ledger.Balances {
		fmt.Println(user+":", balance)
	}
}

func (bc *Blockchain) RebuildLedger() {
	newLedger := ledger.NewLedger()

	newLedger.AddAccount("Alice", 100)
	newLedger.AddAccount("Bob", 50)
	newLedger.AddAccount("Charlie", 75)

	for i, block := range bc.Blocks {
		if i == 0 {
			continue
		}

		for _, tx := range block.Transactions {
			newLedger.ApplyTransaction(tx)
		}
	}

	bc.Ledger = newLedger
}