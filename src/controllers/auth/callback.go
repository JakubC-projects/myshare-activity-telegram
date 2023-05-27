package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/JakubC-projects/myshare-activity-telegram/src/api"
	"github.com/JakubC-projects/myshare-activity-telegram/src/auth"
	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/oauth2"
)

func callbackHandler(c *gin.Context) {

	ctx := c.Request.Context()

	token, err := getCallbackToken(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
	}

	user, err := getCallbackUser(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
	}

	loggedInUser, err := updateUser(ctx, user, token)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
	}

	teams, err := api.GetTeams(ctx, loggedInUser)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	telegram.SendLoggedInMessage(loggedInUser, teams, telegram.Edit)

	// Redirect to logged in page.
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://t.me/%s", config.Get().Telegram.BotName))
}

func getCallbackToken(c *gin.Context) (*oauth2.Token, error) {
	session := sessions.Default(c)
	if c.Query("state") != session.Get("state") {
		return nil, errors.New("invalid state parameter")
	}

	token, err := auth.Auth.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		return nil, errors.New("failed to exchange an authorization code for a token")
	}
	return token, nil
}

func getCallbackUser(c *gin.Context) (models.User, error) {
	session := sessions.Default(c)
	chatId, ok := session.Get("chatId").(int64)
	if !ok {
		return models.User{}, errors.New("missing state: chatId")
	}

	user, err := db.GetUser(c.Request.Context(), chatId)
	if err != nil {
		return models.User{}, fmt.Errorf("cannot get user: %w", err)
	}
	return user, nil
}

func updateUser(ctx context.Context, user models.User, token *oauth2.Token) (models.User, error) {
	idToken := token.Extra("id_token").(string)
	idTokenClaims, err := jwt.Parse([]byte(idToken))
	if err != nil {
		return user, fmt.Errorf("cannot parse id token claims: %w", err)
	}

	userName, ok := idTokenClaims.Get("name")
	if !ok {
		return user, errors.New("cannot find name claim")
	}
	userNameString, ok := userName.(string)
	if !ok {
		return user, errors.New("invalid type of name claim ")
	}

	user.Token = token
	user.DisplayName = userNameString

	if err = db.SaveUser(ctx, user); err != nil {
		return user, fmt.Errorf("cannot save user: %w", err)
	}
	return user, nil
}
