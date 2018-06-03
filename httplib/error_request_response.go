package httplib

/*
 * Created by M. Muflih Kholidin
 * Sun Jun 03 2018 11:30:22
 * mmuflic@gmail.com
 * https://github.com/mmuflih
 **/

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"
)

func ResponseRequestException(w http.ResponseWriter, err []error, code int) {
	exceptions := []ErrorResponse{}
	for _, e := range err {
		_, fn, line, _ := runtime.Caller(1)
		exception := ErrorResponse{
			e.Error() + " on " + fn + ":" + strconv.Itoa(line),
			code,
			"Contact developer or administrator",
			code,
			e.Error(),
		}
		exceptions = append(exceptions, exception)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(
		exceptions,
	)
	return
}
