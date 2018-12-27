package httplib

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type muxRequestReader struct{}

func NewMuxRequestReader() RequestReader {
	return &muxRequestReader{}
}

func (rr *muxRequestReader) GetRouteParam(r *http.Request, name string) string {
	return mux.Vars(r)[name]
}

func (rr *muxRequestReader) GetJsonData(r *http.Request, data interface{}) (err error) {
	err = json.NewDecoder(r.Body).Decode(data)
	return
}

func (rr *muxRequestReader) GetQuery(r *http.Request, key string) string {
	qs := r.URL.Query()
	return qs.Get(key)
}
