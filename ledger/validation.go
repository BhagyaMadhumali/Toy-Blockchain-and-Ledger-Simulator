package ledger

import "errors"

func ValidateTransaction(l *Ledger, tx Transaction) error {
	if tx.Sender == "" || tx.Receiver == "" {
		return errors.New("sender and receiver are required")
	}

	if tx.Sender == tx.Receiver {
		return errors.New("sender and receiver cannot be the same")
	}

	if tx.Amount <= 0 {
		return errors.New("invalid amount")
	}

	if l.GetBalance(tx.Sender) < tx.Amount {
		return errors.New("insufficient balance")
	}

	return nil
}