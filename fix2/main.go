package main

// import (
// 	"fmt"
// 	"log"
// 	"math"
// 	"os"
// 	"os/signal"
// 	"strconv"
// 	"syscall"

// 	binance "github.com/adshao/go-binance/v2"
// 	"github.com/go-resty/resty/v2"
// )

// type Candle struct {
// 	Open, High, Low, Close float64
// 	Timestamp              int64
// }

// type SwingPoint struct {
// 	Index int
// 	Price float64
// 	High  bool // true if swing high, false if swing low
// }

// type FVG struct {
// 	StartIndex, EndIndex int
// 	UpperGap, LowerGap   float64
// 	Valid                bool
// }

// type SignalType string

// const (
// 	BuySignal  SignalType = "BUY"
// 	SellSignal SignalType = "SELL"
// )

// // Globals
// var (
// 	telegramBotToken = "YOUR_TELEGRAM_BOT_TOKEN"
// 	telegramChatID   = "YOUR_TELEGRAM_CHAT_ID"
// )

// func main() {
// 	// Channel to handle graceful shutdown
// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

// 	// Start Binance websocket to get 1m candles for BTCUSDT
// 	doneC, stopC, err := binance.WsKlineServe("btcusdt", binance.KlineInterval1m, handleKline)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("Started Binance WebSocket for BTCUSDT 1m candles...")

// 	select {
// 	case <-stop:
// 		log.Println("Shutting down...")
// 		stopC <- struct{}{}
// 		<-doneC
// 	}
// }

// // buffer of last candles & swings
// var candles []Candle
// var swings []SwingPoint

// func handleKline(event *binance.WsKlineEvent) {
// 	k := event.Kline
// 	candle := Candle{
// 		Open:      mustParseFloat(k.Open),
// 		High:      mustParseFloat(k.High),
// 		Low:       mustParseFloat(k.Low),
// 		Close:     mustParseFloat(k.Close),
// 		Timestamp: k.StartTime,
// 	}

// 	candles = append(candles, candle)
// 	if len(candles) > 1000 {
// 		candles = candles[1:] // keep buffer size manageable
// 	}

// 	// Detect swing points on updated candles
// 	swings = detectSwings(candles)

// 	// Detect FVG
// 	fvgZones := detectFVG(candles)

// 	// Detect BOS and ChoCH
// 	bullishBOS, bearishBOS := detectBOS(candles, swings)
// 	bullishChoCH, bearishChoCH := detectChoCH(swings)

// 	// Detect liquidity zones from swings
// 	liquidityZones := getLiquidityZones(swings)

// 	// Combine all to detect signals
// 	signals := analyzeSignals(candles, swings, fvgZones, liquidityZones, bullishBOS, bearishBOS, bullishChoCH, bearishChoCH)

// 	// Send Telegram alerts for signals
// 	for _, s := range signals {
// 		sendTelegramMessage(fmt.Sprintf("%s signal detected at price %.2f", s.Type, s.EntryPrice))
// 	}
// }

// // Helper to parse string to float64 safely
// func mustParseFloat(str string) float64 {
// 	f, err := strconv.ParseFloat(str, 64)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return f
// }

// // Detect swing points based on local high/lows
// func detectSwings(candles []Candle) []SwingPoint {
// 	var swings []SwingPoint

// 	for i := 2; i < len(candles)-2; i++ {
// 		c := candles[i]

// 		// Check for swing high
// 		if c.High > candles[i-1].High && c.High > candles[i-2].High && c.High > candles[i+1].High && c.High > candles[i+2].High {
// 			swings = append(swings, SwingPoint{Index: i, Price: c.High, High: true})
// 		}

// 		// Check for swing low
// 		if c.Low < candles[i-1].Low && c.Low < candles[i-2].Low && c.Low < candles[i+1].Low && c.Low < candles[i+2].Low {
// 			swings = append(swings, SwingPoint{Index: i, Price: c.Low, High: false})
// 		}
// 	}
// 	return swings
// }

// // Detect Fair Value Gap zones
// func detectFVG(candles []Candle) []FVG {
// 	var fvgList []FVG
// 	for i := 2; i < len(candles); i++ {
// 		// Bullish FVG (gap down)
// 		if candles[i-2].Low > candles[i].High {
// 			fvgList = append(fvgList, FVG{
// 				StartIndex: i - 2,
// 				EndIndex:   i,
// 				UpperGap:   candles[i].High,
// 				LowerGap:   candles[i-2].Low,
// 				Valid:      true,
// 			})
// 		} else if candles[i-2].High < candles[i].Low {
// 			// Bearish FVG (gap up)
// 			fvgList = append(fvgList, FVG{
// 				StartIndex: i - 2,
// 				EndIndex:   i,
// 				UpperGap:   candles[i-2].High,
// 				LowerGap:   candles[i].Low,
// 				Valid:      true,
// 			})
// 		}
// 	}
// 	return fvgList
// }

// // Detect BOS using previous swings and current candle close
// func detectBOS(candles []Candle, swings []SwingPoint) (bool, bool) {
// 	if len(swings) < 2 || len(candles) < 1 {
// 		return false, false
// 	}
// 	lastSwing := swings[len(swings)-1]
// 	prevSwing := swings[len(swings)-2]
// 	lastClose := candles[len(candles)-1].Close

