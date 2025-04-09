package handlers

import (
	"assignment-2/clients"
	"assignment-2/database"
	"assignment-2/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
sets a predefined database pull
*/
func mockGetOneRegistration(id string) (*utils.Dashboard, error) {
	return &utils.Dashboard{
		Id:      id,
		Country: "Norway",
		IsoCode: "NO",
		Features: utils.Features{
			Capital:          true,
			Coordinates:      true,
			Population:       true,
			Area:             true,
			Temperature:      true,
			Precipitation:    true,
			TargetCurrencies: []string{"EUR", "USD"},
		},
		LastChange: time.Now().Format("20060102 15:04"),
	}, nil
}

/*
sets predefined country data
*/
func mockGetCountryData(country, iso string) (*utils.CountryResponse, error) {
	return &utils.CountryResponse{

		Capital:    []string{"Oslo"},
		Latlng:     []float64{62.0, 10.0},
		Population: 5379475,
		Area:       385207.0,
		Currencies: map[string]struct {
			Name   string `json:"name"`
			Symbol string `json:"symbol"`
		}{
			"NOK": {
				Name:   "Norwegian Krone",
				Symbol: "kr",
			},
		},
	}, nil
}

/*
sets predefined weather data
*/
func mockGetWeatherDate(lat float64, lon float64) (*utils.OpenMeteoresponse, error) {
	return &utils.OpenMeteoresponse{
		Daily: struct {
			Temperature   []float64 `json:"temperature_2m_mean"`
			Precipitation []float64 `json:"precipitation_probability_mean"`
		}{
			Temperature:   []float64{2.0, 3.0, 4.0},
			Precipitation: []float64{10.0, 20.0, 30.0},
		},
	}, nil
}

/*
sets predefined weather data
*/
func mockGetCurrencyRates(targets []string, base string) (*utils.CurrencyAPIResult, error) {
	return &utils.CurrencyAPIResult{
		BaseCode:          base,
		TimeLastUpdateUTC: time.Now().Format(time.RFC3339),
		TimeNextUpdateUTC: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
		Rates: []utils.CurrencyResponse{
			{Code: "EUR", Rate: 0.09},
			{Code: "USD", Rate: 0.1},
		},
	}, nil
}

/*
TestDashboardHandler tests the populated dashboard functionality
*/
func TestDashboardHandler(t *testing.T) {
	// Patch actual functions with mocks
	database.GetOneRegistration = mockGetOneRegistration
	clients.GetCountryData = mockGetCountryData
	clients.GetWeatherDate = mockGetWeatherDate
	clients.GetCurrencyRates = mockGetCurrencyRates

	// Create request
	req := httptest.NewRequest("GET", "/dashboard/v1/dashboards/mock-id", nil)
	rec := httptest.NewRecorder()

	DashboardHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rec.Code)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("Error decoding JSON response: %v", err)
	}

	if body["country"] != "Norway" {
		t.Errorf("Expected country Norway, got %v", body["country"])
	}

	features := body["features"].(map[string]interface{})

	capital := features["capital"].([]interface{})
	if capital[0] != "Oslo" {
		t.Errorf("Expected first capital Oslo, got %v", capital[0])
	}

	if features["temperature"] != 3.0 {
		t.Errorf("Expected mean temperature 3.0, got %v", features["temperature"])
	}

	if _, ok := features["targetCurrencies"]; !ok {
		t.Error("Expected targetCurrencies in features, got none")
	}
}
