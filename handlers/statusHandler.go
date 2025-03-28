package handlers

import (
	"assignment-2/config"
	"assignment-2/database"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
Statushandler routes HTTP requests to the appropriate status method handler
*/
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleStatusGetRequest(w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' not supported. "+
			"Currently only '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
		return
	}
}

/*
handleStatusGetRequest checks the status of external APIs and Firestore collection
*/
func handleStatusGetRequest(w http.ResponseWriter, r *http.Request) {

	// Gets the urls and checks the APIs
	countriesAPIStatus := checkAPI(config.RESTCOUNTRIES_ROOT + "alpha/" + config.Testcountry + "/?fields=name")
	currencyAPIStatus := checkAPI(config.CURRENCY_ROOT + config.Testcurrency)
	openmeteoAPIStatus := checkAPI(config.OPENMETEO_ROOT + config.Testweather)

	// Check if we can access dashboards in Firestore
	dashStatusCode := http.StatusOK
	_, dashErr := database.GetAllRegistrations()
	if dashErr != nil {
		dashStatusCode = http.StatusInternalServerError
	}

	// Check if we can access webhooks in Firestore and count the registered webhooks
	notiStatusCode := http.StatusOK
	allHooks, hooksErr := database.GetAllWebhooks(database.Ctx, database.Client)
	if hooksErr != nil {
		notiStatusCode = http.StatusInternalServerError
	}
	totalHooks := len(allHooks)

	// Constructs JSON response using Statusresponse struct
	response := utils.Statusresponse{
		CountriesAPI:         countriesAPIStatus,
		CurrencyAPI:          currencyAPIStatus,
		OpenmeteoAPI:         openmeteoAPIStatus,
		Notificationresponse: notiStatusCode,
		Dashboardresponse:    dashStatusCode,
		Webhookssum:          totalHooks,
		Version:              config.VERSION,
		Uptime:               utils.Gettime(),
	}

	// Convert response to JSON and send to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//sends response to client
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Print("Error encoding status response: ", err)
		http.Error(w, config.ERR_INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
		return
	}
}

// sends a get request to the URL checking for a response
func checkAPI(apiurl string) int {
	resp, err := http.Get(apiurl)
	if err != nil {
		fmt.Println("error checking apiURL:", apiurl, err)
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
