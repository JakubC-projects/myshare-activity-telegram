package webhook

import (
	"context"

	"github.com/JakubC-projects/myshare-activity-telegram/src/log"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func SetupWebhook(token string) (ngrok.Tunnel, error) {
	tun, err := ngrok.Listen(context.Background(),
		config.HTTPEndpoint(),
		ngrok.WithAuthtoken(token),
	)

	if err != nil {
		return nil, err
	}

	log.L.Debug().Str("url", tun.URL()).Msg("Setup ngrok")

	wh, err := tgbotapi.NewWebhook(tun.URL() + "/telegram-update")
	if err != nil {
		return nil, err
	}
	res, err := telegram.Bot.Request(wh)
	if err != nil {
		return nil, err
	}
	log.L.Debug().Interface("response", res).Msg("Setup webhook in telegram")

	return tun, nil
}
