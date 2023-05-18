package auth

import (
	"fmt"
	"net/http"

	"github.com/JakubC-projects/myshare-activity-telegram/src/auth"
	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/samber/lo"
)

func callbackHandler(auth *auth.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		session := sessions.Default(c)
		if c.Query("state") != session.Get("state") {
			c.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx, c.Query("code"))
		if err != nil {
			c.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
			return
		}

		idToken := token.Extra("id_token").(string)
		idTokenClaims, err := jwt.Parse([]byte(idToken))
		if err != nil {
			c.String(http.StatusUnauthorized, "failed to parse id token")
			return
		}

		chatId, ok := session.Get("chatId").(int64)
		if !ok {
			c.String(http.StatusUnauthorized, "Missing state: chatId")
			return
		}
		user, err := db.GetUser(ctx, chatId)
		if err != nil {
			c.String(http.StatusNotFound, "cannot find user")
			return
		}

		user.Token = token
		user.DisplayName = lo.Must(idTokenClaims.Get("name")).(string)
		if err = db.SaveUser(ctx, user); err != nil {
			c.String(http.StatusInternalServerError, "cannot update user")
			return
		}

		telegram.SendLoggedInMessage(user)

		// Redirect to logged in page.
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://t.me/%s", config.Get().Telegram.BotName))
	}
}
