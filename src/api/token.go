package api

import (
	"context"

	"github.com/JakubC-projects/myshare-activity-telegram/src/auth"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"golang.org/x/oauth2"
)

func getTokenSilently(ctx context.Context, u models.User) (*oauth2.Token, error) {
	auth := auth.GetAuthenticator()
	refresher := auth.TokenSource(ctx, u.Token)
	return refresher.Token()
}
