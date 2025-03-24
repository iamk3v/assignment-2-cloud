package main

import (
	"assignment-2/config"
	"assignment-2/handlers"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

func main() {

	// Instantiate global Ctx variable
	config.Ctx = context.Background()

	// get the credentials from file
	sa := option.WithCredentialsFile("./config/service-account.json")

	// Create new app
	app, err := firebase.NewApp(config.Ctx, nil, sa)
	if err != nil {
		log.Println("Error initializing app: " + err.Error())
		return
	}

	// Instantiate client
	config.Client, err = app.Firestore(config.Ctx)
	if err != nil {
		log.Println("Error instantiating app: " + err.Error())
		return
	}

	// Create a new router
	router := http.NewServeMux()

	// Routes
	router.HandleFunc(config.START_URL+"/registrations/", handlers.RegistrationHandler)
	router.HandleFunc(config.START_URL+"/dashboards/", handlers.DashboardHandler)
	router.HandleFunc(config.START_URL+"/notifications/", handlers.NotificationHandler)
	router.HandleFunc(config.START_URL+"/status/", handlers.StatusHandler)

	// Define port
	PORT := "8080"
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}

	// Listen on the designated port for traffic
	log.Println("Starting server on port " + PORT)
	err = http.ListenAndServe(":"+PORT, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
