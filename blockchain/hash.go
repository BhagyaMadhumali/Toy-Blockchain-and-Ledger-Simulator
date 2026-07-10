package blockchain

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

// CalculateHash calculates the SHA-256 hash of a block.
func CalculateHash(block Block) string {
	record := strconv.Itoa(block.Index)
	record += strconv.FormatInt(block.Timestamp, 10)

	for _, tx := range block.Transactions {
		record += tx.Sender
		record += tx.Receiver
		record += strconv.Itoa(tx.Amount)
	}

	record += block.PreviousHash
	record += strconv.Itoa(block.Nonce)

	hash := sha256.Sum256([]byte(record))

	return fmt.Sprintf("%x", hash)
}