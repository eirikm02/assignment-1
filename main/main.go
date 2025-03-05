package main

import (
	"assignment-1/handlers" // Import handlers
	"fmt"
	"net/http"
)

func main() {
	// Define routes
	http.HandleFunc("/countryinfo/v1/info/", handlers.InfoHandler)
	http.HandleFunc("/countryinfo/v1/population/", handlers.PopulationHandler) // Register PopulationHandler
	http.HandleFunc("/countryinfo/v1/status/", handlers.StatusHandler)

	// Start the server
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
