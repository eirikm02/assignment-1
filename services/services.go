package services

import (
	"assignment-1/config" // Importing configuration package
	"assignment-1/models" // Importing models package
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// FetchCountryInfo retrieves country details using the country code from the REST Countries API.
func FetchCountryInfo(countryCode string) (models.CountryInfo, error) {
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

	// Ensure the response contains data
	if len(countries) == 0 {
		return models.CountryInfo{}, fmt.Errorf("no country found for code: %s", countryCode)
	}

	return countries[0], nil // Return the first matched country
}

// FetchCities retrieves a list of cities for a given country code from the CountriesNow API.
func FetchCities(countryCode string) ([]string, error) {
	url := config.CountriesNowAPI + "/countries/cities"

	// Prepare the request body with the country code
	requestBody := fmt.Sprintf(`{"country": "%s"}`, countryCode)
	resp, err := http.Post(url, "application/json", strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response into CitiesResponse structure
	var citiesResponse models.CitiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&citiesResponse); err != nil {
		return nil, err
	}

	return citiesResponse.Data, nil
}

// GetPopulationData fetches and processes population data for a country based on a year range.
func GetPopulationData(countryCode string, startYear int, endYear int) (*models.PopulationResponse, error) {
	// Convert ISO-2 code to ISO-3
	iso3Code, err := FetchISO3FromISO2(countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ISO-2 to ISO-3: %v", err)
	}

	// Fetch population data for the given ISO-3 country code
	populationRecords, err := FetchPopulationData(iso3Code)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch population data: %v", err)
	}

	// Filter records based on year range
	filteredRecords, totalPopulation, count := filterPopulationData(populationRecords, startYear, endYear)
	if count == 0 {
		return nil, fmt.Errorf("404 Not Found: No population data available for the given year range")
	}

	// Calculate mean population
	mean := totalPopulation / count
	return &models.PopulationResponse{
		Mean:   mean,
		Values: filteredRecords,
	}, nil
}

// FetchPopulationData retrieves population data from the CountriesNow API for a given ISO-3 country code.
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

	// Decode API response
	var apiResponse models.CountriesNowPopulation
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding population data: %v", err)
	}

	// Match the requested ISO-3 country code
	for _, country := range apiResponse.Data {
		if strings.EqualFold(country.Iso3, iso3Code) {
			return country.PopulationCounts, nil
		}
	}
	return nil, fmt.Errorf("404 Not Found: Population data not found for country %s", iso3Code)
}

// filterPopulationData filters population data based on a year range and calculates total population.
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

// CheckCountriesNowAPI checks the availability of the CountriesNow API.
func CheckCountriesNowAPI() string {
	url := fmt.Sprintf("%s/countries", config.CountriesNowAPI)

	client := &http.Client{Timeout: 5 * time.Second} // Set timeout to avoid hanging requests
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Sprintf("Unavailable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "OK"
	}
	return fmt.Sprintf("Unavailable: %s", resp.Status)
}

// CheckRestCountriesAPI checks the availability of the REST Countries API.
func CheckRestCountriesAPI() string {
	url := fmt.Sprintf("%s/all", config.RestCountriesAPI)

	client := &http.Client{Timeout: 5 * time.Second} // Set timeout to avoid hanging requests
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Sprintf("Unavailable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "OK"
	}
	return fmt.Sprintf("Unavailable: %s", resp.Status)
}

// FetchISO3FromISO2 converts an ISO-2 country code to an ISO-3 code using the REST Countries API.
func FetchISO3FromISO2(iso2Code string) (string, error) {
	url := fmt.Sprintf("%s/alpha/%s", config.RestCountriesAPI, iso2Code)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch country data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch country data: %s", resp.Status)
	}

	// Decode response into a slice of CountryInfo
	var countries []models.CountryInfo
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return "", fmt.Errorf("failed to decode country data: %v", err)
	}

	if len(countries) == 0 {
		return "", fmt.Errorf("no country found for code: %s", iso2Code)
	}

	iso3Code := countries[0].Cca3
	if iso3Code == "" {
		return "", fmt.Errorf("ISO-3 code not found for country: %s", iso2Code)
	}

	return iso3Code, nil
}
