package api

import (
	"context"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
)

type teamsResponse response[struct {
	UserGroups []models.Team `json:"userGroups"`
}]

func GetTeams(ctx context.Context, u models.User) ([]models.Team, error) {
	res, err := get[teamsResponse](ctx, u, "/Profile/Clubs")
	return res.Data.UserGroups, err
}
