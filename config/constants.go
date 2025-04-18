package config

// The start url for the service
const START_URL = "/dashboard/" + VERSION

// API URLs
const (
	RESTCOUNTRIES_ROOT = "http://129.241.150.113:8080/v3.1/"
	CURRENCY_ROOT      = "http://129.241.150.113:9090/currency/"
	OPENMETEO_ROOT     = "https://api.open-meteo.com/v1/forecast"
)

// API version
const VERSION = "v1"

// used for status testing
const Testcountry = "no"
const Testcurrency = "nok"
const Testweather = "?latitude=52.52&longitude=13.41&hourly=temperature_2m"

// Error messages
const (
	ERR_NOT_FOUND             = "Not found"
	ERR_INTERNAL_SERVER_ERROR = "Internal server error"
	ERR_BAD_REQUEST           = "Bad request"
)

// Database
const PROJECT_ID = "assignment-2-279db"
const DASHBOARD_COLLECTION = "dashboards"
const NOTIFICATION_COLLECTION = "webhooks"
