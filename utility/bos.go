package utility

import (
	"github.com/AndroX7/binance-notifier/candle"
)

// Break Of Structure (BOS) detection (simplified)
// For bullish BOS: price closes above last swing high
// For bearish BOS: price closes below last swing low
func BOS(candles []candle.Candle, swings []SwingPoint) (bool, bool) {
	if len(swings) < 2 {
		return false, false
	}
	prevSwing := swings[len(swings)-2]
	lastClose := candles[len(candles)-1].Close

	bullishBOS := false
	bearishBOS := false

	// Bullish BOS: last close above previous swing high price
	if prevSwing.High && lastClose > prevSwing.Price {
		bullishBOS = true
	}

	// Bearish BOS: last close below previous swing low price
	if !prevSwing.High && lastClose < prevSwing.Price {
		bearishBOS = true
	}

	return bullishBOS, bearishBOS
}
