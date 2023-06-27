package telegram

import (
	"fmt"
	"sort"
	"time"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/microcosm-cc/bluemonday"
	"github.com/samber/lo"
)

const pageSize = 10

func SendShowActivitiesMessage(user models.User, activities []models.MyshareActivity, page int, opts ...Option) (tgbotapi.Message, error) {
	text := "<b>Available activities</b>:\n"
	if len(activities) > pageSize {
		text += fmt.Sprintf("<b>Page</b>: %d / %d\n", page+1, len(activities)/pageSize+1)
	}
	sort.Slice(activities, func(i, j int) bool {
		return activities[j].Start.After(time.Time(activities[i].Start))
	})
	showedActivities := lo.Slice(activities, page*pageSize, (page+1)*pageSize)

	type DayActivities struct {
		Day        string
		Activities []models.MyshareActivity
	}
	dayActivitiesSlice := []DayActivities{}
	dayIndexMap := map[string]int{}
	for _, a := range showedActivities {
		day := a.Start.Format("Monday, 2. January")
		ind, ok := dayIndexMap[day]
		if !ok {
			dayActivitiesSlice = append(dayActivitiesSlice, DayActivities{day, []models.MyshareActivity{a}})
			dayIndexMap[day] = len(dayActivitiesSlice) - 1
		} else {
			dayActivitiesSlice[ind].Activities = append(dayActivitiesSlice[ind].Activities, a)
		}
	}
	for _, dayActivities := range dayActivitiesSlice {
		text += fmt.Sprintf("┌ <b>%s</b>\n", dayActivities.Day)
		for _, activity := range dayActivities.Activities {
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
	return SendMessage(user, text, &buttons, opts...)
}

var sanitizeHtml = bluemonday.StripTagsPolicy()

func SendShowActivityMessage(user models.User, activity models.MyshareActivity, opts ...Option) (tgbotapi.Message, error) {
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

	return SendMessage(user, text, &buttons, opts...)
}

func SendShowShowParticipantsMessage(user models.User, activity models.MyshareActivity, participants []models.Participant, opts ...Option) (tgbotapi.Message, error) {
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

	return SendMessage(user, text, &buttons, opts...)

}

func SendNewActivitiesNotificationMessage(user models.User, activities []models.MyshareActivity, opts ...Option) (tgbotapi.Message, error) {

	buttons := [][]tgbotapi.InlineKeyboardButton{}
	text := "New activities have appeared!\n\n"
	for _, activity := range activities {
		text += fmt.Sprintf("┌ <b>%s</b>\n", activity.Name)
		text += fmt.Sprintf("| 📅 %s\n", activity.Start.Format("02 Jan"))
		text += fmt.Sprintf("| 🕒 %s -> %s\n",
			activity.Start.Format("15.04"),
			activity.Finish.Format("15.04"))
		text += "\n"

		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			{Text: activity.Name, CallbackData: lo.ToPtr(fmt.Sprintf("%s-%d", models.CommandShowActivity, activity.ActivityId))}})
	}
	text += "\n"

	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
		{Text: "Back to menu", CallbackData: &models.CommandShowMenu},
	})

	return SendMessage(user, text, &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons}, opts...)
}
