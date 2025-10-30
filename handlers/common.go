package handlers

import (
	"encoding/json"
	"enterprise-architect-api/models"
	"net/http"
)

// respondWithJSON writes a JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// respondWithError writes an error response
func respondWithError(w http.ResponseWriter, statusCode int, message string, details string) {
	errorResponse := models.ErrorResponse{
		Error:   message,
		Message: details,
	}
	respondWithJSON(w, statusCode, errorResponse)
}

