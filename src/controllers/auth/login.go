package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/JakubC-projects/myshare-activity-telegram/src/auth"
	"github.com/JakubC-projects/myshare-activity-telegram/src/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type queryParams struct {
	ChatId int64 `form:"chatId"`
}

// Handler for our login.
func loginHandler(ctx *gin.Context) {
	state, err := generateRandomState()
	if err != nil {
		errCtx := fmt.Errorf("cannot generate random state: %w", err)
		log.L.Err(err).Send()
		ctx.String(http.StatusInternalServerError, errCtx.Error())
		return
	}

	var qParams queryParams

	err = ctx.BindQuery(&qParams)
	if err != nil {
		errCtx := fmt.Errorf("invalid query parameters: %w", err)
		log.L.Err(err).Send()

		ctx.String(http.StatusBadRequest, errCtx.Error())
		return
	}

	// Save the state inside the session.
	session := sessions.Default(ctx)
	session.Set("state", state)
	session.Set("chatId", qParams.ChatId)

	if err := session.Save(); err != nil {
		errCtx := fmt.Errorf("cannot save session: %w", err)
		log.L.Err(err).Send()

		ctx.String(http.StatusInternalServerError, errCtx.Error())
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, auth.Auth.AuthCodeURL(state))
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
