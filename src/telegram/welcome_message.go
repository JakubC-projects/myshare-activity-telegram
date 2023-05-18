package telegram

import (
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

func SendWelcomeMessage(user models.User) (tgbotapi.Message, error) {
	fmt.Println("Send welcome")
	text := "Welcome to the unofficial MyShare bot\nTo start you need to login below"
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "Login", URL: lo.ToPtr(fmt.Sprintf("%s/login?chatId=%d", config.Get().Server.Host, user.ChatId))}},
		},
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func SendLoggedInMessage(user models.User) (tgbotapi.Message, error) {
	text := fmt.Sprintf("Successfully logged in as %s\nSelect your team:", user.DisplayName)
	return Bot.Send(tgbotapi.NewEditMessageText(user.ChatId, user.LastMessageId, text))
}

// func sendMessage(chatId int64, text string) (tgbotapi.Message, error) {
// 	return Bot.Send(tgbotapi.MessageConfig{
// 		BaseChat: tgbotapi.BaseChat{
// 			ChatID: chatId,
// 		},
// 		Text: text,
// 	})
// }

func sendMessageWithMarkup(chatId int64, text string, replyMarkup tgbotapi.InlineKeyboardMarkup) (tgbotapi.Message, error) {
	return Bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatId,
			ReplyMarkup: replyMarkup,
		},
		Text: text,
	})
}
