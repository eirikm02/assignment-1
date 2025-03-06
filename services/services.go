package services

import (
	"assignment-1/config" // Import config
	"assignment-1/models" // Import models
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func FetchCountryInfo(countryCode string) (models.CountryInfo, error) {
	// Use the centralized URL from config
	url := config.RestCountriesAPI + "/alpha/" + countryCode
	resp, err := http.Get(url)
	if err != nil {
		return models.CountryInfo{}, err
	}
	defer resp.Body.Close()

	// Decode the response into a slice of CountryInfo
	var countries []models.CountryInfo
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return models.CountryInfo{}, err
	}

	// Check if the slice is empty
	if len(countries) == 0 {
		return models.CountryInfo{}, fmt.Errorf("no country found for code: %s", countryCode)
	}

	// Return the first country in the slice
	return countries[0], nil
}

func FetchCities(countryCode string) ([]string, error) {
	// Use the centralized URL from config
	url := config.CountriesNowAPI + "/countries/cities"

	// Create a request body with the country code
	requestBody := fmt.Sprintf(`{"country": "%s"}`, countryCode)
	resp, err := http.Post(url, "application/json", strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response into a CitiesResponse struct
	var citiesResponse models.CitiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&citiesResponse); err != nil {
		return nil, err
	}

	// Return the cities from the response
	return citiesResponse.Data, nil
}

// GetPopulationData processes and filters the population data based on year range and calculates the mean value.
func GetPopulationData(countryCode string, startYear int, endYear int) (*models.PopulationResponse, error) {
	// Convert ISO-2 to ISO-3
	iso3Code, err := FetchISO3FromISO2(countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ISO-2 to ISO-3: %v", err)
	}

	// Fetch population data using the ISO-3 code
	populationRecords, err := FetchPopulationData(iso3Code)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch population data: %v", err)
	}

	// Filter population data by year range
	filteredRecords, totalPopulation, count := filterPopulationData(populationRecords, startYear, endYear)
	if count == 0 {
		return nil, fmt.Errorf("404 Not Found: No population data available for the given year range")
	}

	// Calculate the mean population
	mean := totalPopulation / count
	return &models.PopulationResponse{
		Mean:   mean,
		Values: filteredRecords,
	}, nil
}

func FetchPopulationData(iso3Code string) ([]models.PopulationRecord, error) {
	url := fmt.Sprintf("%s/countries/population", config.CountriesNowAPI)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching population data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch population data: %d", resp.StatusCode)
	}

	var apiResponse models.CountriesNowPopulation
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding population data: %v", err)
	}

	// Iterate through the countries to find a match for the ISO-3 code
	for _, country := range apiResponse.Data {
		if strings.EqualFold(country.Iso3, iso3Code) {
			return country.PopulationCounts, nil
		}
	}
	return nil, fmt.Errorf("404 Not Found: Population data not found for country %s", iso3Code)
}

// filterPopulationData filters the population data based on the specified year range.
func filterPopulationData(populationRecords []models.PopulationRecord, startYear int, endYear int) ([]models.PopulationRecord, int, int) {
	var filteredRecords []models.PopulationRecord
	totalPopulation := 0
	count := 0

	for _, record := range populationRecords {
		if (startYear == -1 || record.Year >= startYear) && (endYear == -1 || record.Year <= endYear) {
			filteredRecords = append(filteredRecords, record)
			totalPopulation += record.Value
			count++
		}
	}

	return filteredRecords, totalPopulation, count
}

// CheckCountriesNowAPI checks the status of the CountriesNow API.
func CheckCountriesNowAPI() string {
	// Define the endpoint to check (e.g., the base URL or a specific endpoint)
	url := fmt.Sprintf("%s/countries", config.CountriesNowAPI)

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 5 * time.Second, // Set a timeout to avoid hanging
	}

	// Make a GET request to the API
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Sprintf("Unavailable: %v", err) // Return the error if the request fails
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		return "OK" // API is working
	}

	// If the status code is not 200, return the status
	return fmt.Sprintf("Unavailable: %s", resp.Status)
}

// CheckRestCountriesAPI checks the status of the REST Countries API.
func CheckRestCountriesAPI() string {
	// Define the endpoint to check (e.g., the base URL or a specific endpoint)
	url := fmt.Sprintf("%s/all", config.RestCountriesAPI)

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 5 * time.Second, // Set a timeout to avoid hanging
	}

	// Make a GET request to the API
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Sprintf("Unavailable: %v", err) // Return the error if the request fails
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		return "OK" // API is working
	}

	// If the status code is not 200, return the status
	return fmt.Sprintf("Unavailable: %s", resp.Status)
}
func FetchISO3FromISO2(iso2Code string) (string, error) {
	// Construct the URL for the REST Countries API
	url := fmt.Sprintf("%s/alpha/%s", config.RestCountriesAPI, iso2Code)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch country data: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch country data: %s", resp.Status)
	}

	// Decode the response into a slice of CountryInfo
	var countries []models.CountryInfo
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return "", fmt.Errorf("failed to decode country data: %v", err)
	}

	// Check if the slice is empty
	if len(countries) == 0 {
		return "", fmt.Errorf("no country found for code: %s", iso2Code)
	}

	// Extract the ISO-3 code from the first country in the slice
	iso3Code := countries[0].Cca3
	if iso3Code == "" {
		return "", fmt.Errorf("ISO-3 code not found for country: %s", iso2Code)
	}

	return iso3Code, nil
}
