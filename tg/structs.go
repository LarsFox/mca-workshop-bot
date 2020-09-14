package tg

import "encoding/json"

// Response is a Telegram response
type Response struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description json.RawMessage `json:"description"` // only if not OK
	ErrorCode   int32           `json:"error_code"`
}

// Update is defined by https://core.telegram.org/bots/api#update.
type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message"`
}

// Message is a Telegram message
type Message struct {
	MessageID int64  `json:"message_id"`
	Text      string `json:"text"`
	Chat      *Chat  `json:"chat"`
}

// Chat is a Telegram chat
type Chat struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type sendMessage struct {
	ChatID      int64  `json:"chat_id"`
	Text        string `json:"text"`
	ParseMode   string `json:"parse_mode"`
	ReplyMarkup string `json:"reply_markup"`
}

// ReplyKeyboardMarkup is a bot Keyboard
type ReplyKeyboardMarkup struct {
	Keyboard        [][]*KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool                `json:"resize_keyboard"`   // optional
	OneTimeKeyboard bool                `json:"one_time_keyboard"` // optional
	Selective       bool                `json:"selective"`         // optional
}

// KeyboardButton one key in ReplyKeyboardMarkup.
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

const noKeyboard = `{"remove_keyboard": true}`
