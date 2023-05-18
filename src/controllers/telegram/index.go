package telegram

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func TelegramUpdateHttpHandler(c *gin.Context) {
	var update tgbotapi.Update
	err := c.ShouldBind(&update)
	if err != nil {
		log.Err(err).Send()
	}
	HandleUpdate(c.Request.Context(), update)
	c.Status(200)
}

func HandleUpdate(ctx context.Context, u tgbotapi.Update) {
	chatId, err := getChatIdFromUpdate(u)
	if err != nil {
		return
	}
	user, err := db.GetOrCreateUser(ctx, chatId)
	if err != nil {
		return
	}
	msg, err := telegram.SendWelcomeMessage(user)
	if err != nil {
		panic(err)
	}
	user.LastMessageId = msg.MessageID
	db.SaveUser(ctx, user)
	fmt.Printf("%+v", u)
}
