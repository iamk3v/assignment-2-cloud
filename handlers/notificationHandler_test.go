package handlers

import (
	"assignment-2/config"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestNotificationHandler_Post creates a new webhook and checks that a valid ID is returned.
func TestNotificationHandler_Post(t *testing.T) {
	// Create a POST request with a valid JSON payload.
	payload := `{"url": "https://example.com/webhook", "country": "NO", "event": "REGISTER"}`
	req := httptest.NewRequest(http.MethodPost, config.START_URL+"/notifications/", strings.NewReader(payload))
	rr := httptest.NewRecorder()

	// Call the NotificationHandler directly.
	NotificationHandler(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	// Verify that the status code is 201 Created.
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}

	// Decode the JSON response.
	var respBody map[string]string
	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check that the response contains a non-empty "id".
	if id, ok := respBody["id"]; !ok || id == "" {
		t.Errorf("Expected a valid webhook id, got %q", respBody["id"])
	}

	// Check for the HTTP Cat URL.
	if cat, ok := respBody["httpCat"]; ok {
		expectedCat := "https://http.cat/201"
		if cat != expectedCat {
			t.Errorf("Expected httpCat URL %q, got %q", expectedCat, cat)
		}
	}
}

// TestNotificationHandler_Get creates a webhook, retrieves it using its ID, and verifies the result.
func TestNotificationHandler_Get(t *testing.T) {
	// Create a webhook using POST.
	postPayload := `{"url": "https://example.com/webhook", "country": "NO", "event": "REGISTER"}`
	reqPost := httptest.NewRequest(http.MethodPost, config.START_URL+"/notifications/", strings.NewReader(postPayload))
	rrPost := httptest.NewRecorder()
	NotificationHandler(rrPost, reqPost)
	resPost := rrPost.Result()
	defer resPost.Body.Close()

	if resPost.StatusCode != http.StatusCreated {
		t.Fatalf("Expected POST status %d, got %d", http.StatusCreated, resPost.StatusCode)
	}

	var postResp map[string]string
	if err := json.NewDecoder(resPost.Body).Decode(&postResp); err != nil {
		t.Fatalf("Failed to decode POST response: %v", err)
	}
	id, ok := postResp["id"]
	if !ok || id == "" {
		t.Fatalf("POST response did not contain a valid id")
	}

	// GET request using the returned id.
	getURL := config.START_URL + "/notifications/" + id
	reqGet := httptest.NewRequest(http.MethodGet, getURL, nil)
	rrGet := httptest.NewRecorder()
	NotificationHandler(rrGet, reqGet)
	resGet := rrGet.Result()
	defer resGet.Body.Close()

	// Expect a 200 OK status.
	if resGet.StatusCode != http.StatusOK {
		t.Errorf("Expected GET status %d, got %d", http.StatusOK, resGet.StatusCode)
	}

	// Decode the GET response.
	var getResp map[string]interface{}
	body, err := io.ReadAll(resGet.Body)
	if err != nil {
		t.Fatalf("Failed to read GET response: %v", err)
	}
	if err := json.Unmarshal(body, &getResp); err != nil {
		t.Fatalf("Failed to unmarshal GET response: %v", err)
	}

	// Verify that the returned id matches.
	if getResp["id"] != id {
		t.Errorf("Expected webhook id %q, got %q", id, getResp["id"])
	}
}

// TestNotificationHandler_Delete creates a webhook, deletes it, then confirms it can no longer be retrieved.
func TestNotificationHandler_Delete(t *testing.T) {
	// Create a webhook via POST.
	payload := `{"url": "https://example.com/webhook", "country": "NO", "event": "REGISTER"}`
	reqPost := httptest.NewRequest(http.MethodPost, config.START_URL+"/notifications/", strings.NewReader(payload))
	rrPost := httptest.NewRecorder()
	NotificationHandler(rrPost, reqPost)
	resPost := rrPost.Result()
	defer resPost.Body.Close()

	if resPost.StatusCode != http.StatusCreated {
		t.Fatalf("Expected POST status %d, got %d", http.StatusCreated, resPost.StatusCode)
	}

	var postResp map[string]string
	if err := json.NewDecoder(resPost.Body).Decode(&postResp); err != nil {
		t.Fatalf("Failed to decode POST response: %v", err)
	}
	id, ok := postResp["id"]
	if !ok || id == "" {
		t.Fatalf("POST response did not contain a valid id")
	}

	// DELETE request for the created webhook.
	delURL := config.START_URL + "/notifications/" + id
	reqDel := httptest.NewRequest(http.MethodDelete, delURL, nil)
	rrDel := httptest.NewRecorder()
	NotificationHandler(rrDel, reqDel)
	resDel := rrDel.Result()
	defer resDel.Body.Close()

	// Expect 204 No Content.
	if resDel.StatusCode != http.StatusNoContent {
		t.Errorf("Expected DELETE status %d, got %d", http.StatusNoContent, resDel.StatusCode)
	}

	// Try to GET the deleted webhook to confirm deletion.
	reqGet := httptest.NewRequest(http.MethodGet, delURL, nil)
	rrGet := httptest.NewRecorder()
	NotificationHandler(rrGet, reqGet)
	resGet := rrGet.Result()
	defer resGet.Body.Close()

	// Expect an error status (either 400 Bad Request or 404 Not Found).
	if resGet.StatusCode != http.StatusNotFound && resGet.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected GET after DELETE to return 404 or 400, got %d", resGet.StatusCode)
	}
}

