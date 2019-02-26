package httplib

import (
	"encoding/json"
	"net/http"
)

type DataPaginate struct {
	Data  interface{} 		`json:"data"`
	Additional interface{}  `json:"additional,omitempty"`
	Count int         		`json:"total"`
	Page  int         		`json:"page"`
	Size  int         		`json:"size"`
	Code  int         		`json:"code"`
}

func NewDataPaginate(data interface{}, count int, page int, size int) DataPaginate {
	dp := DataPaginate {
		data,
		nil,
		count,
		page,
		size,
		http.StatusOK,
	}
	dp.Code = http.StatusOK
	return dp
}

func NewNilDataPaginate() DataPaginate {
	dp := DataPaginate {
		"",
		nil,
		0,
		0,
		0,
		http.StatusInternalServerError,
	}
	dp.Code = http.StatusOK
	return dp
}

func ResponsePaginate(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
	return
}
