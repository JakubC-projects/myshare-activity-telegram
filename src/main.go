package main

import (
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/controllers/auth"
	telegram_ctrl "github.com/JakubC-projects/myshare-activity-telegram/src/controllers/telegram"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	r := gin.Default()
	fmt.Printf("%+v\n", config.Get())

	r.POST("/telegram-update", func(c *gin.Context) {
		var update tgbotapi.Update
		err := c.ShouldBind(&update)
		if err != nil {
			log.Err(err).Send()
		}
		telegram_ctrl.HandleTelegramUpdate(update)
		c.Status(200)
	})

	auth.AddRoutes(r)

	r.Run(fmt.Sprintf(":%d", config.Get().Server.Port))
}
