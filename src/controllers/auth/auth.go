package auth

import (
	"context"
	"errors"
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

// New instantiates the *Authenticator.
func getAuthenticator() (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+config.Get().Oauth.Domain+"/",
	)
	if err != nil {
		return nil, err
	}
	conf := oauth2.Config{
		ClientID:     config.Get().Oauth.ClientID,
		ClientSecret: config.Get().Oauth.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/callback", config.Get().Server.Host),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

// verifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) verifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
