package clients

import (
	"assignment-2/config"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetCurrencyRates(curency []string, countrycode string) ([]utils.CurrencyResponse, error) {
	url := config.CURRENCY_ROOT + countrycode

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rate data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}

	// Define response struct to match the JSON
	var apiResponse struct {
		BaseCode string             `json:"base_code"`
		Rates    map[string]float64 `json:"rates"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	var fullcurrencydata []utils.CurrencyResponse

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

	return fullcurrencydata, nil

}
