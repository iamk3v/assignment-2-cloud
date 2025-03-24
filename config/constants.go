package config

const START_URL = "/dashboard/" + VERSION

// API URLs
const (
	RESTCOUNTRIES_ROOT = "http://129.241.150.113:8080/v3.1/"
	CURRENCY_ROOT      = "http://129.241.150.113:9090/currency/"
	OPENMETEO_ROOT     = "https://api.open-meteo.com/v1/forecast/"
)

// API version
const VERSION = "v1"

//used for status testing
const Testcountry = "no"

// Error messages
const (
	ERR_NOT_FOUND             = "Not found"
	ERR_INTERNAL_SERVER_ERROR = "Internal server error"
	ERR_BAD_REQUEST           = "Bad request"
	// Fill in with more as we go
)
