// package main

// import (
// 	"toy-blockchain/blockchain"
// )

// func main() {
// 	bc := blockchain.NewBlockchain()
// 	bc.PrintChain()
// }




package main

import (
	"fmt"
	"toy-blockchain/blockchain"
	"toy-blockchain/ledger"
)

func main() {

	bc := blockchain.NewBlockchain()

	// setup users
	bc.Ledger.AddAccount("Alice", 100)
	bc.Ledger.AddAccount("Bob", 50)

	fmt.Println("Initial Balances:", bc.Ledger.Balances)

	// transactions
	tx1 := ledger.Transaction{"Alice", "Bob", 20}
	tx2 := ledger.Transaction{"Bob", "Alice", 10}

	bc.AddTransaction(tx1)
	bc.AddTransaction(tx2)

	// mining
	bc.MinePendingTransactions()

	// print
	bc.PrintChain()

	fmt.Println("Final Balances:", bc.Ledger.Balances)
}