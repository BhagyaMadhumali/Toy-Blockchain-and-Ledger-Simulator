package blockchain

const (
	// TargetBlockTime is the desired time between mined blocks.
	TargetBlockTime = int64(30) // seconds

	// Difficulty changes by at most one level for each new block.
	MinDifficulty = 1
	MaxDifficulty = 5
)

// CalculateNextDifficulty deterministically calculates the difficulty for the
// next block from the existing chain's timestamps.
//
// The first two blocks after genesis keep the previous difficulty because a
// normal-to-normal block interval is not available until then. Afterwards:
//   - faster than half the target: increase by one
//   - slower than twice the target: decrease by one
//   - otherwise: keep the current difficulty
func CalculateNextDifficulty(blocks []Block) int {
	if len(blocks) == 0 {
		return DefaultDifficulty
	}

	previousDifficulty := clampDifficulty(blocks[len(blocks)-1].Difficulty)
	if len(blocks) < 3 {
		return previousDifficulty
	}

	previous := blocks[len(blocks)-1]
	beforePrevious := blocks[len(blocks)-2]
	actualBlockTime := previous.Timestamp - beforePrevious.Timestamp

	// A zero or negative interval is treated as extremely fast. A negative
	// timestamp will still be rejected separately by blockchain validation.
	if actualBlockTime < TargetBlockTime/2 {
		return clampDifficulty(previousDifficulty + 1)
	}
	if actualBlockTime > TargetBlockTime*2 {
		return clampDifficulty(previousDifficulty - 1)
	}
	return previousDifficulty
}

func clampDifficulty(difficulty int) int {
	if difficulty < MinDifficulty {
		return MinDifficulty
	}
	if difficulty > MaxDifficulty {
		return MaxDifficulty
	}
	return difficulty
}
