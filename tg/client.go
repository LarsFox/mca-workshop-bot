package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Client works with Telegram.
type Client struct {
	token     string
	lastUpdID int64
}

// NewClient returns a new client to work with Telegram Bot API.
func NewClient(token string) *Client {
	return &Client{token: token}
}

// GetMessagesChan returns a message channel.
func (c *Client) GetMessagesChan() (<-chan *Message, error) {
	msgChan := make(chan *Message)

	go func() {
		for {
			uri := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?timeout=100&offset=%d", c.token, c.lastUpdID+1)
			res, err := http.Get(uri)
			if err != nil {
				log.Println(err)
				continue
			}

			resp := &Response{}
			defer res.Body.Close()
			if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
				log.Println(err)
				continue
			}

			var updates []*Update
			if err := json.Unmarshal(resp.Result, &updates); err != nil {
				log.Println(err)
				continue
			}

			for _, upd := range updates {
				c.lastUpdID = upd.UpdateID
				msgChan <- upd.Message
			}
		}
	}()

	return msgChan, nil
}

// SendMessage sends a message to chat.
func (c *Client) SendMessage(chatID int64, text string) (*Response, error) {
	uri := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&sendMessage{
		ChatID:    chatID,
		ParseMode: "HTML",
		Text:      text,
	}); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return result, nil
}
