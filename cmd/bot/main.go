package main

import (
	"os"

	"github.com/LarsFox/mca-workshop-bot/bot"
	"github.com/LarsFox/mca-workshop-bot/tg"
)

func main() {
	modelAddr := os.Getenv("MCA_WORKSHOP_MODEL_ADDR")
	token := os.Getenv("MCA_WORKSHOP_TG_TOKEN")

	tgClient := tg.NewClient(token)

	bot := bot.New(tgClient, modelAddr)
	bot.Listen()
}
