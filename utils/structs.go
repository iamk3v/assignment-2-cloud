package utils

type Statusresponse struct {
	countriesAPI         int `firestore:"countriesAPI" json:"countriesAPI"`
	currencyAPI          int `firestore:"currencyAPI" json:"currencyAPI"`
	openmeteoAPI         int `firestore:"openmeteoAPI" json:"openmeteoAPI"`
	notificationresponse int `firestore:"notificationresponse" json:"notificationresponse"`
	dashboardresponse    int `firestore:"dashboardresponse" json:"dashboardresponse"`
	webhookssum          int `firestore:"webhookssum" json:"webhookssum"`
	version              int `firestore:"version" json:"version"`
	uptime               int `firestore:"uptime" json:"uptime"`
}
