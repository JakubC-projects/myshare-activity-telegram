package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
)

func GetActivity(ctx context.Context, id int) (models.ContributionsActivity, error) {
	var res models.ContributionsActivity
	doc, err := Activities.Doc(fmt.Sprint(id)).Get(ctx)
	if err != nil {
		return res, fmt.Errorf("cannot fetch user from the database: %w", err)
	}
	err = doc.DataTo(&res)
	return res, err
}

func GetLatestActivity(ctx context.Context) (models.ContributionsActivity, error) {
	var res models.ContributionsActivity

	docIter := Activities.OrderBy("Created", firestore.Desc).Limit(1).Documents(ctx)
	doc, err := docIter.Next()
	if err != nil {
		return res, fmt.Errorf("cannot get latest activity: %w", err)
	}
	err = doc.DataTo(&res)
	return res, err
}

func SaveActivity(ctx context.Context, a models.ContributionsActivity) error {
	_, err := Activities.Doc(fmt.Sprint(a.Id)).Set(ctx, a)
	return err
}
