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
	bc, err := blockchain.LoadBlockchain(dataFile)
	if err != nil {
		fmt.Println("Error loading blockchain:", err)
		return
	}

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) != 5 {
			fmt.Println("Usage: go run cmd/main.go add <sender> <receiver> <amount>")
			return
		}

		amount, err := strconv.Atoi(os.Args[4])
		if err != nil {
			fmt.Println("Invalid amount. Amount must be an integer.")
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

	case "mine":
		fmt.Println("Mining pending transactions...")

		if err := bc.MinePendingTransactions(); err != nil {
			fmt.Println(err)
			return
		}

		if err := bc.SaveBlockchain(dataFile); err != nil {
			fmt.Println("Error saving blockchain:", err)
			return
		}

		fmt.Println("Block mined successfully.")

	case "print":
		bc.PrintChain()

	case "pending":
		bc.PrintPendingTransactions()

	case "validate":
		if bc.ValidateBlockchain() {
			fmt.Println("Blockchain is valid.")
		} else {
			fmt.Println("Blockchain is invalid or has been tampered.")
		}

	case "balance":
		bc.PrintBalances()

	case "save":
		if err := bc.SaveBlockchain(dataFile); err != nil {
			fmt.Println("Error saving blockchain:", err)
			return
		}

		fmt.Println("Blockchain saved successfully.")

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

func printHelp() {
	fmt.Println("Toy Blockchain CLI")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  go run cmd/main.go add <sender> <receiver> <amount>")
	fmt.Println("  go run cmd/main.go mine")
	fmt.Println("  go run cmd/main.go print")
	fmt.Println("  go run cmd/main.go pending")
	fmt.Println("  go run cmd/main.go validate")
	fmt.Println("  go run cmd/main.go balance")
	fmt.Println("  go run cmd/main.go save")
	fmt.Println("  go run cmd/main.go load")
	fmt.Println("  go run cmd/main.go help")
}