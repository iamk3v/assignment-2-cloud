package clients

import (
	"assignment-2/config"
	"assignment-2/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetCountryData(name string, isoCode string) (*utils.CountryResponse, error) {
	var url string

	if isoCode != "" {
		url = fmt.Sprintf("%salpha/%s", config.RESTCOUNTRIES_ROOT, isoCode)
	} else if name != "" {
		url = fmt.Sprintf("%sname/%s", config.RESTCOUNTRIES_ROOT, name)
	} else {
		return nil, errors.New("no country name or isoCode provided")
	}

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
		return nil, fmt.Errorf("Error: Failed to read API response. %w", err)

	}

	// Parse JSON response
	var countryData []utils.CountryResponse
	if err := json.Unmarshal(body, &countryData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Ensure data is available
	if len(countryData) == 0 {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &countryData[0], nil
}
