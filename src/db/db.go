package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	_ "google.golang.org/api/option"
)

const userCollectionName = "users"
const activitiesCollectionName = "activities"

var Db *firestore.Client
var Users *firestore.CollectionRef
var Activities *firestore.CollectionRef

func init() {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: config.Get().Server.GcpProject}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		panic(fmt.Errorf("cannot create firebase app: %w", err))
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		panic(fmt.Errorf("cannot create firestore client: %w", err))
	}
	Db = client
	Users = Db.Collection(userCollectionName)
	Activities = Db.Collection(activitiesCollectionName)
}
