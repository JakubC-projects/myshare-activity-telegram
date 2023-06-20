package api

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
)

func GetActivities(ctx context.Context, u models.User) ([]models.Activity, error) {
	url := fmt.Sprintf("/Activities/AvailableActivities?teamId=%d&groupId=%d", u.Org.TeamId, u.Org.Id)
	res, err := get[response[[]models.Activity]](ctx, u, url)
	return res.Data, err
}

func GetActivity(ctx context.Context, u models.User, activityId int) (models.Activity, error) {
	url := fmt.Sprintf("/Activities/Details/%d", activityId)
	res, err := get[response[models.Activity]](ctx, u, url)
	return res.Data, err
}
