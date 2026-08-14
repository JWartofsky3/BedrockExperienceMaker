package packs

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func WriteJSONError(w http.ResponseWriter, message string, status int) {
	WriteJSON(w, map[string]string{"error": message}, status)
}
