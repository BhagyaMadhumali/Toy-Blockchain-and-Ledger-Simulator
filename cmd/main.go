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

	// Reset must be handled before loading and validating saved data.
	// This allows incompatible data from an older block format to be removed.
	if command == "reset" {
		resetCommand()
		return
	}

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
	case "resolve":
		resolveForkCommand(bc, os.Args[2:])
	case "export":
		exportCommand(bc, os.Args[2:])
	case "balance":
		if err := bc.PrintBalances(); err != nil {
			fmt.Println("Cannot calculate balances:", err)
		}
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
	privateKey, err := ledger.DemoAccountPrivateKey(tx.Sender)
	if err != nil {
		fmt.Println("Cannot sign transaction:", err)
		return
	}
	if err := ledger.SignTransaction(&tx, privateKey); err != nil {
		fmt.Println("Cannot sign transaction:", err)
		return
	}
	if err := bc.AddTransaction(tx); err != nil {
		fmt.Println("Transaction rejected:", err)
		return
	}
	if err := bc.SaveBlockchain(dataFile); err != nil {
		fmt.Println("Error saving blockchain:", err)
		return
	}
	fmt.Println("Transaction signed and added successfully.")
	fmt.Println("Signer key fingerprint:", tx.PublicKeyFingerprint())
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
	fmt.Printf("Workers: %d\n", result.Workers)
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

func resolveForkCommand(bc *blockchain.Blockchain, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: go run ./cmd resolve <candidate-blockchain.json>")
		return
	}

	candidate, err := blockchain.LoadCandidateBlocks(args[0])
	if err != nil {
		fmt.Println("Cannot load candidate blockchain:", err)
		return
	}

	result, err := bc.ResolveFork(candidate)
	if err != nil {
		fmt.Println("Fork resolution failed:", err)
		return
	}

	if !result.Replaced {
		fmt.Printf("Local chain kept. Local length: %d, candidate length: %d.\n", result.OldLength, result.CandidateLength)
		return
	}

	if err := bc.SaveBlockchain(dataFile); err != nil {
		fmt.Println("Fork was resolved but the winning chain could not be saved:", err)
		return
	}

	fmt.Println("Competing chain accepted by the longest-valid-chain rule.")
	fmt.Printf("Old local length: %d\n", result.OldLength)
	fmt.Printf("Winning chain length: %d\n", result.CandidateLength)
	fmt.Printf("Common ancestor block: %d\n", result.CommonAncestorIndex)
	fmt.Printf("Transactions returned to pending: %d\n", result.Requeued)
	fmt.Printf("Duplicate or invalid transactions discarded: %d\n", result.Discarded)
}

func exportCommand(bc *blockchain.Blockchain, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: go run ./cmd export <output-file.json>")
		return
	}
	if err := bc.ExportBlockchain(args[0]); err != nil {
		fmt.Println("Export failed:", err)
		return
	}
	fmt.Println("Blockchain exported to", args[0])
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

	fmt.Println("difficulty,run,workers,nonce,attempts,elapsed_ms,hash")
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
			fmt.Printf("%d,%d,%d,%d,%d,%.3f,%s\n", difficulty, run, result.Workers, block.Nonce, result.Attempts, float64(result.Elapsed.Microseconds())/1000, block.Hash)
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
	fmt.Println("  go run ./cmd resolve <candidate-blockchain.json>")
	fmt.Println("  go run ./cmd export <output-file.json>")
	fmt.Println("  go run ./cmd balance")
	fmt.Println("  go run ./cmd benchmark --min 1 --max 5 --runs 3")
	fmt.Println("  go run ./cmd reset")
}
