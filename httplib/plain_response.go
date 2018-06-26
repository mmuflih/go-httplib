package httplib

import (
	"encoding/json"
	"net/http"
)

func ResponsePlain(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		data,
	)
	return
}
