package ai

import (
	"os"
)

// GetAPIKey retrieves the Groq API key from environment variable or returns default
func GetAPIKey() string {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		// Default API key (should be set via environment variable in production)
		return ""
	}
	return apiKey
}
