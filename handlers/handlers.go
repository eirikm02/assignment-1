package handlers

import (
	"assignment-1/models"   // Import models
	"assignment-1/services" // Import services
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var startTime = time.Now() // Track server start time

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract country code from the URL
	countryCode := strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/info/")

	// Fetch data from external APIs
	countryInfo, err := services.FetchCountryInfo(countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cities, err := services.FetchCities(countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle empty capital array
	capital := ""
	if len(countryInfo.Capital) > 0 {
		capital = countryInfo.Capital[0]
	}

	// Prepare response
	response := map[string]interface{}{
		"name":       countryInfo.Name.Common,
		"continents": countryInfo.Continents,
		"population": countryInfo.Population,
		"languages":  countryInfo.Languages,
		"borders":    countryInfo.Borders,
		"flag":       countryInfo.Flag,
		"capital":    capital, // Use the first capital if available
		"cities":     cities,
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PopulationHandler handles requests to the population endpoint, where you get the population levels for a country.
func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	// Extract country code from URL
	countryCode := extractCountryCode(r.URL.Path)
	if countryCode == "" {
		http.Error(w, "400 Bad Request: No ISO code specified. \n"+
			"Example usage: .../countryinfo/population/no if you want to see the population for Norway. \n"+
			"You can also optionally add for example ?limit=2000-2005 "+
			"if you only want the data for a specific year gap.", http.StatusBadRequest)
		return
	}

	// Extract optional year range from query parameters
	startYear, endYear, err := extractYearRange(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "400 Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch and process population data
	populationResponse, err := services.GetPopulationData(countryCode, startYear, endYear)
	if err != nil {
		handlePopulationError(w, err)
		return
	}

	// Send the response
	sendPopulationResponse(w, populationResponse)
}

// extractCountryCode extracts the country code from the URL path.
func extractCountryCode(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 5 && parts[4] != "" {
		return strings.ToUpper(parts[4])
	}
	return ""
}

// extractYearRange extracts and validates the year range from the query parameter.
func extractYearRange(queryLimit string) (int, int, error) {
	if queryLimit == "" {
		return -1, -1, nil
	}

	years := strings.Split(queryLimit, "-")
	if len(years) != 2 {
		return -1, -1, fmt.Errorf("invalid year range format")
	}

	startYear, err1 := strconv.Atoi(years[0])
	endYear, err2 := strconv.Atoi(years[1])
	if err1 != nil || err2 != nil || startYear > endYear {
		return -1, -1, fmt.Errorf("invalid year range format")
	}

	return startYear, endYear, nil
}

// handlePopulationError handles errors returned by the population data fetching process.
func handlePopulationError(w http.ResponseWriter, err error) {
	if strings.Contains(err.Error(), "404 Not Found") {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// sendPopulationResponse sends the population data as a JSON response.
func sendPopulationResponse(w http.ResponseWriter, populationResponse *models.PopulationResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(populationResponse); err != nil {
		log.Println("Error encoding JSON response:", err)
	}
}

// StatusHandler handles request for API status
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Check external API status
	countriesNowStatus := services.CheckCountriesNowAPI()
	restCountriesStatus := services.CheckRestCountriesAPI()

	// Calculate uptime
	uptime := time.Since(startTime).Seconds()

	// Prepare response
	response := map[string]interface{}{
		"countriesnowapi":  countriesNowStatus,
		"restcountriesapi": restCountriesStatus,
		"version":          "v1",
		"uptime":           uptime,
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
