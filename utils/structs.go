package utils

type Statusresponse struct {
	CountriesAPI         int `firestore:"countriesAPI" json:"countriesAPI"`
	CurrencyAPI          int `firestore:"currencyAPI" json:"currencyAPI"`
	OpenmeteoAPI         int `firestore:"openmeteoAPI" json:"openmeteoAPI"`
	Notificationresponse int `firestore:"notificationresponse" json:"notificationresponse"`
	Dashboardresponse    int `firestore:"dashboardresponse" json:"dashboardresponse"`
	Webhookssum          int `firestore:"webhookssum" json:"webhookssum"`
	Version              int `firestore:"version" json:"version"`
	Uptime               int `firestore:"uptime" json:"uptime"`
}
