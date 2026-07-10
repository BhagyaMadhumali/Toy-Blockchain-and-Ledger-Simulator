package ledger

import "fmt"

// Ledger is a derived view of balances rebuilt from blockchain transactions.
type Ledger struct {
	Balances map[string]int `json:"-"`
}

func NewLedger() *Ledger {
	return &Ledger{Balances: make(map[string]int)}
}

func (l *Ledger) GetBalance(account string) int {
	return l.Balances[account]
}

func (l *Ledger) Credit(account string, amount int) error {
	if account == "" {
		return fmt.Errorf("receiver is required")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	l.Balances[account] += amount
	return nil
}

// ApplyTransaction validates and applies a normal account-to-account transfer.
func (l *Ledger) ApplyTransaction(tx Transaction) error {
	if err := ValidateTransaction(l, tx); err != nil {
		return err
	}
	l.Balances[tx.Sender] -= tx.Amount
	l.Balances[tx.Receiver] += tx.Amount
	return nil
}

func (l *Ledger) Clone() *Ledger {
	clone := NewLedger()
	for account, balance := range l.Balances {
		clone.Balances[account] = balance
	}
	return clone
}
