package telegram

import (
	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func init() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(config.Get().Telegram.ApiKey)
	if err != nil {
		panic(err)
	}
}
