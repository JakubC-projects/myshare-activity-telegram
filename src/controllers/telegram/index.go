package telegram

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/log"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TelegramUpdateHttpHandler(c *gin.Context) {
	var update tgbotapi.Update
	err := c.ShouldBind(&update)
	if err != nil {
		log.L.Err(err).Send()
	}
	err = HandleUpdate(c.Request.Context(), update)

	if err != nil {
		log.L.Err(err).Msg("Telegram update error")
	}

	c.Status(200)
}

func HandleUpdate(ctx context.Context, u tgbotapi.Update) error {
	log.L.Debug().Interface("update", u).Msg("Received message")

	chatId, err := getChatIdFromUpdate(u)
	if err != nil {
		return err
	}

	user, err := db.GetOrCreateUser(ctx, chatId)
	if err != nil {
		return fmt.Errorf("cannot get user: %w", err)
	}
	configured, err := ensureUserIsConfigured(ctx, user, u)
	if !configured || err != nil {
		return err
	}

	return handleUserAction(ctx, user, u)
}
