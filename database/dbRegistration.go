package database

import (
	"assignment-2/config"
	"assignment-2/utils"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"errors"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"time"
)

const collection = "dashboards"

/*
Addregistration reads JSON from the request body, stores it in firebase and writes a JSON response (ID + lastchange)
 */
func AddRegistration(w http.ResponseWriter, r *http.Request) {
	var dash utils.Dashboard
	// Parse JSON body into Dashboard struct
	if err := json.NewDecoder(r.Body).Decode(&dash); err != nil {
		log.Println("Error decoding JSON body: "err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// Update timestamp
	dash.LastChange = time.Now()
	// Adding the document to firestore
	ref, _, err := config.Client.Collection(collection).Add(config.Ctx, dash)
	if err != nil {
		log.Println("Error adding document to Firestore: ", err)
		http.Error(w, "Could not add registration: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Response Object - ID and LastChange
	resp := map[string]string{
		"id": ref.ID,
		"lastChange": dash.LastChange.Format(time.RFC3339),
	}
	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Error encoding AddRegistration response: ", err)
	}
}

/*
DeleteRegistration Deletes a specific registration in Firestore by ID.
 */
func DeleteRegistration(id string, w http.ResponseWriter, r *http.Request) {
	_, err := config.Client.Collection(collection).Doc(id).Delete(config.Ctx)
	if err != nil {
		log.Println("Error deleting document with id: ", id,"->" ,err)
		http.Error(w, "Could not delete registration: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Return status to indicate successful deletion
	w.WriteHeader(http.StatusNoContent)
}

/*
UpdateRegistration replaces an existing registration document with the JSON from the request body
 */
func UpdateRegistration(id string, w http.ResponseWriter, r *http.Request) {
	var dash utils.Dashboard

	if err := json.NewDecoder(r.Body).Decode(&dash); err != nil {
		log.Println("Error decoding JSON body for update: ", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Update timestamp
	dash.LastChange = time.Now()

	// Overwrite the document
	_, err := config.Client.Collection(collection).Doc(id).Set(config.Ctx, dash)
	if err != nil {
		log.Println("Error updating document with id: ", id, "->" ,err)
		http.Error(w, "Could not update registration: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Return status code to indicate success
	w.WriteHeader(http.StatusNoContent)

}

func GetOneRegistration(id string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	// Find the document with specified id
	res := config.Client.Collection(collection).Doc(id)

	// Get the document
	doc, err := res.Get(config.Ctx)

	if err != nil {
		log.Println("Error extracting body of returned document of dashboard " + id + ": " + err.Error())
		http.Error(w, "There was an error getting the document with id "+id, http.StatusInternalServerError)
		return nil, err
	}

	return doc.Data(), nil
}

func GetAllRegistrations(w http.ResponseWriter, r *http.Request) ([]map[string]interface{}, error) {
	// Iterator through all documents
	iter := config.Client.Collection(collection).Documents(config.Ctx)
	var allDocs []map[string]interface{}

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Println("Error iterating dashboards collection: " + err.Error())
			return nil, err
		}

		// Append the document to list
		allDocs = append(allDocs, doc.Data())
	}

	// Return all documents
	return allDocs, nil
}
