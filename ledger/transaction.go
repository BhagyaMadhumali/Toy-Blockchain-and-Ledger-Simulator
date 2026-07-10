package ledger

const SystemAccount = "SYSTEM"

// Transaction represents an integer transfer between two accounts.
type Transaction struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   int    `json:"amount"`
}
