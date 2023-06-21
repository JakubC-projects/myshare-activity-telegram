package telegram

import (
	"fmt"
	"sort"
	"time"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/microcosm-cc/bluemonday"
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
	text := fmt.Sprintf("<b>Successfully logged in:</b> as %s\n<b>Select your org:</b>", user.DisplayName)
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: lo.Map(orgs, func(t models.Org, _ int) []tgbotapi.InlineKeyboardButton {
			callback := fmt.Sprintf("%s-%d", models.CommandChangeOrg, t.Id)
			return []tgbotapi.InlineKeyboardButton{{Text: t.Name, CallbackData: &callback}}
		}),
	}
	if isEdit(opts) {
		msg := tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons)
		msg.ParseMode = tgbotapi.ModeHTML
		return Bot.Send(msg)
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
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
	if isEdit(opts) {
		msg := tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons)
		msg.ParseMode = tgbotapi.ModeHTML
		return Bot.Send(msg)
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

var pageSize = 10

func SendShowActivitiesMessage(user models.User, activities []models.Activity, page int, opts ...Option) (tgbotapi.Message, error) {
	text := "<b>Available activities</b>:\n"
	if len(activities) > pageSize {
		text += fmt.Sprintf("<b>Page</b>: %d / %d\n", page+1, len(activities)/pageSize+1)
	}
	sort.Slice(activities, func(i, j int) bool {
		return activities[j].Start.After(time.Time(activities[i].Start))
	})
	showedActivities := lo.Slice(activities, page*pageSize, (page+1)*pageSize)
	dayActivities := lo.GroupBy(showedActivities, func(a models.Activity) string {
		return a.Start.Format("Monday, 2. January")
	})
	for day, activities := range dayActivities {
		text += fmt.Sprintf("┌ <b>%s</b>\n", day)
		for _, activity := range activities {
			registeredMarker := "⚪"
			if activity.RegistrationStatus == models.RegistrationRegistered {
				registeredMarker = "🟢"
			}
			text += fmt.Sprintf("|  ┌ <b>%s</b> %s\n", activity.Name, registeredMarker)
			text += fmt.Sprintf("|  | 🕒 %s -> %s\n",
				activity.Start.Format("15.04"),
				activity.Finish.Format("15.04"))

			text += fmt.Sprintf("|  | 👥 %d / %d\n", activity.Registrations, activity.NeededRegistrations)
		}
		text += "\n"
	}
	text += fmt.Sprintf("\nRefreshed at: %s", time.Now().Format("2006-01-02 15:04:05"))

	var inlineKeyboard [][]tgbotapi.InlineKeyboardButton

	if len(activities) > pageSize {
		nextPage := page + 1
		previousPage := page - 1
		if previousPage < 0 {
			previousPage = 0
		}
		if nextPage*pageSize > len(activities) {
			nextPage -= 1
		}
		firstPageCallback := fmt.Sprintf("%s-%d", models.CommandShowActivities, 0)
		previousPageCallback := fmt.Sprintf("%s-%d", models.CommandShowActivities, previousPage)
		nextPageCallback := fmt.Sprintf("%s-%d", models.CommandShowActivities, nextPage)
		lastPageCallback := fmt.Sprintf("%s-%d", models.CommandShowActivities, len(activities)/pageSize)

		inlineKeyboard = append(inlineKeyboard, []tgbotapi.InlineKeyboardButton{
			{Text: "<<", CallbackData: &firstPageCallback},
			{Text: "<", CallbackData: &previousPageCallback},
			{Text: ">", CallbackData: &nextPageCallback},
			{Text: ">>", CallbackData: &lastPageCallback},
		})
	}

	for _, a := range showedActivities {
		detailsCallback := fmt.Sprintf("%s-%d", models.CommandShowActivity, a.ActivityId)
		keyboardRow := []tgbotapi.InlineKeyboardButton{{Text: a.Name, CallbackData: &detailsCallback}}
		inlineKeyboard = append(inlineKeyboard, keyboardRow)
	}

	inlineKeyboard = append(inlineKeyboard, []tgbotapi.InlineKeyboardButton{
		{Text: "Go back", CallbackData: &models.CommandShowMenu},
		{Text: "Refresh", CallbackData: &models.CommandShowActivities},
	})

	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}

	if isEdit(opts) {
		msg := tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons)
		msg.ParseMode = tgbotapi.ModeHTML
		return Bot.Send(msg)
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

var sanitizeHtml = bluemonday.StripTagsPolicy()

