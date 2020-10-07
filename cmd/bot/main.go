package main

import (
	"os"

	"github.com/LarsFox/mca-workshop-bot/bot"
	"github.com/LarsFox/mca-workshop-bot/storage"
	"github.com/LarsFox/mca-workshop-bot/tg"
)

func main() {
	token := os.Getenv("MCA_WORKSHOP_TG_TOKEN")

	tgClient := tg.NewClient(token)
	storageClient, err := storage.NewClient("mca_workshop_bot.db")
	if err != nil {
		panic(err)
	}

	bot := bot.New(storageClient, tgClient, map[string]string{
		bot.ModelBert: os.Getenv("MCA_WORKSHOP_MODEL_BERT"),
		bot.ModelElmo: os.Getenv("MCA_WORKSHOP_MODEL_ELMO"),
	})
	bot.Listen()
}
