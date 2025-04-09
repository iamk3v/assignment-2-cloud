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
RegistrationHandler Handles requests sent to the /registrations endpoint, routing the request to
corresponding handle functions based on http methods and if ID was sent or not.
*/
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	basePath := config.START_URL + "/registrations/"
	trimmedPath := strings.TrimPrefix(r.URL.Path, basePath)
	parts := strings.Split(trimmedPath, "/")

	if len(parts) == 1 && parts[0] != "" {
		id := parts[0]
		// ID provided
		switch r.Method {
		case http.MethodGet:
			handleRegGetOneRequest(w, r, id)
		case http.MethodDelete:
			handleRegDeleteRequest(w, r, id)
		case http.MethodPut:
			handleRegPutRequest(w, r, id)
		case http.MethodPatch:
			handleRegPatchRequest(w, r, id)
		case http.MethodHead:
			handleRegHeadRequest(w, r, id)
		default:
			http.Error(w,
				fmt.Sprintf("Method %s not supported on /registrations/{id}", r.Method),
				http.StatusMethodNotAllowed)
			return
		}
	} else {
		// No ID provided
		switch r.Method {
		case http.MethodGet:
			handleRegGetAllRequest(w, r)
		case http.MethodPost:
			handleRegPostRequest(w, r)
		case http.MethodHead:
			handleRegHeadRequest(w, r, "")
		default:
			http.Error(w,
				fmt.Sprintf("Method %s not supported on /registrations/", r.Method),
				http.StatusMethodNotAllowed)
			return
		}
	}
}

/*
handleRegGetOneRequest Gets one registration from the dashboards database based on ID
*/
func handleRegGetOneRequest(w http.ResponseWriter, r *http.Request, id string) {
	rawContent, err := database.GetOneRegistration(id)
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
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	_, err = fmt.Fprintln(w, string(content))
	if err != nil {
		log.Println("Error while writing response body: " + err.Error())
		http.Error(w, "There was an error while writing response body", http.StatusInternalServerError)
		return
	}
}

/*
handleRegGetAllRequest Gets all registrations from the dashboards database
*/
func handleRegGetAllRequest(w http.ResponseWriter, r *http.Request) {
	// If no ID was provided, get all
	rawContent, err := database.GetAllRegistrations()
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
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	_, err = fmt.Fprintln(w, string(content))
	if err != nil {
		log.Println("Error while writing response body: " + err.Error())
		http.Error(w, "There was am error while writing response body", http.StatusInternalServerError)
		return
	}

}

/*
handleRegPostRequest Adds a new registration to the dashboards database
*/
func handleRegPostRequest(w http.ResponseWriter, r *http.Request) {

	// Read the body
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

	dashboard := utils.DashboardPost{}
	// Decode JSON into the dashboard struct
	err = json.Unmarshal(content, &dashboard)
	if err != nil {
		log.Println("Error unmarshalling payload: " + err.Error())
		http.Error(w, "There was an error unmarshalling payload", http.StatusInternalServerError)
		return
	}
	dashboard.LastChange = time.Now().Local().String()

	// Add the dashboard to DB
	id, err := database.AddRegistration(dashboard)
	if err != nil {
		log.Println("Error adding dashboard to database: " + err.Error())
		http.Error(w, "There was an error adding dashboard", http.StatusInternalServerError)
		return
	}

	// Retrieve the new registration and trigger webhooks for REGISTER
	if webhookTrigger != nil {
		newReg, err := database.GetOneRegistration(id)
		if err != nil {
			log.Println("Error retrieving newly added registration: ", err)
		} else if newReg != nil {
			webhookTrigger.TriggerWebhooks("REGISTER", newReg.IsoCode)

		}
	}

	// Create the response struct
	resp := map[string]string{
		"id":         id,
		"lastChange": dashboard.LastChange,
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Error encoding AddRegistration response: ", err)
	}
}

/*
handleRegPutRequest Overwrites an existing registration in the dashboard based on provided ID
*/
func handleRegPutRequest(w http.ResponseWriter, r *http.Request, id string) {
	// If an ID was not provided
	if id == "" {
		http.Error(w, "An ID is required to update a specific dashboard registration", http.StatusBadRequest)
		return
	}

	// Read the body
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

	dashboard := utils.DashboardPost{}
	// Decode JSON into the dashboard struct
	err = json.Unmarshal(content, &dashboard)
	if err != nil {
		log.Println("Error unmarshalling payload: " + err.Error())
		http.Error(w, "There was an error unmarshalling payload", http.StatusInternalServerError)
		return
	}

	// Update timestamp
	dashboard.LastChange = time.Now().Local().String()

	err = database.UpdateRegistration(id, dashboard)
	if err != nil {
		http.Error(w, "Could not update dashboard with id: "+id, http.StatusInternalServerError)
		return
	}

	updateReg, err := database.GetOneRegistration(id)
	if err != nil {
		log.Println("Error retrieving updated registration", err)
	} else if updateReg != nil {
		if webhookTrigger != nil {
			webhookTrigger.TriggerWebhooks("CHANGE", updateReg.IsoCode)
		}
	}

	// Return status code to indicate success
	w.WriteHeader(http.StatusNoContent)
}

