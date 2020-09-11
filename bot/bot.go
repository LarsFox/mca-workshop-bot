package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LarsFox/mca-workshop-bot/tg"
)

// Bot communicates with user via Telegram and with model via HTTP.
type Bot struct {
	modelAddr string
	tgClient  *tg.Client
}

// New returns a new bot.
func New(tgClient *tg.Client, modelAddr string) *Bot {
	return &Bot{
		modelAddr: modelAddr,
		tgClient:  tgClient,
	}
}

// Listen listens the incoming messages.
func (bot *Bot) Listen() {
	log.Println("Listening...")

	msgChan, _ := bot.tgClient.GetMessagesChan()
	for msg := range msgChan {
		if msg != nil {
			go bot.handleMsg(msg)
		}
	}
}

type request struct {
	Model string `json:"model"`
	Text  string `json:"text"`
}

type answer struct {
	Immoral float64 `json:"immoral"`
	Obscene float64 `json:"obscene"`
	Person  float64 `json:"person"`
}

const (
	errMessage   = "Произошла ошибка! Попробуйте повторить запрос чуть позже."
	helpMessage  = "Привет!\nНапиши мне любую реплику, а я оценю, насколько она попадает под критерии оскорбления!"
	okMessage    = "<b>Результаты модели</b>\nАморальность: %.2f\nНаправленность: %.2f\nНецензурность: %.2f"
	startMessage = "Напиши мне любую реплику, а я оценю, насколько она попадает под критерии оскорбления!"
)

func (bot *Bot) handleMsg(msg *tg.Message) {
	if msg.Text == "/start" {
		bot.SendMessage(msg.Chat.ID, startMessage)
		return
	}

	if msg.Text == "/help" {
		bot.SendMessage(msg.Chat.ID, helpMessage)
		return
	}

	b, err := json.Marshal(&request{Text: msg.Text})
	if err != nil {
		log.Println(err)
		bot.SendMessage(msg.Chat.ID, errMessage)
		return
	}

	resp, err := http.Post(bot.modelAddr, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		bot.SendMessage(msg.Chat.ID, errMessage)
		return
	}

	a := &answer{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		log.Println(err)
		bot.SendMessage(msg.Chat.ID, errMessage)
		return
	}

	bot.SendMessage(msg.Chat.ID, fmt.Sprintf(okMessage, a.Immoral, a.Person, a.Obscene))
}

// SendMessage sends a message using tgClient.
func (bot *Bot) SendMessage(chatID int64, text string) {
	resp, err := bot.tgClient.SendMessage(chatID, text)
	if err != nil {
		log.Println(err)
		return
	}

	if !resp.Ok {
		log.Println("not ok!", resp.ErrorCode, string(resp.Description))
	}
}
