package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]string{
		"error": message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Failed to write JSON response: %v\n", err)
	}
}

func ValidationErrors(w http.ResponseWriter, errs error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	validationErrors, ok := errs.(validator.ValidationErrors)
	if !ok {
		fmt.Printf("Error type assertion failed for validation errors\n")
		return
	}

	errors := make([]string, len(validationErrors))

	for i, err := range validationErrors {
		errors[i] = fmt.Sprintf("Field '%s' failed on the '%s' tag", err.Field(), err.Tag())
	}

	response := map[string]interface{}{
		"errors": errors,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Failed to write JSON response: %v\n", err)
	}
}
