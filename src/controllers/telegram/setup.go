package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ensureUserIsConfigured(ctx context.Context, user models.User, u tgbotapi.Update) (bool, error) {
	if user.Token == nil {
		return false, handleWelcomeMessage(ctx, user)
	}

	if user.Org == nil {
		callback := u.CallbackData()
		if callback == "" {
			return false, startChangeOrg(ctx, user)
		}
		command, payload, _ := strings.Cut(callback, "-")
		if command != models.CommandChangeOrg {
			return false, startChangeOrg(ctx, user)
		}
		err := changeOrg(ctx, user, payload)
		if err != nil {
			return false, fmt.Errorf("cannot set org: %w", err)
		}
		return false, nil
	}
	return true, nil
}

func handleWelcomeMessage(ctx context.Context, user models.User) error {
	msg, err := telegram.SendWelcomeMessage(user)
	if err != nil {
		return fmt.Errorf("cannot send welcome message: %w", err)
	}
	user.LastMessageId = msg.MessageID
	err = db.SaveUser(ctx, user)
	if err != nil {
		return fmt.Errorf("cannot save user: %w", err)
	}
	return nil
}
