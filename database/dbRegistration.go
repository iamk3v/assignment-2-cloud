package database

import (
	"assignment-2/config"
	"errors"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
)

const collection = "dashboards"

func AddRegistration(w http.ResponseWriter, r *http.Request) {

}

func DeleteRegistration(id string, w http.ResponseWriter, r *http.Request) {

}

func UpdateRegistration(id string, w http.ResponseWriter, r *http.Request) {

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
