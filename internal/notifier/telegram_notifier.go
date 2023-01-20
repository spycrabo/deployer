package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	apiUrl = "https://api.telegram.org"
)

type TelegramNotifier struct {
	apiUrl   string
	botToken string
	chatId   int64
	params   map[string]interface{}
}

func NewTelegramNotifier(botToken string, chatId int64, params map[string]interface{}) *TelegramNotifier {
	return &TelegramNotifier{
		apiUrl:   strings.TrimSuffix(apiUrl, "/"),
		botToken: botToken,
		chatId:   chatId,
		params:   params,
	}
}

func (n *TelegramNotifier) Notify(message string) error {
	url := fmt.Sprintf("%s/bot%s/sendMessage", n.apiUrl, n.botToken)
	data := map[string]interface{}{
		"chat_id": n.chatId,
		"text":    message,
	}
	for k, v := range n.params {
		if _, ok := data[k]; !ok {
			data[k] = v
		}
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	return nil
}
