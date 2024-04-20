package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ensureUserIsConfigured(ctx context.Context, user models.User, u tgbotapi.Update) (bool, error) {
	if user.Token == nil {
		return false, handleWelcomeMessage(ctx, user)
	}

	if user.Org == nil {
		if u.CallbackQuery == nil {
			return false, startChangeOrg(ctx, user, 0)
		}
		command, payload, _ := strings.Cut(u.CallbackQuery.Data, "-")
		if command != models.CommandChangeOrg {
			return false, startChangeOrg(ctx, user, u.CallbackQuery.Message.MessageID)
		}
		err := changeOrg(ctx, user, payload, u.CallbackQuery.Message.MessageID)
		if err != nil {
			return false, fmt.Errorf("cannot set org: %w", err)
		}
		return false, nil
	}
	return true, nil
}

func handleWelcomeMessage(ctx context.Context, user models.User) error {
	_, err := telegram.SendWelcomeMessage(user, 0)
	if err != nil {
		return fmt.Errorf("cannot send welcome message: %w", err)
	}
	return nil
}
