package clients

import (
	"assignment-2/config"
	"assignment-2/database"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

/*
GetWeatherDate calls the external API, OpenMeteo if the data is not available from cache
*/
var GetWeatherDate = func(latitude float64, longitude float64) (*utils.OpenMeteoresponse, error) {

	// Defines a key for cache based on lat and long
	cacheKey := fmt.Sprintf("Openmeteo_%f_%f", latitude, longitude)

	// Checks if there is cached data
	var weatherData utils.OpenMeteoresponse
	if err := database.GetCachedData(cacheKey, &weatherData); err == nil {
		fmt.Printf("Cache hit for key: %s\n", cacheKey)
		// Trigger webhook event for cache hit
		if webhookTrigger != nil {
			webhookTrigger.TriggerWebhooks("CACHE_HIT", fmt.Sprintf("LAT:%f, LONG:%f", latitude, longitude))
		}
		return &weatherData, nil
	}
	fmt.Printf("Cache miss for key: %s\n", cacheKey)

	// Construct the URL for the API call
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&daily=temperature_2m_mean,precipitation_probability_mean", config.OPENMETEO_ROOT, latitude, longitude)

	// Make the HTTP get request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	// Handle HTTP errors from external API
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenMeteo API returned status %d", resp.StatusCode)
	}

	// Read API response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed to read API response. %w", err)
	}

	// Parse JSON response into weatherData variable
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Ensure data is available
	if len(weatherData.Daily.Precipitation) == 0 {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Cache the retireved data
	if err := database.SetCacheEntry(cacheKey, weatherData); err != nil {
		fmt.Printf("Failed to cache data for key %s: %v\n", cacheKey, err)
	}

	return &weatherData, nil
}

/*
Average returns the average of a slice of float64 numbers
*/
func Average(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}

	var mean float64
	if len(numbers) > 0 {
		mean = sum / float64(len(numbers))
	}

	// Round to 2 decimal places
	return math.Round(mean*100) / 100
}
