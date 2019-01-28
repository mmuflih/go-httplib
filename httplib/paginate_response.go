package httplib

import (
	"encoding/json"
	"net/http"
)

type DataPaginate struct {
	Data  interface{} 		`json:"data"`
	Count int         		`json:"total"`
	Page  int         		`json:"page"`
	Size  int         		`json:"size"`
	Code  int         		`json:"code"`
	Additional interface{}  `json:"additional,omitempty"`
}

func NewDataPaged(data interface{}, count int, page int, size int) *DataPaginate {
	return &DataPaginate {
		data,
		count,
		page,
		size,
		http.StatusOK,
		data,
	}
}

func ResponsePaginate(w http.ResponseWriter, data DataPaginate) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data.Code = http.StatusOK
	json.NewEncoder(w).Encode(data)
	return
}
