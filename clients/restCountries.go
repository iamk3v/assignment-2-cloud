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
GetCountry Retrieves country data using ISO code or country name and then parses the JSON response into CountryResponse
*/
func GetCountry(name string, isoCode string) (*utils.CountryResponse, error) {
	var url string
	// If the Url is ISO code or country name
	if isoCode == "" {
		url = fmt.Sprintf("%s/alpha/%s", config.RESTCOUNTRIES_ROOT, isoCode)
	} else if name != "" {
		url = fmt.Sprintf("%sname%s", config.RESTCOUNTRIES_ROOT, name)
	} else {
		// Error if none is provided
		return nil, errors.New("no country name or isoCode provided")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("REST countries API returned status %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var countries []utils.CountryResponse
	if err := json.Unmarshal(body, &countries); err != nil {
		return nil, err
	}

	// If the country data is returned
	if len(countries) == 0 {
		return nil, errors.New("no country found")
	}

	// Returns the first country
	return &countries[0], nil
}
