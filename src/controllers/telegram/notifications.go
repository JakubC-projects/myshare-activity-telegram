package telegram

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
)

func setNotificationSettings(ctx context.Context, u models.User, enabled bool) error {
	u.NotificationsSettings.Enabled = enabled
	err := db.SaveUser(ctx, u)
	if err != nil {
		return fmt.Errorf("cannot save user: %w", err)

	}
	_, err = telegram.SendMenuMessage(u, telegram.Edit)
	return err
}
