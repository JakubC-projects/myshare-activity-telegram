package telegram

import (
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMenuMessage(user models.User, opts ...Option) (tgbotapi.Message, error) {
	text := fmt.Sprintf("<b>Logged in as</b>: %s\n<b>Selected org:</b> %s\n<b>Choose your action:</b>", user.DisplayName, user.Org.Name)
	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "Show activities", CallbackData: &models.CommandShowActivities}},
			{{Text: "Show status", CallbackData: &models.CommandShowStatus}},
			{{Text: "Change org", CallbackData: &models.CommandStartChangeOrg}},
		},
	}
	if user.NotificationsSettings.Enabled {
		buttons.InlineKeyboard = append(buttons.InlineKeyboard,
			[]tgbotapi.InlineKeyboardButton{{Text: "Disable notifications", CallbackData: &models.CommandDisableNotifications}})
	} else {
		buttons.InlineKeyboard = append(buttons.InlineKeyboard,
			[]tgbotapi.InlineKeyboardButton{{Text: "Enable notifications", CallbackData: &models.CommandEnableNotifications}})
	}

	buttons.InlineKeyboard = append(buttons.InlineKeyboard,
		[]tgbotapi.InlineKeyboardButton{{Text: "Logout", CallbackData: &models.CommandLogout}})

	return SendMessage(user, text, &buttons, opts...)

}
