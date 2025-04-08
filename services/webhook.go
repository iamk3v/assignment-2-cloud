package services

import (
	"assignment-2/database"
	"assignment-2/utils"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
		log.Println("Error retrieving webhooks: " + err.Error())
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
				Time:    time.Now().Local().String(),
			}

			// Payload to JSON
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				log.Println("Error marshalling webhook payload: " + err.Error())
				continue
			}

			// Trigger the webhooks asynchronously
			go func(url string, data []byte) {
				// Send a post request with the webhook URL
				resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
				if err != nil {
					log.Println("Error triggering webhook at: " + url + ": " + err.Error())
					return
				}

				// Log it with the HTTP status code
				defer resp.Body.Close()
				log.Println("Webhook triggered at: " + url + ", status code: " + strconv.Itoa(resp.StatusCode))
			}(hook.URL, payloadBytes)
		}
	}
}
