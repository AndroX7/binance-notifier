package utility

import "github.com/AndroX7/binance-notifier/candle"

// Imbalance zone representation
type Imbalance struct {
	StartIndex int
	EndIndex   int
	Top        float64
	Bottom     float64
	Bullish    bool
}

// Detect imbalances (similar to FVG but more flexible)
func GetImbalance(candles []candle.Candle) []Imbalance {
	var imbalances []Imbalance
	for i := 2; i < len(candles); i++ {
		c1 := candles[i-2]
		c2 := candles[i-1]
		// Bullish imbalance (gap up)
		if c2.Low > c1.High {
			imbalances = append(imbalances, Imbalance{
				StartIndex: i - 2,
				EndIndex:   i - 1,
				Top:        c2.Low,
				Bottom:     c1.High,
				Bullish:    true,
			})
		}
		// Bearish imbalance (gap down)
		if c2.High < c1.Low {
			imbalances = append(imbalances, Imbalance{
				StartIndex: i - 2,
				EndIndex:   i - 1,
				Top:        c1.Low,
				Bottom:     c2.High,
				Bullish:    false,
			})
		}
	}
	return imbalances
}
