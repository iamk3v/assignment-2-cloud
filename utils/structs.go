package utils

import "time"

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

type Dashboard struct {
	Id         string    `firestore:"id" json:"id"`
	Country    string    `firestore:"country" json:"country"`
	IsoCode    string    `firestore:"isoCode" json:"isoCode"`
	Features   Features  `firestore:"features" json:"features"`
	LastChange time.Time `firestore:"lastChange" json:"lastChange"`
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
	ID      string `json:"id"`
	URL     string `json:"url"`
	Country string `json:"country,omitempty"` // if empty, applies to all countries
	Event   string `json:"event"`             // REGISTER, CHANGE, DELETE, INVOKE, ...
}

// WebhookInvocation is the payload we POST to the subscribed URL
type WebhookInvocation struct {
	ID      string `json:"id"`
	Country string `json:"country,omitempty"`
	Event   string `json:"event"`
	Time    string `json:"time"`
}

type CountryResponse struct {
	Population int      `json:"population"`
	Capital    string   `json:"capital"`
	Area       string   `json:"area"`
	Latling    []string `json:"latling"`
}

type OpenMeteoresponse struct {
	Temperature   []int `json:"temperature"`
	Precipitation []int `json:"precipitation"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CurrencyResponse struct {
	TargetCurrencies []string `json:"targetCurrencies"`
}