/*
handleRegDeleteRequest Deletes an existing registration from the dashboard database based on ID
*/
func handleRegDeleteRequest(w http.ResponseWriter, r *http.Request, id string) {
	// If an ID was not provided
	if id == "" {
		http.Error(w, "An ID is required to delete a specific dashboard registration", http.StatusBadRequest)
		return
	}

	// Retrieve registrations to be deleted (for iso code)
	existingReg, err := database.GetOneRegistration(id)
	if err != nil {
		log.Println("Error retrieving registration with id ", id, ": ", err)
		http.Error(w, "Registration was not found", http.StatusNotFound)
		return
	}

	// Try to delete
	err = database.DeleteRegistration(id)
	if err != nil {
		http.Error(w, "There was an error trying to delete that dashboard..", http.StatusInternalServerError)
		return
	}

	// Trigger webhook for delete event
	if webhookTrigger != nil {
		webhookTrigger.TriggerWebhooks("DELETE", existingReg.IsoCode)
	}
	w.WriteHeader(http.StatusNoContent)
}

/*
handleRegPatchRequest Modifies only provided fields in an existing registration in the dashboards database
*/
func handleRegPatchRequest(w http.ResponseWriter, r *http.Request, id string) {
	// If an ID was not provided
	if id == "" {
		http.Error(w, "An ID is required to update a specific dashboard registration", http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Read the body
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

	// Decode Body into an indexable map
	var patchData map[string]interface{}
	err = json.Unmarshal(content, &patchData)
	if err != nil {
		log.Println("Error unmarshalling payload: " + err.Error())
		http.Error(w, "There was an error unmarshalling payload", http.StatusInternalServerError)
		return
	}

	// Get the registration from firebase
	dbData, err := database.GetOneRegistration(id)
	if err != nil {
		log.Println("Error retrieving registration with id " + id + ": " + err.Error())
		http.Error(w, "There was an error retrieving registration with id "+id, http.StatusNotFound)
		return
	}

	// Extract original data into a indexable map
	dbDataJson, err := json.Marshal(dbData)
	if err != nil {
		log.Println("Error marshalling payload: " + err.Error())
		http.Error(w, "There was an error patching registration", http.StatusInternalServerError)
		return
	}

	var originalData map[string]interface{}
	err = json.Unmarshal(dbDataJson, &originalData)
	if err != nil {
		log.Println("Error unmarshalling payload: " + err.Error())
		http.Error(w, "There was an error patching registration", http.StatusInternalServerError)
		return
	}

	// Merge country and isoCode values from the incoming patch with original data
	// Country - If patchData contains a non-empty country, that value is used,
	// otherwise the existing country value is used
	mergedCountry := ""
	if patchVal, ok := patchData["country"].(string); ok && patchVal != "" {
		mergedCountry = patchVal
	} else if origVal, ok := originalData["country"].(string); ok {
		mergedCountry = origVal
	}

	// Also done for iso-code
	mergedIsoCode := ""
	if patchVal, ok := patchData["isoCode"].(string); ok && patchVal != "" {
		mergedIsoCode = patchVal
	} else if origVal, ok := originalData["isoCode"].(string); ok {
		mergedIsoCode = origVal
	}

	// If both are empty, an error is returned.
	if mergedCountry == "" && mergedIsoCode == "" {
		http.Error(w, "Both country code and isoCode cannot be empty", http.StatusBadRequest)
		return
	}

	// Save the original values into the originalData map
	originalData["country"] = mergedCountry
	originalData["isoCode"] = mergedIsoCode

	// Check if incoming patch contains a features field
	if patchVal, exists := patchData["features"]; exists {
		// attempt to assert that the value of features is a map
		patchFeatures, ok := patchVal.(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid format for features", http.StatusBadRequest)
			return
		}
		// check if original data has a features field
		if origVal, exists := originalData["features"]; exists {
			if origFeatures, ok := origVal.(map[string]interface{}); ok {
				for key, value := range patchFeatures {
					origFeatures[key] = value
				}
				originalData["features"] = origFeatures
			} else {
				originalData["features"] = patchFeatures
			}
		} else {
			originalData["features"] = patchFeatures
		}
	}

	// Update timestamp
	originalData["lastChange"] = time.Now().Local().String()

	originalDataJson, err := json.Marshal(originalData)
	if err != nil {
		log.Println("Error marshalling payload: " + err.Error())
		http.Error(w, "There was an error patching registration", http.StatusInternalServerError)
		return
	}

	var updatedData utils.DashboardPost
	err = json.Unmarshal(originalDataJson, &updatedData)
	//Patch with request body
	if err != nil {
		log.Println("Error unmarshalling updated data: " + err.Error())
		http.Error(w, "Error processing updated data", http.StatusInternalServerError)
		return
	}

	err = database.UpdateRegistration(id, updatedData)
	if err != nil {
		http.Error(w, "Could not patch registration with id: "+id+"\nMake sure all fields are valid fields", http.StatusInternalServerError)
		return
	}

	// Trigger Webhook
	if webhookTrigger != nil {
		webhookTrigger.TriggerWebhooks("CHANGE", updatedData.IsoCode)
	}
	// Return status code to indicate success
	w.WriteHeader(http.StatusNoContent)
}

/*
handleRegHeadRequest Provides only the headers for a registration or all registrations in the dashboard database
*/
func handleRegHeadRequest(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" { // No ID provided
		// Get all registrations
		rawContent, err := database.GetAllRegistrations()
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
			return
		}

		// Set and send back headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
		w.WriteHeader(http.StatusNoContent)

	} else { // ID provided
		// Get one registration
		rawContent, err := database.GetOneRegistration(id)
		if err != nil {
			log.Println("Error retrieving registration with id " + id + ": " + err.Error())
			http.Error(w, "There was an error getting the dashboard with id: "+id, http.StatusNotFound)
			return
		}

		// Encode response
		content, err := json.Marshal(rawContent)
		if err != nil {
			log.Println("Error marshalling payload: " + err.Error())
			http.Error(w, "There was an error marshalling payload", http.StatusInternalServerError)
			return
		}

		// Set and send back headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
		w.WriteHeader(http.StatusNoContent)
	}
}
