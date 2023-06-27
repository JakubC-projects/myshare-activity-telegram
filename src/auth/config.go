package auth

import (
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var Auth = GetAuthenticator()

func GetAuthenticator() *oauth2.Config {

	endpoints := oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("https://%s/oauth/authorize?audience=%s", config.Get().Oauth.Domain, config.Get().MyshareAPI.Audience),
		TokenURL: fmt.Sprintf("https://%s/oauth/token", config.Get().Oauth.Domain),
	}

	conf := oauth2.Config{
		ClientID:     config.Get().Oauth.ClientID,
		ClientSecret: config.Get().Oauth.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/callback", config.Get().Server.Host),
		Endpoint:     endpoints,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "offline_access"},
	}

	return &conf
}
