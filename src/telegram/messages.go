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
			{{Text: "Login", URL: &loginUrl}},
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
	buttons.InlineKeyboard = append(buttons.InlineKeyboard, []tgbotapi.InlineKeyboardButton{{Text: "Go back", CallbackData: &models.CommandShowMenu}})
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
		{Text: "Go back", CallbackData: &models.CommandShowMenu},
		{Text: "Refresh", CallbackData: &models.CommandShowActivities},
	})
	if isEdit(opts) {
		return Bot.Send(tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons))
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func SendShowActivityMessage(user models.User, activity models.Activity, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("Name: %s\nDate: %s-%s\nTime: %s-%s\nDescription: %s",
		activity.Name,
		activity.Start.Format("02 Jan"),
		activity.Finish.Format("02 Jan"),
		activity.Start.Format("15.04"),
		activity.Finish.Format("15.04"),
		activity.Description,
	)

	refreshCallback := fmt.Sprintf("%s-%d", models.CommandShowActivity, activity.ActivityId)
	showParticipantsCallback := fmt.Sprintf("%s-%d", models.CommandShowParticipants, activity.ActivityId)
	registerCallback := fmt.Sprintf("%s-%d", models.CommandRegisterActivity, activity.ActivityId)
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "Register", CallbackData: &registerCallback}},
			{{Text: "Show participants", CallbackData: &showParticipantsCallback}},
			{{Text: "Go back", CallbackData: &models.CommandShowActivities},
				{Text: "Refresh", CallbackData: &refreshCallback}},
		},
	}

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
			{{Text: "Show status", CallbackData: &models.CommandShowStatus}},
			{{Text: "Change org", CallbackData: &models.CommandStartChangeOrg}},
			{{Text: "Logout", CallbackData: &models.CommandLogout}},
		},
	}
	if isEdit(opts) {
		return Bot.Send(tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons))
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func SendShowStatusMessage(user models.User, status models.Status, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("My status (refreshed at %s)\n%.2f%%\n%.2f / %.2f %s",
		time.Now().Format("2006-01-02 15:04:05"),
		status.PercentageValue,
		status.TransactionsAmount,
		status.Target,
		status.Currency)
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "Go back", CallbackData: &models.CommandShowMenu},
				{Text: "Refresh", CallbackData: &models.CommandShowStatus}},
		}}
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
