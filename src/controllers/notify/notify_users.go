package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/log"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/iterator"
)

type notifyParams struct {
	PersonID int `form:"personId"`
	TeamId   int `form:"teamId"`
}

type notification struct {
	Text           string
	KeyboardMarkup *tgbotapi.InlineKeyboardMarkup
}

func NotifyUsersHandler(c *gin.Context) {
	var np notifyParams
	var n notification

	apiKey := c.GetHeader("X-Api-Key")
	if apiKey != config.Get().Oauth.Secret {
		err := fmt.Errorf("invalid api key")
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	err := c.BindQuery(&np)
	if err != nil {
		err = fmt.Errorf("invalid query params: %w", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	fmt.Printf(" %+v\n", np)

	err = json.NewDecoder(c.Request.Body).Decode(&n)
	if err != nil {
		err = fmt.Errorf("invalid body: %w", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := handleNotifyUsers(c.Request.Context(), n, np); err != nil {
		err = fmt.Errorf("cannot notify users: %w", err)
		log.L.Err(err).Msg(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func handleNotifyUsers(ctx context.Context, n notification, np notifyParams) error {
	iter := getDbIterator(ctx, np)

	var u models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		err = doc.DataTo(&u)
		if err != nil {
			return fmt.Errorf("cannot map firebase data: %w", err)
		}
		err = sendNotificationToUser(ctx, u, n)
		if err != nil {
			return fmt.Errorf("cannot notify user %d: %w", u.PersonID, err)
		}
	}
	return nil
}

func sendNotificationToUser(ctx context.Context, u models.User, n notification) error {

	_, err := telegram.SendMessage(u, n.Text, n.KeyboardMarkup, 0)
	if err != nil {
		return fmt.Errorf("cannot send message: %w", err)
	}
	return nil
}

func getDbIterator(ctx context.Context, np notifyParams) *firestore.DocumentIterator {
	if np.PersonID != 0 {
		return db.Users.Where("PersonID", "==", np.PersonID).Documents(ctx)
	}
	if np.TeamId != 0 {
		return db.Users.Where("Org.TeamId", "==", np.TeamId).Documents(ctx)
	}
	return db.Users.Documents(ctx)
}
