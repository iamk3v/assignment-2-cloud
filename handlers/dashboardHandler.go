package handlers

import (
	"assignment-2/config"
	"assignment-2/database"
	"assignment-2/utils"
	"encoding/asn1"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleDashGetRequest(w, r)
	default:
		http.Error(w, "REST method '"+r.Method+"' not supported. "+
			"Currently only '"+http.MethodGet+"' is supported.", http.StatusNotImplemented)
		return
	}
}

func handleDashGetRequest(w http.ResponseWriter, r *http.Request) {
	basePath := config.START_URL + "/dashboards/"
	trimmedPath := strings.TrimPrefix(r.URL.Path, basePath)
	parts := strings.Split(trimmedPath, "/")

	// If an ID was provided, get one
	if len(parts) >= 1 && parts[0] != "" {
		http.Error("dasboard ID not provided", http.StatusBadRequest)
		return
	}
	id := parts[0]

	rawContent, err := database.GetOneRegistration(id)
	if err != nil {
		log.Println("Error retrieving dashboard with id " + id + ": " + err.Error())
		http.Error(w, "There was an error getting the dashboard with id: "+id, http.StatusInternalServerError)
		return
	}

	// Extract fields from rawConfig (stored in Firestore)
	country := rawContent.Country
	isoCode := rawContent.IsoCode
	features := rawContent.Features
	lastsignin := rawContent.LastChange

	// Get country info
	countryData, err := database.GetCountryData(country, isoCode)
	if err != nil {
		http.Error(w, "Failed to fetch country data", http.StatusBadGateway)
		return
	}

	// Get weather info
	weatherData, err := database.GetWeatherDate(countryData.Latling[0], countryData.Latling[1])
	if err != nil {
		http.Error(w, "Failed to fetch country data", http.StatusBadGateway)
		return
	}

	precaverage := database.Average(weatherData.Precipitation)
	tempaverage := database.Average(weatherData.Temperature)

	// Prepare response
	response := map[string]interface{}{
		"country": country,
		"isoCode": isoCode,
		"features": utils.Featureseponse{},
		"lastRetrival": lastsignin,

	}

	responseFeatures := response["features"].(map[string]interface{})

	// Capital
	if features.Capital {
		responseFeatures["capital"] = countryData.Capital
	}

	// Coordinates
	if features.Coordinates {
		responseFeatures["coordinates"] = map[string]string{
			"latitude":  countryData.Latling[0],
			"longitude": countryData.Latling[1],
		}
	}

	// Population
	if features.Population {
		responseFeatures["population"] = countryData.Population
	}

	// Area
	if features.Area {
		responseFeatures["area"] = countryData.Area
	}

	if features.Temperature {
		responseFeatures["temperature"] = tempaverage
	}

	if features.Precipitation {
		responseFeatures["precipitation"] = precaverage
	}


	// Currency
	if tcs := features.TargetCurrencies; len(tcs) > 0 {

		rates, err := database.GetCurrencyRates(tcs, countryData.Cca3)
		if err != nil {
			http.Error(w, "Currency API failed", http.StatusBadGateway)
			return
		}
		responseFeatures["targetCurrencies"] = rates
	}

	// Trigger webhooks asynchronously
	go services.TriggerWebhooks("INVOKE", isoCode)

	// Send the final response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
}
