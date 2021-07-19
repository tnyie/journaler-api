package views

import "net/http"

func respondJSON(w http.ResponseWriter, encoded []byte, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(encoded)
}
