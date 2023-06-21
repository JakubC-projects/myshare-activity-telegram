package api

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
)

type orgsResponse response[[]models.Org]

func GetOrgs(ctx context.Context, u models.User) ([]models.Org, error) {
	url := fmt.Sprintf("%s/Profile/Organisations", baseUrl1)
	res, err := get[orgsResponse](ctx, u, url)
	return res.Data, err
}
