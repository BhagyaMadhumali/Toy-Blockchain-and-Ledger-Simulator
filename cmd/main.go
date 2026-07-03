package main

import (
	"toy-blockchain/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	bc.PrintChain()
}