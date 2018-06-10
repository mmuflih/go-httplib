package httplib

import (
	"encoding/json"
	"net/http"
)

type PagedResponse struct {
	Data interface{} `json:"data"`
	Page int         `json:"page"`
	Size int         `json:"size"`
	Code int         `json:"code"`
}

func ResponsePaged(w http.ResponseWriter, page int, size int, data interface{}) {
	exception := PagedResponse{
		data,
		page,
		size,
		http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(
		exception,
	)
	return
}
