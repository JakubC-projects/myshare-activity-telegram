package telegram

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/api"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
)

func showStatus(ctx context.Context, u models.User, editedMessageId int) error {
	userStatus, err := api.GetStatus(ctx, u)
	if err != nil {
		return fmt.Errorf("cannot get status :%w", err)
	}
	_, err = telegram.SendShowStatusMessage(u, userStatus, editedMessageId)
	return err
}
