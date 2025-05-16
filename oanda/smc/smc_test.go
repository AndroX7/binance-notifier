package smc_test

import (
	"testing"

	oanda "github.com/AndroX7/binance-notifier/oanda/client"
	smc "github.com/AndroX7/binance-notifier/oanda/smc"
)

func TestDetectSwingPoints(t *testing.T) {
	candles := []oanda.Candle{
		{Close: 1.0},
		{Close: 1.1},
		{Close: 0.9}, // Low swing
		{Close: 1.3}, // High swing
		{Close: 1.1},
	}

	swings := smc.DetectSwingPoints(candles, 1)
	if len(swings) < 2 {
		t.Errorf("Expected at least 2 swing points, got %d", len(swings))
	}
}

func TestDetectBOS(t *testing.T) {
	candles := []oanda.Candle{
		{Close: 1.1}, {Close: 1.2}, {Close: 1.3},
	}
	swings := []smc.SwingPoint{
		{Index: 0, Price: 1.2, High: true},
		{Index: 1, Price: 1.25, High: true},
	}

	bull, _ := smc.DetectBOS(candles, swings)
	if !bull {
		t.Error("Expected bullish BOS to be true")
	}
}

func TestDetectCHoCH(t *testing.T) {
	swings := []smc.SwingPoint{
		{Index: 0, Price: 1.0, High: false},
		{Index: 1, Price: 1.1, High: false},
		{Index: 2, Price: 1.2, High: false},
		{Index: 3, Price: 1.3, High: false},
	}

	bull, _ := smc.DetectCHoCH(swings)
	if !bull {
		t.Error("Expected bullish CHoCH to be detected")
	}
}

func TestDetectFVG(t *testing.T) {
	candles := []oanda.Candle{
		{High: 1.1, Low: 1.0},
		{High: 1.3, Low: 1.2},
		{High: 1.5, Low: 1.4},
	}
	fvgs := smc.DetectFVG(candles)
	if len(fvgs) == 0 {
		t.Error("Expected FVG to be detected")
	}
}

func TestDetectImbalance(t *testing.T) {
	candles := []oanda.Candle{
		{High: 1.2, Low: 1.0},
		{High: 1.4, Low: 1.3},
		{High: 1.5, Low: 1.4},
	}
	imbalances := smc.DetectImbalance(candles)
	if len(imbalances) == 0 {
		t.Error("Expected imbalance to be detected")
	}
}

func TestDetectLiquiditySweeps(t *testing.T) {
	candles := []oanda.Candle{
		{High: 1.2, Low: 1.1},
		{High: 1.3, Low: 1.2},
		{High: 1.4, Low: 1.3},
	}
	liqs := smc.DetectLiquiditySweeps(candles)
	if len(liqs) == 0 {
		t.Error("Expected liquidity sweep to be detected")
	}
}
