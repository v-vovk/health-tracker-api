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

	json.NewEncoder(w).Encode(response)
}

func ValidationErrors(w http.ResponseWriter, errs error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	validationErrors := errs.(validator.ValidationErrors)
	errors := make([]string, len(validationErrors))

	for i, err := range validationErrors {
		errors[i] = fmt.Sprintf("Field '%s' failed on the '%s' tag", err.Field(), err.Tag())
	}

	response := map[string]interface{}{
		"errors": errors,
	}

	json.NewEncoder(w).Encode(response)
}
