package httplib

import (
	"encoding/json"
	"net/http"
)

type DataPaginator struct {
	Data  interface{} `json:"data"`
	Count int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Code  int         `json:"code"`
}

type pageResponse struct {
	Data DataPaginator `json:"data"`
	Code int           `json:"code"`
}

func ResponsePaged(w http.ResponseWriter, data DataPaginator) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data.Code = http.StatusOK
	json.NewEncoder(w).Encode(data)
	return
}
