# Toy Blockchain Research Report

**Student:** Bhagya Madhumali  
**Project:** Toy Blockchain and Ledger Simulator  
**Language:** Go

## 1. Introduction

This report evaluates the integrity and proof-of-work behaviour of my toy blockchain implementation. The experiments were run against the implementation itself rather than using estimated values. The blockchain records initial currency issuance in a fixed genesis block, derives balances by replaying transactions, links blocks through previous hashes, and requires proof-of-work for every block.

> Before submission, run the commands shown below on the final project and replace the table placeholders with the exact output from that run.

## 2. Tamper-evidence experiment

### 2.1 Procedure

I first created and mined a valid transaction:

```bash
go run ./cmd reset
go run ./cmd add Alice Bob 20
go run ./cmd mine
go run ./cmd validate
```

### 2.2 Output before tampering

```text
Blockchain is valid.
```

I then opened `data/blockchain.json` and changed the first normal transaction amount from `20` to `999` without updating the block hash.

### 2.3 Output after tampering

```text
Error loading blockchain: saved blockchain is invalid: block 1: stored hash does not match block contents
```

### 2.4 Why validation failed

A block hash commits to its index, timestamp, transactions, previous hash, difficulty and nonce. Changing the amount changes the serialized block data. Recalculating SHA-256 therefore produces a different value from the hash stored in the block. Validation detects the first invalid block at the hash-integrity check.

A more advanced attacker could recalculate or re-mine the edited block. However, the next block would still contain the old hash in its `previous_hash` field, breaking the chain link. Re-mining all later blocks would also be required. In addition, replay validation rejects transactions that spend more than the sender owns, so changing the amount to an unaffordable value is detected even after re-mining.

## 3. Difficulty versus mining effort

### 3.1 Method

I used the independent benchmark command, which does not modify the real blockchain:

```bash
go run ./cmd benchmark --min 1 --max 5 --runs 3
```

For each difficulty, the miner repeatedly changed the nonce and calculated a SHA-256 hash until the hash began with the required number of hexadecimal zeroes. I recorded the attempts and elapsed time.

### 3.2 Results

Fill this table using the averages from the actual benchmark output.

| Difficulty | Average attempts | Average elapsed time (ms) |
| ---------: | ---------------: | ------------------------: |
|          1 |          REPLACE |                   REPLACE |
|          2 |          REPLACE |                   REPLACE |
|          3 |          REPLACE |                   REPLACE |
|          4 |          REPLACE |                   REPLACE |
|          5 |          REPLACE |                   REPLACE |

### 3.3 Discussion

The growth is not linear. Each hexadecimal hash digit has 16 possible values, and only one value is zero. Therefore, requiring one additional leading zero multiplies the expected work by approximately 16. The expected attempts are roughly `16^difficulty`, although individual runs can be much lower or higher because mining is probabilistic.

Elapsed time generally follows the attempt count because each attempt performs the same serialization and SHA-256 calculation. Small timing variations can occur because of CPU scheduling and other programs running on the computer.

## 4. Hashing design

The hash input is serialized as JSON and contains:

1. Block index
2. Timestamp
3. Ordered transaction list
4. Previous block hash
5. Difficulty
6. Nonce

I used structured JSON serialization instead of joining values directly. Direct concatenation can be ambiguous: index `1` and timestamp `23` produce the same joined text as index `12` and timestamp `3`. JSON preserves field boundaries and transaction structure, so those blocks produce different hash inputs.

The hash field itself is excluded from the hash input because it is the output being calculated. Any change to a committed field changes the calculated SHA-256 digest.

## 5. How validation guarantees chain integrity

Validation begins with a trusted, deterministic genesis block. Its timestamp, previous hash, allocations and difficulty are fixed in the program. This prevents a fresh run from creating a different blockchain and prevents an edited genesis allocation from becoming trusted.

For every block, validation checks sequential indexes, non-decreasing timestamps, the recalculated hash, the expected consensus difficulty, and proof-of-work. For every block after genesis, it checks that `previous_hash` equals the actual hash of the preceding block.

Validation also rebuilds a new ledger from zero by replaying all transactions in order. Genesis `SYSTEM` transactions create the initial balances. Normal transactions are rejected if names are missing, the sender and receiver match, the amount is non-positive, a reserved system account is used, or the sender lacks funds. This ensures the chain cannot validate if replay would create a negative balance.

The saved ledger is not trusted or persisted. Consequently, manually adding a balance section to JSON has no effect. Balances shown by the CLI always come from replaying the validated chain.

## 6. Discussion questions

### How does the previous-hash link make old-block tampering difficult?

Changing an old block changes its hash. The next block still points to the original hash, so validation detects a broken link. To hide the modification, an attacker must update and re-mine that block and every later block. As the chain grows and difficulty increases, the required work grows substantially.

### Why is proof-of-work useful for tamper evidence?

Without proof-of-work, an attacker could recalculate hashes very quickly after editing blocks. Proof-of-work requires repeated hashing until a rare prefix is found. Rewriting history therefore requires repeating the computational work for the changed block and all following blocks.

### Why must balances be rebuilt from the blockchain?

If a separate balance map is trusted, someone can edit it without changing any block. Replaying transactions makes the blockchain the only source of truth. Every displayed balance then has a traceable origin in genesis issuance and later transfers.

### Why are pending transactions validated again during mining?

Persisted pending data can be manually edited, and earlier pending transactions affect the funds available to later ones. Replaying them sequentially against a fresh chain-derived ledger prevents injected overspending and pending-pool double-spends from being mined.

### What are the limitations of this implementation?

This is an educational single-node blockchain. It has no networking, digital signatures, public/private keys, distributed consensus, transaction fees, mining rewards, fork handling or peer-to-peer synchronization. The trusted difficulty is fixed in code, and JSON files do not provide concurrent multi-process storage safety.

## 7. Conclusion

The experiments demonstrate that the blockchain detects direct transaction tampering and that mining effort grows exponentially in expectation as difficulty increases. The final design strengthens integrity by using a deterministic genesis block, unambiguous hashing, trusted proof-of-work rules, previous-hash links and complete transaction replay. Most importantly, account balances are derived entirely from the chain rather than trusted from editable runtime state.
