package auth

import (
	"fmt"
	"net/http"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/log"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	"github.com/gin-gonic/gin"
)

// Handler for our login.
func logoutHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var qParams queryParams

	err := c.BindQuery(&qParams)
	if err != nil {
		errCtx := fmt.Errorf("invalid query parameters: %w", err)
		log.L.Err(err).Send()

		c.String(http.StatusBadRequest, errCtx.Error())
		return
	}

	u, err := db.GetUser(ctx, qParams.ChatId)
	if err != nil {
		errCtx := fmt.Errorf("cannot find user: %w", err)
		log.L.Err(err).Send()

		c.String(http.StatusBadRequest, errCtx.Error())
		return
	}

	u.Token = nil

	err = db.SaveUser(ctx, u)
	if err != nil {
		errCtx := fmt.Errorf("remove user session: %w", err)
		log.L.Err(err).Send()

		c.String(http.StatusBadRequest, errCtx.Error())
		return
	}

	telegram.SendWelcomeMessage(u, 0)

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://t.me/%s", config.Get().Telegram.BotName))
}
