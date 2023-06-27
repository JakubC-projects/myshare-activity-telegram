package apim2m

import (
	"context"
	"net/http"
	"net/url"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var client = getClient()

func getClient() *http.Client {
	creds := clientcredentials.Config{
		ClientID:       config.Get().Oauth.ClientID,
		ClientSecret:   config.Get().Oauth.ClientSecret,
		TokenURL:       "https://" + config.Get().Oauth.Domain + "/oauth/token",
		AuthStyle:      oauth2.AuthStyleInHeader,
		Scopes:         []string{"contributions#read"},
		EndpointParams: url.Values{"audience": []string{config.Get().ContributionsAPI.Audience}},
	}
	return creds.Client(context.Background())
}
