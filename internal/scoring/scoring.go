package scoring

import "time"

// Calculate computes the score for a debugging session.
//
// Formula:
//   score = 100
//   - deduct 10 per hint used
//   - deduct 5 per unnecessary file changed (beyond expected)
//   - deduct 1 per minute over the par time
//
// Score is clamped to a minimum of 0.
func Calculate(timeTaken time.Duration, parTime time.Duration, hintsUsed int, filesChanged int, expectedFiles int) int {
	score := 100

	// Deduct for hints
	score -= 10 * hintsUsed

	// Deduct for unnecessary file changes
	extraFiles := filesChanged - expectedFiles
	if extraFiles > 0 {
		score -= 5 * extraFiles
	}

	// Deduct for overtime
	overtimeMinutes := int(timeTaken.Minutes()) - int(parTime.Minutes())
	if overtimeMinutes > 0 {
		score -= overtimeMinutes
	}

	// Clamp to 0
	if score < 0 {
		score = 0
	}

	return score
}
