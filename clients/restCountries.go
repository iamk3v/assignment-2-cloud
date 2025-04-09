package clients

import (
	"assignment-2/config"
	"assignment-2/database"
	"assignment-2/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

/*
GetCountryData Retrieves data for countries by country name or ISO code.
It attempts to load the data from cache. If it fails the external api is called.
*/
var GetCountryData = func(name string, isoCode string) (*utils.CountryResponse, error) {
	var url string
	var cacheKey string

	// Determine URL and cache key via either country name or ISO code
	if isoCode != "" {
		url = fmt.Sprintf("%salpha/%s", config.RESTCOUNTRIES_ROOT, isoCode)
		cacheKey = fmt.Sprintf("Country_alpha_%s", isoCode)
	} else if name != "" {
		url = fmt.Sprintf("%sname/%s", config.RESTCOUNTRIES_ROOT, name)
		cacheKey = fmt.Sprintf("Country_alpha_%s", name)
	} else {
		return nil, errors.New("no country name or isoCode provided")
	}

	// Response to hold the API response data
	var countryData []utils.CountryResponse

	// Tries to get a cache hit using the cache key
	if err := database.GetCachedData(cacheKey, &countryData); err == nil {
		fmt.Printf("Cache hit for key: %s\n", cacheKey)
		// Trigger webhook notification for the cache hit
		if webhookTrigger != nil {
			if isoCode != "" {
				webhookTrigger.TriggerWebhooks("CACHE_HIT", isoCode)
			} else {
				webhookTrigger.TriggerWebhooks("CACHE_HIT", name)
			}
		}
		// Returns the first entry if successful
		return &countryData[0], nil
	}
	fmt.Printf("Cache miss for key: %s\n", cacheKey)

	// Calls the API
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch country data: %w", err)
	}
	defer resp.Body.Close()

	// Handle HTTP errors from external API
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("REST Countries API returned status %d", resp.StatusCode)
	}

	// Read API response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error: Failed to read API response: %w", err)
	}

	// Unmarshal the response
	if err := json.Unmarshal(body, &countryData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Ensure data is available
	if len(countryData) == 0 {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// The retrieved result is cached
	if err := database.SetCacheEntry(cacheKey, countryData); err != nil {
		fmt.Printf("Failed to cache data for key %s: %v\n", cacheKey, err)
	}

	return &countryData[0], nil
}
