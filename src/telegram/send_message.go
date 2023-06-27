package telegram

import (
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

type Option interface {
	IsOption()
}

type editMessage struct{}

func (editMessage) IsOption() {}

var Edit = editMessage{}

func isEdit(opts []Option) bool {
	_, foundEdit := lo.Find(opts, func(p Option) bool {
		_, ok := p.(editMessage)
		return ok
	})
	return foundEdit
}

func sendMessage(user models.User, text string, replyMarkup tgbotapi.InlineKeyboardMarkup, opts ...Option) (tgbotapi.Message, error) {
	if isEdit(opts) {
		msg := tgbotapi.NewEditMessageTextAndMarkup(user.ChatId, user.LastMessageId, text, replyMarkup)
		msg.ParseMode = tgbotapi.ModeHTML
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
