package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/JakubC-projects/myshare-activity-telegram/src/auth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type queryParams struct {
	ChatId int64 `form:"chatId"`
}

// Handler for our login.
func loginHandler(auth *auth.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state, err := generateRandomState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		var qParams queryParams

		err = ctx.BindQuery(&qParams)
		if err != nil {
			errMsg := fmt.Sprintf("invalid query parameters: %s", err.Error())
			ctx.String(http.StatusBadRequest, errMsg)
			return
		}

		fmt.Println("params", qParams)

		// Save the state inside the session.
		session := sessions.Default(ctx)
		session.Set("state", state)
		session.Set("chatId", qParams.ChatId)

		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
	}
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
