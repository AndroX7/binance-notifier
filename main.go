package main

import (
	"fmt"
	"log"
	"math"

	telegram "github.com/AndroX7/binance-notifier/notify/telegram"
	"github.com/AndroX7/binance-notifier/utility"
)

const (
	Symbol   = "BTCUSDT"
	Interval = "15m"
	Limit    = 100
)

func main() {

	candles, err := utility.FetchKlines(Symbol, Interval, Limit)
	if err != nil {
		log.Fatal(err)
	}

	// 1. Find swing points
	swings := utility.GetSwingPoint(candles, 3)

	// 2. Detect BOS
	bullBOS, bearBOS := utility.BOS(candles, swings)

	// 3. Detect ChoCH
	bullChoCH, bearChoCH := utility.DetectChoCH(swings)

	// 4. Find imbalances (FVGs & gaps)
	imbalances := utility.GetImbalance(candles)

	// 5. Find liquidity pools (last 10 swings)
	highPools, lowPools := utility.FindLiquidityPools(swings, 10)

	lastClose := candles[len(candles)-1].Close

	// Entry logic example
	entryMsg := ""

	if bullBOS && bullChoCH {
		// Check if price is retesting any bullish imbalance or liquidity pool
		for _, imb := range imbalances {
			if imb.Bullish && utility.PriceInZone(lastClose, imb.Top, imb.Bottom) {
				entryMsg = fmt.Sprintf("Bullish Entry Zone: Price %.2f is retesting imbalance zone between %.2f - %.2f", lastClose, imb.Bottom, imb.Top)
				break
			}
		}

		// Check liquidity pools for retest
		for _, lp := range lowPools {
			if math.Abs(lp-lastClose) < 0.005*lastClose { // within 0.5%
				entryMsg = fmt.Sprintf("Bullish Entry Zone: Price %.2f near liquidity pool %.2f", lastClose, lp)
				break
			}
		}
	}

	if bearBOS && bearChoCH {
		for _, imb := range imbalances {
			if !imb.Bullish && utility.PriceInZone(lastClose, imb.Top, imb.Bottom) {
				entryMsg = fmt.Sprintf("Bearish Entry Zone: Price %.2f is retesting imbalance zone between %.2f - %.2f", lastClose, imb.Bottom, imb.Top)
				break
			}
		}

		for _, hp := range highPools {
			if math.Abs(hp-lastClose) < 0.005*lastClose { // within 0.5%
				entryMsg = fmt.Sprintf("Bearish Entry Zone: Price %.2f near liquidity pool %.2f", lastClose, hp)
				break
			}
		}
	}

	if entryMsg != "" {
		fmt.Println(entryMsg)
		err := telegram.Notify(entryMsg)
		if err != nil {
			log.Println("Telegram send error:", err)
		}
	} else {
		fmt.Println("No sharp entry conditions met now.")
	}

}
