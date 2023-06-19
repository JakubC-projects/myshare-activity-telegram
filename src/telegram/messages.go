package telegram

import (
	"fmt"
	"time"

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
			{{Text: "Login", WebApp: &tgbotapi.WebAppInfo{URL: loginUrl}}},
		},
	}

	if isEdit(opts) {
		return Bot.Send(tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons))
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func SendLoggedInMessage(user models.User, orgs []models.Org, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("Successfully logged in as %s\nSelect your org:", user.DisplayName)
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: lo.Map(orgs, func(t models.Org, _ int) []tgbotapi.InlineKeyboardButton {
			callback := fmt.Sprintf("%s-%d", models.CommandChangeOrg, t.Id)
			return []tgbotapi.InlineKeyboardButton{{Text: t.Name, CallbackData: &callback}}
		}),
	}
	if isEdit(opts) {
		return Bot.Send(tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons))
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func SendChangeOrgMessage(user models.User, orgs []models.Org, opts ...Option) (tgbotapi.Message, error) {
	text := "Change your org:"
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: lo.Map(orgs, func(t models.Org, _ int) []tgbotapi.InlineKeyboardButton {
			callback := fmt.Sprintf("%s-%d", models.CommandChangeOrg, t.Id)
			return []tgbotapi.InlineKeyboardButton{{Text: t.Name, CallbackData: &callback}}
		}),
	}
	buttons.InlineKeyboard = append(buttons.InlineKeyboard, []tgbotapi.InlineKeyboardButton{{Text: "Back to menu", CallbackData: &models.CommandShowMenu}})
	if isEdit(opts) {
		return Bot.Send(tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons))
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func SendShowActivitiesMessage(user models.User, activities []models.Activity, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("Available activities (refreshed at %s):", time.Now().Format("2006-01-02 15:04:05"))
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: lo.Map(activities, func(a models.Activity, _ int) []tgbotapi.InlineKeyboardButton {
			text := fmt.Sprintf("%s %s %d/%d", a.Name, a.Start.Format("2006-01-02 15:04:05"), a.Registrations, a.NeededRegistrations)
			callback := fmt.Sprintf("%s-%d", models.CommandShowActivity, a.ActivityId)
			return []tgbotapi.InlineKeyboardButton{{Text: text, CallbackData: &callback}}
		}),
	}
	buttons.InlineKeyboard = append(buttons.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		{Text: "Back to menu", CallbackData: &models.CommandShowMenu},
		{Text: "Refresh", CallbackData: &models.CommandShowActivities},
	})
	if isEdit(opts) {
		return Bot.Send(tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons))
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func SendMenuMessage(user models.User, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("Logged in as %s\nSelected org: %s\nChoose your action:", user.DisplayName, user.Org.Name)
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "Show activities", CallbackData: &models.CommandShowActivities}},
			{{Text: "Change org", CallbackData: &models.CommandStartChangeOrg}},
			{{Text: "Logout", CallbackData: &models.CommandLogout}},
		},
	}
	if isEdit(opts) {
		return Bot.Send(tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons))
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
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
