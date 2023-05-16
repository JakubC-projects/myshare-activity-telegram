package telegram

import (
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

func SendWelcomeMessage(chatId int64) {
	Bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatId,
			ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{{Text: "Login", URL: lo.ToPtr(fmt.Sprintf("%s/login", config.Get().Server.Host))}},
				},
			},
		},
		Text: "Hello world",
	})
}
