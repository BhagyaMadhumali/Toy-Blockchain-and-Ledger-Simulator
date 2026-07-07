package blockchain

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"toy-blockchain/ledger"
)

func CalculateHash(block Block) string {
	record := strconv.Itoa(block.Index) +
		strconv.FormatInt(block.Timestamp, 10)

	for _, tx := range block.Transactions {
		record += tx.Sender + tx.Receiver + strconv.Itoa(tx.Amount)
	}

	record += block.PreviousHash
	record += strconv.Itoa(block.Nonce)

	hash := sha256.Sum256([]byte(record))
	return fmt.Sprintf("%x", hash)
}