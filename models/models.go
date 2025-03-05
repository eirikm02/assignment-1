package models

type PopulationData struct {
	Mean   int                `json:"mean"`
	Values []YearlyPopulation `json:"values"`
}

type YearlyPopulation struct {
	Year  int `json:"year"`
	Value int `json:"value"`
}

type CountryInfo struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Continents []string          `json:"continents"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flag"`
	Capital    []string          `json:"capital"` // Updated to handle array of strings
}

type PopulationYear struct {
	Year  int `json:"year"`
	Value int `json:"value"`
}

type CitiesResponse struct {
	Data []string `json:"data"`
}

type PopulationResponse struct {
	Data []PopulationYear `json:"data"`
}
