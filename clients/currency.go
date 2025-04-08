package clients

import (
	"assignment-2/config"
	"assignment-2/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetCurrencyRates(curency []string, countrycode string) (*utils.CurrencyAPIResult, error) {
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

	//returns data as seperate variables
	return &utils.CurrencyAPIResult{
		BaseCode:          apiResponse.BaseCode,
		TimeLastUpdateUTC: apiResponse.TimeLastUpdateUTC,
		TimeNextUpdateUTC: apiResponse.TimeNextUpdateUTC,
		Rates:             fullcurrencydata,
	}, nil

}
