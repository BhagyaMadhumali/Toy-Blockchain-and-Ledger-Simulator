package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
	"toy-blockchain/blockchain"
	"toy-blockchain/ledger"
)

const dataFile = "data/blockchain.json"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]
	if command == "benchmark" {
		benchmarkCommand(os.Args[2:])
		return
	}

	bc, err := blockchain.LoadBlockchain(dataFile)
	if err != nil {
		fmt.Println("Error loading blockchain:", err)
		return
	}

	switch command {
	case "add":
		addTransactionCommand(bc, os.Args[2:])
	case "mine":
		mineCommand(bc)
	case "print":
		bc.PrintChain()
	case "pending":
		bc.PrintPendingTransactions()
	case "validate":
		validateCommand(bc)
	case "balance":
		if err := bc.PrintBalances(); err != nil {
			fmt.Println("Cannot calculate balances:", err)
		}
	case "reset":
		resetCommand()
	case "help":
		printHelp()
	default:
		fmt.Println("Unknown command:", command)
		printHelp()
	}
}

func addTransactionCommand(bc *blockchain.Blockchain, args []string) {
	if len(args) != 3 {
		fmt.Println("Usage: go run ./cmd add <sender> <receiver> <amount>")
		return
	}
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Invalid amount. Amount must be an integer.")
		return
	}
	tx := ledger.Transaction{Sender: args[0], Receiver: args[1], Amount: amount}
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

func mineCommand(bc *blockchain.Blockchain) {
	result, err := bc.MinePendingTransactions()
	if err != nil {
		fmt.Println("Mining failed:", err)
		return
	}
	if err := bc.SaveBlockchain(dataFile); err != nil {
		fmt.Println("Error saving blockchain:", err)
		return
	}
	block := bc.Blocks[len(bc.Blocks)-1]
	fmt.Println("Block mined successfully.")
	fmt.Printf("Difficulty: %d\n", block.Difficulty)
	fmt.Printf("Nonce: %d\n", block.Nonce)
	fmt.Printf("Hashes attempted: %d\n", result.Attempts)
	fmt.Printf("Elapsed time: %s\n", result.Elapsed)
	fmt.Printf("Hash: %s\n", block.Hash)
}

func validateCommand(bc *blockchain.Blockchain) {
	if err := bc.ValidateBlockchain(); err != nil {
		fmt.Println("Blockchain is invalid:", err)
		return
	}
	fmt.Println("Blockchain is valid.")
}

// benchmark mines independent sample blocks and does not modify the real chain.
func benchmarkCommand(args []string) {
	flags := flag.NewFlagSet("benchmark", flag.ContinueOnError)
	minDifficulty := flags.Int("min", 1, "minimum difficulty")
	maxDifficulty := flags.Int("max", 5, "maximum difficulty")
	runs := flags.Int("runs", 1, "runs per difficulty")
	if err := flags.Parse(args); err != nil {
		return
	}
	if *minDifficulty < 1 || *maxDifficulty > 6 || *minDifficulty > *maxDifficulty || *runs < 1 {
		fmt.Println("Use: benchmark --min 1 --max 6 --runs 1")
		return
	}

	fmt.Println("difficulty,run,nonce,attempts,elapsed_ms,hash")
	for difficulty := *minDifficulty; difficulty <= *maxDifficulty; difficulty++ {
		for run := 1; run <= *runs; run++ {
			block := blockchain.Block{
				Index:        1,
				Timestamp:    time.Now().UnixNano(),
				PreviousHash: "benchmark",
				Transactions: []ledger.Transaction{{Sender: "Alice", Receiver: "Bob", Amount: 1}},
			}
			result, err := blockchain.MineBlock(&block, difficulty)
			if err != nil {
				fmt.Println("Benchmark failed:", err)
				return
			}
			fmt.Printf("%d,%d,%d,%d,%.3f,%s\n", difficulty, run, block.Nonce, result.Attempts, float64(result.Elapsed.Microseconds())/1000, block.Hash)
		}
	}
}

func resetCommand() {
	if err := os.Remove(dataFile); err != nil && !os.IsNotExist(err) {
		fmt.Println("Reset failed:", err)
		return
	}
	fmt.Println("Runtime blockchain data removed. The next command will create the fixed genesis chain.")
}

func printHelp() {
	fmt.Println("Toy Blockchain CLI")
	fmt.Println("  go run ./cmd add <sender> <receiver> <amount>")
	fmt.Println("  go run ./cmd mine")
	fmt.Println("  go run ./cmd print")
	fmt.Println("  go run ./cmd pending")
	fmt.Println("  go run ./cmd validate")
	fmt.Println("  go run ./cmd balance")
	fmt.Println("  go run ./cmd benchmark --min 1 --max 5 --runs 3")
	fmt.Println("  go run ./cmd reset")
}
