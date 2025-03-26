package database

import (
	"assignment-2/config"
	"assignment-2/utils"
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
func AddRegistration(dash utils.Dashboard) (string, error) {
	// Adding the document to firestore
	ref, _, err := config.Client.Collection(collection).Add(config.Ctx, dash)
	if err != nil {
		log.Println("Error adding document to Firestore: " + err.Error())
		return "", err
	}
	// If nothing went wrong
	return ref.ID, nil
}

/*
DeleteRegistration Deletes a specific registration in Firestore by ID.
*/
func DeleteRegistration(id string) error {
	_, err := config.Client.Collection(collection).Doc(id).Delete(config.Ctx)
	if err != nil {
		log.Println("Error deleting document with id " + id + ": " + err.Error())
		return err
	}

	// Return nil as there was no error
	return nil
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
		log.Println("Error updating document with id: ", id, "->", err)
		http.Error(w, "Could not update registration: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Return status code to indicate success
	w.WriteHeader(http.StatusNoContent)

}

func GetOneRegistration(id string) (map[string]interface{}, error) {
	// Find the document with specified id
	res := config.Client.Collection(collection).Doc(id)

	// Get the document
	doc, err := res.Get(config.Ctx)

	if err != nil {
		log.Println("Error extracting body of returned document of dashboard " + id + ": " + err.Error())
		return nil, err
	}

	return doc.Data(), nil
}

func GetAllRegistrations() ([]map[string]interface{}, error) {
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
