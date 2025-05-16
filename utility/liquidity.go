package utility

// Find liquidity pools - recent swing highs and lows (areas likely to attract stops)
func FindLiquidityPools(swings []SwingPoint, lookback int) (highPools, lowPools []float64) {
	n := len(swings)
	if n == 0 {
		return nil, nil
	}

	start := n - lookback
	if start < 0 {
		start = 0
	}

	for i := start; i < n; i++ {
		if swings[i].High {
			highPools = append(highPools, swings[i].Price)
		} else {
			lowPools = append(lowPools, swings[i].Price)
		}
	}

	return highPools, lowPools
}
