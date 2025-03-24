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
