package services

import (
	"assignment-2/database"
	"assignment-2/utils"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

/*
TriggerWebhooks Checks for registered webhooks that match the given event country and sends
a post notification
*/
func TriggerWebhooks(event string, country string) {
	// Retrieve all webhooks from the database
	hooks, err := database.GetAllWebhooks()
	if err != nil {
		log.Println("Error retrieving webhooks: ", err)
		return
	}

	// Looping through all webhooks
	for _, hook := range hooks {
		// Check if the webhook event matches
		if hook.Event == event && (hook.Country == "" || hook.Country == country) {
			// Create the payload for the webhook invocation plus a time stamp
			payload := utils.WebhookInvocation{
				ID:      hook.ID,
				Country: hook.Country,
				Event:   event,
				Time:    time.Now().Format(time.RFC3339),
			}

			// Payload to JSON
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				log.Println("Error marshalling webhook payload: ", err)
				continue
			}

			// Trigger the webhooks asynchronously
			go func(url string, data []byte) {
				// Send a post request with the webhook URL
				resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
				if err != nil {
					log.Println("Error triggering webhook at: ", url, ":", err)
					return
				}

				// Log it with the HTTP status code
				defer resp.Body.Close()
				log.Println("Webhook triggered at: %s, status code: %s", url, resp.StatusCode)
			}(hook.URL, payloadBytes)
		}
	}
}
