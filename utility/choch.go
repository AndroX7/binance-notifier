package utility

// Change Of Character (ChoCH) detection (simplified)
// Detect if the last swing breaks structure and changes trend direction
func DetectChoCH(swings []SwingPoint) (bool, bool) {
	if len(swings) < 4 {
		return false, false
	}

	// last 4 swings for pattern recognition
	last := swings[len(swings)-1]
	prev := swings[len(swings)-2]
	prev3 := swings[len(swings)-4]

	bullishChoCH := false
	bearishChoCH := false

	// Bullish ChoCH example: last two lows are higher than previous lows
	if !prev3.High && !prev.Low && last.Low && prev.Price > prev3.Price && last.Price > prev.Price {
		bullishChoCH = true
	}

	// Bearish ChoCH example: last two highs lower than previous highs
	if prev3.High && prev.High && last.High && prev.Price < prev3.Price && last.Price < prev.Price {
		bearishChoCH = true
	}

	return bullishChoCH, bearishChoCH
}
