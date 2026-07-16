package ledger

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const SystemAccount = "SYSTEM"

// Transaction represents a signed integer transfer between two accounts.
// Genesis SYSTEM transactions are trusted allocations and remain unsigned.
type Transaction struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Amount    int    `json:"amount"`
	PublicKey string `json:"public_key,omitempty"`
	Signature string `json:"signature,omitempty"`
}

// SigningBytes returns the exact canonical data covered by the signature.
// PublicKey and Signature are excluded so the signature cannot sign itself.
func (tx Transaction) SigningBytes() ([]byte, error) {
	payload := struct {
		Sender   string `json:"sender"`
		Receiver string `json:"receiver"`
		Amount   int    `json:"amount"`
	}{
		Sender:   tx.Sender,
		Receiver: tx.Receiver,
		Amount:   tx.Amount,
	}
	return json.Marshal(payload)
}

// SignTransaction attaches the public key and Ed25519 signature to a normal
// transaction. The caller must first ensure the private key belongs to Sender.
func SignTransaction(tx *Transaction, privateKey ed25519.PrivateKey) error {
	if tx == nil {
		return fmt.Errorf("transaction is required")
	}
	if tx.Sender == SystemAccount {
		return fmt.Errorf("SYSTEM genesis transactions are not signed")
	}
	if len(privateKey) != ed25519.PrivateKeySize {
		return fmt.Errorf("invalid Ed25519 private key")
	}

	message, err := tx.SigningBytes()
	if err != nil {
		return fmt.Errorf("serialize transaction for signing: %w", err)
	}
	publicKey := privateKey.Public().(ed25519.PublicKey)
	tx.PublicKey = base64.StdEncoding.EncodeToString(publicKey)
	tx.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, message))
	return nil
}

// VerifySignature verifies that the signature covers the transaction fields.
func (tx Transaction) VerifySignature() error {
	if tx.Sender == SystemAccount {
		return fmt.Errorf("SYSTEM transactions are allowed only in the genesis block")
	}
	if tx.PublicKey == "" || tx.Signature == "" {
		return fmt.Errorf("public key and signature are required")
	}

	publicKey, err := base64.StdEncoding.DecodeString(tx.PublicKey)
	if err != nil || len(publicKey) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid Ed25519 public key")
	}
	signature, err := base64.StdEncoding.DecodeString(tx.Signature)
	if err != nil || len(signature) != ed25519.SignatureSize {
		return fmt.Errorf("invalid Ed25519 signature encoding")
	}
	message, err := tx.SigningBytes()
	if err != nil {
		return fmt.Errorf("serialize transaction for verification: %w", err)
	}
	if !ed25519.Verify(ed25519.PublicKey(publicKey), message, signature) {
		return fmt.Errorf("digital signature verification failed")
	}
	return nil
}

// PublicKeyFingerprint returns a short printable identifier for a key.
func (tx Transaction) PublicKeyFingerprint() string {
	publicKey, err := base64.StdEncoding.DecodeString(tx.PublicKey)
	if err != nil {
		return "invalid"
	}
	hash := sha256.Sum256(publicKey)
	return fmt.Sprintf("%x", hash[:6])
}
