package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	// Encode before sending headers so encoding failures can return a 500.
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("Encode response: %v", err)
		status = http.StatusInternalServerError
		data = []byte(`{"error":"Internal server error"}`)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(append(data, '\n')); err != nil {
		log.Printf("Write response: %v", err)
	}
}
