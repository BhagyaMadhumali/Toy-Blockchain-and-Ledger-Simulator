package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// CalculateHash hashes the block header. Transactions are represented only by
// their Merkle root rather than by serializing the complete transaction list.
func CalculateHash(block Block) string {
	input := struct {
		Index        int    `json:"index"`
		Timestamp    int64  `json:"timestamp"`
		MerkleRoot   string `json:"merkle_root"`
		PreviousHash string `json:"previous_hash"`
		Difficulty   int    `json:"difficulty"`
		Nonce        uint64 `json:"nonce"`
	}{
		Index:        block.Index,
		Timestamp:    block.Timestamp,
		MerkleRoot:   block.MerkleRoot,
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
