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
	if u.CallbackQuery != nil {
		command, payload, _ := strings.Cut(u.CallbackQuery.Data, "-")
		switch command {
		case models.CommandStartChangeOrg:
			return startChangeOrg(ctx, user, u.CallbackQuery.Message.MessageID)
		case models.CommandChangeOrg:
			return changeOrg(ctx, user, payload, u.CallbackQuery.Message.MessageID)
		case models.CommandShowActivities:
			return showActivities(ctx, user, payload, u.CallbackQuery.Message.MessageID)
		case models.CommandShowActivity:
			return ShowActivity(ctx, user, payload, u.CallbackQuery.Message.MessageID)
		case models.CommandShowParticipants:
			return showParticipants(ctx, user, payload, u.CallbackQuery.Message.MessageID)
		case models.CommandRegisterActivity:
			return registerActivity(ctx, user, payload, u.CallbackQuery.Message.MessageID)
		case models.CommandUnregisterActivity:
			return unregisterActivity(ctx, user, payload, u.CallbackQuery.Message.MessageID)
		case models.CommandShowStatus:
			return showStatus(ctx, user, u.CallbackQuery.Message.MessageID)
		case models.CommandShowMenu:
			return showMenu(ctx, user, u.CallbackQuery.Message.MessageID)
		case models.CommandEnableNotifications:
			return setNotificationSettings(ctx, user, true, u.CallbackQuery.Message.MessageID)
		case models.CommandDisableNotifications:
			return setNotificationSettings(ctx, user, false, u.CallbackQuery.Message.MessageID)
		case models.CommandLogout:
			return logoutUser(ctx, user, u.CallbackQuery.Message.MessageID)
		}
	}

	_, err := telegram.SendMenuMessage(user, 0)
	if err != nil {
		return fmt.Errorf("cannot send menu message: %w", err)
	}

	return nil
}

func showMenu(ctx context.Context, u models.User, editedMessageId int) error {
	_, err := telegram.SendMenuMessage(u, editedMessageId)
	if err != nil {
		return fmt.Errorf("cannot send menu message: %w", err)
	}
	return nil
}

func logoutUser(ctx context.Context, u models.User, editedMessageId int) error {
	u.Token = nil
	u.Org = nil

	if err := db.SaveUser(ctx, u); err != nil {
		return fmt.Errorf("cannot remove user session: %w", err)
	}

	_, err := telegram.DeleteMessage(u, editedMessageId)
	if err != nil {
		return fmt.Errorf("cannot send welcome message: %w", err)
	}

	return nil
}