// 	bullishBOS := false
// 	bearishBOS := false

// 	// Bullish BOS: last close above previous swing high price
// 	if prevSwing.High && lastClose > prevSwing.Price {
// 		bullishBOS = true
// 	}

// 	// Bearish BOS: last close below previous swing low price
// 	if !prevSwing.High && lastClose < prevSwing.Price {
// 		bearishBOS = true
// 	}

// 	return bullishBOS, bearishBOS
// }

// // Detect ChoCH (Change of Character)
// func detectChoCH(swings []SwingPoint) (bool, bool) {
// 	if len(swings) < 4 {
// 		return false, false
// 	}

// 	last := swings[len(swings)-1]
// 	prev := swings[len(swings)-2]
// 	prev2 := swings[len(swings)-3]
// 	prev3 := swings[len(swings)-4]

// 	bullishChoCH := false
// 	bearishChoCH := false

// 	// Bullish ChoCH: last three lows progressively higher
// 	if !prev3.High && !prev2.High && !prev.High && !last.High &&
// 		prev2.Price > prev3.Price &&
// 		prev.Price > prev2.Price &&
// 		last.Price > prev.Price {
// 		bullishChoCH = true
// 	}

// 	// Bearish ChoCH: last three highs progressively lower
// 	if prev3.High && prev2.High && prev.High && last.High &&
// 		prev2.Price < prev3.Price &&
// 		prev.Price < prev2.Price &&
// 		last.Price < prev.Price {
// 		bearishChoCH = true
// 	}

// 	return bullishChoCH, bearishChoCH
// }

// // Get liquidity zones from swing points
// func getLiquidityZones(swings []SwingPoint) []float64 {
// 	var zones []float64
// 	for _, sp := range swings {
// 		zones = append(zones, sp.Price)
// 	}
// 	return zones
// }

// // Signal struct for detected signals
// type Signal struct {
// 	Type       SignalType
// 	EntryPrice float64
// 	StopLoss   float64
// 	TakeProfit float64
// }

// // Analyze all conditions and generate signals
// func analyzeSignals(
// 	candles []Candle,
// 	swings []SwingPoint,
// 	fvgZones []FVG,
// 	liquidityZones []float64,
// 	bullishBOS, bearishBOS, bullishChoCH, bearishChoCH bool,
// ) []Signal {

// 	var signals []Signal
// 	lastClose := candles[len(candles)-1].Close

// 	// Loop FVG zones to find price proximity & confluence with BOS, ChoCH and liquidity
// 	for _, fvg := range fvgZones {
// 		distToLower := math.Abs(lastClose - fvg.LowerGap)
// 		distToUpper := math.Abs(lastClose - fvg.UpperGap)

// 		// Threshold proximity to FVG
// 		proximityThreshold := 0.5 // adjust to your preference

// 		if bullishBOS && bullishChoCH && distToLower < proximityThreshold {
// 			// Find closest liquidity zone below price for stop loss
// 			sl := findClosestLiquidityBelow(liquidityZones, lastClose)
// 			tp := lastClose + (lastClose-sl)*2 // 2:1 RR ratio example

// 			signals = append(signals, Signal{
// 				Type:       BuySignal,
// 				EntryPrice: lastClose,
// 				StopLoss:   sl,
// 				TakeProfit: tp,
// 			})
// 		} else if bearishBOS && bearishChoCH && distToUpper < proximityThreshold {
// 			// Find closest liquidity zone above price for stop loss
// 			sl := findClosestLiquidityAbove(liquidityZones, lastClose)
// 			tp := lastClose - (sl-lastClose)*2 // 2:1 RR ratio example

// 			signals = append(signals, Signal{
// 				Type:       SellSignal,
// 				EntryPrice: lastClose,
// 				StopLoss:   sl,
// 				TakeProfit: tp,
// 			})
// 		}
// 	}

// 	return signals
// }

// // Find closest liquidity zone below price (for SL)
// func findClosestLiquidityBelow(zones []float64, price float64) float64 {
// 	closest := price * 0.99 // default 1% below
// 	minDiff := math.MaxFloat64
// 	for _, zone := range zones {
// 		if zone < price && (price-zone) < minDiff {
// 			minDiff = price - zone
// 			closest = zone
// 		}
// 	}
// 	return closest
// }

// // Find closest liquidity zone above price (for SL)
// func findClosestLiquidityAbove(zones []float64, price float64) float64 {
// 	closest := price * 1.01 // default 1% above
// 	minDiff := math.MaxFloat64
// 	for _, zone := range zones {
// 		if zone > price && (zone-price) < minDiff {
// 			minDiff = zone - price
// 			closest = zone
// 		}
// 	}
// 	return closest
// }

// // Send message to Telegram bot
// func sendTelegramMessage(msg string) {
// 	client := resty.New()
// 	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramBotToken)
// 	_, err := client.R().
// 		SetQueryParams(map[string]string{
// 			"chat_id": telegramChatID,
// 			"text":    msg,
// 		}).
// 		Get(url)
// 	if err != nil {
// 		log.Println("Error sending Telegram message:", err)
// 	}
// }
