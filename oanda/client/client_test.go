package client_test

// import (
// 	"os"
// 	"testing"
// 	oanda "github.com/AndroX7/binance-notifier/oanda/client"
// )

// func TestFetchCandles(t *testing.T) {
// 	token := os.Getenv("OANDA_API_TOKEN")
// 	accountID := os.Getenv("OANDA_ACCOUNT_ID")
// 	if token == "" || accountID == "" {
// 		t.Skip("OANDA credentials not set")
// 	}

// 	oanda.SetCredentials(token, accountID)

// 	candles, err := oanda.FetchCandles("EUR_USD", "M15", 10)
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}
// 	if len(candles) == 0 {
// 		t.Fatal("Expected candles to be fetched")
// 	}
// }
