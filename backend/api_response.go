package main

import (
	"encoding/json"
	"net/http"
)

// P is shortcut for map[string]any
type P map[string]any

// WriteJSON writes json response `data` to `w` with `status` status code.
// Returns `error` if either marshaling or writing returns an error.
func WriteJSON(w http.ResponseWriter, status int, data P) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}
