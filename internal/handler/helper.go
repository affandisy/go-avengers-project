package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", slog.Any("error", err))
	}
}

func writeError(w http.ResponseWriter, status int, message string, errors any) {
	writeJSON(w, status, Response{
		Message: message,
		Errors:  errors,
	})
}

func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := strings.ToLower(e.Field())

			switch e.Tag() {
			case "required":
				errors[field] = field + "is required"
			case "email":
				errors[field] = field + " must be a valid email address"
			case "min":
				errors[field] = field + " must be at least " + e.Param() + " characters"
			case "max":
				errors[field] = field + " must be at most " + e.Param() + " characters"
			case "gte":
				errors[field] = field + " must be greater than or equal to " + e.Param()
			case "lte":
				errors[field] = field + " must be less than or equal to " + e.Param()
			case "gt":
				errors[field] = field + " must be greater than " + e.Param()
			case "oneof":
				errors[field] = field + " must be one of: " + e.Param()
			default:
				errors[field] = field + "is invalid"
			}
		}
	}

	return errors
}
