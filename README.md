# Toy Blockchain and Ledger Simulator

A simplified blockchain implementation built from scratch using **Go**. This project is being developed as part of a backend engineering learning exercise to understand the core concepts of blockchain technology and Go programming.

---

# Project Status

**Current Progress:** ✅ Day 1 Completed

The project currently implements the foundation of a blockchain system, including the blockchain structure, genesis block creation, deterministic SHA-256 hashing, and blockchain printing.

---

# Features Implemented (Day 1)

- Go project initialization
- Clean project structure
- Block data structure
- Blockchain data structure
- Genesis block creation
- SHA-256 hash generation
- Blockchain printing through the command line

---

# Project Structure

```text
toy-blockchain/
│
├── cmd/
│   └── main.go
│
├── blockchain/
│   ├── block.go
│   ├── blockchain.go
│   └── hash.go
│
├── ledger/
│   └── transaction.go
│
├── storage/
│
├── data/
│
├── go.mod
└── README.md
```

---

# Implemented Components

## Block

A block represents one unit of the blockchain.

Current fields:

- Index
- Timestamp
- Data
- Previous Hash
- Hash

---

## Blockchain

The blockchain maintains an ordered collection of blocks.

Currently, it contains only the Genesis Block.

---

## Genesis Block

The Genesis Block is the first block in the blockchain.

Properties:

- Block Index = 0
- Previous Hash = `"0000000000"`
- Contains the text `"Genesis Block"`
- Hash generated using SHA-256

---

## Hashing

The project uses Go's standard library:

```go
crypto/sha256
```

The hash is generated using:

- Block Index
- Timestamp
- Data
- Previous Hash

The generated SHA-256 hash uniquely identifies each block.

---

# Current Workflow

```
Start Program
      │
      ▼
Create Genesis Block
      │
      ▼
Calculate SHA-256 Hash
      │
      ▼
Store Genesis Block in Blockchain
      │
      ▼
Print Blockchain
```

---

# How to Run

### Clone the project

```bash
git clone <repository-url>
```

### Navigate to the project

```bash
cd toy-blockchain
```

### Run the application

```bash
go run cmd/main.go
```

---

# Example Output

```text
---------------
Index: 0
Timestamp: 1783059283
Data: Genesis Block
Previous Hash: 0000000000
Hash: 2853e4167e9bed843dff7500cfe5025003158af4ae5612ce943cb6a30493e01d
```

---

# Concepts Learned

During Day 1, the following blockchain concepts were implemented and understood:

- Blockchain architecture
- Genesis Block
- Block structure
- SHA-256 hashing
- Deterministic hashing
- Blockchain initialization
- Go packages and project organization

---

# Technologies Used

- Go 1.26.4
- Go Standard Library
  - `crypto/sha256`
  - `time`
  - `fmt`
  - `strconv`

---

# Next Steps (Day 2)

The following features will be implemented next:

- Transaction model
- Pending transaction pool
- Ledger (account balances)
- Transaction validation
- Proof of Work (Mining)
- Mining difficulty
- Adding new blocks to the blockchain

---

# Learning Objective

This project is being developed incrementally to understand blockchain concepts from first principles while improving Go programming skills. Each feature is implemented one step at a time to ensure a clear understanding of the underlying design and logic before moving on to more advanced blockchain functionality.

# Toy Blockchain - Day 2

## Overview

Day 2 extends the basic blockchain created on Day 1 by introducing a transaction system, account ledger, transaction validation, pending transaction pool, and Proof of Work (PoW) mining. The blockchain can now process transactions between users, mine new blocks, and update account balances.

---

## Features Implemented

- Transaction system
- Account ledger with balances
- Transaction validation
- Pending transaction pool
- Proof of Work (PoW)
- Nonce-based mining
- Difficulty-based hash generation
- Automatic ledger updates after mining
- Blockchain growth with newly mined blocks

---

## Project Structure

```
toy-blockchain/
│
├── cmd/
│   └── main.go
│
├── blockchain/
│   ├── block.go
│   ├── blockchain.go
│   ├── hash.go
│   ├── mining.go
│   └── pow.go
│
├── ledger/
│   ├── ledger.go
│   ├── transaction.go
│   └── validation.go
│
├── go.mod
└── README.md
```

---

## Components

### Transaction (`ledger/transaction.go`)

Defines the structure of a transaction.

```go
type Transaction struct {
    Sender   string
    Receiver string
    Amount   int
}
```

Each transaction contains:

- Sender
- Receiver
- Amount

---

### Block (`blockchain/block.go`)

Each block stores:

- Block Index
- Timestamp
- Transactions
- Previous Block Hash
- Nonce
- Current Block Hash

---

### Hashing (`blockchain/hash.go`)

Each block's hash is generated using SHA-256.

The hash is calculated from:

- Block Index
- Timestamp
- Transactions
- Previous Hash
- Nonce

Any modification to block data generates a completely different hash.

---

