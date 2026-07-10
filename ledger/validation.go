package ledger

import "errors"

// ValidateTransaction validates a transaction
// using the current ledger balance.
func ValidateTransaction(
	l *Ledger,
	tx Transaction,
) error {
	currentBalance := l.GetBalance(tx.Sender)

	return ValidateTransactionWithBalance(
		tx,
		currentBalance,
	)
}

// ValidateTransactionWithBalance validates a transaction
// using the available balance provided by the blockchain.
func ValidateTransactionWithBalance(
	tx Transaction,
	availableBalance int,
) error {
	if tx.Sender == "" || tx.Receiver == "" {
		return errors.New(
			"sender and receiver are required",
		)
	}

	if tx.Sender == tx.Receiver {
		return errors.New(
			"sender and receiver cannot be the same",
		)
	}

	if tx.Amount <= 0 {
		return errors.New(
			"amount must be greater than zero",
		)
	}

	if availableBalance < tx.Amount {
		return errors.New(
			"insufficient available balance",
		)
	}

	return nil
}