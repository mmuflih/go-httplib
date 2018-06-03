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
)

func ResponseRequestException(w http.ResponseWriter, err []error, code int) {
	messages := ""
	for _, e := range err {
		messages += e.Error() + ","
	}
	messages = messages[:len(messages)-1]
	exception := ErrorResponse{
		messages,
		code,
		"Contact developer or administrator",
		code,
		messages,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(
		exception,
	)
	return
}
