package httplib

/*
 * Created by M. Muflih Kholidin
 * Sun Jun 03 2018 11:30:22
 * mmuflic@gmail.com
 * https://github.com/mmuflih
 **/

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

type ErrorResponse struct {
	DeveloperMessage string `json:"developer_message"`
	ErrorCode        int    `json:"error_code"`
	MoreInfo         string `json:"more_info"`
	Status           int    `json:"status"`
	UserMessage      string `json:"user_message"`
}

func ResponseException(w http.ResponseWriter, err error, code int) {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("[error] %s:%d %v on %s", fn, line, err, pc)
	exception := ErrorResponse{
		err.Error() + " on " + fn + ":" + strconv.Itoa(line),
		code,
		"Contact developer or administrator",
		code,
		err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(
		exception,
	)
	return
}
