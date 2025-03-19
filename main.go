package main

import (
	"assignment-2/config"
	"assignment-2/handlers"
	"log"
	"net/http"
	"os"
)

func main() {

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
	err := http.ListenAndServe(":"+PORT, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
