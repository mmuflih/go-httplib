package httplib

import (
	"encoding/json"
	"net/http"
)

func ResponseError(w http.ResponseWriter, payload interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(
		payload,
	)
	return
}
