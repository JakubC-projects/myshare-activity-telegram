package main

import (
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/controllers/auth"
	"github.com/JakubC-projects/myshare-activity-telegram/src/controllers/telegram"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/telegram-update", telegram.TelegramUpdateHttpHandler)

	auth.AddRoutes(r)

	r.Run(fmt.Sprintf(":%d", config.Get().Server.Port))
}
