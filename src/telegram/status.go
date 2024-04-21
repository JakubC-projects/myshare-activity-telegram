package telegram

import (
	"fmt"
	"time"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

func SendShowStatusMessage(user models.User, status models.Status, editedMessageId int) (tgbotapi.Message, error) {
	text := fmt.Sprintf(`<b>Your Status</b>
%.2f%%
%.2f / %.2f %s`,
		status.PercentageValue,
		status.TransactionsAmount,
		status.Target,
		status.Currency)

	text += getPeacefulRoadStatus(status)

	text += fmt.Sprintf("\n\nRefreshed at: %s", time.Now().Format("2006-01-02 15:04:05"))

	buttons := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{{Text: "Go back", CallbackData: &models.CommandShowMenu},
				{Text: "Refresh", CallbackData: &models.CommandShowStatus}},
		}}

	return SendMessage(user, text, &buttons, editedMessageId)
}

var peacefulRoadStartTime = lo.Must(time.Parse(time.RFC3339, "2024-03-06T00:00:00Z"))
var peacefulRoadEndTime = lo.Must(time.Parse(time.RFC3339, "2024-07-11T23:59:59Z"))
var peacefulRoadStartPercentage = 35.0
var peacefulRoadEndPercentage = 70.0

func getPeacefulRoadStatus(status models.Status) string {
	now := time.Now()
	if now.Before(peacefulRoadStartTime) || now.After(peacefulRoadEndTime) {
		return ""
	}

	actionDuration := float64(peacefulRoadEndTime.Sub(peacefulRoadStartTime))
	elapsedDuration := float64(now.Truncate(time.Hour * 24).Sub(peacefulRoadStartTime))

	percentForNow := elapsedDuration/actionDuration*(peacefulRoadEndPercentage-peacefulRoadStartPercentage) + peacefulRoadStartPercentage

	var statusMessage string
	if status.PercentageValue > percentForNow {
		statusMessage = "You're on the peaceful road! 🕊"
	} else {
		statusMessage = "Keep on working 🛠"
	}

	return fmt.Sprintf(`
	
<i>Peaceful Road</i>
Status for Today: %.2f%%
%s`,
		percentForNow,
		statusMessage)
}
