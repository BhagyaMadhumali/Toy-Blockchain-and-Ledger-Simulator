package ledger

type Ledger struct {
	Balances map[string]int `json:"balances"`
}

func NewLedger() *Ledger {
	return &Ledger{
		Balances: make(map[string]int),
	}
}

func (l *Ledger) AddAccount(user string, balance int) {
	l.Balances[user] = balance
}

func (l *Ledger) GetBalance(user string) int {
	return l.Balances[user]
}

func (l *Ledger) UpdateBalance(sender, receiver string, amount int) {
	l.Balances[sender] -= amount
	l.Balances[receiver] += amount
}

func (l *Ledger) ApplyTransaction(tx Transaction) {
	l.UpdateBalance(tx.Sender, tx.Receiver, tx.Amount)
}