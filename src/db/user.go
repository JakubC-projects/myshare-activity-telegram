package db

import (
	"context"
	"fmt"

	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetUser(ctx context.Context, chatId int64) (models.User, error) {
	var res models.User
	doc, err := Users.Doc(fmt.Sprint(chatId)).Get(ctx)
	if err != nil {
		return res, fmt.Errorf("cannot fetch user from the database: %w", err)
	}
	err = doc.DataTo(&res)
	return res, err
}

func GetUsersToNotify(ctx context.Context, teamId int) ([]models.User, error) {
	docs, err := Users.Where("Org.TeamId", "==", teamId).Where("NotificationsSettings.Enabled", "==", true).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch user from the database: %w", err)
	}

	res := make([]models.User, len(docs))
	for i, doc := range docs {
		err = doc.DataTo(&res[i])
		if err != nil {
			return nil, fmt.Errorf("cannot parse user: %w", err)
		}
	}
	return res, nil
}

func GetOrCreateUser(ctx context.Context, chatId int64) (models.User, error) {
	user, err := GetUser(ctx, chatId)
	if err == nil {
		return user, nil
	}
	if status.Code(err) != codes.NotFound {
		return user, fmt.Errorf("error getting user: %w", err)
	}
	user = models.User{ChatId: chatId}
	err = SaveUser(ctx, user)
	if err != nil {
		return user, fmt.Errorf("cannot save user: %w", err)
	}
	return user, err
}

func SaveUser(ctx context.Context, user models.User) error {
	_, err := Users.Doc(fmt.Sprint(user.ChatId)).Set(ctx, user)
	return err
}
