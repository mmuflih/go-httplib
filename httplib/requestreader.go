package httplib

import "net/http"

type RequestReader interface {
	GetRouteParam(r *http.Request, name string) string
	GetJsonData(r *http.Request, data interface{}) (err error)
	GetQuery(r *http.Request, query string) string
}
