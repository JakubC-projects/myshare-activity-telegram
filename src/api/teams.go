package api

import (
	"context"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
)

type orgsResponse response[[]models.Org]

func GetOrgs(ctx context.Context, u models.User) ([]models.Org, error) {
	res, err := get[orgsResponse](ctx, u, "/Profile/Organisations")
	return res.Data, err
}
