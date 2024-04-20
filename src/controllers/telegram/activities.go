package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/JakubC-projects/myshare-activity-telegram/src/api"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/JakubC-projects/myshare-activity-telegram/src/telegram"
)

func showActivities(ctx context.Context, u models.User, pagePayload string, editedMessageId int) error {
	page := 0
	if pagePayload != "" {
		var err error
		page, err = strconv.Atoi(pagePayload)
		if err != nil {
			return fmt.Errorf("cannot parse page payload")
		}
	}

	availableActivities, err := api.GetActivities(ctx, u)
	if err != nil {
		return fmt.Errorf("cannot get activities :%w", err)
	}
	_, err = telegram.SendShowActivitiesMessage(u, availableActivities, page, editedMessageId)
	return err
}

func ShowActivity(ctx context.Context, u models.User, activityIdString string, editedMessageId int) error {
	activityId, err := strconv.Atoi(activityIdString)
	if err != nil {
		return fmt.Errorf("invalid activity id %s: %w", activityIdString, err)
	}
	activity, err := api.GetActivity(ctx, u, activityId)
	if err != nil {
		return fmt.Errorf("cannot get activity :%w", err)
	}
	_, err = telegram.SendShowActivityMessage(u, activity, editedMessageId)
	return err
}

func showParticipants(ctx context.Context, u models.User, activityIdString string, editedMessageId int) error {
	activityId, err := strconv.Atoi(activityIdString)
	if err != nil {
		return fmt.Errorf("invalid activity id %s: %w", activityIdString, err)
	}
	activity, err := api.GetActivity(ctx, u, activityId)
	if err != nil {
		return fmt.Errorf("cannot get activity :%w", err)
	}
	participants, err := api.GetParticipants(ctx, u, activityId)
	if err != nil {
		return fmt.Errorf("cannot get participants :%w", err)
	}
	_, err = telegram.SendShowShowParticipantsMessage(u, activity, participants, editedMessageId)

	return err
}

func registerActivity(ctx context.Context, u models.User, activityIdString string, editedMessageId int) error {
	activityId, err := strconv.Atoi(activityIdString)
	if err != nil {
		return fmt.Errorf("invalid activity id %s: %w", activityIdString, err)
	}

	registrationReq := models.Registration{
		ActivityId:          activityId,
		UserId:              u.PersonID,
		IsSwipeRegistration: false,
		ShowComments:        false,
	}
	_, err = api.RegisterActivity(ctx, u, registrationReq)
	if err != nil {
		return fmt.Errorf("cannot register activity :%w", err)
	}
	activity, err := api.GetActivity(ctx, u, activityId)
	if err != nil {
		return fmt.Errorf("cannot get activity :%w", err)
	}
	_, err = telegram.SendShowActivityMessage(u, activity, editedMessageId)

	return err
}

func unregisterActivity(ctx context.Context, u models.User, activityIdString string, editedMessageId int) error {
	activityId, err := strconv.Atoi(activityIdString)
	if err != nil {
		return fmt.Errorf("invalid activity id %s: %w", activityIdString, err)
	}
	activity, err := api.GetActivity(ctx, u, activityId)
	if err != nil {
		return fmt.Errorf("cannot get activity :%w", err)
	}
	if activity.RegistrationId == 0 {
		return fmt.Errorf("missing registration id")
	}

	registrationReq := models.Registration{
		RegistrationId: activity.RegistrationId,
		ActivityId:     activityId,
	}
	err = api.UnregisterActivity(ctx, u, registrationReq)
	if err != nil {
		return fmt.Errorf("cannot unregister activity :%w", err)
	}
	activity, err = api.GetActivity(ctx, u, activityId)
	if err != nil {
		return fmt.Errorf("cannot get activity :%w", err)
	}
	_, err = telegram.SendShowActivityMessage(u, activity, editedMessageId)

	return err
}
