package handlers

import (
	"assignment-2/config"
	"assignment-2/utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Store the testId for the registration sent in TestPostRegistration
var testId string

/*
TestPostRegistration creates a test to add a new registration to the dashboard DB, expected result: ok
*/
func TestPostRegistration(t *testing.T) {
	// Define the test post body
	testBody := utils.DashboardPost{
		Country: "Norway",
		IsoCode: "NO",
		Features: utils.Features{
			Temperature:      true,
			Precipitation:    false,
			Capital:          true,
			Coordinates:      true,
			Population:       false,
			Area:             false,
			TargetCurrencies: []string{"EUR", "SEK", "DEK"},
		},
		LastChange: time.Now().Local().String(),
	}

	// Encode the body
	postData, err := json.Marshal(testBody)
	if err != nil {
		t.Error("Failed to marshal test post data", err.Error())
	}

	// Create the request
	req := httptest.NewRequest("POST", "/registration", bytes.NewBuffer(postData))
	w := httptest.NewRecorder()

	// Send request to the handler
	RegistrationHandler(w, req)

	// Capture result
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Error("Expected status code 201, got", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error("Failed to read response body", err.Error())
	}

	// Define response struct from a successful post
	type respStruct struct {
		Id         string `json:"id"`
		LastChange string `json:"lastChange"`
	}
	var result respStruct
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Error("Failed to unmarshal response body", err.Error())
	}
	testId = result.Id
}

/*
TestGetOne creates a test to get one spesific dashboard, expected result: ok
*/
func TestGetOne(t *testing.T) {
	// Create the request
	req := httptest.NewRequest(http.MethodGet, config.START_URL+"/registrations/"+testId, nil)
	w := httptest.NewRecorder()

	// Send request to the handler
	RegistrationHandler(w, req)

	// Get the response
	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Store response
	var gotten utils.Dashboard
	err = json.Unmarshal(body, &gotten)
	if err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}

	// Build the expected struct
	expected := utils.Dashboard{
		Id:      testId,
		Country: "Norway",
		IsoCode: "NO",
		Features: utils.Features{
			Temperature:      true,
			Precipitation:    false,
			Capital:          true,
			Coordinates:      true,
			Population:       false,
			Area:             false,
			TargetCurrencies: []string{"EUR", "SEK", "DEK"},
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

/*
TestGetAll creates a test to get all dashboards, expected result: ok
*/
func TestGetAll(t *testing.T) {
	// Create the request
	req := httptest.NewRequest(http.MethodGet, config.START_URL+"/registrations/", nil)
	w := httptest.NewRecorder()

	// Send request to the handler
	RegistrationHandler(w, req)

	// Get the response
	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(body) == 0 || body == nil {
		t.Fatalf("Registration response body is empty")
	}
}

/*
TestDeleteRegistration deletes the test post registration sent in TestPostRegistration
*/
func TestDeleteRegistration(t *testing.T) {
	// Create the request
	req := httptest.NewRequest(http.MethodDelete, config.START_URL+"/registrations/"+testId, nil)
	w := httptest.NewRecorder()

	// Send request to the handler
	RegistrationHandler(w, req)

	// Get the response
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}
