package utils

import (
	"log"
	"os"
)

// DefaultPort Default port for the server
const DefaultPort = "8080"

// GetPort Get the port from the environment variable, or use the default port
func GetPort() string {
	// Get the PORT environment variable
	port := os.Getenv("PORT")

	// Use default Port variable if not provided
	if port == "" {
		log.Println("$PORT has not been set. Default: " + DefaultPort)
		port = DefaultPort
	}

	return port
}
