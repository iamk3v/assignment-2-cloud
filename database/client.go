package database

import (
	"assignment-2/config"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)

var Client *firestore.Client
var Ctx context.Context

func init() {
	var err error
	Client, err = initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
}

func initDatabase() (*firestore.Client, error) {
	Ctx = context.Background()
	// get the credentials from file
	sa := option.WithCredentialsFile("config/service-account.json")
	dbConfig := &firebase.Config{
		ProjectID: config.PROJECT_ID,
	}

	// Create new app
	app, err := firebase.NewApp(Ctx, dbConfig, sa)
	if err != nil {
		log.Println("Error initializing app: " + err.Error())
		return nil, err
	}

	// Instantiate client
	client, err := app.Firestore(Ctx)
	if err != nil {
		log.Println("Error instantiating app: " + err.Error())
		return nil, err
	}

	return client, nil
}
