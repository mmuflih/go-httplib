package httplib

import (
	"encoding/json"
	"net/http"
)

type DataResponse struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

func ResponseData(w http.ResponseWriter, data interface{}) {
	exception := DataResponse{
		data,
		http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		exception,
	)
	return
}