### Ledger (`ledger/ledger.go`)

The ledger maintains account balances.

Functions:

- Create new ledger
- Add accounts
- Retrieve balances
- Update balances after successful mining

Example:

```
Alice : 100
Bob   : 50
```

---

### Transaction Validation (`ledger/validation.go`)

Transactions are validated before being accepted.

Validation checks:

- Amount must be greater than zero
- Sender must have sufficient balance

Returns an error if validation fails.

---

### Blockchain (`blockchain/blockchain.go`)

The blockchain stores:

- All mined blocks
- Pending transactions
- Mining difficulty
- Ledger

Responsibilities:

- Create Genesis Block
- Add transactions
- Store pending transactions
- Print the blockchain

---

### Proof of Work (`blockchain/pow.go`)

Implements the mining difficulty.

A valid block hash must begin with a predefined number of leading zeros.

Example (Difficulty = 4):

```
0000ab3d8d92...
```

If the generated hash does not satisfy the difficulty requirement, the miner increases the nonce and tries again.

---

### Mining (`blockchain/mining.go`)

Mining performs the following steps:

1. Collect pending transactions.
2. Create a new block.
3. Generate hashes repeatedly.
4. Increment the nonce until a valid hash is found.
5. Store the mined block.
6. Update account balances.
7. Clear pending transactions.

---

## Program Flow

```
Start Program
      │
      ▼
Create Blockchain
      │
      ▼
Create Genesis Block
      │
      ▼
Initialize Ledger
      │
      ▼
Create User Accounts
      │
      ▼
Create Transactions
      │
      ▼
Add to Pending Pool
      │
      ▼
Mine Pending Transactions
      │
      ▼
Proof of Work
      │
      ▼
Generate Valid Hash
      │
      ▼
Update Ledger
      │
      ▼
Append Block to Blockchain
      │
      ▼
Print Blockchain
      │
      ▼
Display Final Balances
```

---

## Expected Output

```
Initial Balances:
map[Alice:100 Bob:50]

---------------
Index: 0
Time: ...
Prev: 0000
Nonce: 0
Hash: ...

---------------
Index: 1
Time: ...
Prev: ...
Nonce: 48291
Hash: 0000....

   Alice -> Bob 20
   Bob -> Alice 10

Final Balances:
map[Alice:90 Bob:60]
```

(The nonce and hashes will vary on each run.)

---

## Learning Outcomes

By completing Day 2, you learned how to:

- Design a transaction model
- Maintain account balances using a ledger
- Validate transactions
- Implement a pending transaction pool
- Understand SHA-256 hashing
- Implement Proof of Work
- Mine new blocks using a nonce
- Link blocks using hashes
- Update the blockchain after mining

---

## How to Run

### 1. Clone the repository

```bash
git clone <repository-url>
```

### 2. Navigate to the project

```bash
cd toy-blockchain
```

### 3. Run the application

```bash
go run ./cmd
```

---

## Day 2 Summary

Day 2 transforms the blockchain from a static chain of blocks into a functional blockchain capable of handling transactions. The implementation introduces account management through a ledger, validates transactions, mines blocks using the Proof of Work consensus mechanism, updates balances after successful mining, and appends new blocks to the chain, providing a strong foundation for future enhancements such as digital signatures, peer-to-peer networking, wallets, and consensus algorithms.

# Toy Blockchain Simulator (Day 3)

## Project Overview

The **Toy Blockchain Simulator** is a Go-based application developed to demonstrate the core concepts of blockchain technology. It simulates how transactions are validated, grouped into blocks, mined using the Proof of Work algorithm, and securely linked together to form a blockchain.

Day 3 extends the project by adding blockchain validation, tampering detection, a command-line interface (CLI), JSON persistence, and unit testing.

---

# Features

### Day 1 Features

- Blockchain data structure
- Genesis block creation
- SHA-256 hashing
- Block structure
- Blockchain printing

### Day 2 Features

- Transaction model
- Pending transaction pool
- Ledger for account balances
- Transaction validation
- Proof of Work (PoW) mining
- Mining difficulty
- Block creation from transactions

### Day 3 Features

- Blockchain validation
- Tampering detection
- Command Line Interface (CLI)
- Save blockchain to JSON
- Load blockchain from JSON
- Automatic blockchain persistence
- Unit testing
- Project documentation

---

# Project Structure

```text
toy-blockchain/
│
├── cmd/
│   └── main.go
│
├── blockchain/
│   ├── block.go
│   ├── blockchain.go
│   ├── hash.go
│   ├── mining.go
│   ├── persistence.go
│   ├── pow.go
│   └── validation.go
│
├── ledger/
│   ├── ledger.go
│   ├── transaction.go
│   └── validation.go
│
├── tests/
│   ├── blockchain_test.go
│   ├── hash_test.go
│   ├── mining_test.go
│   ├── tamper_test.go
│   └── transaction_test.go
│
├── data/
│   └── blockchain.json
│
├── go.mod
└── README.md
```

