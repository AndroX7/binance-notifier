package telegram

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	TelegramBotToken = "YOUR_TELEGRAM_BOT_TOKEN"
	TelegramChatID   = "YOUR_TELEGRAM_CHAT_ID"
)

// Send Telegram message
func Notify(msg string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TelegramBotToken)

	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {
			TelegramChatID},
		"text":       {msg},
		"parse_mode": {"Markdown"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("telegram error: %s", string(body))
	}

	return nil
}
