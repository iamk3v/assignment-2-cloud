package handlers

import (
	"assignment-2/config"
	"assignment-2/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleRegGetRequest(w, r)
	case http.MethodPost:
		handleRegPostRequest(w, r)
	case http.MethodDelete:
		handleRegDeleteRequest(w, r)
	case http.MethodPut:
		handleRegPutRequest(w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' not supported. "+
			"Currently only '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
		return
	}
}

func handleRegGetRequest(w http.ResponseWriter, r *http.Request) {
	basePath := config.START_URL + "/registrations/"
	trimmedPath := strings.TrimPrefix(r.URL.Path, basePath)
	parts := strings.Split(trimmedPath, "/")

	// If an ID was provided, get one
	log.Println(len(parts))
	if len(parts) >= 1 && parts[0] != "" {
		id := parts[0]
		rawContent, err := database.GetOneRegistration(id, w, r)
		if err != nil {
			log.Println("Error retrieving registration with id " + id + ": " + err.Error())
			http.Error(w, "There was an error getting the dashboard with id: "+id, http.StatusInternalServerError)
			return
		}

		// Encode response
		content, err := json.Marshal(rawContent)
		if err != nil {
			log.Println("Error marshalling payload: " + err.Error())
			http.Error(w, "There was an error marshalling payload", http.StatusInternalServerError)
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		_, err = fmt.Fprintln(w, string(content))
		if err != nil {
			log.Println("Error while writing response body: " + err.Error())
			http.Error(w, "There was am error while writing response body", http.StatusInternalServerError)
		}
	} else {
		// If no ID was provided, get all
		rawContent, err := database.GetAllRegistrations(w, r)
		if err != nil {
			log.Println("Error retrieving all dashboards: " + err.Error())
			http.Error(w, "There was an error retrieving all dashboards", http.StatusInternalServerError)
			return
		}

		// Encode response
		content, err := json.Marshal(rawContent)
		if err != nil {
			log.Println("Error marshalling payload: " + err.Error())
			http.Error(w, "There was an error marshalling payload", http.StatusInternalServerError)
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		_, err = fmt.Fprintln(w, string(content))
		if err != nil {
			log.Println("Error while writing response body: " + err.Error())
			http.Error(w, "There was am error while writing response body", http.StatusInternalServerError)
		}
	}

}

func handleRegPostRequest(w http.ResponseWriter, r *http.Request) {

}

func handleRegDeleteRequest(w http.ResponseWriter, r *http.Request) {

}

func handleRegPutRequest(w http.ResponseWriter, r *http.Request) {

}
