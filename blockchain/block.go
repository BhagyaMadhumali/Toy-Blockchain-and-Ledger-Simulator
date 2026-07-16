package blockchain

import "toy-blockchain/ledger"

// Block represents one immutable unit in the blockchain.
type Block struct {
	Index        int                  `json:"index"`
	Timestamp    int64                `json:"timestamp"`
	Transactions []ledger.Transaction `json:"transactions"`
	MerkleRoot   string               `json:"merkle_root"`
	PreviousHash string               `json:"previous_hash"`
	Difficulty   int                  `json:"difficulty"`
	Nonce        uint64               `json:"nonce"`
	Hash         string               `json:"hash"`
}
