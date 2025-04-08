package database

import (
	"assignment-2/config"
	"assignment-2/utils"
	"errors"
	"google.golang.org/api/iterator"
	"log"
)

/*
AddRegistration Adds a specific registration in Firestore by ID
*/
func AddRegistration(dash utils.DashboardPost) (string, error) {
	// Adding the document to firestore
	ref, _, err := Client.Collection(config.DASHBOARD_COLLECTION).Add(Ctx, dash)
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
	_, err := Client.Collection(config.DASHBOARD_COLLECTION).Doc(id).Delete(Ctx)
	if err != nil {
		log.Println("Error deleting document with id " + id + ": " + err.Error())
		return err
	}

	// Return nil as there was no error
	return nil
}

/*
UpdateRegistration Updates a specific registration in Firestore by ID
*/
func UpdateRegistration(id string, dash utils.DashboardPost) error {

	// Overwrite the document
	_, err := Client.Collection(config.DASHBOARD_COLLECTION).Doc(id).Set(Ctx, dash)
	if err != nil {
		log.Println("Error updating document with id: " + id + ": " + err.Error())
		return err
	}
	return nil
}

/*
GetOneRegistration Gets a specific registration in Firestore by ID
*/
func GetOneRegistration(id string) (*utils.Dashboard, error) {
	// Find the document with specified id
	res := Client.Collection(config.DASHBOARD_COLLECTION).Doc(id)

	// Get the document
	doc, err := res.Get(Ctx)
	if err != nil {
		log.Println("hei")
		log.Println("Error extracting body of returned document of dashboard " + id + ": " + err.Error())
		return nil, err
	}

	// Convert the firebase document into a dashboard struct
	var dashboard utils.Dashboard
	err = doc.DataTo(&dashboard)
	if err != nil {
		return nil, err
	}

	dashboard.Id = doc.Ref.ID

	return &dashboard, nil
}

/*
GetAllRegistrations Gets all currently stored registrations from Firestore
*/
func GetAllRegistrations() ([]utils.Dashboard, error) {
	// Iterator through all documents
	iter := Client.Collection(config.DASHBOARD_COLLECTION).Documents(Ctx)
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
