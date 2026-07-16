package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"toy-blockchain/ledger"
)

// CalculateMerkleRoot summarizes all transactions in a block.
// Each leaf is SHA-256(canonical transaction JSON). Parent nodes are
// SHA-256(left child bytes || right child bytes). When a level has an odd
// number of nodes, its final node is duplicated.
func CalculateMerkleRoot(transactions []ledger.Transaction) string {
	if len(transactions) == 0 {
		empty := sha256.Sum256(nil)
		return hex.EncodeToString(empty[:])
	}

	level := make([][]byte, 0, len(transactions))
	for _, tx := range transactions {
		encoded, err := json.Marshal(tx)
		if err != nil {
			panic("transaction serialization for Merkle tree failed: " + err.Error())
		}
		hash := sha256.Sum256(encoded)
		node := make([]byte, len(hash))
		copy(node, hash[:])
		level = append(level, node)
	}

	for len(level) > 1 {
		if len(level)%2 != 0 {
			duplicate := append([]byte(nil), level[len(level)-1]...)
			level = append(level, duplicate)
		}

		nextLevel := make([][]byte, 0, len(level)/2)
		for i := 0; i < len(level); i += 2 {
			combined := make([]byte, 0, len(level[i])+len(level[i+1]))
			combined = append(combined, level[i]...)
			combined = append(combined, level[i+1]...)
			parent := sha256.Sum256(combined)
			node := make([]byte, len(parent))
			copy(node, parent[:])
			nextLevel = append(nextLevel, node)
		}
		level = nextLevel
	}

	return hex.EncodeToString(level[0])
}
