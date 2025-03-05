package handlers

import (
	"assignment-1/models"   // Import models
	"assignment-1/services" // Import services
	"encoding/json"
	"net/http"
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

// Handler for population data
func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	// Extract country code from URL
	countryCode := r.URL.Path[len("/countryinfo/v1/population/"):]
	if countryCode == "" {
		http.Error(w, "Missing country code", http.StatusBadRequest)
		return
	}

	// Mock population data
	response := models.PopulationData{
		Mean: 5044396,
		Values: []models.YearlyPopulation{
			{Year: 2010, Value: 4889252},
			{Year: 2011, Value: 4953088},
			{Year: 2012, Value: 5018573},
			{Year: 2013, Value: 5079623},
			{Year: 2014, Value: 5137232},
			{Year: 2015, Value: 5188607},
		},
	}

	// Return JSON response
	jsonResponse(w, response)
}

// Helper function to send JSON response
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

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

// Helper function to calculate mean population
func calculateMean(populationData []models.PopulationYear) float64 {
	sum := 0
	for _, data := range populationData {
		sum += data.Value
	}
	return float64(sum) / float64(len(populationData))
}
