package telegram

import (
	"fmt"
	"time"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

	return SendMessage(user, text, &buttons, opts...)
}
