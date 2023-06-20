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

func handleUserAction(ctx context.Context, user models.User, u tgbotapi.Update) error {
	callback := u.CallbackData()
	if callback != "" {
		command, payload, _ := strings.Cut(callback, "-")
		switch command {
		case models.CommandStartChangeOrg:
			return startChangeOrg(ctx, user)
		case models.CommandChangeOrg:
			return changeOrg(ctx, user, payload)
		case models.CommandShowActivities:
			return showActivities(ctx, user)
		case models.CommandShowActivity:
			return showActivity(ctx, user, payload)

		case models.CommandShowStatus:
			return showStatus(ctx, user)
		case models.CommandShowMenu:
			return showMenu(ctx, user)
		case models.CommandLogout:
			return logoutUser(ctx, user)
		}
	}

	msg, err := telegram.SendMenuMessage(user)
	if err != nil {
		return fmt.Errorf("cannot send menu message: %w", err)
	}
	user.LastMessageId = msg.MessageID

	if err = db.SaveUser(ctx, user); err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}
	return nil
}

func showMenu(ctx context.Context, u models.User) error {
	_, err := telegram.SendMenuMessage(u, telegram.Edit)
	if err != nil {
		return fmt.Errorf("cannot send menu message: %w", err)
	}
	return nil
}

func logoutUser(ctx context.Context, u models.User) error {
	msg, err := telegram.SendWelcomeMessage(u)
	if err != nil {
		return fmt.Errorf("cannot send welcome message: %w", err)
	}

	u.Token = nil
	u.Org = nil
	u.LastMessageId = msg.MessageID

	if err = db.SaveUser(ctx, u); err != nil {
		return fmt.Errorf("cannot remove user session: %w", err)
	}
	return nil
}
