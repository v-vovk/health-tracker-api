package errors

import (
	"encoding/json"
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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status_code": statusCode,
		"message":     message,
	})
}
