package clients

import (
	"assignment-2/config"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

func GetWeatherDate(latitude float64, longitude float64) (*utils.OpenMeteoresponse, error) {

	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&daily=temperature_2m_mean,precipitation_probability_mean", config.OPENMETEO_ROOT, latitude, longitude)

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
	var weatherData utils.OpenMeteoresponse

	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Ensure data is available
	if len(weatherData.Daily.Precipitation) == 0 {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &weatherData, nil
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

	// Round to 2 decimal places
	return math.Round(mean*100) / 100
}
