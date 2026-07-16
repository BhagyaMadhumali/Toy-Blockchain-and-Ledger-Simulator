package ledger

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// DemoAccountPrivateKey creates stable keys for the built-in educational
// accounts. This keeps the CLI easy to use and the genesis block reproducible.
// Never use deterministic, source-code-derived private keys in a real system.
func DemoAccountPrivateKey(account string) (ed25519.PrivateKey, error) {
	switch account {
	case "Alice", "Bob", "Charlie":
		seed := sha256.Sum256([]byte("toy-blockchain-demo-key:" + account))
		return ed25519.NewKeyFromSeed(seed[:]), nil
	default:
		return nil, fmt.Errorf("unknown signing account %q; available demo accounts: Alice, Bob, Charlie", account)
	}
}

func DemoAccountPublicKey(account string) (string, error) {
	privateKey, err := DemoAccountPrivateKey(account)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public().(ed25519.PublicKey)
	return base64.StdEncoding.EncodeToString(publicKey), nil
}

// VerifySenderIdentity prevents a valid key owned by one account from being
// used to spend another built-in account's balance.
func VerifySenderIdentity(tx Transaction) error {
	expectedPublicKey, err := DemoAccountPublicKey(tx.Sender)
	if err != nil {
		return err
	}
	if tx.PublicKey != expectedPublicKey {
		return fmt.Errorf("public key does not belong to sender %s", tx.Sender)
	}
	return nil
}
