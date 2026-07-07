package ledger

import "errors"

func ValidateTransaction(l *Ledger, tx Transaction) error {
	if tx.Amount <= 0 {
		return errors.New("invalid amount")
	}

	if l.GetBalance(tx.Sender) < tx.Amount {
		return errors.New("insufficient balance")
	}

	return nil
}