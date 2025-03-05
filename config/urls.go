package config

// BaseURL is the base URL for the local server.
const BaseURL = "http://localhost:8080"

// External API URLs
const (
	CountriesNowAPI  = "http://129.241.150.113:3500/api/v0.1" // Base URL for CountriesNow API
	RestCountriesAPI = "http://129.241.150.113:8080/v3.1"     // Base URL for REST Countries API
)

// Application endpoints
const (
	Root               = "/"                           // Root endpoint
	InfoEndpoint       = "/countryinfo/v1/info/"       // Endpoint for country information
	PopulationEndpoint = "/countryinfo/v1/population/" // Endpoint for population data
	StatusEndpoint     = "/countryinfo/v1/status/"     // Endpoint for service status
)

// HTTP methods
const (
	GET  = "GET"  // HTTP GET method
	POST = "POST" // HTTP POST method
)

// API response status codes
const (
	StatusOK                  = 200 // Success status code
	StatusBadRequest          = 400 // Bad request status code
	StatusNotFound            = 404 // Resource not found status code
	StatusInternalServerError = 500 // Internal server error status code
)

// API rate limits (if applicable)
const (
	RateLimitPerMinute = 60 // Maximum allowed requests per minute
)

// Mock data paths (for development/testing)
const (
	MockCountryInfoPath = "mocks/country_info.json"    // Path to mock country info data
	MockPopulationPath  = "mocks/population_data.json" // Path to mock population data
	MockCitiesPath      = "mocks/cities.json"          // Path to mock cities data
)
