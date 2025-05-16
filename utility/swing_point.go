package utility

import "github.com/AndroX7/binance-notifier/candle"

type SwingPoint struct {
	Index int
	Price float64
	High  bool
	Low   bool
}

// Find swing highs and lows (simple method)
// windowSize defines how many candles to look left and right for swing points
func GetSwingPoint(candles []candle.Candle, windowSize int) []SwingPoint {
	sp := []SwingPoint{}
	for i := windowSize; i < len(candles)-windowSize; i++ {
		isHigh := true
		isLow := true
		for j := i - windowSize; j <= i+windowSize; j++ {
			if candles[j].High > candles[i].High {
				isHigh = false
			}
			if candles[j].Low < candles[i].Low {
				isLow = false
			}
		}
		if isHigh {
			sp = append(sp, SwingPoint{Index: i, Price: candles[i].High, High: true, Low: false})
		}
		if isLow {
			sp = append(sp, SwingPoint{Index: i, Price: candles[i].Low, High: false, Low: true})
		}
	}
	return sp
}
