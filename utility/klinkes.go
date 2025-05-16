package utility

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/AndroX7/binance-notifier/candle"
)

const (
	BinanceKlineURL = "https://api.binance.com/api/v3/klines"
)

func FetchKlines(symbol string, interval string, limit int) ([]candle.Candle, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("interval", interval)
	params.Set("limit", fmt.Sprintf("%d", limit))

	resp, err := http.Get(BinanceKlineURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data [][]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	candles := make([]candle.Candle, 0, len(data))
	for _, k := range data {
		openTime := int64(k[0].(float64))
		open, _ := ParseFloat(k[1])
		high, _ := ParseFloat(k[2])
		low, _ := ParseFloat(k[3])
		closePrice, _ := ParseFloat(k[4])
		closeTime := int64(k[6].(float64))

		candles = append(candles, candle.Candle{
			OpenTime:  openTime,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			CloseTime: closeTime,
		})
	}

	return candles, nil
}
