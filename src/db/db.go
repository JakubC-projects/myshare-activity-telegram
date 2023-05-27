package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	_ "google.golang.org/api/option"
)

const userCollectionName = "users"

var Db *firestore.Client
var Users *firestore.CollectionRef

func init() {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "myshare-telegram-notifications"}
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
}
