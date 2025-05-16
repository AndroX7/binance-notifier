package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/AndroX7/binance-notifier/utility"
)

const (
	baseURL = "https://api-fxpractice.oanda.com/v3"
)

type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}

type oandaResponse struct {
	Candles []struct {
		Time string `json:"time"`
		Mid  struct {
			O string `json:"o"`
			H string `json:"h"`
			L string `json:"l"`
			C string `json:"c"`
		} `json:"mid"`
		Volume   int  `json:"volume"`
		Complete bool `json:"complete"`
	} `json:"candles"`
}

// FetchCandles fetches historical candles from OANDA
func FetchCandles(token, instrument, granularity string, count int) ([]Candle, error) {
	url := fmt.Sprintf("%s/instruments/%s/candles?count=%d&granularity=%s&price=M", baseURL, instrument, count, granularity)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result oandaResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var candles []Candle
	for _, c := range result.Candles {
		if !c.Complete {
			continue
		}
		t, _ := time.Parse(time.RFC3339, c.Time)
		open, _ := utility.ParseFloat(c.Mid.O)
		high, _ := utility.ParseFloat(c.Mid.H)
		low, _ := utility.ParseFloat(c.Mid.L)
		close, _ := utility.ParseFloat(c.Mid.C)
		candles = append(candles, Candle{
			Time:   t,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: c.Volume,
		})
	}

	return candles, nil
}
