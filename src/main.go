package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/controllers/auth"
	"github.com/JakubC-projects/myshare-activity-telegram/src/controllers/telegram"
	"github.com/JakubC-projects/myshare-activity-telegram/src/log"
	"github.com/JakubC-projects/myshare-activity-telegram/src/webhook"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.POST("/telegram-update", telegram.TelegramUpdateHttpHandler)

	auth.AddRoutes(r)

	if _, found := os.LookupEnv("NGROK_AUTHTOKEN"); found {
		tun, err := webhook.SetupWebhook()
		if err != nil {
			log.L.Fatal().AnErr("err", err).Send()
		}
		go http.Serve(tun, r)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Get().Server.Port), r); err != nil {
		log.L.Fatal().AnErr("err", err).Send()
	}
}
