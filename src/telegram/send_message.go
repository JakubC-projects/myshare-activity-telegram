package telegram

import (
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessage(user models.User, text string, replyMarkup *tgbotapi.InlineKeyboardMarkup, editedMessageId int) (tgbotapi.Message, error) {
	if editedMessageId != 0 {
		msg := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      user.ChatId,
				MessageID:   editedMessageId,
				ReplyMarkup: replyMarkup,
			},
			Text:      text,
			ParseMode: tgbotapi.ModeHTML,
		}
		return Bot.Send(msg)
	}

	return Bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      user.ChatId,
			ReplyMarkup: replyMarkup,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	})
}

func DeleteMessage(user models.User, messageId int) (tgbotapi.Message, error) {
	return Bot.Send(tgbotapi.NewDeleteMessage(user.ChatId, messageId))
}
