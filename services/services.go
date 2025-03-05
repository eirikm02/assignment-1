package services

import (
	"assignment-1/config" // Import config
	"assignment-1/models" // Import models
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func FetchPopulationData(countryCode, limit string) ([]models.PopulationYear, error) {
	// Use the centralized URL from config
	url := config.CountriesNowAPI + "/countries/population/" + countryCode

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the raw response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("Raw response body:", string(body)) // Debug log

	// Decode the response into a PopulationResponse struct
	var populationResponse models.PopulationResponse
	if err := json.Unmarshal(body, &populationResponse); err != nil {
		return nil, err
	}

	// Debug log: Print the parsed response
	fmt.Println("Population response:", populationResponse)

	// Return the population data from the response
	return populationResponse.Data, nil
}

func CheckCountriesNowAPI() string {
	// Implement API status check logic
	return "OK"
}

func CheckRestCountriesAPI() string {
	// Implement API status check logic
	return "OK"
}