// TestNotificationHandler_Patch tests the PATCH endpoint for updating an existing webhook.
func TestNotificationHandler_Patch(t *testing.T) {
	// Create a webhook via POST.
	postPayload := `{"url": "https://example.com/webhook", "country": "NO", "event": "REGISTER"}`
	reqPost := httptest.NewRequest(http.MethodPost, config.START_URL+"/notifications/", strings.NewReader(postPayload))
	rrPost := httptest.NewRecorder()
	NotificationHandler(rrPost, reqPost)
	resPost := rrPost.Result()
	defer resPost.Body.Close()

	if resPost.StatusCode != http.StatusCreated {
		t.Fatalf("Expected POST status %d, got %d", http.StatusCreated, resPost.StatusCode)
	}

	var postResp map[string]string
	if err := json.NewDecoder(resPost.Body).Decode(&postResp); err != nil {
		t.Fatalf("Failed to decode POST response: %v", err)
	}

	id, ok := postResp["id"]
	if !ok || id == "" {
		t.Fatalf("POST response did not contain a valid id")
	}

	// Create a PATCH request to update the webhook.
	patchPayload := `{"url": "https://updated-example.com/webhook"}`
	patchURL := config.START_URL + "/notifications/" + id
	reqPatch := httptest.NewRequest(http.MethodPatch, patchURL, strings.NewReader(patchPayload))
	rrPatch := httptest.NewRecorder()
	NotificationHandler(rrPatch, reqPatch)
	resPatch := rrPatch.Result()
	defer resPatch.Body.Close()

	// Expect a 204 No Content status on successful patch.
	if resPatch.StatusCode != http.StatusNoContent {
		t.Errorf("Expected PATCH status %d, got %d", http.StatusNoContent, resPatch.StatusCode)
	}

	// Issue a GET request to verify that the update took effect.
	reqGet := httptest.NewRequest(http.MethodGet, patchURL, nil)
	rrGet := httptest.NewRecorder()
	NotificationHandler(rrGet, reqGet)
	resGet := rrGet.Result()
	defer resGet.Body.Close()

	if resGet.StatusCode != http.StatusOK {
		t.Fatalf("Expected GET status %d, got %d", http.StatusOK, resGet.StatusCode)
	}

	// Decode the GET response.
	var getResp map[string]interface{}
	body, err := io.ReadAll(resGet.Body)
	if err != nil {
		t.Fatalf("Failed to read GET response: %v", err)
	}
	if err := json.Unmarshal(body, &getResp); err != nil {
		t.Fatalf("Failed to unmarshal GET response: %v", err)
	}

	// Verify that the "url" field was updated.
	updatedURL, exists := getResp["url"].(string)
	if !exists {
		t.Errorf("GET response does not include 'url' field")
	} else if updatedURL != "https://updated-example.com/webhook" {
		t.Errorf("Expected updated url %q, got %q", "https://updated-example.com/webhook", updatedURL)
	}
}
