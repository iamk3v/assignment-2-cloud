package clients

import (
	"assignment-2/config"
	"assignment-2/database"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/*
GetCurrencyRates Retrieves the currency API result from cache or the external API
*/
func GetCurrencyRates(curency []string, countrycode string) (*utils.CurrencyAPIResult, error) {
	// Create a unique cache key via the country code
	cacheKey := fmt.Sprintf("currency_%s", countrycode)

	var result utils.CurrencyAPIResult

	// Retrieve cached data
	if err := database.GetCachedData(cacheKey, &result); err == nil {
		fmt.Printf("Cache hit for key: %s\n", cacheKey)
		return &result, nil
	}
	fmt.Printf("Cache miss for key: %s", cacheKey)

	// Build the API url
	url := config.CURRENCY_ROOT + countrycode

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rate data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}

	// Define response struct to match the JSON
	var apiResponse struct {
		BaseCode          string             `json:"base_code"`
		TimeLastUpdateUTC string             `json:"time_last_update_utc"`
		TimeNextUpdateUTC string             `json:"time_next_update_utc"`
		Rates             map[string]float64 `json:"rates"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	var fullcurrencydata []utils.CurrencyResponse

	//extracts the currency rates based on the currency code
	for _, code := range curency {
		rate, exists := apiResponse.Rates[code]
		if !exists {
			return nil, fmt.Errorf("currency code %s not found in API response", code)
		}

		fullcurrencydata = append(fullcurrencydata, utils.CurrencyResponse{
			Code: code,
			Rate: rate,
		})
	}

	// Creating the result
	result = utils.CurrencyAPIResult{
		BaseCode:          apiResponse.BaseCode,
		TimeLastUpdateUTC: apiResponse.TimeLastUpdateUTC,
		TimeNextUpdateUTC: apiResponse.TimeNextUpdateUTC,
		Rates:             fullcurrencydata,
	}

	// Cache the result for future calls with the same key
	if err := database.SetCacheEntry(cacheKey, result); err != nil {
		fmt.Printf("Failed to cache data for key %s: %v\n", cacheKey, err)
	}

	return &result, nil
}
