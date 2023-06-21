package api

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
)

func GetActivities(ctx context.Context, u models.User) ([]models.Activity, error) {
	url := fmt.Sprintf("%s/Activities/AvailableActivities?teamId=%d&groupId=%d", baseUrl1, u.Org.TeamId, u.Org.Id)
	res, err := get[response[[]models.Activity]](ctx, u, url)
	return res.Data, err
}

func GetActivity(ctx context.Context, u models.User, activityId int) (models.Activity, error) {
	url := fmt.Sprintf("%s/Activities/Details/%d", baseUrl1, activityId)
	res, err := get[response[models.Activity]](ctx, u, url)
	return res.Data, err
}

func GetParticipants(ctx context.Context, u models.User, activityId int) ([]models.Participant, error) {
	url := fmt.Sprintf("%s/Registrations/%d/Participants", baseUrl2, activityId)
	res, err := get[response[struct{ Participants []models.Participant }]](ctx, u, url)
	return res.Data.Participants, err
}

func RegisterActivity(ctx context.Context, u models.User, registration models.Registration) (models.Registration, error) {
	url := fmt.Sprintf("%s/Registrations/", baseUrl2)
	res, err := post[response[models.Registration]](ctx, u, url, registration)
	return res.Data, err
}

func UnregisterActivity(ctx context.Context, u models.User, registration models.Registration) error {
	url := fmt.Sprintf("%s/Registrations/%d/Unregistration/%d", baseUrl2, registration.RegistrationId, registration.ActivityId)
	fmt.Println(url)
	_, err := post[response[any]](ctx, u, url, nil)
	return err
}
