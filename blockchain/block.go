package blockchain

import "toy-blockchain/ledger"

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []ledger.Transaction
	PreviousHash string
	Nonce        int
	Hash         string
}