package utils

type Statusresponse struct {
	CountriesAPI         int    `firestore:"countriesAPI" json:"countriesAPI"`
	CurrencyAPI          int    `firestore:"currencyAPI" json:"currencyAPI"`
	OpenmeteoAPI         int    `firestore:"openmeteoAPI" json:"openmeteoAPI"`
	Notificationresponse int    `firestore:"notificationresponse" json:"notificationresponse"`
	Dashboardresponse    int    `firestore:"dashboardresponse" json:"dashboardresponse"`
	Webhookssum          int    `firestore:"webhookssum" json:"webhookssum"`
	Version              string `firestore:"version" json:"version"`
	Uptime               string `firestore:"uptime" json:"uptime"`
}

type DashboardPost struct {
	Country    string   `firestore:"country" json:"country"`
	IsoCode    string   `firestore:"isoCode" json:"isoCode"`
	Features   Features `firestore:"features" json:"features"`
	LastChange string   `firestore:"lastChange" json:"lastChange"`
}

type Dashboard struct {
	Id         string   `firestore:"id" json:"id"`
	Country    string   `firestore:"country" json:"country"`
	IsoCode    string   `firestore:"isoCode" json:"isoCode"`
	Features   Features `firestore:"features" json:"features"`
	LastChange string   `firestore:"lastChange" json:"lastChange"`
}
type Features struct {
	Temperature      bool     `firestore:"temperature" json:"temperature"`
	Precipitation    bool     `firestore:"precipitation" json:"precipitation"`
	Capital          bool     `firestore:"capital" json:"capital"`
	Coordinates      bool     `firestore:"coordinates" json:"coordinates"`
	Population       bool     `firestore:"population" json:"population"`
	Area             bool     `firestore:"area" json:"area"`
	TargetCurrencies []string `firestore:"targetCurrencies" json:"targetCurrencies"`
}

type Webhook struct {
	ID      string `firestore:"id" json:"id"`
	URL     string `firestore:"url" json:"url"`
	Country string `firestore:"country" json:"country,omitempty"` // if empty, applies to all countries
	Event   string `firestore:"event" json:"event"`               // REGISTER, CHANGE, DELETE, INVOKE, ...
}

// WebhookInvocation is the payload we POST to the subscribed URL
type WebhookInvocation struct {
	ID      string `json:"id"`
	Country string `json:"country,omitempty"`
	Event   string `json:"event"`
	Time    string `json:"time"`
}

type CountryResponse struct {
	Population int       `json:"population"`
	Capital    []string  `json:"capital"`
	Area       float64   `json:"area"`
	Latlng     []float64 `json:"latlng"`
	Cca3       string    `json:"cca3"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

type OpenMeteoresponse struct {
	Daily struct {
		Temperature   []float64 `json:"temperature_2m_mean"`
		Precipitation []float64 `json:"precipitation_probability_mean"`
	}
}

type CurrencyAPIResult struct {
	BaseCode          string
	TimeLastUpdateUTC string
	TimeNextUpdateUTC string
	Rates             []CurrencyResponse
}
type CurrencyResponse struct {
	Code string  `json:"code"`
	Rate float64 `json:"rate"`
}

type GroupedCurrencyResponse struct {
	BaseCode               string             `json:"base_code"`
	TimeLastCurrencyUpdate string             `json:"time_last_update_utc"`
	TimeNextCurrencyUpdate string             `json:"time_next_update_utc"`
	Rates                  []CurrencyResponse `json:"rates"`
}
