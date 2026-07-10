package ledger

// Ledger stores account names and balances.
type Ledger struct {
	Balances map[string]int `json:"balances"`
}

// NewLedger creates an empty ledger.
func NewLedger() *Ledger {
	return &Ledger{
		Balances: make(map[string]int),
	}
}

// AddAccount creates an account or replaces
// the balance of an existing account.
func (l *Ledger) AddAccount(user string, balance int) {
	l.Balances[user] = balance
}

// GetBalance returns an account balance.
//
// If the account does not exist, Go returns the
// zero value for int, which is 0.
func (l *Ledger) GetBalance(user string) int {
	return l.Balances[user]
}

// UpdateBalance moves an amount from sender to receiver.
func (l *Ledger) UpdateBalance(
	sender string,
	receiver string,
	amount int,
) {
	l.Balances[sender] -= amount
	l.Balances[receiver] += amount
}

// ApplyTransaction updates balances using one transaction.
func (l *Ledger) ApplyTransaction(tx Transaction) {
	l.UpdateBalance(
		tx.Sender,
		tx.Receiver,
		tx.Amount,
	)
}