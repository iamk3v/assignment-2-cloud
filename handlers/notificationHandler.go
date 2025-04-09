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
	// Split the URL path based on the notifications endpoint to extract potential ID.
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, config.START_URL+"/notifications"), "/")

	// Check if there is an ID present (second element in the split path).
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
			// Return an error if the HTTP method is not supported.
			http.Error(w,
				fmt.Sprintf("Method %s not supported on /notifications/{id}", r.Method),
				http.StatusMethodNotAllowed)
		}
	} else {
		// No specific ID is provided; operate on the entire collection.
		switch r.Method {
		case http.MethodGet:
			handleNotiGetAllRequest(w, r)
		case http.MethodPost:
			handleNotiPostRequest(w, r)
		case http.MethodHead:
			handleNotiHeadRequest(w, r, "")
		default:
			// Return an error if the HTTP method is not supported.
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
	// Retrieve all registered webhooks from the database.
	hooks, err := database.GetAllWebhooks()
	if err != nil {
		// Log the error and return a 500 Internal Server Error if retrieval fails.
		log.Println("Error retrieving webhooks: " + err.Error())
		http.Error(w, config.ERR_INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}
	// Set the response header to indicate that the content is in JSON format.
	w.Header().Set("Content-Type", "application/json")
	// Encode the slice of webhooks into JSON and write it to the response.
	json.NewEncoder(w).Encode(hooks)
}

/*
handleNotiGetOneRequest handles GET requests for a specific webhook registration.
It fetches the webhook identified by the provided id and returns it as JSON.
*/
func handleNotiGetOneRequest(w http.ResponseWriter, r *http.Request, id string) {
	// Retrieve the webhook associated with the given ID.
	hook, err := database.GetWebhook(id)
	if err != nil {
		// Log the error and return a 404 Not Found if the webhook does not exist.
		log.Println("Error retrieving webhook: " + err.Error())
		http.Error(w, config.ERR_NOT_FOUND, http.StatusNotFound)
		return
	}
	// Set the Content-Type to JSON for the response.
	w.Header().Set("Content-Type", "application/json")
	// Encode the webhook data into JSON and send it in the response.
	json.NewEncoder(w).Encode(hook)
}

/*
handleNotiPostRequest handles POST requests to create a new webhook registration.
It decodes the webhook from the request body, stores it in the database, and returns
the new webhook's ID along with an HTTP cat URL as JSON.
*/
func handleNotiPostRequest(w http.ResponseWriter, r *http.Request) {
	var hook utils.Webhook
	// Decode the incoming JSON body into a webhook struct.
	if err := json.NewDecoder(r.Body).Decode(&hook); err != nil {
		// Log and respond with an error if JSON decoding fails.
		log.Println("Error decoding webhook body: " + err.Error())
		http.Error(w, config.ERR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	// Create the webhook entry in the database.
	id, err := database.CreateWebhook(hook)
	if err != nil {
		// Log and respond with an error if creation fails.
		log.Println("Error creating webhook: " + err.Error())
		http.Error(w, config.ERR_INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}
	// Prepare the response payload containing the new webhook's ID and an HTTP cat image URL.
	resp := map[string]string{
		"id":      id,
		"httpCat": "https://http.cat/201",
	}
	// Set the Content-Type header and the HTTP status code indicating resource creation.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Encode the response map into JSON and write it to the response.
	json.NewEncoder(w).Encode(resp)
}

/*
handleNotiDeleteRequest handles DELETE requests to remove a webhook registration.
It deletes the webhook identified by id from the database and returns a 204 No Content status.
*/
func handleNotiDeleteRequest(w http.ResponseWriter, r *http.Request, id string) {
	// Attempt to delete the webhook from the database.
	err := database.DeleteWebhook(id)
	if err != nil {
		// If deletion fails, log the error and return a 404 Not Found status.
		log.Println("Error deleting webhook: " + err.Error())
		http.Error(w, config.ERR_NOT_FOUND, http.StatusNotFound)
		return
	}
	// Log the successful deletion and return a 204 No Content response.
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

	// Ensure that the request contains a body; otherwise, return an error.
	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Read the request body.
	content, err := io.ReadAll(r.Body)
	if err != nil {
		// Log and respond with an error if reading the payload fails.
		log.Println("Error reading request body: " + err.Error())
		http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
		return
	}
	// Ensure the request body is closed after reading.
	defer r.Body.Close()

	// Check if the payload is empty and return an error if it is.
	if len(content) == 0 {
		http.Error(w, "Your payload appears to be empty.", http.StatusBadRequest)
		return
	}

	// Decode the JSON payload into a map for patching.
	var patchData map[string]interface{}
	err = json.Unmarshal(content, &patchData)
	if err != nil {
		// Log and return an error if JSON unmarshalling fails.
		log.Println("Error unmarshalling payload: " + err.Error())
		http.Error(w, "There was an error unmarshalling payload", http.StatusInternalServerError)
		return
	}

	// Retrieve the existing webhook from Firestore.
	existingHook, err := database.GetWebhook(id)
	if err != nil {
		// Log error and notify the client if the webhook cannot be retrieved.
		log.Println("Error retrieving webhook with id " + id + ": " + err.Error())
		http.Error(w, "Error retrieving webhook with id "+id, http.StatusInternalServerError)
		return
	}

	// Marshal the existing webhook to JSON, then unmarshal into a map.
	existingJSON, err := json.Marshal(existingHook)
	if err != nil {
		// Log and return an error if marshalling fails.
		log.Println("Error marshalling existing webhook: " + err.Error())
		http.Error(w, "Error patching webhook", http.StatusInternalServerError)
		return
	}

	// Convert the JSON back into a map so that it can be merged with the patch data.
	var originalData map[string]interface{}
	err = json.Unmarshal(existingJSON, &originalData)
	if err != nil {
		// Log and return an error if the conversion fails.
		log.Println("Error unmarshalling existing webhook: " + err.Error())
		http.Error(w, "Error patching webhook", http.StatusInternalServerError)
		return
	}

	// Merge the patch data into the original data.
	for key, value := range patchData {
		originalData[key] = value
	}

	// Update the "lastChange" field with the current local timestamp.
	originalData["lastChange"] = time.Now().Local().String()

	// Apply the updated data map to the webhook document in the database.
	err = database.UpdateWebhook(id, originalData)
	if err != nil {
		http.Error(w, "Could not patch webhook with id: "+id, http.StatusInternalServerError)
		return
	}

	// Respond with a 204 No Content status to indicate the patch was successful.
	w.WriteHeader(http.StatusNoContent)
}

/*
handleNotiHeadRequest processes HEAD requests for webhook registrations.
For a specific webhook (by ID), it retrieves the document and returns only the headers.
If no ID is provided, it applies to the entire collection of webhooks.
*/
func handleNotiHeadRequest(w http.ResponseWriter, r *http.Request, id string) {

	if id == "" { // No ID provided
		// Attempt to retrieve all webhooks to ensure the collection is accessible.
		_, err := database.GetAllWebhooks()
		if err != nil {
			// Log error and respond with an internal server error if retrieval fails.
			log.Println("Error retrieving webhooks: " + err.Error())
			http.Error(w, config.ERR_INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
			return
		}

		// Set the header to indicate JSON content, and return no content.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

	} else { // Specific webhook ID provided.
		// Attempt to retrieve the specific webhook.
		_, err := database.GetWebhook(id)
		if err != nil {
			// Log error and return a not found status if the webhook does not exist.
			log.Println("Error retrieving webhook: " + err.Error())
			http.Error(w, config.ERR_NOT_FOUND, http.StatusNotFound)
			return
		}

		// Set the response header to JSON and return a no content status.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}
