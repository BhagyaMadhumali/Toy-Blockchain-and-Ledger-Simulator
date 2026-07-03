package blockchain

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

func CalculateHash(block Block) string {

	record :=
		strconv.Itoa(block.Index) +
			strconv.FormatInt(block.Timestamp, 10) +
			block.Data +
			block.PreviousHash

	hash := sha256.Sum256([]byte(record))
	return fmt.Sprintf("%x", hash)
}