package handlers

import (
	"assignment-2/config"
	"assignment-2/database"
	"assignment-2/services"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
DashboardHandler Handles GET requests for a populated dashboard
*/
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

/*
handleDashGetRequest Extracts the dashboard ID, gets configuration, fetches external data and
sends the dashboard response
*/
func handleDashGetRequest(w http.ResponseWriter, r *http.Request) {
	basePath := config.START_URL + "/dashboards/"
	trimmedPath := strings.TrimPrefix(r.URL.Path, basePath)
	parts := strings.Split(trimmedPath, "/")

	// Check if ID is provided
	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "Dashboard ID not provided", http.StatusBadRequest)
		return
	}
	id := parts[0]

	// Retrieve the dashboard configuration from firestore
	reg, err := database.GetOneRegistration(id)
	if err != nil {
		log.Println("Error retrieving dashboard with id " + id + ": " + err.Error())
		http.Error(w, "There was an error getting the dashboard with id: "+id, http.StatusInternalServerError)
		return
	}

	// Extract fields from the registration
	country := reg.Country
	isoCode := reg.IsoCode
	features := reg.Features

	// Get country info from the REST Countries API
	countryData, err := database.GetCountryData(country, isoCode)
	if err != nil {
		http.Error(w, "Failed to fetch country data", http.StatusBadGateway)
		return
	}

	// Get weather info from the Open-Meteo API
	weatherData, err := database.GetWeatherDate(countryData.Latling[0], countryData.Latling[1])
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusBadGateway)
		return
	}

	tempAverage := weatherData.Temperature
	precAverage := weatherData.Precipitation

	// Assemble the features based on the configuration
	featuresMap := make(map[string]interface{})

	// Capital
	if features.Capital {
		featuresMap["capital"] = countryData.Capital
	}

	// Coordinates
	if features.Coordinates {
		featuresMap["coordinates"] = map[string]string{
			"latitude":  countryData.Latling[0],
			"longitude": countryData.Latling[1],
		}
	}

	// Population
	if features.Population {
		featuresMap["population"] = countryData.Population
	}

	// Area
	if features.Area {
		featuresMap["area"] = countryData.Area
	}

	// Temperature
	if features.Temperature {
		featuresMap["temperature"] = tempAverage
	}

	// Precipitation
	if features.Precipitation {
		featuresMap["precipitation"] = precAverage
	}

	// Currency
	if len(features.TargetCurrencies) > 0 {
		rates, err := database.GetCurrencyRates(features.TargetCurrencies, countryData.Cca3)
		if err != nil {
			http.Error(w, "Currency API failed", http.StatusBadGateway)
			return
		}
		featuresMap["targetCurrencies"] = rates
	}

	// Building the response
	response := map[string]interface{}{
		"country":       country,
		"isoCode":       isoCode,
		"features":      featuresMap,
		"lastRetrieval": time.Now().Format(time.RFC3339),
	}

	// Trigger webhooks asynchronously
	go services.TriggerWebhooks("INVOKE", isoCode)

	// Send the final response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
