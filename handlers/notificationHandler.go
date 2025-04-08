package handlers

import (
	"assignment-2/config"
	"assignment-2/database"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
NotificationHandler handles requests to the /notifications endpoint.
It routes the request to the appropriate sub-handler based on whether an ID
is provided in the URL and the HTTP method used.
*/
func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	// /dashboard/v1/notifications/{id} or /dashboard/v1/notifications/
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, config.START_URL+"/notifications"), "/")

	if len(pathParts) > 1 && pathParts[1] != "" {
		id := pathParts[1]
		switch r.Method {
		case http.MethodGet:
			handleNotiGetOneRequest(w, r, id)
		case http.MethodDelete:
			handleNotiDeleteRequest(w, r, id)
		case http.MethodPatch:
			handleNotiPatchRequest(w, r, id)
		case http.MethodHead:
			handleNotiHeadRequest(w, r, id)
		default:
			http.Error(w,
				fmt.Sprintf("Method %s not supported on /notifications/{id}", r.Method),
				http.StatusMethodNotAllowed)
		}
	} else {
		// collection
		switch r.Method {
		case http.MethodGet:
			handleNotiGetAllRequest(w, r)
		case http.MethodPost:
			handleNotiPostRequest(w, r)
		case http.MethodHead:
			handleNotiHeadRequest(w, r, "")
		default:
			http.Error(w,
				fmt.Sprintf("Method %s not supported on /notifications/", r.Method),
				http.StatusMethodNotAllowed)
		}
	}
}

/*
handleNotiGetAllRequest handles GET requests to retrieve all webhook registrations.
It fetches all webhooks from the database and returns them in JSON format.
*/
func handleNotiGetAllRequest(w http.ResponseWriter, r *http.Request) {
	hooks, err := database.GetAllWebhooks()
	if err != nil {
		log.Println("Error retrieving webhooks: " + err.Error())
		http.Error(w, config.ERR_INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hooks)
}

/*
handleNotiGetOneRequest handles GET requests for a specific webhook registration.
It fetches the webhook identified by the provided id and returns it as JSON.
*/
func handleNotiGetOneRequest(w http.ResponseWriter, r *http.Request, id string) {
	hook, err := database.GetWebhook(id)
	if err != nil {
		log.Println("Error retrieving webhook: " + err.Error())
		http.Error(w, config.ERR_NOT_FOUND, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hook)
}

/*
handleNotiPostRequest handles POST requests to create a new webhook registration.
It decodes the webhook from the request body, stores it in the database, and returns
the new webhook's ID along with an HTTP cat URL as JSON.
*/
func handleNotiPostRequest(w http.ResponseWriter, r *http.Request) {
	var hook utils.Webhook
	if err := json.NewDecoder(r.Body).Decode(&hook); err != nil {
		log.Println("Error decoding webhook body: " + err.Error())
		http.Error(w, config.ERR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	id, err := database.CreateWebhook(hook)
	if err != nil {
		log.Println("Error creating webhook: " + err.Error())
		http.Error(w, config.ERR_INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}
	resp := map[string]string{
		"id":      id,
		"httpCat": "https://http.cat/201",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

/*
handleNotiDeleteRequest handles DELETE requests to remove a webhook registration.
It deletes the webhook identified by id from the database and returns a 204 No Content status.
*/
func handleNotiDeleteRequest(w http.ResponseWriter, r *http.Request, id string) {
	err := database.DeleteWebhook(id)
	if err != nil {
		log.Println("Error deleting webhook: " + err.Error())
		http.Error(w, config.ERR_NOT_FOUND, http.StatusNotFound)
		return
	}
	log.Println("Deleted webhook: " + id)
	w.WriteHeader(http.StatusNoContent)
}

/*
handleNotiPatchRequest processes PATCH requests to update a webhook registration partially.
It reads the request body, merges the provided patch data with the existing webhook data,
updates the lastChange timestamp, and writes the updated document to the database.
*/
func handleNotiPatchRequest(w http.ResponseWriter, r *http.Request, id string) {
	// Check if ID is provided
	if id == "" {
		http.Error(w, "An ID is required to update a specific webhook", http.StatusBadRequest)
		return
	}

	// Ensure the request body is not nil.
	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Read the request body.
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body: " + err.Error())
		http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(content) == 0 {
		http.Error(w, "Your payload appears to be empty.", http.StatusBadRequest)
		return
	}

	// Decode the patch data into a map.
	var patchData map[string]interface{}
	err = json.Unmarshal(content, &patchData)
	if err != nil {
		log.Println("Error unmarshalling payload: " + err.Error())
		http.Error(w, "There was an error unmarshalling payload", http.StatusInternalServerError)
		return
	}

	// Retrieve the existing webhook from Firestore.
	existingHook, err := database.GetWebhook(id)
	if err != nil {
		log.Println("Error retrieving webhook with id " + id + ": " + err.Error())
		http.Error(w, "Error retrieving webhook with id "+id, http.StatusInternalServerError)
		return
	}

	// Marshal the existing webhook to JSON, then unmarshal into a map.
	existingJSON, err := json.Marshal(existingHook)
	if err != nil {
		log.Println("Error marshalling existing webhook: " + err.Error())
		http.Error(w, "Error patching webhook", http.StatusInternalServerError)
		return
	}

	var originalData map[string]interface{}
	err = json.Unmarshal(existingJSON, &originalData)
	if err != nil {
		log.Println("Error unmarshalling existing webhook: " + err.Error())
		http.Error(w, "Error patching webhook", http.StatusInternalServerError)
		return
	}

	// Merge the patch data into the original data.
	for key, value := range patchData {
		originalData[key] = value
	}

	originalData["lastChange"] = time.Now().Local().String()

	// Update the webhook document in Firestore using the UpdateWebhook function.
	err = database.UpdateWebhook(id, originalData)
	if err != nil {
		http.Error(w, "Could not patch webhook with id: "+id, http.StatusInternalServerError)
		return
	}

	// Return no content to indicate success.
	w.WriteHeader(http.StatusNoContent)
}

/*
handleNotiHeadRequest processes HEAD requests for webhook registrations.
For a specific webhook (by ID), it retrieves the document and returns only the headers.
If no ID is provided, it applies to the entire collection of webhooks.
*/
func handleNotiHeadRequest(w http.ResponseWriter, r *http.Request, id string) {

	if id == "" { // No ID provided
		// Get all webhooks
		_, err := database.GetAllWebhooks()
		if err != nil {
			log.Println("Error retrieving webhooks: " + err.Error())
			http.Error(w, config.ERR_INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
			return
		}

		// Set and send back headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

	} else { // ID provided
		_, err := database.GetWebhook(id)
		if err != nil {
			log.Println("Error retrieving webhook: " + err.Error())
			http.Error(w, config.ERR_NOT_FOUND, http.StatusNotFound)
			return
		}

		// Set and send back headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}
