package api

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
)

func GetStatus(ctx context.Context, u models.User) (models.Status, error) {
	url := fmt.Sprintf("%s/TargetStatus/%d/Member/%d", baseUrl1, u.Org.Id, u.PersonID)
	res, err := get[response[models.Status]](ctx, u, url)
	return res.Data, err
}
