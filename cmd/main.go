package main

import (
	"fmt"
	"os"
	"strconv"
	"toy-blockchain/blockchain"
	"toy-blockchain/ledger"
)

const dataFile = "data/blockchain.json"

func main() {
	// Load the existing blockchain.
	// A new blockchain is created if the file does not exist.
	bc, err := blockchain.LoadBlockchain(dataFile)

	if err != nil {
		fmt.Println("Error loading blockchain:", err)
		return
	}

	// The program requires at least one command.
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		addTransactionCommand(bc)

	case "mine":
		mineCommand(bc)

	case "print":
		bc.PrintChain()

	case "pending":
		bc.PrintPendingTransactions()

	case "validate":
		validateCommand(bc)

	case "balance":
		bc.PrintBalances()

	case "save":
		saveCommand(bc)

	case "load":
		fmt.Println("Blockchain loaded successfully.")
		bc.PrintChain()

	case "help":
		printHelp()

	default:
		fmt.Println("Unknown command:", command)
		printHelp()
	}
}

// addTransactionCommand handles the add command.
func addTransactionCommand(bc *blockchain.Blockchain) {
	if len(os.Args) != 5 {
		fmt.Println(
			"Usage: go run cmd/main.go add <sender> <receiver> <amount>",
		)
		return
	}

	amount, err := strconv.Atoi(os.Args[4])

	if err != nil {
		fmt.Println(
			"Invalid amount. Amount must be an integer.",
		)
		return
	}

	tx := ledger.Transaction{
		Sender:   os.Args[2],
		Receiver: os.Args[3],
		Amount:   amount,
	}

	if err := bc.AddTransaction(tx); err != nil {
		fmt.Println("Transaction rejected:", err)
		return
	}

	if err := bc.SaveBlockchain(dataFile); err != nil {
		fmt.Println("Error saving blockchain:", err)
		return
	}

	fmt.Println("Transaction added successfully.")
}

// mineCommand mines all pending transactions.
func mineCommand(bc *blockchain.Blockchain) {
	fmt.Println("Mining pending transactions...")

	if err := bc.MinePendingTransactions(); err != nil {
		fmt.Println("Mining failed:", err)
		return
	}

	if err := bc.SaveBlockchain(dataFile); err != nil {
		fmt.Println("Error saving blockchain:", err)
		return
	}

	fmt.Println("Block mined successfully.")
}

// validateCommand validates the blockchain.
func validateCommand(bc *blockchain.Blockchain) {
	if bc.ValidateBlockchain() {
		fmt.Println("Blockchain is valid.")
		return
	}

	fmt.Println(
		"Blockchain is invalid or has been tampered.",
	)
}

// saveCommand manually saves the blockchain.
func saveCommand(bc *blockchain.Blockchain) {
	if err := bc.SaveBlockchain(dataFile); err != nil {
		fmt.Println("Error saving blockchain:", err)
		return
	}

	fmt.Println("Blockchain saved successfully.")
}

// printHelp displays available commands.
func printHelp() {
	fmt.Println("Toy Blockchain CLI")
	fmt.Println()
	fmt.Println("Commands:")

	fmt.Println(
		"  go run cmd/main.go add <sender> <receiver> <amount>",
	)

	fmt.Println(
		"  go run cmd/main.go mine",
	)

	fmt.Println(
		"  go run cmd/main.go print",
	)

	fmt.Println(
		"  go run cmd/main.go pending",
	)

	fmt.Println(
		"  go run cmd/main.go validate",
	)

	fmt.Println(
		"  go run cmd/main.go balance",
	)

	fmt.Println(
		"  go run cmd/main.go save",
	)

	fmt.Println(
		"  go run cmd/main.go load",
	)

	fmt.Println(
		"  go run cmd/main.go help",
	)
}