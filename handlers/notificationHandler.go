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
		default:
			http.Error(w,
				fmt.Sprintf("Method %s not supported on /notifications/", r.Method),
				http.StatusMethodNotAllowed)
		}
	}
}

func handleNotiGetAllRequest(w http.ResponseWriter, r *http.Request) {
	hooks, err := database.GetAllWebhooks()
	if err != nil {
		log.Println("Error retrieving webhooks:", err)
		http.Error(w, config.ERR_INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hooks)
}

func handleNotiGetOneRequest(w http.ResponseWriter, r *http.Request, id string) {
	hook, err := database.GetWebhook(id)
	if err != nil {
		log.Println("Error retrieving webhook:", err)
		http.Error(w, config.ERR_NOT_FOUND, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hook)
}

func handleNotiPostRequest(w http.ResponseWriter, r *http.Request) {
	var hook utils.Webhook
	if err := json.NewDecoder(r.Body).Decode(&hook); err != nil {
		log.Println("Error decoding webhook body:", err)
		http.Error(w, config.ERR_BAD_REQUEST, http.StatusBadRequest)
		return
	}
	id, err := database.CreateWebhook(hook)
	if err != nil {
		log.Println("Error creating webhook:", err)
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

func handleNotiDeleteRequest(w http.ResponseWriter, r *http.Request, id string) {
	err := database.DeleteWebhook(id)
	if err != nil {
		log.Println("Error deleting webhook:", err)
		http.Error(w, config.ERR_NOT_FOUND, http.StatusNotFound)
		return
	}
	log.Println("Deleted webhook:", id)
	w.WriteHeader(http.StatusNoContent)
}

// handleNotiPatchRequest processes PATCH requests to update a webhook partially.
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
		log.Println("Error reading request body:", err)
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
		log.Println("Error unmarshalling payload:", err)
		http.Error(w, "There was an error unmarshalling payload", http.StatusInternalServerError)
		return
	}

	// Retrieve the existing webhook from Firestore.
	existingHook, err := database.GetWebhook(id)
	if err != nil {
		log.Println("Error retrieving webhook with id", id, ":", err)
		http.Error(w, "Error retrieving webhook with id "+id, http.StatusInternalServerError)
		return
	}

	// Marshal the existing webhook to JSON, then unmarshal into a map.
	existingJSON, err := json.Marshal(existingHook)
	if err != nil {
		log.Println("Error marshalling existing webhook:", err)
		http.Error(w, "Error patching webhook", http.StatusInternalServerError)
		return
	}

	var originalData map[string]interface{}
	err = json.Unmarshal(existingJSON, &originalData)
	if err != nil {
		log.Println("Error unmarshalling existing webhook:", err)
		http.Error(w, "Error patching webhook", http.StatusInternalServerError)
		return
	}

	// Merge the patch data into the original data.
	for key, value := range patchData {
		originalData[key] = value
	}

	originalData["lastChange"] = time.Now().Format(time.RFC3339)

	// Update the webhook document in Firestore using the UpdateWebhook function.
	err = database.UpdateWebhook(id, originalData)
	if err != nil {
		http.Error(w, "Could not patch webhook with id: "+id, http.StatusInternalServerError)
		return
	}

	// Return no content to indicate success.
	w.WriteHeader(http.StatusNoContent)
}
