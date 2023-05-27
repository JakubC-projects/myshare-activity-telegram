package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/JakubC-projects/myshare-activity-telegram/src/api"
	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/log"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

func TelegramUpdateHttpHandler(c *gin.Context) {
	var update tgbotapi.Update
	err := c.ShouldBind(&update)
	if err != nil {
		log.L.Err(err).Send()
	}
	err = HandleUpdate(c.Request.Context(), update)

	if err != nil {
		log.L.Err(err).Msg("Telegram update error")
	}

	c.Status(200)
}

func HandleUpdate(ctx context.Context, u tgbotapi.Update) error {
	log.L.Debug().Interface("update", u).Msg("Received message")

	chatId, err := getChatIdFromUpdate(u)
	if err != nil {
		return err
	}
	user, err := db.GetOrCreateUser(ctx, chatId)
	if err != nil {
		return fmt.Errorf("cannot get user: %w", err)
	}

	if user.Token == nil {
		handleWelcomeMessage(ctx, user)
		return nil
	}

	if user.Team == nil && u.CallbackQuery != nil {
		err := handleSetTeam(ctx, user, u.CallbackQuery)
		if err != nil {
			return fmt.Errorf("cannot set team: %w", err)
		}
		return nil
	}

	_, err = telegram.SendMenuMessage(user)
	if err != nil {
		return fmt.Errorf("cannot send menu message: %w", err)
	}
	return nil
}

func handleWelcomeMessage(ctx context.Context, user models.User) error {
	msg, err := telegram.SendWelcomeMessage(user)
	if err != nil {
		return fmt.Errorf("cannot send welcome message: %w", err)
	}
	user.LastMessageId = msg.MessageID
	err = db.SaveUser(ctx, user)
	if err != nil {
		return fmt.Errorf("cannot save user: %w", err)
	}
	return nil
}

func handleSetTeam(ctx context.Context, user models.User, callback *tgbotapi.CallbackQuery) error {
	if callback.Message.MessageID != user.LastMessageId {
		return fmt.Errorf("invalid callback Id: expected: %d, found: %d", user.LastMessageId, callback.Message.MessageID)
	}

	userTeams, err := api.GetTeams(ctx, user)
	if err != nil {
		return fmt.Errorf("cannot get teams :%w", err)
	}

	selectedTeamId := lo.Must(strconv.Atoi(callback.Data))

	selectedTeam, found := lo.Find(userTeams, func(t models.Team) bool {
		return t.TeamId == selectedTeamId
	})
	if !found {
		return fmt.Errorf("selected invalid team: %d", selectedTeamId)
	}

	user.Team = &selectedTeam
	err = db.SaveUser(ctx, user)
	if err != nil {
		return fmt.Errorf("cannot save user: %w", err)

	}
	_, err = telegram.SendMenuMessage(user, telegram.Edit)
	if err != nil {
		return fmt.Errorf("cannot send menu message: %w", err)
	}

	return nil
}
