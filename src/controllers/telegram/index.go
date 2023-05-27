package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/JakubC-projects/myshare-activity-telegram/src/api"
	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func TelegramUpdateHttpHandler(c *gin.Context) {
	var update tgbotapi.Update
	err := c.ShouldBind(&update)
	if err != nil {
		log.Err(err).Send()
	}
	HandleUpdate(c.Request.Context(), update)
	c.Status(200)
}

func HandleUpdate(ctx context.Context, u tgbotapi.Update) {
	chatId, err := getChatIdFromUpdate(u)
	if err != nil {
		return
	}
	user, err := db.GetOrCreateUser(ctx, chatId)
	if err != nil {
		return
	}

	if user.Token == nil {
		handleWelcomeMessage(ctx, user)
		return
	}

	if user.Team == nil && u.CallbackQuery != nil {
		handleSetTeam(ctx, user, u.CallbackQuery)
		return
	}

	telegram.SendMenuMessage(user)

	fmt.Printf("%+v", u)
}

func handleWelcomeMessage(ctx context.Context, user models.User) {
	msg, err := telegram.SendWelcomeMessage(user)
	if err != nil {
		panic(err)
	}
	user.LastMessageId = msg.MessageID
	db.SaveUser(ctx, user)
}

func handleSetTeam(ctx context.Context, user models.User, callback *tgbotapi.CallbackQuery) {
	if callback.Message.MessageID != user.LastMessageId {
		panic("invalid callback")
	}

	userTeams, err := api.GetTeams(ctx, user)
	if err != nil {
		panic(err)
	}

	selectedTeamId := lo.Must(strconv.Atoi(callback.Data))

	selectedTeam, found := lo.Find(userTeams, func(t models.Team) bool {
		return t.TeamId == selectedTeamId
	})
	if !found {
		panic("invalid team selected")
	}

	user.Team = &selectedTeam
	err = db.SaveUser(ctx, user)
	if err != nil {
		panic(err)
	}
	_, err = telegram.SendMenuMessage(user)
	if err != nil {
		panic(err)
	}
}
