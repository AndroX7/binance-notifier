package smc

import (
	"math"

	oanda "github.com/AndroX7/binance-notifier/oanda/client"
)

type SwingPoint struct {
	Index int
	Price float64
	High  bool
}

func DetectSwingPoints(candles []oanda.Candle, period int) []SwingPoint {
	var swings []SwingPoint
	for i := period; i < len(candles)-period; i++ {
		high := true
		low := true
		for j := i - period; j <= i+period; j++ {
			if candles[i].Close < candles[j].Close {
				high = false
			}
			if candles[i].Close > candles[j].Close {
				low = false
			}
		}
		if high {
			swings = append(swings, SwingPoint{Index: i, Price: candles[i].Close, High: true})
		}
		if low {
			swings = append(swings, SwingPoint{Index: i, Price: candles[i].Close, High: false})
		}
	}
	if len(swings) == 0 {
		return nil
	}
	return swings
}

func DetectBOS(candles []oanda.Candle, swings []SwingPoint) (bool, bool) {
	if len(swings) < 2 {
		return false, false
	}
	lastClose := candles[len(candles)-1].Close
	prevSwing := swings[len(swings)-2]

	bullish := false
	bearish := false

	if prevSwing.High && lastClose > prevSwing.Price {
		bullish = true
	}
	if !prevSwing.High && lastClose < prevSwing.Price {
		bearish = true
	}

	return bullish, bearish
}

func DetectCHoCH(swings []SwingPoint) (bool, bool) {
	if len(swings) < 4 {
		return false, false
	}

	last := swings[len(swings)-1]
	prev := swings[len(swings)-2]
	prev3 := swings[len(swings)-4]

	bullish := false
	bearish := false

	if !prev3.High && !prev.High && last.High && prev.Price > prev3.Price && last.Price > prev.Price {
		bullish = true
	}
	if prev3.High && prev.High && last.High && prev.Price < prev3.Price && last.Price < prev.Price {
		bearish = true
	}

	return bullish, bearish
}

type FVG struct {
	Index     int
	Direction string // "buy" or "sell"
	From      float64
	To        float64
}

func DetectFVG(candles []oanda.Candle) []FVG {
	var fvgs []FVG
	for i := 2; i < len(candles); i++ {
		c1 := candles[i-2]
		c3 := candles[i]
		if c1.High < c3.Low {
			// Bullish FVG
			fvgs = append(fvgs, FVG{Index: i, Direction: "buy", From: c1.High, To: c3.Low})
		} else if c1.Low > c3.High {
			// Bearish FVG
			fvgs = append(fvgs, FVG{Index: i, Direction: "sell", From: c3.High, To: c1.Low})
		}
	}
	return fvgs
}

func DetectImbalance(candles []oanda.Candle) []FVG {
	var imbalances []FVG
	for i := 1; i < len(candles); i++ {
		diff := math.Abs(candles[i].Close - candles[i].Open)
		body := math.Abs(candles[i].High - candles[i].Low)
		if body > 2*diff {
			dir := "sell"
			if candles[i].Close > candles[i].Open {
				dir = "buy"
			}
			imbalances = append(imbalances, FVG{
				Index:     i,
				Direction: dir,
				From:      candles[i].Low,
				To:        candles[i].High,
			})
		}
	}
	return imbalances
}

func DetectLiquiditySweeps(candles []oanda.Candle) []int {
	var sweepIndices []int
	for i := 2; i < len(candles); i++ {
		// Simple logic: sweep when current high is slightly higher than previous highs (and reverses)
		if candles[i].High > candles[i-1].High && candles[i].High > candles[i-2].High &&
			candles[i].Close < candles[i].Open {
			sweepIndices = append(sweepIndices, i)
		}
	}
	return sweepIndices
}
