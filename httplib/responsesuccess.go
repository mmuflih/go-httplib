package httplib

import (
	"encoding/json"
	"net/http"
)

func ResponseSuccess(w http.ResponseWriter, payload interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(
		payload,
	)
	return
}
