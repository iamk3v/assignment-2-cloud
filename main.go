package main

import (
	"assignment-2/config"
	"assignment-2/database"
	_ "assignment-2/database"
	"assignment-2/handlers"
	"assignment-2/utils"
	"log"
	"net/http"
	"os"
)

/*
main The main function of the service which starts it up and routes endpoint requests to corresponding handlers
*/
func main() {

	//start uptime timer
	utils.StartTime()
	log.Println("Uptime timer started:", utils.GetTime())

	// Create a new router
	router := http.NewServeMux()

	// Routes
	router.HandleFunc(config.START_URL+"/registrations/", handlers.RegistrationHandler)
	router.HandleFunc(config.START_URL+"/registrations", handlers.RegistrationHandler)
	router.HandleFunc(config.START_URL+"/dashboards/", handlers.DashboardHandler)
	router.HandleFunc(config.START_URL+"/dashboards", handlers.DashboardHandler)
	router.HandleFunc(config.START_URL+"/notifications/", handlers.NotificationHandler)
	router.HandleFunc(config.START_URL+"/notifications", handlers.NotificationHandler)
	router.HandleFunc(config.START_URL+"/status/", handlers.StatusHandler)
	router.HandleFunc(config.START_URL+"/status", handlers.StatusHandler)

	//Handle all 404 if no match found
	router.HandleFunc("/", handlers.NotFoundHandler)

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

	// Close the client when service shuts down
	defer func() {
		errClose := database.Client.Close()
		if errClose != nil {
			log.Fatal("Closing of the Firebase client failed. Error: " + errClose.Error())
		}
	}()
}
