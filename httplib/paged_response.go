package httplib

import (
	"encoding/json"
	"net/http"
)

type DataPaginator struct {
	Data  interface{} `json:"items"`
	Count int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

type pageResponse struct {
	data DataPaginator `json:"data"`
	code int           `json:"code"`
}

func ResponsePaged(w http.ResponseWriter, data DataPaginator) {
	resp := pageResponse{
		data,
		http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(resp)
	return
}
