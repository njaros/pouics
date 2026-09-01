package common

import (
	"encoding/json"
	"net/http"
)

// USUAL RESPONSE

// Common function to write any response
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// ERRORS

// Common function to write any error
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// Unauthorized répond avec un 401 au format JSON, cohérent avec les handlers.
func Unauthorized(w http.ResponseWriter) {
	WriteError(w, http.StatusUnauthorized, "unauthorized")
}

// InternalServerError répond avec un 500 au format JSON, cohérent avec les handlers.
func InternalServerError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, "internal server error")
}
