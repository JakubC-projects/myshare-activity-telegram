package main

import (
	"fmt"
	"net/http"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/controllers/auth"
	checkactivites "github.com/JakubC-projects/myshare-activity-telegram/src/controllers/check_activities"
	"github.com/JakubC-projects/myshare-activity-telegram/src/controllers/notify"
	"github.com/JakubC-projects/myshare-activity-telegram/src/controllers/telegram"
	"github.com/JakubC-projects/myshare-activity-telegram/src/log"
	"github.com/JakubC-projects/myshare-activity-telegram/src/webhook"
	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Get()

	r := gin.New()

	r.POST("/telegram-update", telegram.TelegramUpdateHttpHandler)
	r.POST("/check-activities", checkactivites.CheckActivitiesHandler)
	r.POST("/notify", notify.NotifyUsersHandler)
	auth.AddRoutes(r)

	if config.Ngrok.AuthToken != "" {
		tun, err := webhook.SetupWebhook(config.Ngrok.AuthToken, config.Ngrok.Domain)
		if err != nil {
			log.L.Fatal().AnErr("err", err).Send()
		}
		http.Serve(tun, r)
	} else {
		r.Run(fmt.Sprintf(":%d", config.Server.Port))
	}
}
