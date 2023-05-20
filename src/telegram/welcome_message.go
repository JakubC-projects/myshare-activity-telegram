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

func SendLoggedInMessage(user models.User, teams []models.Team) (tgbotapi.Message, error) {
	text := fmt.Sprintf("Successfully logged in as %s\nSelect your team:", user.DisplayName)
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: lo.Map(teams, func(t models.Team, _ int) []tgbotapi.InlineKeyboardButton {
			return []tgbotapi.InlineKeyboardButton{{Text: t.Name, CallbackData: lo.ToPtr(fmt.Sprint(t.TeamId))}}
		}),
	}
	return Bot.Send(tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons))
}

func sendMessageWithMarkup(chatId int64, text string, replyMarkup tgbotapi.InlineKeyboardMarkup) (tgbotapi.Message, error) {
	return Bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatId,
			ReplyMarkup: replyMarkup,
		},
		Text: text,
	})
}
