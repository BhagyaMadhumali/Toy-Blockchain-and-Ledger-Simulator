package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// CalculateHash serializes every consensus field as JSON before hashing.
// Structured serialization prevents ambiguous field concatenation.
func CalculateHash(block Block) string {
	input := struct {
		Index        int         `json:"index"`
		Timestamp    int64       `json:"timestamp"`
		Transactions interface{} `json:"transactions"`
		PreviousHash string      `json:"previous_hash"`
		Difficulty   int         `json:"difficulty"`
		Nonce        uint64      `json:"nonce"`
	}{
		Index:        block.Index,
		Timestamp:    block.Timestamp,
		Transactions: block.Transactions,
		PreviousHash: block.PreviousHash,
		Difficulty:   block.Difficulty,
		Nonce:        block.Nonce,
	}

	encoded, err := json.Marshal(input)
	if err != nil {
		panic("block hash serialization failed: " + err.Error())
	}

	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}
