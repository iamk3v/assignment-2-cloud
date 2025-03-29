package clients

import (
	"assignment-2/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

/*
CurrencyAPIResponse Contains a map of currency codes to their corresponding exchange rates
*/
type CurrencyAPIResponse struct {
	Rates map[string]float64 `json:"rates"`
}

/*
GetRates Retrieves currency exchange rates for the provided base currency and returns them
*/
func GetRates(baseCurrency string) (map[string]float64, error) {
	// Construct the API endpoint Url
	url := fmt.Sprintf("%s%s", config.CURRENCY_ROOT, baseCurrency)
	// Fetch the currency data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprint("Currency API returned: %d", resp.StatusCode))
	}
	// Reading the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp CurrencyAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}
	// return the map of exchange rates
	return apiResp.Rates, nil
}
