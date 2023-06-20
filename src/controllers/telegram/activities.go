package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/JakubC-projects/myshare-activity-telegram/src/api"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
)

func showActivities(ctx context.Context, u models.User) error {
	availableActivities, err := api.GetActivities(ctx, u)
	if err != nil {
		return fmt.Errorf("cannot get activities :%w", err)
	}
	_, err = telegram.SendShowActivitiesMessage(u, availableActivities, telegram.Edit)
	return err
}

func showActivity(ctx context.Context, u models.User, activityIdString string) error {
	activityId, err := strconv.Atoi(activityIdString)
	if err != nil {
		return fmt.Errorf("invalid activity id %s: %w", activityIdString, err)
	}
	activity, err := api.GetActivity(ctx, u, activityId)
	if err != nil {
		return fmt.Errorf("cannot get activity :%w", err)
	}
	_, err = telegram.SendShowActivityMessage(u, activity, telegram.Edit)
	return err
}
