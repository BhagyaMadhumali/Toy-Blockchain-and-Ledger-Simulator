package ledger

import "fmt"

func ValidateTransaction(l *Ledger, tx Transaction) error {
	return ValidateTransactionWithBalance(tx, l.GetBalance(tx.Sender))
}

func ValidateTransactionWithBalance(tx Transaction, availableBalance int) error {
	if tx.Sender == "" || tx.Receiver == "" {
		return fmt.Errorf("sender and receiver are required")
	}
	if tx.Sender == SystemAccount {
		return fmt.Errorf("SYSTEM transactions are allowed only in the genesis block")
	}
	if tx.Receiver == SystemAccount {
		return fmt.Errorf("SYSTEM cannot receive normal transactions")
	}
	if tx.Sender == tx.Receiver {
		return fmt.Errorf("sender and receiver cannot be the same")
	}
	if tx.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if availableBalance < tx.Amount {
		return fmt.Errorf("insufficient available balance: have %d, need %d", availableBalance, tx.Amount)
	}
	return nil
}