func SendShowActivityMessage(user models.User, activity models.Activity, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("<b>%s</b>\n\n", activity.Name)
	text += fmt.Sprintf("📅 %s -> %s\n\n",
		activity.Start.Format("02 Jan"),
		activity.Finish.Format("02 Jan"))

	text += fmt.Sprintf("🕒 %s -> %s\n\n",
		activity.Start.Format("15.04"),
		activity.Finish.Format("15.04"))

	text += fmt.Sprintf("🌍 %s\n\n",
		activity.ActivityLocation)

	text += fmt.Sprintf("👥 %d / %d \n\n", activity.Registrations, activity.NeededRegistrations)

	regStatus := "👤⚪"
	if activity.RegistrationStatus == models.RegistrationRegistered {
		regStatus = "👤🟢"
	}
	text += regStatus + "\n\n"

	if activity.Description != "" {
		description := sanitizeHtml.Sanitize(activity.Description)
		text += fmt.Sprintf("📜 %s\n", description)
	}

	text += fmt.Sprintf("\nRefreshed at: %s", time.Now().Format("2006-01-02 15:04:05"))

	inlineButtons := [][]tgbotapi.InlineKeyboardButton{}

	if activity.RegistrationStatus == models.RegistrationNotRegistered {
		registerCallback := fmt.Sprintf("%s-%d", models.CommandRegisterActivity, activity.ActivityId)
		inlineButtons = append(inlineButtons, []tgbotapi.InlineKeyboardButton{{Text: "Register", CallbackData: &registerCallback}})
	} else if activity.RegistrationStatus == models.RegistrationRegistered {
		unregisterCallback := fmt.Sprintf("%s-%d", models.CommandUnregisterActivity, activity.ActivityId)
		inlineButtons = append(inlineButtons, []tgbotapi.InlineKeyboardButton{{Text: "Unregister", CallbackData: &unregisterCallback}})
	}

	refreshCallback := fmt.Sprintf("%s-%d", models.CommandShowActivity, activity.ActivityId)
	showParticipantsCallback := fmt.Sprintf("%s-%d", models.CommandShowParticipants, activity.ActivityId)

	inlineButtons = append(inlineButtons, []tgbotapi.InlineKeyboardButton{{Text: "Show participants", CallbackData: &showParticipantsCallback}})
	inlineButtons = append(inlineButtons, []tgbotapi.InlineKeyboardButton{
		{Text: "Go back", CallbackData: &models.CommandShowActivities},
		{Text: "Refresh", CallbackData: &refreshCallback},
	})

	buttons := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineButtons}

	if isEdit(opts) {
		msg := tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons)
		msg.ParseMode = tgbotapi.ModeHTML
		return Bot.Send(msg)
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func SendShowShowParticipantsMessage(user models.User, activity models.Activity, participants []models.Participant, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("<b>%s</b>\n", activity.Name)
	text += fmt.Sprintf("👥 %d / %d\n", activity.Registrations, activity.NeededRegistrations)
	text += "\n"

	for _, p := range participants {
		text += p.DisplayName
		if p.Comments != "" {
			text += fmt.Sprintf(" - %s", p.Comments)
		}
		text += "\n"
	}

	text += fmt.Sprintf("\nRefreshed at: %s", time.Now().Format("2006-01-02 15:04:05"))

	refreshCallback := fmt.Sprintf("%s-%d", models.CommandShowParticipants, activity.ActivityId)
	goBackCallback := fmt.Sprintf("%s-%d", models.CommandShowActivity, activity.ActivityId)

	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "Go back", CallbackData: &goBackCallback},
				{Text: "Refresh", CallbackData: &refreshCallback}},
		},
	}

	if isEdit(opts) {
		msg := tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons)
		msg.ParseMode = tgbotapi.ModeHTML
		return Bot.Send(msg)
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func SendMenuMessage(user models.User, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("<b>Logged in as</b>: %s\n<b>Selected org:</b> %s\n<b>Choose your action:</b>", user.DisplayName, user.Org.Name)
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "Show activities", CallbackData: &models.CommandShowActivities}},
			{{Text: "Show status", CallbackData: &models.CommandShowStatus}},
			{{Text: "Change org", CallbackData: &models.CommandStartChangeOrg}},
			{{Text: "Logout", CallbackData: &models.CommandLogout}},
		},
	}
	if isEdit(opts) {
		msg := tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons)
		msg.ParseMode = tgbotapi.ModeHTML
		return Bot.Send(msg)
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func SendShowStatusMessage(user models.User, status models.Status, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("My status\n%.2f%%\n%.2f / %.2f %s",
		status.PercentageValue,
		status.TransactionsAmount,
		status.Target,
		status.Currency)

	text += fmt.Sprintf("\nRefreshed at: %s", time.Now().Format("2006-01-02 15:04:05"))

	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "Go back", CallbackData: &models.CommandShowMenu},
				{Text: "Refresh", CallbackData: &models.CommandShowStatus}},
		}}
	if isEdit(opts) {
		msg := tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, buttons)
		msg.ParseMode = tgbotapi.ModeHTML
		return Bot.Send(msg)
	}
	return sendMessageWithMarkup(user.ChatId, text, buttons)
}

func sendMessageWithMarkup(chatId int64, text string, replyMarkup tgbotapi.InlineKeyboardMarkup) (tgbotapi.Message, error) {
	return Bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatId,
			ReplyMarkup: replyMarkup,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	})
}
