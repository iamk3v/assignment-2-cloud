package handlers

import (
	"assignment-2/config"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

func handleStatusGetRequest(w http.ResponseWriter, r *http.Request) {

	//gets the urls and checks the APIs
	countriesAPIStatus := checkAPI(config.RESTCOUNTRIES_ROOT + "alpha/" + config.Testcountry + "/?fields=name")
	currencyAPIStatus := checkAPI(config.CURRENCY_ROOT + config.Testcurrency)
	openmeteoAPIStatus := checkAPI(config.OPENMETEO_ROOT + config.Testweather)

	//constructs JSON response
	response := utils.Statusresponse{
		CountriesAPI: countriesAPIStatus,
		CurrencyAPI:  currencyAPIStatus,
		OpenmeteoAPI: openmeteoAPIStatus,
		//Notificationresponse:,
		//Webhookssum:,
		Version: config.VERSION,
		Uptime:  utils.Gettime(),
	}

	// Convert response to JSON and send to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//sends response to client
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Print("error when trying to encode and send response to client" + err.Error())
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
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