---

# Technologies Used

- Go (Golang)
- SHA-256 Hashing
- JSON
- Standard Go Libraries

---

# Blockchain Workflow

```
Start Program
      │
      ▼
Load blockchain.json
      │
      ▼
If file not found
Create Genesis Block
      │
      ▼
Wait for CLI Command
      │
 ┌──────────────┬───────────────┬─────────────┐
 │              │               │             │
 ▼              ▼               ▼             ▼
Add          Mine           Validate      Print
 │              │               │             │
 ▼              ▼               ▼             ▼
Pending     Proof of Work    Check Hashes  Display
Pool         Mining          Previous Hash Blockchain
 │              │               │
 ▼              ▼               ▼
Ledger      Save JSON      Blockchain Valid
Update
```

---

# Block Structure

Each block contains:

- Block Index
- Timestamp
- Transactions
- Previous Hash
- Current Hash
- Nonce

---

# Transaction Flow

1. User creates a transaction.
2. Transaction is validated.
3. Valid transactions are stored in the pending transaction pool.
4. Mining collects pending transactions into a new block.
5. Proof of Work generates a valid hash.
6. The mined block is added to the blockchain.
7. Ledger balances are updated.
8. Blockchain is saved as a JSON file.

---

# Proof of Work

Mining repeatedly changes the **Nonce** until the generated SHA-256 hash satisfies the required mining difficulty.

Example:

```
Difficulty = 2

Invalid:
3af129...

Invalid:
91ab23...

Valid:
00fd823ab...
```

---

# Blockchain Validation

The validation process checks:

- Every block hash is correct.
- Previous hash links are valid.
- Proof of Work requirement is satisfied.
- The blockchain has not been modified after mining.

If all checks pass, the blockchain is considered valid.

---

# Tampering Detection

The simulator can detect blockchain tampering.

Examples of tampering:

- Changing transaction amounts
- Editing sender or receiver
- Modifying previous hash
- Editing mined block data

If any block is modified after mining, validation fails because the recalculated hash no longer matches the stored hash.

---

# JSON Persistence

The blockchain can be stored in a JSON file.

Functions:

- SaveBlockchain()
- LoadBlockchain()

Benefits:

- Data remains available after the program exits.
- Blockchain automatically reloads when the application starts.
- If no blockchain exists, a new genesis block is created.

Saved file:

```
data/blockchain.json
```

---

# Command Line Interface (CLI)

## Add Transaction

```bash
go run cmd/main.go add Alice Bob 20
```

## Mine Block

```bash
go run cmd/main.go mine
```

## Print Blockchain

```bash
go run cmd/main.go print
```

## Validate Blockchain

```bash
go run cmd/main.go validate
```

## Show Account Balances

```bash
go run cmd/main.go balance
```

## Save Blockchain

```bash
go run cmd/main.go save
```

## Load Blockchain

```bash
go run cmd/main.go load
```

## Help

```bash
go run cmd/main.go help
```

---

# Unit Testing

The project includes automated tests for the main blockchain components.

## Hash Tests

- Hash generation
- Hash consistency
- Different blocks generate different hashes

## Mining Tests

- Proof of Work
- Difficulty verification
- Nonce increment
- Valid mined hash

## Blockchain Tests

- Genesis block creation
- Transaction addition
- Mining
- Blockchain validation

## Tampering Tests

- Modify transaction amount
- Modify previous hash
- Detect blockchain tampering

## Transaction Tests

- Valid transactions
- Negative amounts
- Zero amounts
- Insufficient balance

Run all tests:

```bash
go test ./...
```

---

# Sample Execution

```
> go run cmd/main.go add Alice Bob 20

Transaction added to pending pool.

> go run cmd/main.go mine

Mining block...
Block mined successfully.

> go run cmd/main.go print

Block 0
Hash: ...

Block 1
Alice -> Bob : 20

> go run cmd/main.go balance

Alice : 80
Bob : 70
Charlie : 75

> go run cmd/main.go validate

Blockchain is valid.

> go run cmd/main.go save

Blockchain saved successfully.
```

---

# Learning Outcomes

After completing this project, the following concepts were implemented and understood:

- Blockchain architecture
- Genesis block creation
- SHA-256 hashing
- Transactions and ledgers
- Pending transaction pools
- Proof of Work mining
- Mining difficulty
- Blockchain validation
- Tampering detection
- JSON data persistence
- Command Line Interface (CLI)
- Unit testing in Go
- Modular software design
- File handling in Go

---

# Future Improvements

Possible enhancements include:

- Mining rewards
- Digital signatures
- Wallet generation
- Peer-to-peer networking
- REST API
- Web dashboard
- Merkle Trees
- Consensus algorithms
- Multi-node blockchain simulation
- Docker deployment

---

# Author

**Toy Blockchain Simulator**

Developed as a learning project to understand the fundamentals of blockchain technology using Go (Golang).
