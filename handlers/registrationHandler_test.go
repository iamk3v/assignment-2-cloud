package handlers

import (
	"assignment-2/config"
	"assignment-2/utils"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Function creates a test to get one spesific dashboard, expect ok
func TestGetOne(t *testing.T) {
	testId := "sSzvOGs4zEuRPPCbwyVT"
	req := httptest.NewRequest(http.MethodGet, config.START_URL+"/registrations/"+testId, nil)
	w := httptest.NewRecorder()

	RegistrationHandler(w, req) // Directly call the handler

	// Get the response
	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	var gotten utils.Dashboard
	err = json.Unmarshal(body, &gotten)
	if err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}

	// Build the expected struct
	expected := utils.Dashboard{
		Id:      "sSzvOGs4zEuRPPCbwyVT",
		Country: "Norway",
		IsoCode: "",
		Features: utils.Features{
			Temperature:      true,
			Precipitation:    true,
			Capital:          true,
			Coordinates:      false,
			Population:       true,
			Area:             false,
			TargetCurrencies: []string{"EUR", "USD", "SEK"},
		},
	}

	// Test a few attributes to see if they are expected
	if expected.Id != gotten.Id {
		t.Errorf("Expected ID %s, got ID %s", expected.Id, gotten.Id)
	}

	if expected.Country != gotten.Country {
		t.Errorf("Expected Country %s, got Country %s", expected.Country, gotten.Country)
	}

	if expected.Features.Temperature != gotten.Features.Temperature {
		t.Errorf("Expected value %t, got value %t", expected.Features.Temperature, gotten.Features.Temperature)
	}
}

// Function creates a test to get all dashboards, expect ok
func TestGetAll(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, config.START_URL+"/registrations/", nil)
	w := httptest.NewRecorder()

	RegistrationHandler(w, req) // Directly call the handler

	// Get the response
	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	if len(body) == 0 || body == nil {
		t.Fatalf("Registration response body is empty")
	}
}
