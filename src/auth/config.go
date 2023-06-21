package auth

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

var Auth *Authenticator

func init() {
	a, err := GetAuthenticator()
	if err != nil {
		panic(err)
	}
	Auth = a
}

// New instantiates the *Authenticator.
func GetAuthenticator() (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+config.Get().Oauth.Domain+"/",
	)
	if err != nil {
		return nil, err
	}

	endpoints := provider.Endpoint()
	endpoints.AuthURL += "?audience=" + config.Get().MyshareAPI.Audience
	conf := oauth2.Config{
		ClientID:     config.Get().Oauth.ClientID,
		ClientSecret: config.Get().Oauth.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/callback", config.Get().Server.Host),
		Endpoint:     endpoints,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "offline_access"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}
