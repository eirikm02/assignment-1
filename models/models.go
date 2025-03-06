package models

// PopulationResponse represents the final response structure for population data.
type PopulationResponse struct {
	Mean   int                `json:"mean"`   // Mean population value
	Values []PopulationRecord `json:"values"` // List of population records
}

// PopulationRecord represents the population data for a specific year.
type PopulationRecord struct {
	Year  int `json:"year"`  // Year of the population record
	Value int `json:"value"` // Population value for the year
}

// CountriesNowPopulation represents the structure of the population data from the CountriesNow API.
type CountriesNowPopulation struct {
	Error bool                    `json:"error"` // Indicates if there was an error
	Msg   string                  `json:"msg"`   // Message from the API
	Data  []CountryPopulationData `json:"data"`  // List of country population data
}

// CountryPopulationData represents the population data for a specific country.
type CountryPopulationData struct {
	Country          string             `json:"country"`          // Name of the country
	Code             string             `json:"code"`             // ISO-2 country code
	Iso3             string             `json:"iso3"`             // ISO-3 country code
	PopulationCounts []PopulationRecord `json:"populationCounts"` // List of population records for the country
}

// CountryInfo represents the structure of the country data from the REST Countries API.
type CountryInfo struct {
	Name struct {
		Common   string `json:"common"`   // Common name of the country
		Official string `json:"official"` // Official name of the country
	} `json:"name"`
	Cca2       string            `json:"cca2"`       // ISO-2 country code
	Cca3       string            `json:"cca3"`       // ISO-3 country code
	Population int               `json:"population"` // Population of the country
	Languages  map[string]string `json:"languages"`  // Languages spoken in the country
	Continents []string          `json:"continents"` // Continents the country belongs to
	Borders    []string          `json:"borders"`    // Bordering countries
	Flag       string            `json:"flag"`       // URL to the country's flag
	Capital    []string          `json:"capital"`    // Capital cities of the country
}

// CombinedInfo represents the combined country and city information.
type CombinedInfo struct {
	Name       string            `json:"name"`       // Common name of the country
	Continent  string            `json:"continent"`  // Continent of the country
	Population int               `json:"population"` // Population of the country
	Languages  map[string]string `json:"languages"`  // Languages spoken in the country
	Borders    []string          `json:"borders"`    // Bordering countries
	Flag       string            `json:"flag"`       // URL to the country's flag
	Capital    string            `json:"capital"`    // Capital city of the country
	Cities     []string          `json:"cities"`     // List of cities in the country
}

// CitiesResponse represents the structure of the cities data from the CountriesNow API
type CitiesResponse struct {
	Data []string `json:"data"`
}
