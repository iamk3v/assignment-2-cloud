package main

import (
	"assignment-2/clients"
	"assignment-2/config"
	"assignment-2/database"
	_ "assignment-2/database"
	"assignment-2/handlers"
	"assignment-2/services"
	"assignment-2/utils"
	"log"
	"net/http"
	"os"
	"time"
)

/*
main The main function of the service which starts it up and routes endpoint requests to corresponding handlers
*/
func main() {

	//start uptime timer
	utils.StartTime()
	log.Println("Uptime timer started:", utils.GetTime())

	// Set the webhook trigger implementation
	database.SetDBWebhookTrigger(services.WebhookService{})
	clients.SetClientWebhookTrigger(services.WebhookService{})
	handlers.SetHandlerWebhookTrigger(services.WebhookService{})

	// Purge cached entries at startup
	if err := database.PurgeExpiredCacheEntries(database.Ctx); err != nil {
		log.Printf("Error purging expired cache entries at startup: %v\n", err)
	} else {
		log.Println("Successfully purged expired cache entries at startup")
	}

	// STARTING background routine for purging expired cache - Checks every hour
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			<-ticker.C
			err := database.PurgeExpiredCacheEntries(database.Ctx)
			if err != nil {
				log.Printf("Error purging expired cache entries: %v\n", err)
			} else {
				log.Println("Expired cache entries purged successfully.")
			}
		}
	}()

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
