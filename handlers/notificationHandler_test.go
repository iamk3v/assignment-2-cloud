package handlers

import (
	"assignment-2/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestNotificationHandler_Post tests the POST endpoint for creating a new webhook.
func TestNotificationHandler_Post(t *testing.T) {
	// Create a POST request with a valid JSON payload.
	payload := `{"url": "https://example.com/webhook", "country": "NO", "event": "REGISTER"}`
	req := httptest.NewRequest(http.MethodPost, config.START_URL+"/notifications/", strings.NewReader(payload))
	rr := httptest.NewRecorder()

	// Directly call NotificationHandler.
	NotificationHandler(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	// Check that the status code is 201 Created.
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}

	// Decode the JSON response.
	var respBody map[string]string
	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify that the response contains a non-empty "id".
	if id, ok := respBody["id"]; !ok || id == "" {
		t.Errorf("Expected a non-empty webhook id, got %q", id)
	}

	// Verify that the response includes the HTTP Cat URL for 201.
	if cat, ok := respBody["httpCat"]; ok {
		expectedCat := "https://http.cat/201"
		if cat != expectedCat {
			t.Errorf("Expected httpCat URL %q, got %q", expectedCat, cat)
		}
	}
}
