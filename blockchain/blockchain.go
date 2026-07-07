// package blockchain

// import (
// 	"fmt"
// 	"time"
// )

// type Blockchain struct {
// 	Blocks []Block
// }

// func NewBlockchain() *Blockchain {

// 	genesisBlock := Block{
// 		Index:        0,
// 		Timestamp:    time.Now().Unix(),
// 		Data:         "Genesis Block",
// 		PreviousHash: "0000000000",
// 	}

// 	genesisBlock.Hash = CalculateHash(genesisBlock)

// 	return &Blockchain{
// 		Blocks: []Block{genesisBlock},
// 	}
// }

// func (bc *Blockchain) PrintChain() {
// 	for _, block := range bc.Blocks {
// 		fmt.Println("---------------")
// 		fmt.Println("Index:", block.Index)
// 		fmt.Println("Timestamp:", block.Timestamp)
// 		fmt.Println("Data:", block.Data)
// 		fmt.Println("Previous Hash:", block.PreviousHash)
// 		fmt.Println("Hash:", block.Hash)
// 	}
// }




package blockchain

import (
	"fmt"
	"time"
	"toy-blockchain/ledger"
)

type Blockchain struct {
	Blocks               []Block
	PendingTransactions  []ledger.Transaction
	Difficulty           int
	Ledger               *ledger.Ledger
}

func NewBlockchain() *Blockchain {
	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []ledger.Transaction{},
		PreviousHash: "0000",
		Nonce:        0,
	}

	genesis.Hash = CalculateHash(genesis)

	return &Blockchain{
		Blocks:              []Block{genesis},
		Difficulty:          4,
		Ledger:              ledger.NewLedger(),
		PendingTransactions: []ledger.Transaction{},
	}
}

func (bc *Blockchain) AddTransaction(tx ledger.Transaction) {
	bc.PendingTransactions = append(bc.PendingTransactions, tx)
}

func (bc *Blockchain) PrintChain() {
	for _, b := range bc.Blocks {
		fmt.Println("---------------")
		fmt.Println("Index:", b.Index)
		fmt.Println("Time:", b.Timestamp)
		fmt.Println("Prev:", b.PreviousHash)
		fmt.Println("Nonce:", b.Nonce)
		fmt.Println("Hash:", b.Hash)

		for _, tx := range b.Transactions {
			fmt.Println("  ", tx.Sender, "->", tx.Receiver, tx.Amount)
		}
	}
}