package blockchain

import (
	"fmt"
	"time"
)

type Blockchain struct {
	Blocks []Block
}

func NewBlockchain() *Blockchain {

	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Data:         "Genesis Block",
		PreviousHash: "0000000000",
	}

	genesisBlock.Hash = CalculateHash(genesisBlock)

	return &Blockchain{
		Blocks: []Block{genesisBlock},
	}
}

func (bc *Blockchain) PrintChain() {
	for _, block := range bc.Blocks {
		fmt.Println("---------------")
		fmt.Println("Index:", block.Index)
		fmt.Println("Timestamp:", block.Timestamp)
		fmt.Println("Data:", block.Data)
		fmt.Println("Previous Hash:", block.PreviousHash)
		fmt.Println("Hash:", block.Hash)
	}
}