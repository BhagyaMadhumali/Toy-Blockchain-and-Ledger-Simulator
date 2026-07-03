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
