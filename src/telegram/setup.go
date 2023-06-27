package telegram

import (
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

func SendWelcomeMessage(user models.User, opts ...Option) (tgbotapi.Message, error) {
	text := "Welcome to the unofficial MyShare bot\nTo start you need to login below"

	loginUrl := fmt.Sprintf("%s/login?chatId=%d", config.Get().Server.Host, user.ChatId)
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "Login", URL: &loginUrl}},
		},
	}

	return sendMessage(user, text, buttons, opts...)
}

func SendLoggedInMessage(user models.User, orgs []models.Org, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("<b>Successfully logged in:</b> as %s\n<b>Select your org:</b>", user.DisplayName)
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: lo.Map(orgs, func(t models.Org, _ int) []tgbotapi.InlineKeyboardButton {
			callback := fmt.Sprintf("%s-%d", models.CommandChangeOrg, t.Id)
			return []tgbotapi.InlineKeyboardButton{{Text: t.Name, CallbackData: &callback}}
		}),
	}

	return sendMessage(user, text, buttons, opts...)
}

func SendChangeOrgMessage(user models.User, orgs []models.Org, opts ...Option) (tgbotapi.Message, error) {
	text := "<b>Change your org:</b>"
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: lo.Map(orgs, func(t models.Org, _ int) []tgbotapi.InlineKeyboardButton {
			callback := fmt.Sprintf("%s-%d", models.CommandChangeOrg, t.Id)
			return []tgbotapi.InlineKeyboardButton{{Text: t.Name, CallbackData: &callback}}
		}),
	}
	buttons.InlineKeyboard = append(buttons.InlineKeyboard, []tgbotapi.InlineKeyboardButton{{Text: "Go back", CallbackData: &models.CommandShowMenu}})

	return sendMessage(user, text, buttons, opts...)
}
