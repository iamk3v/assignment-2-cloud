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
DashboardHandler Handles requests sent to the /dashboards endpoint, routing the request to
corresponding handle functions based on http methods.
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
handleDashGetRequest gets configuration, fetches external data and sends the dashboard response
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
		log.Println("failed to fetch country data: " + err.Error())
		http.Error(w, "Failed to fetch country data", http.StatusBadGateway)
		return
	}

	currencyCode := []string{}
	for code := range countryData.Currencies {
		currencyCode = append(currencyCode, code)
	}
	// Check if no currency codes were found
	if len(currencyCode) == 0 {
		http.Error(w, "no currency codes found for country", http.StatusInternalServerError)
		return
	}

	// Get weather info from the Open-Meteo API
	weatherData, err := clients.GetWeatherDate(countryData.Latlng[0], countryData.Latlng[1])
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusBadGateway)
		return
	}

	tempAverage := clients.Average(weatherData.Daily.Temperature)
	precAverage := clients.Average(weatherData.Daily.Precipitation)

	// Assemble the features based on the configuration in the database
	featuresMap := make(map[string]interface{})

	if features.Capital {
		featuresMap["capital"] = countryData.Capital
	}

	if features.Coordinates {
		featuresMap["coordinates"] = map[string]float64{
			"latitude":  countryData.Latlng[0],
			"longitude": countryData.Latlng[1],
		}
	}

	if features.Population {
		featuresMap["population"] = countryData.Population
	}

	if features.Area {
		featuresMap["area"] = countryData.Area
	}

	if features.Temperature {
		featuresMap["temperature"] = tempAverage
	}

	if features.Precipitation {
		featuresMap["precipitation"] = precAverage
	}

	if len(features.TargetCurrencies) > 0 {
		for currency := range currencyCode {
			//get currency data from the currency API
			result, err := clients.GetCurrencyRates(features.TargetCurrencies, currencyCode[currency])
			if err != nil {
				http.Error(w, "Currency API failed", http.StatusBadGateway)
				return
			}
			// Initialize if needed to avoid panic
			if featuresMap["targetCurrencies"] == nil {
				featuresMap["targetCurrencies"] = []utils.GroupedCurrencyResponse{}
			}

			existingGroups := featuresMap["targetCurrencies"].([]utils.GroupedCurrencyResponse)

			groupExists := false

			// Check if group exists, and if exits it becomes an array
			for i, group := range existingGroups {
				if group.BaseCode == result.BaseCode {
					existingGroups[i].Rates = append(existingGroups[i].Rates, result.Rates...)
					groupExists = true
					break
				}
			}

			//if group dosnt exist create feature dashboard
			if !groupExists {
				existingGroups = append(existingGroups, utils.GroupedCurrencyResponse{
					BaseCode:               result.BaseCode,
					TimeLastCurrencyUpdate: result.TimeLastUpdateUTC,
					TimeNextCurrencyUpdate: result.TimeNextUpdateUTC,
					Rates:                  result.Rates,
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
		"lastRetrieval": time.Now().Local().String(),
	}

	// Trigger webhooks asynchronously
	go services.TriggerWebhooks("INVOKE", isoCode)

	// Send the final response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response: " + err.Error())
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

/*
handleDashHeadRequest Handles HEAD requests sent to the dashboard handler
*/
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
