package telegram_ctrl

import (
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleTelegramUpdate(u tgbotapi.Update) {
	chatId := u.Message.Chat.ID
	telegram.SendWelcomeMessage(chatId)
	fmt.Printf("%+v", u)
}
