package blockchain

import (
	"strconv"
	"strings"
)

func (bc *Blockchain) isValidHash(hash string) bool {
	prefix := strings.Repeat("0", bc.Difficulty)
	return hash[:bc.Difficulty] == prefix
}