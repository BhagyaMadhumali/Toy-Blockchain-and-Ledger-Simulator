package blockchain

import "toy-blockchain/ledger"

// Block represents one block in the blockchain.
type Block struct {
	Index        int                  `json:"index"`
	Timestamp    int64                `json:"timestamp"`
	Transactions []ledger.Transaction `json:"transactions"`
	PreviousHash string               `json:"previous_hash"`
	Nonce        int                  `json:"nonce"`
	Hash         string               `json:"hash"`
}