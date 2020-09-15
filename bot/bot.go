package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LarsFox/mca-workshop-bot/storage"
	"github.com/LarsFox/mca-workshop-bot/tg"
)

const (
	errMessage      = "Произошла ошибка! Попробуйте повторить запрос чуть позже."
	okMessage       = "<b>Результаты модели %s</b>\nАморальность: %.2f\nНаправленность: %.2f\nНецензурность: %.2f"
	selectMessage   = "В какую модель отправлять реплики на анализ?"
	selectedMessage = "Понял. Все последующие реплики буду отправлять в %s\n\nСмело отправляй реплики на оценку!"
	modelAll        = "Во все сразу 📊"
)

var singleCommands = map[string]string{
	"/start":   "Привет!\nНапиши мне любую реплику, а я оценю, насколько она попадает под критерии оскорбления.\n\nПо умолчанию анализ проходит по модели Элмо. Если хочешь выбрать другую модель, отправь /select",
	"/help":    "Напиши мне любую реплику, а я оценю, насколько она попадает под критерии оскорбления.\n\nПо умолчанию для оценки используется модель Элмо — отправь /select, чтобы выбрать другую",
	textCancel: "Окей.\nНапиши мне любую реплику, а я оценю, насколько она попадает под критерии оскорбления",
}

// Available models.
const (
	ModelBert     = "bert"
	ModelElmo     = "elmo"
	ModelFasttext = "fasttext"
)

const (
	textBert     = "Берт 🅱️"
	textElmo     = "Элмо 🍪"
	textFasttext = "Фасттекст 🔤"
	textCancel   = "Отмена 🔙"
)

var modelsTexts = map[string]string{
	ModelBert:     textBert,
	ModelElmo:     textElmo,
	ModelFasttext: textFasttext,
}

type request struct {
	Text string `json:"text"`
}

type answer struct {
	Immoral float64 `json:"immoral"`
	Obscene float64 `json:"obscene"`
	Person  float64 `json:"person"`
}

// Bot communicates with user via Telegram and with model via HTTP.
type Bot struct {
	modelAddrs    map[string]string
	storageClient *storage.Client
	tgClient      *tg.Client
}

// New returns a new bot.
func New(storageClient *storage.Client, tgClient *tg.Client, modelAddrs map[string]string) *Bot {
	return &Bot{
		modelAddrs:    modelAddrs,
		storageClient: storageClient,
		tgClient:      tgClient,
	}
}

// Listen listens the incoming messages.
func (bot *Bot) Listen() {
	log.Println("Listening...")

	msgChan, err := bot.tgClient.GetMessagesChan()
	if err != nil {
		log.Println("Error with the message channel, try again")
		return
	}

	for msg := range msgChan {
		if msg != nil {
			go bot.handleMsg(msg)
		}
	}
}

func (bot *Bot) handleMsg(msg *tg.Message) {
	for command, text := range singleCommands {
		if msg.Text != command {
			continue
		}

		bot.sendMessage(msg.Chat.ID, text, nil)
		return
	}

	if msg.Text == "/select" {
		bot.sendMessage(msg.Chat.ID, selectMessage, &tg.ReplyKeyboardMarkup{
			Keyboard: [][]*tg.KeyboardButton{
				{{Text: textBert}},
				{{Text: textFasttext}},
				{{Text: textElmo}},
				{{Text: textCancel}},
			},
			OneTimeKeyboard: true,
		})
		return
	}

	for model, text := range modelsTexts {
		if msg.Text != text {
			continue
		}

		if err := bot.storageClient.SaveUserModel(msg.Chat.ID, model); err != nil {
			log.Println(err)
			bot.sendErrorMessage(msg.Chat.ID)
			return
		}

		bot.tgClient.SendMessage(msg.Chat.ID, fmt.Sprintf(selectedMessage, text), nil)
		return
	}

	log.Println("message from:", msg.Chat.ID, msg.Chat.Username, msg.Text)
	for _, p := range msg.Photo {
		log.Println("photo from:", msg.Chat.ID, msg.Chat.Username, p.FileID, p.FileUniqueID)
	}

	model, err := bot.storageClient.GetUserModel(msg.Chat.ID)
	if err != nil {
		log.Println(err)
		bot.sendErrorMessage(msg.Chat.ID)
		return
	}

	if model == "" {
		model = ModelElmo
	}

	addr, ok := bot.modelAddrs[model]
	if !ok {
		log.Println("no addr for model", model)
		bot.sendErrorMessage(msg.Chat.ID)
		return
	}

	a, err := bot.callModel(addr, msg.Text)
	if err != nil {
		bot.sendErrorMessage(msg.Chat.ID)
		return
	}
	bot.sendMessage(msg.Chat.ID, fmt.Sprintf(okMessage, modelsTexts[model], a.Immoral, a.Person, a.Obscene), nil)
}

func (bot *Bot) callModel(addr, text string) (*answer, error) {
	b, err := json.Marshal(&request{Text: text})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(b)) // #nosec 117
	if err != nil {
		return nil, err
	}

	a := &answer{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		return nil, err
	}
	return a, nil
}

// sendMessage sends a message using tgClient.
func (bot *Bot) sendMessage(chatID int64, text string, keyboard *tg.ReplyKeyboardMarkup) {
	resp, err := bot.tgClient.SendMessage(chatID, text, keyboard)
	if err != nil {
		log.Println(err)
		return
	}

	if !resp.Ok {
		log.Println("not ok!", resp.ErrorCode, string(resp.Description))
	}
}

// sendErrorMessage sends an error message using tgClient.
func (bot *Bot) sendErrorMessage(chatID int64) {
	bot.sendMessage(chatID, errMessage, nil)
}
