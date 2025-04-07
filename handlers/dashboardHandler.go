package handlers

import (
	"assignment-2/clients"
	"assignment-2/config"
	"assignment-2/database"
	"assignment-2/services"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
DashboardHandler Handles GET requests for a populated dashboard
*/
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	basePath := config.START_URL + "/dashboards/"
	trimmedPath := strings.TrimPrefix(r.URL.Path, basePath)
	parts := strings.Split(trimmedPath, "/")
	id := parts[0]
	// Check if ID is provided

	if len(parts) < 1 || id == "" {
		http.Error(w, "Dashboard ID not provided", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleDashGetRequest(w, r, id)
	case http.MethodHead:
		handleDashHeadRequest(w, r, id)
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
func handleDashGetRequest(w http.ResponseWriter, r *http.Request, id string) {

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
	countryData, err := clients.GetCountryData(country, isoCode)
	if err != nil {
		log.Println("failed to fetch country data:" + err.Error())
		http.Error(w, "Failed to fetch country data", http.StatusBadGateway)
		return
	}

	currencyCode := []string{}
	for code := range countryData.Currencies {
		currencyCode = append(currencyCode, code)
	}

	// Get weather info from the Open-Meteo API
	weatherData, err := clients.GetWeatherDate(countryData.Latlng[0], countryData.Latlng[1])
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusBadGateway)
		return
	}

	tempAverage := clients.Average(weatherData.Daily.Temperature)
	precAverage := clients.Average(weatherData.Daily.Precipitation)

	// Assemble the features based on the configuration
	featuresMap := make(map[string]interface{})

	// Capital
	if features.Capital {
		featuresMap["capital"] = countryData.Capital
	}

	// Coordinates
	if features.Coordinates {
		featuresMap["coordinates"] = map[string]float64{
			"latitude":  countryData.Latlng[0],
			"longitude": countryData.Latlng[1],
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
		for currency := range currencyCode {
			rates, err := clients.GetCurrencyRates(features.TargetCurrencies, currencyCode[currency])
			if err != nil {
				http.Error(w, "Currency API failed", http.StatusBadGateway)
				return
			}
			// Initialize if needed to avoid panic
			if featuresMap["targetCurrencies"] == nil {
				featuresMap["targetCurrencies"] = []utils.GroupedCurrencyRates{}
			}

			existingGroups := featuresMap["targetCurrencies"].([]utils.GroupedCurrencyRates)

			base := currencyCode[currency] // e.g., "BND", "SGD"

			groupExists := false

			// Check if group exists
			for i, group := range existingGroups {
				if group.BaseCode == base {
					existingGroups[i].Rates = append(existingGroups[i].Rates, rates...)
					groupExists = true
					break
				}
			}

			if !groupExists {
				existingGroups = append(existingGroups, utils.GroupedCurrencyRates{
					BaseCode: base,
					Rates:    rates,
				})
			}

			featuresMap["targetCurrencies"] = existingGroups
		}
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

func handleDashHeadRequest(w http.ResponseWriter, r *http.Request, id string) {
	// Get one dashboard
	rawContent, err := database.GetOneRegistration(id)
	if err != nil {
		log.Println("Error retrieving dashboard with id " + id + ": " + err.Error())
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

	// Set and send back headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	w.WriteHeader(http.StatusNoContent)
}
