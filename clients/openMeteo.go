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

/*
openMeteoAPIResponse Represents the structure of the Open-Meteo API response.
*/
type openMeteoAPIResponse struct {
	Hourly struct {
		Temperature2m []float64 `json:"temperature_2m"`
		Precipitation []float64 `json:"precipitation"`
	} `json:"hourly"`
}

/*
average Computes the mean value of a float number, and returns 0 if it is empty
*/
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	// Loop through each value in the slice and add it to the sum
	for _, v := range values {
		sum += v
	}
	// Returning the average
	return sum / float64(len(values))
}

/*
GetWeather Retrieves weather data for a given latitude and longitude from the Open-Meteo API.
It computes the average temperature and percipitation from the hourly data
*/
func GetWeather(latitude, longitude string) (*utils.OpenMeteoresponse, error) {
	// Construct the Url
	url := fmt.Sprintf("%s?latitude=%s&longitude=%s&hourly=temperature_2m,precipitation", config.OPENMETEO_ROOT, latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("OpenMeteo API returned status %d", resp.StatusCode))
	}
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var apiResp openMeteoAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}
	// Compute the average values for temp and precipitation
	weather := &utils.OpenMeteoresponse{
		Temperature:   average(apiResp.Hourly.Temperature2m),
		Precipitation: average(apiResp.Hourly.Precipitation),
	}
	return weather, nil
}
