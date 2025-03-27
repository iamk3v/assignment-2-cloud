package database

import (
	"assignment-2/config"
	"assignment-2/utils"
	"errors"
	"google.golang.org/api/iterator"
	"log"
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
func UpdateRegistration(id string, dash utils.Dashboard) error {

	// Overwrite the document
	_, err := config.Client.Collection(collection).Doc(id).Set(config.Ctx, dash)
	if err != nil {
		log.Println("Error updating document with id: " + id + ": " + err.Error())
		return err
	}
	return nil
}

func GetOneRegistration(id string) (utils.Dashboard, error) {
	// Find the document with specified id
	res := config.Client.Collection(collection).Doc(id)

	// Get the document
	doc, err := res.Get(config.Ctx)

	if err != nil {
		log.Println("Error extracting body of returned document of dashboard " + id + ": " + err.Error())
		return utils.Dashboard{}, err
	}

	// Convert the firebase document into a dashboard struct
	var dashboard utils.Dashboard
	err = doc.DataTo(&dashboard)
	if err != nil {
		return utils.Dashboard{}, err
	}

	dashboard.Id = doc.Ref.ID

	return dashboard, nil
}

func GetAllRegistrations() ([]utils.Dashboard, error) {
	// Iterator through all documents
	iter := config.Client.Collection(collection).Documents(config.Ctx)
	var allDashboards []utils.Dashboard

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Println("Error iterating dashboards collection: " + err.Error())
			return nil, err
		}

		// Convert the firebase document into a dashboard struct
		var dashboard utils.Dashboard
		err = doc.DataTo(&dashboard)
		if err != nil {
			return nil, err
		}

		dashboard.Id = doc.Ref.ID
		// Append the document to list
		allDashboards = append(allDashboards, dashboard)
	}

	// Return all documents
	return allDashboards, nil
}
