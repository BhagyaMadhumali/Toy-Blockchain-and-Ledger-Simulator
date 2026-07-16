package blockchain

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type MiningResult struct {
	Attempts uint64
	Elapsed  time.Duration
	Workers  int
}

type miningWinner struct {
	nonce uint64
	hash  string
}

func HasValidProof(hash string, difficulty int) bool {
	if difficulty < 1 || difficulty > 64 {
		return false
	}
	return strings.HasPrefix(hash, strings.Repeat("0", difficulty))
}

// MineBlock searches the nonce space concurrently using one worker per logical
// CPU. Each worker checks a separate nonce sequence and all workers stop when
// the first valid proof is found.
func MineBlock(block *Block, difficulty int) (MiningResult, error) {
	return MineBlockWithWorkers(block, difficulty, runtime.NumCPU())
}

// MineBlockWithWorkers is the configurable concurrent miner. Worker 0 checks
// 0, workers, 2*workers...; worker 1 checks 1, workers+1... and so on.
// A value of 1 provides deterministic sequential mining, which is used for the
// fixed trusted genesis block.
func MineBlockWithWorkers(block *Block, difficulty, workers int) (MiningResult, error) {
	if difficulty < 1 || difficulty > 6 {
		return MiningResult{}, fmt.Errorf("difficulty must be between 1 and 6")
	}
	if workers < 1 {
		return MiningResult{}, fmt.Errorf("workers must be at least 1")
	}

	block.Difficulty = difficulty
	block.MerkleRoot = CalculateMerkleRoot(block.Transactions)
	block.Nonce = 0
	block.Hash = ""

	baseBlock := *block
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	winnerChannel := make(chan miningWinner, 1)
	var waitGroup sync.WaitGroup
	var attempts atomic.Uint64

	for workerID := 0; workerID < workers; workerID++ {
		waitGroup.Add(1)

		go func(startNonce uint64) {
			defer waitGroup.Done()

			candidate := baseBlock
			step := uint64(workers)

			for nonce := startNonce; ; nonce += step {
				select {
				case <-ctx.Done():
					return
				default:
				}

				candidate.Nonce = nonce
				hash := CalculateHash(candidate)
				attempts.Add(1)

				if HasValidProof(hash, difficulty) {
					select {
					case winnerChannel <- miningWinner{nonce: nonce, hash: hash}:
						cancel()
					default:
					}
					return
				}

				// Prevent uint64 wraparound from causing an endless repeated search.
				if nonce > ^uint64(0)-step {
					return
				}
			}
		}(uint64(workerID))
	}

	workersDone := make(chan struct{})
	go func() {
		waitGroup.Wait()
		close(workersDone)
	}()

	var winner miningWinner
	select {
	case winner = <-winnerChannel:
		<-workersDone
	case <-workersDone:
		return MiningResult{}, fmt.Errorf("nonce space exhausted without finding a valid proof")
	}

	block.Nonce = winner.nonce
	block.Hash = winner.hash

	return MiningResult{
		Attempts: attempts.Load(),
		Elapsed:  time.Since(start),
		Workers:  workers,
	}, nil
}
