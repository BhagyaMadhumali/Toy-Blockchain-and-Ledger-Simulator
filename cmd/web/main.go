package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sync"

	"toy-blockchain/blockchain"
	"toy-blockchain/ledger"
)

const dataFile = "data/blockchain.json"

var (
	bc  *blockchain.Blockchain
	mux sync.Mutex
)

type PageData struct {
	Message  string
	Valid    bool
	Chain    *blockchain.Blockchain
	Balances map[string]int
}

func main() {
	var err error

	bc, err = blockchain.LoadBlockchain(dataFile)
	if err != nil {
		fmt.Println("Error loading blockchain:", err)
		return
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/transfer", transferHandler)
	http.HandleFunc("/mine", mineHandler)
	http.HandleFunc("/tamper", tamperHandler)
	http.HandleFunc("/validate", validateHandler)
	http.HandleFunc("/reset", resetHandler)

	fmt.Println("UI running at: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "Welcome to Toy Blockchain UI")
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	mux.Lock()
	defer mux.Unlock()

	sender := r.FormValue("sender")
	receiver := r.FormValue("receiver")

	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		renderPage(w, "Invalid amount")
		return
	}

	tx := ledger.Transaction{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}

	if err := bc.AddTransaction(tx); err != nil {
		renderPage(w, "Transaction rejected: "+err.Error())
		return
	}

	bc.SaveBlockchain(dataFile)

	renderPage(w, "Transaction added to pending pool. Mine it to complete transfer.")
}

func mineHandler(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()

	if err := bc.MinePendingTransactions(); err != nil {
		renderPage(w, err.Error())
		return
	}

	bc.SaveBlockchain(dataFile)

	renderPage(w, "Block mined successfully. Money transfer completed.")
}

func tamperHandler(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()

	if len(bc.Blocks) < 2 {
		renderPage(w, "No mined block available to tamper.")
		return
	}

	if len(bc.Blocks[1].Transactions) == 0 {
		renderPage(w, "Block has no transaction to tamper.")
		return
	}

	bc.Blocks[1].Transactions[0].Amount = 999
	bc.SaveBlockchain(dataFile)

	renderPage(w, "Blockchain has been tampered. Amount changed to 999.")
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if bc.ValidateBlockchain() {
		renderPage(w, "Blockchain is valid.")
	} else {
		renderPage(w, "Blockchain is invalid or has been tampered.")
	}
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()

	bc = blockchain.NewBlockchain()
	bc.SaveBlockchain(dataFile)

	renderPage(w, "Blockchain reset successfully.")
}

func renderPage(w http.ResponseWriter, message string) {
	data := PageData{
		Message:  message,
		Valid:    bc.ValidateBlockchain(),
		Chain:    bc,
		Balances: bc.Ledger.Balances,
	}

	tmpl := template.Must(template.New("page").Parse(pageHTML))
	tmpl.Execute(w, data)
}

const pageHTML = `
<!DOCTYPE html>
<html>
<head>
	<title>Toy Blockchain UI</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			background: #f4f6f8;
			margin: 0;
			padding: 20px;
		}

		.container {
			max-width: 1100px;
			margin: auto;
		}

		h1 {
			text-align: center;
			color: #222;
		}

		.card {
			background: white;
			padding: 20px;
			margin-bottom: 20px;
			border-radius: 12px;
			box-shadow: 0 2px 8px rgba(0,0,0,0.1);
		}

		input, select, button {
			padding: 10px;
			margin: 5px;
			border-radius: 6px;
			border: 1px solid #ccc;
		}

		button {
			background: #2563eb;
			color: white;
			border: none;
			cursor: pointer;
		}

		button:hover {
			background: #1d4ed8;
		}

		.danger {
			background: #dc2626;
		}

		.danger:hover {
			background: #b91c1c;
		}

		.success {
			color: green;
			font-weight: bold;
		}

		.error {
			color: red;
			font-weight: bold;
		}

		table {
			width: 100%;
			border-collapse: collapse;
		}

		th, td {
			padding: 10px;
			border-bottom: 1px solid #ddd;
			text-align: left;
		}

		.block {
			background: #f9fafb;
			padding: 15px;
			margin-bottom: 15px;
			border-left: 5px solid #2563eb;
			border-radius: 8px;
			word-break: break-all;
		}

		.hash {
			font-size: 13px;
			color: #555;
		}

		.actions a {
			text-decoration: none;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>Toy Blockchain Money Transfer UI</h1>

		<div class="card">
			<p><strong>Message:</strong> {{.Message}}</p>

			{{if .Valid}}
				<p class="success">Status: Blockchain Valid</p>
			{{else}}
				<p class="error">Status: Blockchain Invalid / Tampered</p>
			{{end}}
		</div>

		<div class="card">
			<h2>Transfer Money</h2>

			<form action="/transfer" method="POST">
				<select name="sender">
					<option value="Alice">Alice</option>
					<option value="Bob">Bob</option>
					<option value="Charlie">Charlie</option>
				</select>

				<select name="receiver">
					<option value="Alice">Alice</option>
					<option value="Bob">Bob</option>
					<option value="Charlie">Charlie</option>
				</select>

				<input type="number" name="amount" placeholder="Amount" required>

				<button type="submit">Add Transaction</button>
			</form>

			<p>Note: After adding transaction, click <strong>Mine Block</strong> to complete transfer.</p>
		</div>

		<div class="card actions">
			<h2>Actions</h2>

			<a href="/mine"><button>Mine Block</button></a>
			<a href="/validate"><button>Validate Blockchain</button></a>
			<a href="/tamper"><button class="danger">Tamper Blockchain</button></a>
			<a href="/reset"><button class="danger">Reset Blockchain</button></a>
		</div>

		<div class="card">
			<h2>Account Balances</h2>

			<table>
				<tr>
					<th>User</th>
					<th>Balance</th>
				</tr>

				{{range $user, $balance := .Balances}}
				<tr>
					<td>{{$user}}</td>
					<td>{{$balance}}</td>
				</tr>
				{{end}}
			</table>
		</div>

		<div class="card">
			<h2>Pending Transactions</h2>

			{{if .Chain.PendingTransactions}}
				<table>
					<tr>
						<th>Sender</th>
						<th>Receiver</th>
						<th>Amount</th>
					</tr>

					{{range .Chain.PendingTransactions}}
					<tr>
						<td>{{.Sender}}</td>
						<td>{{.Receiver}}</td>
						<td>{{.Amount}}</td>
					</tr>
					{{end}}
				</table>
			{{else}}
				<p>No pending transactions.</p>
			{{end}}
		</div>

		<div class="card">
			<h2>Blockchain</h2>

			{{range .Chain.Blocks}}
			<div class="block">
				<h3>Block {{.Index}}</h3>
				<p><strong>Timestamp:</strong> {{.Timestamp}}</p>
				<p><strong>Nonce:</strong> {{.Nonce}}</p>
				<p class="hash"><strong>Previous Hash:</strong> {{.PreviousHash}}</p>
				<p class="hash"><strong>Hash:</strong> {{.Hash}}</p>

				<h4>Transactions</h4>

				{{if .Transactions}}
					<ul>
					{{range .Transactions}}
						<li>{{.Sender}} → {{.Receiver}} : {{.Amount}}</li>
					{{end}}
					</ul>
				{{else}}
					<p>No transactions</p>
				{{end}}
			</div>
			{{end}}
		</div>
	</div>
</body>
</html>
`