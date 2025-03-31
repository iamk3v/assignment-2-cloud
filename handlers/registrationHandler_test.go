package handlers

import (
	"assignment-2/config"
	"assignment-2/utils"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetOne(t *testing.T) {
	id := "5NgZBNfIJA93lquV1n9R"
	req := httptest.NewRequest(http.MethodGet, config.START_URL+"/registrations/"+id, nil)
	w := httptest.NewRecorder()

	RegistrationHandler(w, req) // Directly call the handler

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	expected := utils.Dashboard{
		Id:      "5NgZBNfIJA93lquV1n9R",
		Country: "Norway",
		IsoCode: "NO",
		Features: utils.Features{
			Temperature:      true,
			Precipitation:    true,
			Capital:          true,
			Coordinates:      false,
			Population:       true,
			Area:             false,
			TargetCurrencies: []string{"EUR", "USD", "SEK"},
		},
		LastChange: time.Now(),
	}

	// Marshal it to JSON to compare
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Error marshaling expected dashboard: %v", err)
	}

	if string(body) != string(expectedJSON) {
		t.Errorf("Expected %s, got %s", string(expectedJSON), string(body))
	}
}
