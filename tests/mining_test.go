package tests

import (
	"strings"
	"testing"
	"toy-blockchain/blockchain"
)

func TestMiningCreatesValidHash(t *testing.T) {
	block := blockchain.Block{
		Index:        1,
		Timestamp:    1000,
		PreviousHash: "abc",
		Nonce:        0,
	}

	difficulty := 2

	blockchain.MineBlock(
		&block,
		difficulty,
	)

	requiredPrefix := strings.Repeat(
		"0",
		difficulty,
	)

	if !strings.HasPrefix(
		block.Hash,
		requiredPrefix,
	) {
		t.Errorf(
			"mined hash does not match difficulty",
		)
	}
}

func TestMiningChangesNonce(t *testing.T) {
	block := blockchain.Block{
		Index:        1,
		Timestamp:    1000,
		PreviousHash: "abc",
		Nonce:        0,
	}

	blockchain.MineBlock(
		&block,
		2,
	)

	if block.Nonce == 0 {
		t.Errorf(
			"expected nonce to change during mining",
		)
	}
}