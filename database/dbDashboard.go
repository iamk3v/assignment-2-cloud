package database

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
		url = fmt.Sprintf("%s/alpha/%s", config.RESTCOUNTRIES_ROOT, isoCode)
	} else if name != "" {
		url = fmt.Sprintf("%s/name/%s", config.RESTCOUNTRIES_ROOT, name)
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

func GetWeatherDate(latitude string, longitude string) (*utils.OpenMeteoresponse, error) {

	url := fmt.Sprintf("%s?latitude=%s&longtitude=%s", config.RESTCOUNTRIES_ROOT, latitude, longitude)

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

	// Parse JSON response
	var weatherData []utils.OpenMeteoresponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Ensure data is available
	if len(weatherData) == 0 {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &weatherData[0], nil
}

func Average(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}

	var mean float64
	if len(numbers) > 0 {
		mean = sum / float64(len(numbers))
	}

	return mean
}

func GetCurrencyRates(curency []string, countrycode string) (*utils.CurrencyResponse, error) {
	var fullcurrencydata []utils.CurrencyResponse

	for _, cur := range curency {
		url := config.RESTCOUNTRIES_ROOT + countrycode

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

		// Parse JSON response
		var CurrencyData []utils.CurrencyResponse
		if err := json.Unmarshal(body, &CurrencyData); err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}

		// Ensure data is available
		if len(CurrencyData) == 0 {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}
		fullcurrencydata[cur] = CurrencyData
	}

	return &fullcurrencydata[], nil
}
