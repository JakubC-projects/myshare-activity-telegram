package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getChatIdFromUpdate(u tgbotapi.Update) (int64, error) {
	if u.Message != nil && u.Message.Chat != nil {
		return u.Message.Chat.ID, nil
	}

	if u.CallbackQuery != nil && u.CallbackQuery.Message != nil && u.CallbackQuery.Message.Chat != nil {
		return u.CallbackQuery.Message.Chat.ID, nil
	}

	return 0, errors.New("cannot read chat id")
}
