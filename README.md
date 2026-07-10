# Toy Blockchain and Ledger Simulator

A command-line blockchain simulator written in Go. The blockchain is the only source of truth: account balances are derived by replaying validated transactions from the fixed genesis block through the latest block.

## Features

- Integer-only transactions
- Fixed deterministic genesis block
- Genesis coin allocations recorded as `SYSTEM` transactions
- Pending-pool double-spend prevention
- Transaction revalidation before mining
- SHA-256 hashing using structured JSON serialization
- Proof-of-work with trusted consensus difficulty
- Nonce, hash-attempt and elapsed-time reporting
- Validation of hashes, links, proof-of-work, indexes, timestamps and balances
- First-offending-block validation errors
- JSON persistence without saving derived ledger balances
- Difficulty benchmark command for research experiments
- Unit tests for hashing, mining, tampering, persistence and transaction rules

## Project structure

```text
toy-blockchain/
├── blockchain/
│   ├── block.go
│   ├── blockchain.go
│   ├── hash.go
│   ├── mining.go
│   ├── persistence.go
│   ├── pow.go
│   └── validation.go
├── cmd/
│   └── main.go
├── ledger/
│   ├── ledger.go
│   ├── transaction.go
│   └── validation.go
├── tests/
├── research-report.md
├── .gitattributes
├── .gitignore
├── go.mod
└── README.md
```

## Source-of-truth design

Initial funds are created in the deterministic genesis block:

- Alice: 100
- Bob: 50
- Charlie: 75

They are represented as transactions from the reserved `SYSTEM` account. The ledger is marked `json:"-"` and is never stored in `data/blockchain.json`. On every load, the application validates and replays the chain to rebuild balances. Editing or injecting a separate balance map therefore cannot change account balances.

## Hashing scheme

The hash includes these fields in this order through JSON serialization:

1. Index
2. Timestamp
3. Transactions
4. Previous hash
5. Difficulty
6. Nonce

Structured serialization avoids collisions caused by direct string concatenation, such as index `1` with timestamp `23` and index `12` with timestamp `3` both becoming `"123"`.

## Validation rules

Validation stops at and reports the first offending block. It checks:

- The genesis block exactly matches fixed trusted constants
- Block indexes are sequential
- Timestamps never move backward
- Stored hashes match recalculated hashes
- Previous-hash links are correct
- Every block uses the trusted difficulty
- Every hash satisfies proof-of-work
- `SYSTEM` issuance occurs only in the genesis block
- Transaction fields and amounts are valid
- Replay never creates a negative balance

Pending transactions are also replayed sequentially when loading and immediately before mining.

## Commands

Run commands from the project root.

```bash
go run ./cmd help
go run ./cmd balance
go run ./cmd add Alice Bob 20
go run ./cmd pending
go run ./cmd mine
go run ./cmd print
go run ./cmd validate
go run ./cmd benchmark --min 1 --max 5 --runs 3
go run ./cmd reset
```

`add` and `mine` automatically save state. Manual `save` and `load` commands are intentionally omitted because loading happens at startup and state-changing commands save automatically.

## Example mining output

```text
Block mined successfully.
Difficulty: 3
Nonce: 4321
Hashes attempted: 4322
Elapsed time: 8.42ms
Hash: 000...
```

Actual results vary because proof-of-work is probabilistic.

## Running tests and checks

```bash
gofmt -w .
go test ./...
go vet ./...
```

## Runtime data

`data/blockchain.json` is generated during use and excluded from Git. Delete it manually or run `go run ./cmd reset` to start again from the fixed genesis block.

## Research report

See [`research-report.md`](research-report.md). Run the benchmark command and replace the marked example table cells with measurements from the current machine before final submission.
