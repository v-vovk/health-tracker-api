package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

// HTTPError represents a structured API error response
type HTTPError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// WriteHTTPError writes an HTTPError to the response writer
func WriteHTTPError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"status_code": statusCode,
		"message":     message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error if encoding fails
		log.Printf("Failed to write HTTP error response: %v", err)
	}
}
