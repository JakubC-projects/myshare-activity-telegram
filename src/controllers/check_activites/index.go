package checkactivites

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JakubC-projects/myshare-activity-telegram/src/api"
	"github.com/JakubC-projects/myshare-activity-telegram/src/apim2m"
	"github.com/JakubC-projects/myshare-activity-telegram/src/db"
	"github.com/JakubC-projects/myshare-activity-telegram/src/log"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckActivitiesHandler(c *gin.Context) {
	if err := handleActivitiesCheck(c.Request.Context()); err != nil {
		log.L.Err(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func handleActivitiesCheck(ctx context.Context) error {
	newActivities, err := getNewActivities(ctx)
	if err != nil {
		return fmt.Errorf("cannot get new activities: %w", err)
	}

	teamActivities := lo.GroupBy(newActivities, func(a models.ContributionsActivity) int {
		return a.TeamId
	})
	for team, activities := range teamActivities {
		err := notifyTeamOfActivities(ctx, team, activities)
		if err != nil {
			return fmt.Errorf("cannot sync activities for team %d: %w", team, err)
		}
	}
	for _, a := range newActivities {
		err = db.SaveActivity(ctx, a)
		if err != nil {
			return fmt.Errorf("cannot save activity to db: %w", err)
		}
	}
	return nil
}

func getNewActivities(ctx context.Context) ([]models.ContributionsActivity, error) {
	refreshTime := time.Now().Add(-100 * time.Hour)

	freshActivities, err := apim2m.GetActivities(ctx, apim2m.ActivitiesQueryParams{
		Filter: &apim2m.ActivitiesFilter{
			Created: &apim2m.Filter[time.Time]{
				Gt: &refreshTime,
			},
			Start: &apim2m.Filter[time.Time]{
				Gt: lo.ToPtr(time.Now()),
			},
		},
	})

	var newActivities []models.ContributionsActivity
	for _, a := range freshActivities {
		_, err := db.GetActivity(ctx, a.Id)
		if status.Code(err) == codes.NotFound {
			newActivities = append(newActivities, a)
		} else if err != nil {
			return nil, fmt.Errorf("cannot check if activity already exists: %w", err)
		}
	}

	return newActivities, err
}

func notifyTeamOfActivities(ctx context.Context, teamId int, activities []models.ContributionsActivity) error {
	users, err := db.GetUsersToNotify(ctx, teamId)
	if err != nil {
		return fmt.Errorf("cannot fetch team users: %w", err)
	}
	for _, u := range users {
		err := notifyUserOfActivities(ctx, u, activities)
		if err != nil {
			return fmt.Errorf("cannot send message to user: %w", err)
		}
	}

	return nil
}

func notifyUserOfActivities(ctx context.Context, u models.User, newActivities []models.ContributionsActivity) error {
	userActivities, err := api.GetActivities(ctx, u)
	if err != nil {
		return fmt.Errorf("cannot get activity :%w", err)
	}
	newActivitiesForUser := lo.Filter(userActivities, func(ua models.MyshareActivity, _ int) bool {
		return lo.ContainsBy(newActivities, func(na models.ContributionsActivity) bool {
			return na.Id == ua.ActivityId
		})
	})
	msg, err := telegram.SendNewActivitiesNotificationMessage(u, newActivitiesForUser)
	u.LastMessageId = msg.MessageID
	db.SaveUser(ctx, u)
	return err
}
